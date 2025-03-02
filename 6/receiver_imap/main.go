package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Config struct {
	SMTPHost      string `json:"smtpHost"`
	SMTPPort      string `json:"smtpPort"`
	IMAPServer    string `json:"imapServer"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	From          string `json:"from"`
	To            string `json:"to"`
	LabName       string `json:"labName"`
	SenderName    string `json:"senderName"`
	SenderGroup   string `json:"senderGroup"`
	ReceiverName  string `json:"receiverName"`
	ReceiverGroup string `json:"receiverGroup"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func fetchEmails(config *Config) error {
	c, err := client.DialTLS(config.IMAPServer, nil)
	if err != nil {
		return err
	}
	defer c.Logout()

	if err := c.Login(config.Username, config.Password); err != nil {
		return err
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return err
	}

	fmt.Printf("Количество писем в ящике: %d\n", mbox.Messages)

	seqset := new(imap.SeqSet)
	if mbox.Messages > 0 {
		var fromSeq uint32 = 1
		if mbox.Messages > 5 {
			fromSeq = mbox.Messages - 5 + 1
		}
		seqset.AddRange(fromSeq, mbox.Messages)
	}

	messages := make(chan *imap.Message, 5)
	go func() {
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("\nЗаголовки последних писем:")
	count := 0
	for msg := range messages {
		if msg.Envelope != nil {
			fmt.Printf("От: %v, Тема: %v, Дата: %v\n", msg.Envelope.From, msg.Envelope.Subject, msg.Envelope.Date)
		}
		count++
		if count >= 5 {
			break
		}
	}

	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("From", config.To)
	uids, err := c.Search(criteria)
	if err != nil {
		return err
	}

	if len(uids) == 0 {
		fmt.Println("\nПисьмо от одногруппника не найдено.")
	} else {
		fmt.Println("\nПисьмо от одногруппника найдено. Сохраняем его...")
		seqset = new(imap.SeqSet)
		seqset.AddNum(uids[0])
		section := &imap.BodySectionName{}
		msgChan := make(chan *imap.Message, 1)
		go func() {
			if err := c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, msgChan); err != nil {
				log.Fatal(err)
			}
		}()

		msg := <-msgChan
		if msg == nil {
			fmt.Println("Не удалось получить письмо.")
		} else {
			r := msg.GetBody(section)
			if r == nil {
				fmt.Println("Тело письма пустое.")
			} else {
				f, err := os.Create("classmate_email.eml")
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := io.Copy(f, r); err != nil {
					return err
				}
				fmt.Println("Письмо от одногруппника сохранено в файл classmate_email.eml")
			}
		}
	}
	return nil
}

func main() {
	config, err := loadConfig("../config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	if err := fetchEmails(config); err != nil {
		log.Fatal("Ошибка при получении писем:", err)
	}
}

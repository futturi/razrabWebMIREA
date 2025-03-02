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
	"github.com/emersion/go-message/mail"
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

func fetchEmailWithAttachment(config *Config) error {
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

	var seqStart uint32 = 1
	if mbox.Messages > 10 {
		seqStart = mbox.Messages - 10 + 1
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(seqStart, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	section := &imap.BodySectionName{}
	go func() {
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			continue
		}
		mr, err := mail.CreateReader(r)
		if err != nil {
			continue
		}
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				break
			}
			switch h := p.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				f, err := os.Create(filename)
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := io.Copy(f, p.Body); err != nil {
					return err
				}
				fmt.Printf("Вложение сохранено в файл: %s\n", filename)
				return nil
			}
		}
	}
	fmt.Println("Письмо с вложением не найдено среди последних писем.")
	return nil
}

func main() {
	config, err := loadConfig("../config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	if err := fetchEmailWithAttachment(config); err != nil {
		log.Fatal("Ошибка при получении вложения:", err)
	}
}

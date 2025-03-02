package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"time"
)

// Config описывает структуру конфигурационного файла.
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

func sendEmail(config *Config) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	now := time.Now().Format(time.RFC1123Z)
	subject := config.LabName

	msg := "From: " + config.From + "\r\n" +
		"To: " + config.To + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"X-Priority: 1\r\n" +
		"Importance: High\r\n" +
		"\r\n" +
		"Отправитель: " + config.SenderName + ", группа " + config.SenderGroup + "\r\n" +
		"Получатель: " + config.ReceiverName + ", группа " + config.ReceiverGroup + "\r\n" +
		"Время отправки: " + now + "\r\n"

	if err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.From, []string{config.To}, []byte(msg)); err != nil {
		return err
	}
	fmt.Println("Письмо отправлено успешно!")
	return nil
}

func main() {
	config, err := loadConfig("../config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	if err := sendEmail(config); err != nil {
		log.Fatal("Ошибка при отправке письма:", err)
	}
}

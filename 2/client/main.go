package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	ServerPort    int    `json:"server_port"`
	SendInterval  int    `json:"send_interval"`
	Name          string `json:"name"`
	Group         string `json:"group"`
	LogFile       string `json:"log_file"`
}

type LogEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	EventType     string    `json:"event_type"`
	Details       string    `json:"details"`
	ServerAddress string    `json:"server_address"`
}

func main() {
	config, err := loadConfig("./2/client/config.json")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}

	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка при открытии лог-файла: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "CLIENT: ", log.Ldate|log.Ltime|log.Lshortfile)

	serverAddress := fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к серверу %s: %v", serverAddress, err)
	}
	defer conn.Close()

	logEvent(logger, "connect", "Успешное подключение", config.ServerAddress)

	ticker := time.NewTicker(time.Duration(config.SendInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		message := fmt.Sprintf("%s %s", config.Name, config.Group)

		_, err := fmt.Fprintf(conn, message+"\n")
		if err != nil {
			logEvent(logger, "error", fmt.Sprintf("Ошибка при отправке сообщения: %v", err), config.ServerAddress)
			continue
		}
		logEvent(logger, "send", fmt.Sprintf("Отправлено сообщение: %s", message), config.ServerAddress)

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logEvent(logger, "error", fmt.Sprintf("Ошибка при получении ответа: %v", err), config.ServerAddress)
			continue
		}
		logEvent(logger, "receive", fmt.Sprintf("Получено сообщение: %s", strings.TrimSpace(response)), config.ServerAddress)
	}
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func logEvent(logger *log.Logger, eventType string, details string, serverAddress string) {
	entry := LogEntry{
		Timestamp:     time.Now(),
		EventType:     eventType,
		Details:       details,
		ServerAddress: serverAddress,
	}

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("Ошибка при сериализации в JSON: %v", err)
		logger.Printf("%s: %s - %s - %s", entry.Timestamp.Format(time.RFC3339), eventType, details, serverAddress)
		return
	}

	logger.Printf("%s\n", string(jsonEntry))
}

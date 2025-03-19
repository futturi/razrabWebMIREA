package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("./2/server/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	serverStart := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Сервер запущен: %s", serverStart)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен и слушает порт 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при подключении: %v", err)
			continue
		}

		log.Printf("Клиент подключился: %s", time.Now().Format("2006-01-02 15:04:05"))
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		log.Printf("Клиент отключился: %s", time.Now().Format("2006-01-02 15:04:05"))
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка при чтении: %v", err)
			return
		}
		message = message[:len(message)-1]

		log.Printf("Получено сообщение в %s: %s", time.Now().Format("2006-01-02 15:04:05"), message)

		response := fmt.Sprintf("Привет, %s!", message)
		_, err = fmt.Fprintf(conn, response+"\n")
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			return
		}
		log.Printf("Отправлено сообщение в %s: %s", time.Now().Format("2006-01-02 15:04:05"), response)
	}
}

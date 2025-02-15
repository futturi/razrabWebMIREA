package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка при чтении: %v", err)
			return
		}
		message = message[:len(message)-1]

		response := fmt.Sprintf("Привет, %s!", message)
		fmt.Fprintf(conn, response+"\n")
		log.Printf("Получено: %s, Отправлено: %s", message, response)
	}
}

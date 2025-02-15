package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Многопроцессный сервер запущен и слушает порт 8080...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Printf("Получен сигнал: %v. Завершение работы...", sig)
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при подключении: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			pid, err := forkProcess(conn)
			if err != nil {
				log.Printf("Ошибка при создании дочернего процесса: %v", err)
				return
			}
			log.Printf("Создан дочерний процесс с PID: %d", pid)

			handleConnection(conn)
		}(conn)
	}
}

func forkProcess(conn net.Conn) (int, error) {
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess(os.Args[0], os.Args, attr)
	if err != nil {
		return 0, err
	}

	return process.Pid, nil
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка при чтении: %v", err)
			return
		}
		message = message[:len(message)-1]

		response := fmt.Sprintf("Привет от процесса %d, %s!", os.Getpid(), message)
		fmt.Fprintf(conn, response+"\n")
		log.Printf("Процесс %d, Получено: %s, Отправлено: %s", os.Getpid(), message, response)
	}
}

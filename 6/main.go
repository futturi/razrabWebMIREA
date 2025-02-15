package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Многопроцессный web-сервер запущен и слушает порт 8080...")

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
		request, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Ошибка при чтении запроса: %v", err)
			return
		}

		request = strings.TrimSpace(request)

		if request != "" {
			response := generateHTMLResponse("Кручинин Иван", getPrimes(100)) // Замените на свои ФИО
			conn.Write([]byte(response))
			log.Printf("Процесс %d, Отправлен ответ клиенту", os.Getpid())
		}

	}
}

func generateHTMLResponse(name string, primes []int) string {
	html := `<!DOCTYPE html>
<html>
<head>
<title>Простые числа</title>
</head>
<body>
<h1>ФИО: ` + name + `</h1>
<h2>Простые числа от 0 до 100:</h2>
<p>` + strings.Join(strings.Split(fmt.Sprint(primes), " "), ", ") + `</p>
</body>
</html>`
	return "HTTP/1.1 200 OK\r\nContent-Type: text/html; charset=utf-8\r\nConnection: close\r\n\r\n" + html
}

func getPrimes(limit int) []int {
	primes := []int{}
	for i := 2; i <= limit; i++ {
		isPrime := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, i)
		}
	}
	return primes
}

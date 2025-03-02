package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении из соединения:", err)
		return
	}
	html := generateHTML()

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n", len(html)) +
		"\r\n" + html

	conn.Write([]byte(response))
}

func generateHTML() string {
	studentName := "Кручинин Иван Юрьевич"
	primes := getPrimes(0, 100)

	html := "<html><head><title>Простые числа</title></head><body>"
	html += "<h1>ФИО студента: " + studentName + "</h1>"
	html += "<h2>Простые числа от 0 до 100:</h2><ul>"
	for _, prime := range primes {
		html += fmt.Sprintf("<li>%d</li>", prime)
	}
	html += "</ul></body></html>"
	return html
}

func getPrimes(start, end int) []int {
	var primes []int
	for i := start; i <= end; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

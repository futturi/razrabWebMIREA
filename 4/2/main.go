package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "www.yandex.ru:80")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		os.Exit(1)
	}
	defer conn.Close()

	request := "GET / HTTP/1.1\r\nHost: www.yandex.ru\r\nConnection: close\r\n\r\n"
	fmt.Fprintf(conn, request)

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "set-cookie:") || strings.HasPrefix(line, "Set-Cookie:") {
			fmt.Println("Cookie:", line)
			continue
		}
		fmt.Print(line)
	}
}

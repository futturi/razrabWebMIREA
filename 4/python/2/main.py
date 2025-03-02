import socket

def main():
    host = "www.yandex.ru"
    port = 80

    try:
        conn = socket.create_connection((host, port))
    except Exception as e:
        print("Ошибка подключения:", e)
        exit(1)
    request = (
        "GET / HTTP/1.1\r\n"
        "Host: www.yandex.ru\r\n"
        "Connection: close\r\n"
        "\r\n"
    )
    conn.sendall(request.encode('utf-8'))

    with conn.makefile("r", encoding="utf-8") as response:
        for line in response:
            if line.lower().startswith("set-cookie:"):
                print("Cookie:", line.strip())
                continue
            print(line, end="")

    conn.close()

if __name__ == "__main__":
    main()

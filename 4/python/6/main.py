import socket
import threading


def is_prime(n):
    if n < 2:
        return False
    for i in range(2, int(n ** 0.5) + 1):
        if n % i == 0:
            return False
    return True


def get_primes(start, end):
    return [i for i in range(start, end + 1) if is_prime(i)]


def generate_html():
    student_name = "Кручинин Иван Юрьевич" #todo тут менять фио
    primes = get_primes(0, 100)

    html = "<html><head><title>Простые числа</title></head><body>"
    html += f"<h1>ФИО студента: {student_name}</h1>"
    html += "<h2>Простые числа от 0 до 100:</h2><ul>"
    for prime in primes:
        html += f"<li>{prime}</li>"
    html += "</ul></body></html>"
    return html


def handle_connection(conn, addr):
    print(f"Подключение от {addr}")
    try:
        conn.recv(1024)

        html = generate_html()
        response = (
                "HTTP/1.1 200 OK\r\n"
                "Content-Type: text/html; charset=utf-8\r\n"
                f"Content-Length: {len(html)}\r\n"
                "\r\n" + html
        )
        conn.sendall(response.encode('utf-8'))
    except Exception as e:
        print("Ошибка при обработке соединения:", e)
    finally:
        conn.close()


def main():
    host = ""
    port = 8080
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((host, port))
    server_socket.listen(5)
    print(f"Сервер запущен на порту {port}...")

    while True:
        conn, addr = server_socket.accept()
        threading.Thread(target=handle_connection, args=(conn, addr), daemon=True).start()


if __name__ == "__main__":
    main()

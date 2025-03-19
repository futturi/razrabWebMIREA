import socket
import threading
import logging
from datetime import datetime

logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s: %(message)s', handlers=[logging.FileHandler('server.log', encoding="utf-8", mode="a")])

def handle_connection(conn, addr):
    logging.info("Клиент подключился: %s", addr)
    with conn:
        sock_file = conn.makefile('rwb')
        while True:
            try:
                line = sock_file.readline()
                if not line:
                    logging.info("Клиент отключился: %s", addr)
                    break

                message = line.decode('utf-8').rstrip('\n')
                logging.info("Получено сообщение в %s от %s: %s", datetime.now().strftime('%Y-%m-%d %H:%M:%S'), addr, message)

                response = f"Привет, {message}!\n"
                sock_file.write(response.encode('utf-8'))
                sock_file.flush()

                logging.info("Отправлено сообщение в %s к %s: %s", datetime.now().strftime('%Y-%m-%d %H:%M:%S'), addr, response.strip())
            except Exception as e:
                logging.error("Ошибка при работе с клиентом %s: %s", addr, e)
                break

def main():
    host = ''
    port = 8080

    logging.info("Сервер запущен: %s", datetime.now().strftime('%Y-%m-%d %H:%M:%S'))

    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((host, port))
    server_socket.listen(5)
    logging.info("Сервер слушает порт %d...", port)

    while True:
        try:
            conn, addr = server_socket.accept()
            threading.Thread(target=handle_connection, args=(conn, addr), daemon=True).start()
        except Exception as e:
            logging.error("Ошибка при подключении: %s", e)

if __name__ == "__main__":
    main()

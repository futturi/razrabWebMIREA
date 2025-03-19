import socket
import logging
from multiprocessing import Process

def handle_connection(conn, addr):
    with conn:
        logging.info("Новое соединение от %s", addr)
        sock_file = conn.makefile('rwb')
        while True:
            try:
                line = sock_file.readline()
                if not line:
                    logging.info("Соединение закрыто от %s", addr)
                    break
                message = line.decode('utf-8').rstrip('\n')
                response = f"Привет, {message}!\n"
                sock_file.write(response.encode('utf-8'))
                sock_file.flush()
                logging.info("Получено: %s, Отправлено: %s", message, response.strip())
            except Exception as e:
                logging.error("Ошибка при обработке соединения %s: %s", addr, e)
                break

def main():
    logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s: %(message)s', handlers=[logging.FileHandler('server.log', encoding="utf-8", mode="a")])
    host = ''
    port = 8080
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((host, port))
    server_socket.listen(5)
    logging.info("Многопроцессный сервер запущен на порту %d", port)

    while True:
        try:
            conn, addr = server_socket.accept()
            process = Process(target=handle_connection, args=(conn, addr))
            process.daemon = True
            process.start()
            conn.close()
        except Exception as e:
            logging.error("Ошибка при принятии соединения: %s", e)

if __name__ == '__main__':
    main()

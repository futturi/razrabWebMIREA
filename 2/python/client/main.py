import json
import socket
import time
import logging

def load_config(filename):
    with open(filename, "r", encoding="utf-8") as f:
        return json.load(f)

def log_event(logger, event_type, details, server_address):
    entry = {
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%S", time.localtime()),
        "event_type": event_type,
        "details": details,
        "server_address": server_address
    }
    try:
        json_entry = json.dumps(entry, ensure_ascii=False)
        logger.info(json_entry)
    except Exception as e:
        logger.error("Ошибка при сериализации в JSON: %s", e)
        logger.info(f"{entry['timestamp']}: {event_type} - {details} - {server_address}")

def main():
    config_path = "config.json"
    try:
        config = load_config(config_path)
    except Exception as e:
        print(f"Ошибка при чтении конфигурации: {e}")
        return

    log_file = config.get("log_file", "client.log")
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s: %(message)s",
        handlers=[logging.FileHandler(log_file, encoding="utf-8", mode="a")]
    )
    logger = logging.getLogger("CLIENT")

    server_address = config["server_address"]
    server_port = config["server_port"]

    try:
        sock = socket.create_connection((server_address, server_port))
    except Exception as e:
        logger.error("Ошибка при подключении к серверу %s:%s: %s", server_address, server_port, e)
        return

    log_event(logger, "connect", "Успешное подключение", server_address)
    send_interval = config["send_interval"]

    while True:
        message = f"{config['name']} {config['group']}\n"
        logger.info("Отправляем: " + message.strip())
        try:
            sock.sendall(message.encode('utf-8'))
            log_event(logger, "send", f"Отправлено сообщение: {message.strip()}", server_address)
        except Exception as e:
            log_event(logger, "error", f"Ошибка при отправке сообщения: {e}", server_address)
            time.sleep(send_interval)
            continue

        try:
            response_bytes = b""
            while not response_bytes.endswith(b"\n"):
                chunk = sock.recv(1024)
                if not chunk:
                    raise Exception("Нет ответа от сервера")
                response_bytes += chunk
            response = response_bytes.decode('utf-8').strip()
            log_event(logger, "receive", f"Получено сообщение: {response}", server_address)
        except Exception as e:
            log_event(logger, "error", f"Ошибка при получении ответа: {e}", server_address)
        time.sleep(send_interval)

if __name__ == "__main__":
    main()

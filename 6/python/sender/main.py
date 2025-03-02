import json
import smtplib
import ssl
import datetime
from email.message import EmailMessage

def load_config(filename):
    with open(filename, "r", encoding="utf-8") as f:
        return json.load(f)

def send_email(config):
    msg = EmailMessage()
    msg["From"] = config["from"]
    msg["To"] = config["to"]
    msg["Subject"] = config["labName"]
    msg["X-Priority"] = "1"
    msg["Importance"] = "High"

    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    body = (
        f"Отправитель: {config['senderName']}, группа {config['senderGroup']}\n"
        f"Получатель: {config['receiverName']}, группа {config['receiverGroup']}\n"
        f"Время отправки: {now}\n"
    )
    msg.set_content(body)

    smtp_host = config["smtpHost"]
    smtp_port = int(config["smtpPort"])
    username = config["username"]
    password = config["password"]

    context = ssl.create_default_context()
    with smtplib.SMTP(smtp_host, smtp_port) as server:
        server.starttls(context=context)
        server.login(username, password)
        server.send_message(msg)
    print("Письмо отправлено успешно!")

if __name__ == '__main__':
    config = load_config("../../config.json")
    send_email(config)

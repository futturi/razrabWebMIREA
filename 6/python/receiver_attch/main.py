import json
import imaplib
import email
import os

def load_config(filename="../../config.json"):
    with open(filename, "r", encoding="utf-8") as f:
        return json.load(f)

def fetch_attachment(config):
    imap_server = config["imapServer"]
    username = config["username"]
    password = config["password"]

    if ":" in imap_server:
        host, port = imap_server.split(":")
        port = int(port)
    else:
        host = imap_server
        port = 993

    mail_server = imaplib.IMAP4_SSL(host, port)
    mail_server.login(username, password)
    mail_server.select("INBOX")

    typ, data = mail_server.search(None, "ALL")
    if typ != "OK":
        print("Ошибка поиска писем")
        return

    msg_nums = data[0].split()
    last_ten = msg_nums[-10:] if len(msg_nums) >= 10 else msg_nums

    attachment_found = False
    for num in last_ten:
        typ, msg_data = mail_server.fetch(num, "(RFC822)")
        if typ != "OK":
            print("Ошибка получения письма", num)
            continue
        msg = email.message_from_bytes(msg_data[0][1])
        for part in msg.walk():
            content_disposition = part.get("Content-Disposition", "")
            if "attachment" in content_disposition.lower():
                filename = part.get_filename()
                if filename:
                    with open(filename, "wb") as f:
                        f.write(part.get_payload(decode=True))
                    print(f"Вложение сохранено в файл: {filename}")
                    attachment_found = True
                    break
        if attachment_found:
            break

    if not attachment_found:
        print("Письмо с вложением не найдено среди последних писем.")
    mail_server.logout()

if __name__ == '__main__':
    config = load_config()
    fetch_attachment(config)

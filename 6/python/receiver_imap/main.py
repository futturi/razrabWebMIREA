import json
import imaplib
import email
import os
from email.header import decode_header

def load_config(filename="../../config.json"):
    with open(filename, "r", encoding="utf-8") as f:
        return json.load(f)

def decode_mime_words(s):
    if s:
        decoded_fragments = decode_header(s)
        return ''.join(
            fragment.decode(encoding if encoding else 'utf-8', errors='replace')
            if isinstance(fragment, bytes) else fragment
            for fragment, encoding in decoded_fragments
        )
    return s

def fetch_emails(config):
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
    print(f"Количество писем в ящике: {len(msg_nums)}")

    last_five = msg_nums[-5:] if len(msg_nums) >= 5 else msg_nums
    print("\nЗаголовки последних писем:")
    for num in last_five:
        typ, msg_data = mail_server.fetch(num, "(BODY.PEEK[HEADER])")
        if typ != "OK":
            print("Ошибка получения письма", num)
            continue
        msg = email.message_from_bytes(msg_data[0][1])
        from_header = decode_mime_words(msg.get("From"))
        subject_header = decode_mime_words(msg.get("Subject"))
        date_header = decode_mime_words(msg.get("Date"))
        print("От:", from_header, ", Тема:", subject_header, ", Дата:", date_header)

    search_criteria = f'(FROM "{config["to"]}")'
    typ, msg_nums = mail_server.search(None, search_criteria)
    if typ != "OK" or not msg_nums[0]:
        print("\nПисьмо от одногруппника не найдено.")
    else:
        msg_num = msg_nums[0].split()[0]
        typ, msg_data = mail_server.fetch(msg_num, "(RFC822)")
        if typ != "OK":
            print("Ошибка получения письма от одногруппника")
    mail_server.logout()

if __name__ == '__main__':
    config = load_config()
    fetch_emails(config)

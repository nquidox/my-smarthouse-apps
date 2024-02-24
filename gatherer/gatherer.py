import ast
import os
import sqlite3
import time
import paho.mqtt.subscribe as subscribe

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DB = os.path.join(BASE_DIR, 'database.db')

broker = "hostname"
topic = "topic"


def db_worker(op: str, sql: str, values: tuple = None) -> None:
    db = sqlite3.Connection(database=DB)
    cursor = db.cursor()

    match op:
        case "init":
            cursor.execute(sql)

        case "insert":
            cursor.execute(sql, values)

    db.commit()
    cursor.close()
    db.close()


def get_data_from_broker() -> list:
    try:
        msg = subscribe.simple(topic, hostname=broker)
        if msg is not None:
            return ast.literal_eval(msg.payload.decode("utf-8"))

    except Exception as e:
        pass


def main() -> None:
    while True:
        data = get_data_from_broker()
        for item in data:
            sql = '''INSERT INTO sensors(hex_id, temp_value, seconds, date_time) VALUES (?, ?, ?, ?)'''
            values = (item["id"], item["temp"], time.time(), time.strftime("%d-%m-%Y %H:%M:%S", time.localtime()))
            db_worker("insert", sql, values)
        time.sleep(1)


if __name__ == "__main__":
    db_worker("init", '''CREATE TABLE IF NOT EXISTS sensors(
                "id" INTEGER,
                "hex_id" TEXT,
                "temp_value" REAL,
                "seconds" REAL,
                "date_time" TEXT,
                PRIMARY KEY ("id" AUTOINCREMENT)
            )''')
    main()

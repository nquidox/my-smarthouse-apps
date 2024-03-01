import ast
import time
import psycopg2
import paho.mqtt.subscribe as subscribe

broker = "192.168.0.201"
topic = "esp32/banya"


def get_data_from_broker() -> list:
    try:
        msg = subscribe.simple(topic, hostname=broker)
        if msg is not None:
            return ast.literal_eval(msg.payload.decode("utf-8"))

    except Exception as e:
        pass


def main():
    conn = psycopg2.connect(
        database="",
        user="",
        password="",
        host="",
        port="5432"
    )
    cursor = conn.cursor()

    while True:
        for item in get_data_from_broker():
            cursor.execute(
                "INSERT INTO bathhouse_sensors(hex_id, temp_value, seconds, date_time) VALUES (%s, %s, %s, %s)",
                (item["id"], item["temp"], time.time(), time.strftime("%d-%m-%Y %H:%M:%S", time.localtime()))
            )
        conn.commit()
        time.sleep(1)

    cursor.close()
    conn.close()

if __name__ == "__main__":
    main()


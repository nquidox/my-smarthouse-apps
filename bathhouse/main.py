import time
import machine
import onewire, ds18x20
from umqtt.robust import MQTTClient

esp_name = ""
broker = ""
topic = "" 
client = MQTTClient(esp_name, broker)
client.connect()

dat = machine.Pin(12)
ds = ds18x20.DS18X20(onewire.OneWire(dat))
roms = ds.scan()

def send_data():
    res = dict()
    ds.convert_temp()
    for rom in roms:
        temp = ds.read_temp(rom)
        res.update({"id": rom.hex(), "temp": temp})
        client.publish(topic, str(res))

while True:
    send_data()
    time.sleep(1)
    
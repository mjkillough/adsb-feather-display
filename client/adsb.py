import machine
import ujson
import network
import time

import textwrap
import ssd1306
import uwebsockets.client

WIFI_SSID = # FILL IN
WIFI_PASS = # FILL IN

WS_URL = 'ws://192.168.86.21:8080/ws'

OLED_WIDTH = 14
OLED_LINE_HEIGHT = 10

BUTTON_GPIO_PIN = 2

def wifi_connect(ssid, password):
    sta_if = network.WLAN(network.STA_IF)
    if not sta_if.isconnected():
        sta_if.active(True)
        sta_if.connect(ssid, password)
        while not sta_if.isconnected():
            pass

class OLED:
    def __init__(self, width=OLED_WIDTH):
        i2c = machine.I2C(-1, machine.Pin(5), machine.Pin(4))
        self.oled = ssd1306.SSD1306_I2C(128, 32, i2c)
        self.width = width

    def wrap(self, text, column):
        count = 0
        output = ''
        for part in text.split():
            if count > 0 and (count + len(part)) > column:
                output += '\n'
                count = 0
            output += part + ' '
            count += len(part)
        return output

    def show_text(self, text):
        if isinstance(text, str):
            text = self.wrap(text, self.width).split('\n')

        self.oled.fill(0)
        y = 0
        for line in text:
            print(line)
            self.oled.text(line, 0, y)
            y += OLED_LINE_HEIGHT
        self.oled.show()

class Button:
    def __init__(self, pin):
        self.pin = machine.Pin(pin, machine.Pin.IN)

    def pressed(self):
        return not bool(self.pin.value())

def display_aircraft(oled, aircraft):
    callsign = aircraft['Callsign']
    from_iata = aircraft['Route']['From']['Iata']
    to_iata = aircraft['Route']['To']['Iata']
    altitude = int(aircraft['Altitude'])
    velocity = (aircraft['Velocity'] * 60 * 60) // 1000 # m/s -> km/h

    oled.show_text([
        '{} to {}'.format(from_iata, to_iata),
        '{}m {} km/h'.format(altitude, velocity),
        '{}'.format(callsign),
    ])

def display_route(oled, route):
    airport = route['Name']
    country = route['Country']

    oled.show_text('{} ({})'.format(airport, country))


def report_error(oled, error):
    oled.show_text('Error: {}'.format(error))

def loop(oled, button):
    with uwebsockets.client.connect(WS_URL) as websocket:
        websocket.settimeout(1) # seconds

        aircraft = None
        displaying_route = False
        force = False

        while True:
            if button.pressed() and (not displaying_route or force):
                displaying_route = True
                if aircraft is not None:
                    display_route(oled, aircraft['Route']['From'])
            elif not button.pressed() and (displaying_route or force):
                displaying_route = False
                force = False
                if aircraft is not None:
                    display_aircraft(oled, aircraft)

            try:
                resp = ujson.loads(websocket.recv())
            except OSError:
                # Timeout
                pass
            except Exception as e:
                report_error(oled, e)
            else:
                if resp is None:
                    continue
                if not resp:
                    oled.show_text("The sky is empty!")
                    continue
                aircraft = resp[0]
                force = True

def main():
    oled = OLED()
    oled.show_text('Connecting to WiFi...')

    wifi_connect(WIFI_SSID, WIFI_PASS)

    oled.show_text('Waiting for data...')

    button = Button(BUTTON_GPIO_PIN)

    while True:
        try:
            loop(oled, button)
        except Exception as e:
            report_error(oled, e)
        time.sleep(1)

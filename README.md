# adsb-feather-display

A simple client (in MicroPython) and server (in Go) for showing overhead planes
on a Adafruit FEATHER (ESP8266) with FeatherWing OLED display.

The server uses the OpenSky API to find planes in the given bounding box, then
uses the VirtualRadar database (locally) to determine the route of the aircraft
from its callsign.

The server exposes accepts web socket connections from the client, and will
periodically poll the OpenSky API as long as there is a client connected. Any
overhead planes are sent to the client as JSON documents over the web socket
connection. An empty list indicates there are no planes overhead.

## Example

It looks like this on the FEATHER:

```
CHQ to LHR
1000m 300 km/h
BAW661
```

With additional information being shown about the origin airport available on a
button press:

```
Chania (Greece)
```

## Server

The server is a standard Go application and can be developed as usual with the
Go tools.

It requires the VirtualRadar database to be available to run, which can be
downloaded with:

```
make data
```

### Deployment

Building:

```
docker build -t adsb-server .
```

Running:

```
docker run -d -p 8080:8080 adsb-server
```

## Client

The client is a simple MicroPython application. Developing the client requires
MicroPython to be installed on the host machine.

The application can be deployed to the FEATHER with:

```
# assumes /dev/ttyUSB0
make deploy
```

## License

MIT

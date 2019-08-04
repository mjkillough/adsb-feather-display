package opensky

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const url = "https://opensky-network.org/api/states/all?lamin=%v&lomin=%v&lamax=%v&lomax=%v"

// State is a aircraft returned by the OpenSky API.
type State struct {
	// Unique ICAO 24-bit address of the transponder in hex string representation.
	Icao24 string
	// Callsign of the vehicle (8 chars). Nil if no callsign has been received.
	Callsign *string
	// Country name inferred from the ICAO 24-bit address.
	OriginCountry string
	// Unix timestamp (seconds) for the last position update. Nil if no position report was
	// received by OpenSky within the past 15s.
	TimePosition *uint64
	// Unix timestamp (seconds) for the last update in general. This field is updated for any new,
	// valid message received from the transponder.
	LastContact uint64
	// WGS-84 longitude in decimal degrees.
	Longitude *float64
	// WGS-84 latitude in decimal degrees.
	Latitude *float64
	// Barometric altitude in meters.
	BarometricAltitude *float64
	// Boolean value which indicates if the position was retrieved from a surface position report.
	OnGround bool
	// Velocity over ground in m/s.
	Velocity *float64
	// True track in decimal degrees clockwise from north (north=0°).
	TrueTrack *float64
	// Vertical rate in m/s. A positive value indicates that the airplane is climbing, a negative
	// value indicates that it descends.
	VerticalRate *float64
	// IDs of the receivers which contributed to this state vector. Nil if no filtering for
	// sensor was used in the request.
	Sensors []uint64
	// Geometric altitude in meters.
	GeometricAltitude *float64
	// The transponder code aka Squawk.
	Squak *string
	// Whether flight status indicates special purpose indicator.
	Spi bool
	// Origin of this state’s position: 0 = ADS-B, 1 = ASTERIX, 2 = MLAT
	PositionSource uint64
}

// UnmarshalJSON is a custom implementation to cope with the fact that the API returns an array
// rather than an object.
func (s *State) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&s.Icao24, &s.Callsign, &s.OriginCountry, &s.TimePosition, &s.LastContact, &s.Longitude,
		&s.Latitude, &s.BarometricAltitude, &s.OnGround, &s.Velocity, &s.TrueTrack,
		&s.VerticalRate, &s.Sensors, &s.GeometricAltitude, &s.Squak, &s.Spi, &s.PositionSource,
	}
	want := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if got := len(tmp); got != want {
		return errors.Errorf("opensky: wrong number of fields: %d != %d", got, want)
	}
	return nil

}

// Client is an OpenSky API client.
type Client struct {
	http                             http.Client
	minLat, maxLat, minLong, maxLong float64
}

// New creates a new OpenSky API client.
func New(minLat, maxLat, minLong, maxLong float64) Client {
	return Client{
		minLat:  minLat,
		maxLat:  maxLat,
		minLong: minLong,
		maxLong: maxLong,
	}
}

// Fetch fetches the current OpenSky states from the API.
func (c *Client) Fetch() ([]State, error) {
	resp, err := c.http.Get(fmt.Sprintf(url, c.minLat, c.minLong, c.maxLat, c.maxLong))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.Errorf("opensky: http %v (%v)", resp.StatusCode, resp.Status)
	}

	var apiResp struct {
		States []State
	}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&apiResp); err != nil {
		return nil, errors.WithStack(err)
	}

	return apiResp.States, nil
}

package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mjkillough/adsb-display/server/opensky"
	"github.com/mjkillough/adsb-display/server/virtualradar"
	"github.com/pkg/errors"
)

type Aircraft struct {
	Callsign string
	Route    virtualradar.Route
	Altitude *float64
	Velocity *float64
}

type Fetcher struct {
	Reporter     interface{ Report(error) }
	OpenSky      opensky.Client
	VirtualRadar virtualradar.Database
}

func (f *Fetcher) Fetch() ([]byte, error) {
	fmt.Println("Fetching")

	states, err := f.OpenSky.Fetch()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	aircraft := []Aircraft{}
	for _, state := range states {
		if state.Callsign == nil {
			continue
		}

		callsign := strings.TrimSpace(*state.Callsign)
		route, err := f.VirtualRadar.FindRoute(callsign)
		if err != nil {
			log.Printf("Missing route for callsign: %v\n", callsign)
			continue
		}

		aircraft = append(aircraft, Aircraft{
			callsign,
			route,
			state.GeometricAltitude,
			state.Velocity,
		})
	}

	bytes, err := json.Marshal(aircraft)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bytes, nil
}

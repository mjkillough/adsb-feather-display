package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mjkillough/adsb-display/server/fetcher"
	"github.com/mjkillough/adsb-display/server/opensky"
	"github.com/mjkillough/adsb-display/server/virtualradar"
	"github.com/mjkillough/adsb-display/server/ws"
)

var (
	addr = flag.String("addr", ":8080", "http service address")

	minLat  = flag.Float64("min-lat", 51.417680, "min lat of bounding box to listen for planes in")
	maxLat  = flag.Float64("max-lat", 51.495647, "max lat of bounding box to listen for planes in")
	minLong = flag.Float64("min-long", -0.233946, "min long of bounding box to listen for planes in")
	maxLong = flag.Float64("max-long", -0.102311, "max long of bounding box to listen for planes in")

	dbPath = flag.String("db", "data/StandingData.sqb", "path to VirtualRadar database")
)

// Reporter is a simple error-reporter that reports to stdout.
type Reporter struct{}

// Report reports an error.
func (r *Reporter) Report(err error) {
	log.Printf("error: %+v\n", err)
}

func main() {
	flag.Parse()

	reporter := &Reporter{}

	fmt.Printf(
		"Listening for planes in bounding box: (%v, %v) (%v, %v)\n",
		*minLat, *minLong, *maxLat, *maxLong,
	)
	fmt.Printf("Loading database: %v\n", *dbPath)

	client := opensky.New(*minLat, *maxLat, *minLong, *maxLong)
	db, err := virtualradar.New(*dbPath)
	if err != nil {
		reporter.Report(err)
	}
	defer db.Close()

	fetcher := &fetcher.Fetcher{
		Reporter:     reporter,
		OpenSky:      client,
		VirtualRadar: db,
	}
	pollInterval := 10 * time.Second
	hub := ws.New(fetcher, reporter, pollInterval)

	go hub.Run()

	fmt.Printf("Listening on %v\n", *addr)
	http.HandleFunc("/ws", hub.Handler)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		reporter.Report(err)
	}
}

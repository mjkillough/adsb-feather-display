package virtualradar

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Airport struct {
	Iata, Name, Country string
}

type Route struct {
	From, To Airport
}

type Database struct {
	sqlite *sql.DB
	stmt   *sql.Stmt
}

func New(path string) (Database, error) {
	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return Database{}, errors.WithStack(err)
	}

	stmt, err := sqlite.Prepare(`
		select
			FromAirportIata, FromAirportName, FromAirportCountry,
			ToAirportIata, ToAirportName, ToAirportCountry
		from RouteView
		where Callsign = ?
	`)
	if err != nil {
		return Database{}, errors.WithStack(err)
	}

	return Database{sqlite, stmt}, nil
}

func (d *Database) Close() {
	d.stmt.Close()
	d.sqlite.Close()
}

func (d *Database) FindRoute(callsign string) (Route, error) {
	var (
		FromAirportIata, FromAirportName, FromAirportCountry string
		ToAirportIata, ToAirportName, ToAirportCountry       string
	)
	err := d.stmt.QueryRow(callsign).Scan(
		&FromAirportIata, &FromAirportName, &FromAirportCountry,
		&ToAirportIata, &ToAirportName, &ToAirportCountry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Route{}, errors.Errorf("virtualradar: missing route for callsign: %v", callsign)
		}
		return Route{}, errors.WithStack(err)
	}

	return Route{
		From: Airport{
			Iata:    FromAirportIata,
			Name:    FromAirportName,
			Country: FromAirportCountry,
		},
		To: Airport{
			Iata:    ToAirportIata,
			Name:    ToAirportName,
			Country: ToAirportCountry,
		},
	}, nil
}

package firefox

import (
	"database/sql"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/adamdecaf/browsehist"
	"github.com/adamdecaf/browsehist/browser"
	_ "github.com/mattn/go-sqlite3"
)

var (
	home = os.Getenv("HOME")

	placesFilepaths = []string{
		filepath.Join(home, `/Library/Application\ Support/Firefox/Profiles/rrdlhe7o.default/places.sqlite`), // OSX
	}
)

type Firefox struct {
	browser.Browser
}

func (ff Firefox) List() ([]*browsehist.HistoryItem, error) {
	where, err := ff.findPlacesDB()
	if err != nil {
		return nil, err
	}
	if where == "" {
		return nil, errors.New("Unable to find places.sqlite file")
	}

	db, err := ff.driver(where)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Grab urls from places.sqlite file
	return ff.readHistoryItems(db)
}

func (ff Firefox) readHistoryItems(db *sql.DB) ([]*browsehist.HistoryItem, error) {
	rows, err := db.Query(`select url, last_visit_date from moz_places where last_visit_date is not null;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*browsehist.HistoryItem, 0)
	for rows.Next() {
		var where string // url
		var when int64   // last_visit_date
		err = rows.Scan(&where, &when)
		if err != nil {
			return nil, err
		}
		item, err := ff.parse(where, when)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (ff Firefox) parse(where string, when int64) (*browsehist.HistoryItem, error) {
	u, err := url.Parse(where)
	if err != nil {
		return nil, err
	}
	t := time.Unix(int64(when/1e6), 0).UTC() // throw away nsec
	return &browsehist.HistoryItem{
		Address:    *u,
		AccessTime: t,
	}, nil
}

func (ff Firefox) findPlacesDB() (string, error) {
	for i := range placesFilepaths {
		s, err := os.Stat(placesFilepaths[i])
		if err == nil && s.Size() > 0 {
			return placesFilepaths[i], nil
		}
	}
	return "", errors.New("No places.sqlite file found")
}

// Callers are expected to Close() the *sql.DB returned
func (ff Firefox) driver(where string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", where)
	if err != nil {
		return nil, err
	}
	return db, err
}

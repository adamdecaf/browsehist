package chrome

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

	historyFilepaths = []string{
		filepath.Join(home, `/Library/Application Support/Google/Chrome/Default/History`), // OSX
	}

	// Linux: /home/$USER/.config/google-chrome/
	// Linux: /home/$USER/.config/chromium/
	// Windows Vista (and Win 7): C:\Users\[USERNAME]\AppData\Local\Google\Chrome\
	// Windows XP: C:\Documents and Settings\[USERNAME]\Local Settings\Application Data\Google\Chrome\

	// midnight UTC of 1 January 1601
	jan11601 = time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC)
	// unix EPOCH
	jan11970 = time.Unix(0, 0)
)

// Docs
// - https://digital-forensics.sans.org/blog/2010/01/21/google-chrome-forensics/

type Chrome struct {
	browser.Browser
}

func (c Chrome) List() ([]*browsehist.HistoryItem, error) {
	where, err := c.findHistoryDB()
	if err != nil {
		return nil, err
	}
	if where == "" {
		return nil, errors.New("Unable to find Histiry file")
	}

	db, err := c.driver(where)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return c.readHistoryItems(db)
}

func (c Chrome) readHistoryItems(db *sql.DB) ([]*browsehist.HistoryItem, error) {
	rows, err := db.Query(`select url,last_visit_time from urls where last_visit_time is not null;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*browsehist.HistoryItem, 0)
	for rows.Next() {
		var where string // url
		var when int64   // last_visit_time
		err = rows.Scan(&where, &when)
		if err != nil {
			return nil, err
		}
		item, err := c.parse(where, when)
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

func (c Chrome) parse(where string, when int64) (*browsehist.HistoryItem, error) {
	u, err := url.Parse(where)
	if err != nil {
		return nil, err
	}

	// "timestamp in the visit table is formatted as the number of microseconds since midnight UTC of 1 January 1601"
	// https://digital-forensics.sans.org/blog/2010/01/21/google-chrome-forensics/
	t := time.Date(1601, time.January, 1, 0, 0, int(when/1e6), 0, time.UTC)

	return &browsehist.HistoryItem{
		Address:    *u,
		AccessTime: t,
	}, nil
}

func (c Chrome) findHistoryDB() (string, error) {
	for i := range historyFilepaths {
		s, err := os.Stat(historyFilepaths[i])
		if err == nil && s.Size() > 0 {
			return historyFilepaths[i], nil
		}
	}
	return "", errors.New("No History file found")
}

// Callers are expected to Close() the *sql.DB returned
func (c Chrome) driver(where string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", where)
	if err != nil {
		return nil, err
	}
	return db, err
}

package browser

import (
	"net/url"
	"time"
)

type Browser interface {
	// List returns the given history items from a browser
	List() ([]*HistoryItem, error)
}

type History struct {
	Items []HistoryItem
}

type HistoryItem struct {
	Address    url.URL
	AccessTime time.Time
}

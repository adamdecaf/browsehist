package browsehist

import (
	"net/url"
	"time"
)

type History struct {
	Items []HistoryItem
}

type HistoryItem struct {
	Address url.URL
	AccessTime time.Time
}

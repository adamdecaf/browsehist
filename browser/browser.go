package browser

import (
	"github.com/adamdecaf/browsehist"
)

type Browser interface {
	// List returns the given history items from a browser
	List() ([]*browsehist.HistoryItem, error)
}

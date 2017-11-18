package chrome

import (
	"github.com/adamdecaf/browsehist"
	"github.com/adamdecaf/browsehist/browser"
)

type Chrome struct {
	browser.Browser
}

func (ff Chrome) List() ([]*browsehist.HistoryItem, error) {
	return nil, nil
}

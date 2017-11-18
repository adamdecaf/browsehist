package firefox

import (
	"github.com/adamdecaf/browsehist"
	"github.com/adamdecaf/browsehist/browser"
)

type Firefox struct {
	browser.Browser
}

func (ff Firefox) List() ([]*browsehist.HistoryItem, error) {
	return nil, nil
}

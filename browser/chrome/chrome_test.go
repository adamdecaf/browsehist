package chrome

import (
	"testing"
)

func init() {
	// Sneak in and change where we're searching for places.sqlite files
	historyFilepaths = []string{
		"../../testdata/History",
	}
}

func TestChrome__findSqliteDB(t *testing.T) {
	c := Chrome{}
	where, err := c.findHistoryDB()
	if err != nil {
		t.Fatal(err)
	}
	if where == "" {
		t.Error("no History file found..")
	}
}

func TestChrome__parse(t *testing.T) {
	c := Chrome{}
	item, err := c.parse(
		"https://www.lastpass.com/how-it-works/",
		13156450184247212,
	)
	if err != nil {
		t.Error(err)
	}
	if item.Address.String() != "https://www.lastpass.com/how-it-works/" {
		t.Errorf("got: %s", item.Address.String())
	}
	if item.AccessTime.String() != "2017-11-29 17:29:44 +0000 UTC" {
		t.Errorf("got: %s", item.AccessTime.String())
	}
}

func TestChrome__readHistoryItems(t *testing.T) {
	c := Chrome{}
	db, err := c.driver(historyFilepaths[0])
	if err != nil {
		t.Fatal(err)
	}
	items, err := c.readHistoryItems(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items instead of 2", len(items))
	}

	urls := make(map[string]int, 0)
	urls[`http://google.com/`] = 0
	urls[`http://www.google.com/`] = 0
	urls[`https://www.google.com/`] = 0

	whens := make(map[string]int, 0)
	whens[`2017-11-29 18:58:13 +0000 UTC`] = 0

	// Mark what we've see
	for i := range items {
		urls[items[i].Address.String()] += 1
		whens[items[i].AccessTime.String()] += 1
	}

	// Compare
	for k, v := range urls {
		if v <= 0 {
			t.Errorf("%s should have been visited", k)
		}
	}
	for k, v := range whens {
		if v <= 0 {
			t.Errorf("%s should have been parsed", k)
		}
	}
	if len(urls) != 3 || len(whens) != 1 {
		t.Errorf("got %d of 3 urls, %d of 1 whens", len(urls), len(whens))
		for k, v := range urls {
			t.Errorf(" urls: %s, %d", k, v)
		}
		for k, v := range whens {
			t.Errorf(" whens: %s, %d", k, v)
		}
	}
}

package firefox

import (
	"testing"
)

func init() {
	// Sneak in and change where we're searching for places.sqlite files
	placesFilepaths = []string{
		"../../testdata/places.sqlite",
	}
}

func TestFirefox__findSqliteDB(t *testing.T) {
	ff := Firefox{}
	where, err := ff.findPlacesDB()
	if err != nil {
		t.Fatal(err)
	}
	if where == "" {
		t.Error("no places.sqlite file found..")
	}
}

func TestFirefox__parse(t *testing.T) {
	ff := Firefox{}
	item, err := ff.parse(
		"https://mozilla.org",
		1509307351773044,
	)
	if err != nil {
		t.Error(err)
	}
	if item.Address.String() != "https://mozilla.org" {
		t.Errorf("got: %s", item.Address.String())
	}
	if item.AccessTime.String() != "2017-10-29 15:02:31 -0500 CDT" {
		t.Errorf("got: %s", item.AccessTime.String())
	}
}

func TestFirefox__readHistoryItems(t *testing.T) {
	ff := Firefox{}
	db, err := ff.driver(placesFilepaths[0])
	if err != nil {
		t.Fatal(err)
	}
	items, err := ff.readHistoryItems(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items instead of 2", len(items))
	}

	urls := make(map[string]int, 0)
	urls[`https://www.mozilla.org/en-US/firefox/57.0/firstrun/`] = 0
	urls[`https://www.mozilla.org/privacy/firefox/`] = 0
	urls[`https://www.mozilla.org/en-US/privacy/firefox/`] = 0

	whens := make(map[string]int, 0)
	whens[`2017-11-24 12:52:35 -0600 CST`] = 0
	whens[`2017-11-24 12:52:36 -0600 CST`] = 0

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
	if len(urls) != 3 || len(whens) != 2 {
		t.Errorf("got %d of 3 urls, %d of 2 whens", len(urls), len(whens))
	}
}

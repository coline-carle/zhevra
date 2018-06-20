package curseforge

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wow-sweetlie/zhevra/storage"
)

func TestIntegrationDecode(t *testing.T) {
	fixture, err := os.Open("./fixtures/clientdb/complete.json")
	if err != nil {
		t.Fatalf("error loading fixture: %s", err)
	}
	db, err := DecodeDB(fixture)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	const minAddonCount = 5000
	if len(db.Addons) < minAddonCount {
		log.Fatalf("expected number of addon, got > %d, got: %d\n",
			minAddonCount, len(db.Addons))
	}
}

func TestDecode(t *testing.T) {
	date, err := time.Parse(dateLayout, "2017-09-20T13:00:12.89")
	if err != nil {
		t.Fatalf("error parsing date: %s", date)
	}

	modules := []Module{
		Module{
			Foldername:  "Gatherer",
			Fingerprint: 3785330930,
		},
		Module{
			Foldername:  "SlideBar",
			Fingerprint: 1977511019,
		},
		Module{
			Foldername:  "!Swatter",
			Fingerprint: 1122147061,
		},
		Module{
			Foldername:  "Gatherer_HUD",
			Fingerprint: 2170801839,
		},
	}

	latest := Release{
		ID:          2482102,
		Filename:    "7.3.1",
		CreatedAt:   dbDate{date},
		URL:         "https://files.forgecdn.net/files/2482/102/Gatherer-7.3.1.zip",
		GameVersion: []string{"7.3.0"},
		Modules:     modules,
	}
	gatherer := Addon{
		ID:            32,
		Name:          "Gatherer",
		URL:           "https://www.curseforge.com/wow/addons/gatherer",
		Summary:       "Helps track the closest plants, deposits and treasure locations on your minimap.",
		DownloadCount: 15465595,
		Releases:      []Release{latest},
	}

	gathererDB := &ClientDB{
		Timestamp: 1528482271834,
		Addons:    []Addon{gatherer},
	}
	fixture, err := os.Open("./fixtures/clientdb/gatherer.json")
	if err != nil {
		t.Fatalf("error loading fixture: %s", err)
	}
	db, err := DecodeDB(fixture)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, gathererDB, db)

	// now test convertion into our models

	curseRelease := storage.CurseRelease{
		ID:          2482102,
		Filename:    "7.3.1",
		CreatedAt:   date,
		URL:         "https://files.forgecdn.net/files/2482/102/Gatherer-7.3.1.zip",
		GameVersion: "7.3.0",
	}
	expectedAddon := storage.CurseAddon{
		ID:            32,
		Name:          "Gatherer",
		URL:           "https://www.curseforge.com/wow/addons/gatherer",
		Summary:       "Helps track the closest plants, deposits and treasure locations on your minimap.",
		DownloadCount: 15465595,
		Releases:      []storage.CurseRelease{curseRelease},
	}

	curseAddon := NewCurseAddon(gatherer)
	assert.Equal(t, expectedAddon, curseAddon)

}

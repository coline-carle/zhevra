package curseforge

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationDecode(t *testing.T) {
	t.SkipNow()
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
	latest := Release{
		ID:          2482102,
		Filename:    "7.3.1",
		Date:        dbDate{date},
		DownloadURL: "https://files.forgecdn.net/files/2482/102/Gatherer-7.3.1.zip",
		GameVersion: []string{"7.3.0"},
	}
	gatherer := Addon{
		ID:            32,
		Name:          "Gatherer",
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

}

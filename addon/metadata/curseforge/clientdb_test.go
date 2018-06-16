package curseforge

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	latest := LatestFile{
		ID:          2482102,
		Filename:    "7.3.1",
		Date:        "2017-09-20T13:00:12.89",
		DownloadURL: "https://files.forgecdn.net/files/2482/102/Gatherer-7.3.1.zip",
		GameVersion: []string{"7.3.0"},
	}
	gatherer := Addon{
		ID:            32,
		Name:          "Gatherer",
		Summary:       "Helps track the closest plants, deposits and treasure locations on your minimap.",
		DownloadCount: 15465595,
		LatestFiles:   []LatestFile{latest},
	}

	gathererDB := &ClientDB{
		Timestamp: 1528482271834,
		Data:      []Addon{gatherer},
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

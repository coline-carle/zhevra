package curseforge

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

// ClientDB database root
type ClientDB struct {
	Addons    []Addon `json:"data"`
	Timestamp int64   `json:"timestamp"`
}

// Addon Curseforge addon struct
type Addon struct {
	ID            int64     `json:"Id"`
	Name          string    `json:"Name"`
	URL           string    `json:"WebSiteIURL"`
	Summary       string    `json:"Summary"`
	DownloadCount float64   `json:"DownloadCount"`
	Releases      []Release `json:"LatestFiles"`
}

// Release Curseforge last release struct
type Release struct {
	ID          int64    `json:"Id"`
	Filename    string   `json:"Filename"`
	Date        dbDate   `json:"Filedate"`
	DownloadURL string   `json:"DownloadURL"`
	GameVersion []string `json:"GameVersion"`
}

type dbDate struct {
	time.Time
}

const dateLayout = "2006-01-02T15:04:05"

func (d *dbDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	d.Time, err = time.Parse(dateLayout, s)
	return err
}

// DecodeDB decode client database
func DecodeDB(r io.Reader) (*ClientDB, error) {
	db := &ClientDB{}
	err := json.NewDecoder(r).Decode(db)
	return db, err
}

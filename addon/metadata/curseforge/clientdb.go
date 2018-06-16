package curseforge

import (
	"encoding/json"
	"io"
)

// ClientDB database root
type ClientDB struct {
	Data      []Addon `json:"data"`
	Timestamp int64   `json:"timestamp"`
}

// Addon Curseforge addon struct
type Addon struct {
	ID            int64        `json:"Id"`
	Name          string       `json:"Name"`
	URL           string       `json:"WebSiteIURL"`
	Summary       string       `json:"Summary"`
	DownloadCount int64        `json:"DownloadCount"`
	LatestFiles   []LatestFile `json:"LatestFiles"`
}

// LatestFile Curseforge last release struct
type LatestFile struct {
	ID          int64    `json:"Id"`
	Filename    string   `json:"Filename"`
	Date        string   `json:"Filedate"`
	DownloadURL string   `json:"DownloadURL"`
	GameVersion []string `json:"GameVersion"`
}

// DecodeDB decode client database
func DecodeDB(r io.Reader) (*ClientDB, error) {
	db := &ClientDB{}
	err := json.NewDecoder(r).Decode(db)
	return db, err
}

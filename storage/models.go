package store

import "time"

type Addon struct {
	ID   int64
	Name string
}

type CurseProvider struct {
	ID            int64
	URL           string
	Summary       string
	DownloadCount float64
	Releases      []CurseRelease
}

type CurseRelease struct {
	ID          int64
	Filename    string
	Date        time.Time
	DownloadURL string
	GameVersion string
}

package storage

import (
	"errors"
	"time"
)

type CurseAddon struct {
	ID            int64
	Name          string
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
	AddonID     int64
	addon       *CurseAddon
}

var (
	ErrCurseAddonDoesNotExists = errors.New("curse addon does not exist")
)

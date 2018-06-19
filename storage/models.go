package storage

import (
	"errors"
	"time"
)

// CurseAddon represent a given addon handled by curse provider
type CurseAddon struct {
	ID            int64
	Name          string
	URL           string
	Summary       string
	DownloadCount int64
	Releases      []CurseRelease
}

// CurseRelease represent a release of an addon handled by curse provider
type CurseRelease struct {
	ID          int64
	Filename    string
	CreatedAt   time.Time
	URL         string
	GameVersion string
	AddonID     int64
	addon       *CurseAddon
}

var (
	// ErrCurseAddonDoesNotExists is the error returned when no addon row match
	ErrCurseAddonDoesNotExists = errors.New("curse addon does not exist")
	// ErrCurseReleaseDoesNotExists is the error returned when no release row match
	ErrCurseReleaseDoesNotExists = errors.New("curse addon release does not exist")
)

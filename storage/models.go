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

// CurseFolder is a folder in addons directory
type CurseFolder struct {
	ReleaseID int64
	Name      string
}

// CurseRelease represent a release of an addon handled by curse provider
type CurseRelease struct {
	ID          int64
	Filename    string
	CreatedAt   time.Time
	URL         string
	GameVersion string
	AddonID     int64
	IsAlternate bool
	Directories []string
}

var (
	// ErrCurseAddonDoesNotExists is the error returned when no addon row match
	ErrCurseAddonDoesNotExists = errors.New("curse addon does not exist")
	// ErrCurseReleaseDoesNotExists is the error returned when no release row match
	ErrCurseReleaseDoesNotExists = errors.New("curse addon release does not exist")
)

// MainReleases return release that are not flagged as alternate by curse
// maxVersion is the max version (avoid betas for a standard release)
func (a *CurseAddon) MainReleases(minVersion string, maxVersion string) ([]CurseRelease, error) {
	numMaxVersion, err := VersionToInt(maxVersion)
	if err != nil {
		return nil, err
	}

	numMinVersion, err := VersionToInt(minVersion)
	if err != nil {
		return nil, err
	}

	releases := []CurseRelease{}
	for _, release := range a.Releases {
		if !release.IsAlternate {
			version, err := release.NumericGameVersion()
			if err != nil {
				return nil, err
			}
			if version >= numMinVersion && version <= numMaxVersion {
				releases = append(releases, release)
			}
		}
	}

	return releases, nil
}

// NumericGameVersion return a numeric representation of game version for
// comparing
func (c *CurseRelease) NumericGameVersion() (int, error) {
	return VersionToInt(c.GameVersion)
}

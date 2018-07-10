package curseforge

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/coline-carle/zhevra/storage/model"
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
	URL           string    `json:"WebSiteURL"`
	Summary       string    `json:"Summary"`
	DownloadCount float64   `json:"DownloadCount"`
	Releases      []Release `json:"LatestFiles"`
}

// Release Curseforge last release struct
type Release struct {
	ID          int64    `json:"Id"`
	Filename    string   `json:"Filename"`
	CreatedAt   dbDate   `json:"Filedate"`
	URL         string   `json:"DownloadURL"`
	GameVersion []string `json:"GameVersion"`
	Modules     []Module `json:"Modules"`
	IsAlternate bool     `json:"IsAlternate"`
}

// Module repesent the different directories in addons folder
type Module struct {
	Fingerprint int64
	Foldername  string
}

type dbDate struct {
	time.Time
}

// NewCurseRelease transform curse database struct into app entity
func NewCurseRelease(release Release) model.CurseRelease {
	versions := make([]int, 0, len(release.GameVersion))
	for _, gameVersionStr := range release.GameVersion {
		version, err := model.VersionToInt(gameVersionStr)

		// FIXME: not the better place to fail
		if err != nil {
			log.Fatalf("error parsing gameversion (%s): %s", gameVersionStr, err)
		}
		versions = append(versions, version)
	}
	curseRelease := model.CurseRelease{
		ID:           release.ID,
		Filename:     release.Filename,
		CreatedAt:    release.CreatedAt.Time,
		URL:          release.URL,
		GameVersions: versions,
		Directories:  make([]string, 0, len(release.Modules)),
	}
	for _, module := range release.Modules {
		curseRelease.Directories = append(curseRelease.Directories, module.Foldername)
	}
	return curseRelease
}

// NewCurseAddon transform curse database struct into app entity
func NewCurseAddon(addon Addon) model.CurseAddon {
	curseAddon := model.CurseAddon{
		ID:            addon.ID,
		Name:          addon.Name,
		URL:           addon.URL,
		Summary:       addon.Summary,
		DownloadCount: int64(addon.DownloadCount),
		Releases:      []model.CurseRelease{},
	}

	for _, release := range addon.Releases {
		curseRelease := NewCurseRelease(release)
		curseAddon.Releases = append(curseAddon.Releases, curseRelease)
	}
	return curseAddon
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

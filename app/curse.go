package app

import (
	"compress/bzip2"
	"database/sql"
	"io"
	"net/http"

	"github.com/coline-carle/zhevra/addon/metadata/curseforge"
	"github.com/pkg/errors"
)

// ErrWowNotDetected is the error returned when the wow folder can not be found
var ErrWowNotDetected = errors.New("world of warcraft not detected")

const curseDatabaseURL = "http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2"

// DownloadDatabase download and decode a curseforge database
func (a *App) DownloadDatabase(url string) (
	*curseforge.ClientDB,
	error) {
	resp, err := http.Get(url)
	if err != nil {
		errors.Wrap(err, "Download Database network error")
	}
	defer resp.Body.Close()
	curseDB := bzip2.NewReader(resp.Body)
	return curseforge.DecodeDB(curseDB)
}

// ImportCurseDB import complete curse database
// curseDB reader of curse json database
func (a *App) ImportCurseDB(curseDB io.Reader) error {
	jsonDB, err := curseforge.DecodeDB(curseDB)
	if err != nil {
		return errors.Wrap(err, "ImportDB")
	}
	err = a.storage.Tx(func(tg *sql.Tx) error {
		for _, addon := range jsonDB.Addons {
			curseAddon := curseforge.NewCurseAddon(addon)
			err = a.storage.Tx(func(tx *sql.Tx) error {
				return a.storage.CreateCurseAddon(tx, curseAddon)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "ImportDB migration")
	}
	return nil
}

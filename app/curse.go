package app

import (
	"compress/bzip2"
	"database/sql"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/addon/metadata/curseforge"
)

const curseDatabaseURL = "http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2"

// DownloadDatabase download and decode a curseforge database
func (a *App) DownloadDatabase(url string) (
	*curseforge.ClientDB,
	error) {
	resp, err := http.Get(url)
	if err != nil {
		err.Wrap("DownloadDabase network error", err)
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
			err = a.Storage.Tx(func(tx *sql.Tx) error {
				return a.Storage.CreateCurseAddon(tx, curseAddon)
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
}

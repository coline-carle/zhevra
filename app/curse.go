package app

import (
	"compress/bzip2"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/coline-carle/zhevra/addon/metadata/curseforge"
	"github.com/pkg/errors"
)

const url = "http://clientupdate-v6.cursecdn.com/feed/addons/1/v10/complete.json.bz2"

// DownloadDatabase download and decode a curseforge database
func (a *App) DownloadDatabase() (*curseforge.ClientDB, error) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		errors.Wrap(err, "Download Database network error")
	}
	defer resp.Body.Close()
	curseDB := bzip2.NewReader(resp.Body)
	clientdb, err := curseforge.DecodeDB(curseDB)
	elapsed := time.Now().Sub(start)
	fmt.Printf("downloaded database in %v\n", elapsed)
	return clientdb, err
}

// ImportCurseDB import complete curse database
// curseDB reader of curse json database
func (a *App) ImportCurseDB(curseDB *curseforge.ClientDB) error {
	var err error
	err = a.storage.Tx(func(tg *sql.Tx) error {
		for _, addon := range curseDB.Addons {
			curseAddon := curseforge.NewCurseAddon(addon)
			err = a.storage.CreateCurseAddon(tg, curseAddon)
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

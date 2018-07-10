package app

import (
	"compress/bzip2"
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/coline-carle/zhevra/addon/metadata/curseforge"
	"github.com/pkg/errors"
)

// ErrWowNotDetected is the error returned when the wow folder can not be found
var ErrWowNotDetected = errors.New("world of warcraft not detected")

const curseDatabaseURL = "http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2"
const wow6432RegKey = `SOFTWARE\WOW6432Node\Blizzard Entertainment\World of Warcraft`
const wowRegKey = `SOFTWARE\WOW6432Node\Blizzard Entertainment\World of Warcraft`
const pathKey = "InstallPath"
const defaultInstall64 = `%PROGRAMFILES(x86)%\World of Warcraft`
const defaultInstall32 = `%PROGRAMFILES%\World of Warcraft`

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

// GuessWowDir return wow directory name
func (a *App) GuessWowDir() (string, error) {
	var path string
	var err error
	path, err = guessFromReg(wow6432RegKey)
	if err == nil && checkDir(path) {
		return path, nil
	}
	path, err = guessFromReg(wow6432RegKey)
	if err == nil && checkDir(path) {
		return path, nil
	}
	if checkDir(defaultInstall64) {
		return defaultInstall64, nil
	}
	if checkDir(defaultInstall32) {
		return defaultInstall32, nil
	}
	return "", errors.New("world of warcraft not detected")
}

func pathExists(path string, filename string) bool {
	fullpath := filepath.Join(path, filename)
	if _, err := os.Stat(fullpath); err == nil {
		return true
	}
	return false
}
func checkDir(dir string) bool {
	anyPath := []string{
		"Wow-64.exe",
		"Wow.exe",
	}
	for _, filename := range anyPath {
		if pathExists(dir, filename) {
			return true
		}
	}
	return false
}

func guessFromReg(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, key, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	path, _, err := k.GetStringValue(pathKey)
	if err != nil {
		return "", err
	}
	return path, nil
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

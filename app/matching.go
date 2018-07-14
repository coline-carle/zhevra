package app

import (
	"database/sql"
	"log"

	"github.com/coline-carle/zhevra/addondir"

	"github.com/coline-carle/zhevra/storage/model"
)

func (app *App) directoryCoverage() error {
	if len(app.directories) == 0 {
		err := app.ScanAddonsDirectories()
		if err != nil {
			return err
		}
	}
	dirMap := make(map[string][]model.CurseAddon)
	for directory := range app.directories {
		var matchingAddons []model.CurseAddon
		err := app.storage.Tx(func(tx *sql.Tx) error {
			var err error
			matchingAddons, err = app.storage.FindAddonsWithDirectoryName(tx, directory)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, addon := range matchingAddons {
			mainReleases, err := addon.MainReleases(app.wowVersion)
			if err != nil {
				log.Fatal(err)
			}
			for _, release := range mainReleases {
				if allDirectoriesExists(release.Directories, app.directories) {
					for _, addonDir := range release.Directories {
						dirMap[addonDir] = append(dirMap[addonDir], addon)
					}
				}
			}
		}
	}
	return nil
}

func allDirectoriesExists(directories []string, allDirectories map[string]addondir.Hashmap) bool {
	for _, directory := range directories {
		if _, ok := allDirectories[directory]; !ok {
			return false
		}
	}
	return true
}

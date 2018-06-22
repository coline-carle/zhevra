package cmd

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/wow-sweetlie/zhevra/storage"
	"github.com/wow-sweetlie/zhevra/storage/sqlite"
)

var (
	errNoMatchingAddon = errors.New("no matching addon")
)

var wowDir string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan wow folder and identify addons",
	Run:   runScanCmd,
}

func setScanCmdFlags() {
	scanCmd.Flags().StringVarP(&wowDir, "wow", "w", "", "wow directory")
	scanCmd.MarkFlagRequired("wow")
}

func runScanCmd(cmd *cobra.Command, args []string) {
	addondb, err := sqlite.NewStorage("addons.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(wowDir)
	if err != nil {
		log.Fatal(err)
	}

	var dirMap map[string](map[int64]storage.CurseAddon)
	dirMap = make(map[string](map[int64]storage.CurseAddon))

	addonsDirs := []string{}

	for _, f := range files {
		if f.IsDir() {
			dirMap[f.Name()] = make(map[int64]storage.CurseAddon)
			addonsDirs = append(addonsDirs, f.Name())
		}
	}
	for _, directory := range addonsDirs {
		var matchingAddons []storage.CurseAddon
		err := addondb.Tx(func(tx *sql.Tx) error {
			matchingAddons, err := addondb.FindAddonsWithDirectoryName(tx, directory)
			return err
		})
		for _, addon := range matchingAddons {
			if len(addon.Releases) > 0 {
				release := addon.Releases[0]
				for _, addonDir := range release.Directories {
					if _, ok := dirMap[addonDir]; ok {
						dirMap[addonDir][addon.ID] = addon
					}
				}
			}
		}
	}

}

func FiltrableAddon(
	db *sqlite.Storage, addonsDirs []string) (*storage.CurseAddon, error) {
	var matchingAddons []storage.CurseAddon
	for _, directory := range addonsDirs {
		if err != nil {
			return nil, err
		}
		if len(matchingAddons) == 0 {
			return nil, errors.New("addon not found")
		}
		if len(matchinAddons) == 1 {
			return matchingAddons[0], nil
		}
	}
	return nil, errors.New("no addon matching")
}

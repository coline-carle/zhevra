package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

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
	maxVersion := "7.3.5"
	minVersion := "6.0.0"
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(wowDir)
	if err != nil {
		log.Fatal(err)
	}

	var dirMap map[string](map[int64]storage.CurseAddon)
	dirMap = make(map[string](map[int64]storage.CurseAddon))

	addonsDirs := make(map[string]bool)

	for _, f := range files {
		if f.IsDir() {
			tocFilename := fmt.Sprintf("%s.toc", f.Name())
			tocPath := path.Join(wowDir, f.Name(), tocFilename)
			if _, err := os.Stat(tocPath); os.IsNotExist(err) {
				continue
			}
			dirMap[f.Name()] = make(map[int64]storage.CurseAddon)
			addonsDirs[f.Name()] = true
		}
	}

	for directory := range addonsDirs {
		var matchingAddons []storage.CurseAddon
		err := addondb.Tx(func(tx *sql.Tx) error {
			matchingAddons, err = addondb.FindAddonsWithDirectoryName(tx, directory)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, addon := range matchingAddons {
			mainReleases, err := addon.MainReleases(minVersion, maxVersion)
			if err != nil {
				log.Fatal(err)
			}
			for _, release := range mainReleases {
				if directory == "AllTheThings" {
					fmt.Printf("%+v\n", addon)
				}
				if allDirectoriesExists(release.Directories, addonsDirs) {
					for _, addonDir := range release.Directories {
						dirMap[addonDir][addon.ID] = addon
					}
				}
			}
		}
	}

	fmt.Printf("%+v", dirMap["AllTheThings"])

	coveredDirs := make(map[string]bool)
	validatedAddons := make(map[int64]storage.CurseAddon)

	for dir, addons := range dirMap {
		if _, covered := coveredDirs[dir]; !covered && len(addons) == 1 {
			var addon storage.CurseAddon
			for _, addon = range addons {
			}
			mainReleases, err := addon.MainReleases(minVersion, maxVersion)
			if err != nil {
				log.Fatal(err)
			}
			for _, release := range mainReleases {
				if allDirectoriesExists(release.Directories, addonsDirs) {
					validatedAddons[addon.ID] = addon
					for _, releaseDir := range release.Directories {
						coveredDirs[releaseDir] = true
					}
					break
				}
			}
		}
	}
	fmt.Printf("\n")

	for directory := range addonsDirs {
		if _, ok := coveredDirs[directory]; !ok {
			fmt.Printf("not covered: %s\n", directory)
		}
	}
}

func allDirectoriesExists(directories []string, allDirectories map[string]bool) bool {
	for _, directory := range directories {
		if _, ok := allDirectories[directory]; !ok {
			return false
		}
	}
	return true
}
package cmd

import (
	"compress/bzip2"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/coline-carle/zhevra/addon/metadata/curseforge"
	"github.com/coline-carle/zhevra/storage/sqlite"
	"github.com/spf13/cobra"
)

const url = "http://clientupdate-v6.cursecdn.com/feed/addons/1/v10/complete.json.bz2"

var filename string

var importdbCmd = &cobra.Command{
	Use:   "importdb",
	Short: "import a curse database for the skee of debugging",
	Run:   runImportdbCmd,
}

func setImportdbCmdFlags() {
	importdbCmd.Flags().StringVarP(
		&filename, "filename", "f", "", "import cursedb filename")
	importdbCmd.PersistentFlags().Bool("online", true, "fetch the database from internet")
}

func runImportdbCmd(cmd *cobra.Command, args []string) {
	var curseDB *curseforge.ClientDB
	var err error
	if len(filename) > 0 {
		curseDBFile, err := os.Open(filename)
		defer curseDBFile.Close()
		if err != nil {
			log.Fatalf("error loading cursedb: %s", err)
		}
		curseDB, err = curseforge.DecodeDB(curseDBFile)
	} else {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("error downloading cursedb: %s", err)
		}
		defer resp.Body.Close()
		reader := bzip2.NewReader(resp.Body)
		curseDB, err = curseforge.DecodeDB(reader)
	}

	if err != nil {
		log.Fatalf("unexpected error decoding the database: %s", err)
	}

	programDB, err := sqlite.NewStorage("addons.sqlite")
	if err != nil {
		log.Fatalf("unexpected error creating the database: %s", err)
	}
	programDB.Migrate()
	if err != nil {
		log.Fatalf("unexpected while migrating: %s", err)
	}
	defer programDB.Close()

	err = programDB.Tx(func(tg *sql.Tx) error {
		for _, addon := range curseDB.Addons {
			curseAddon := curseforge.NewCurseAddon(addon)
			err = programDB.Tx(func(tx *sql.Tx) error {
				return programDB.CreateCurseAddon(tx, curseAddon)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error while inserting records: %s", err)
	}
}

package cmd

import (
	"database/sql"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wow-sweetlie/zhevra/addon/metadata/curseforge"
	"github.com/wow-sweetlie/zhevra/storage/sqlite"
)

var cursedb string

var importdbCmd = &cobra.Command{
	Use:   "importdb",
	Short: "import a curse database for the skee of debugging",
	Run:   runImportdbCmd,
}

func setImportdbCmdFlags() {
	importdbCmd.Flags().StringVarP(
		&cursedb, "cursedb", "c", "", "cursedb file (required)")
	importdbCmd.MarkFlagRequired("cursedb")
}

func runImportdbCmd(cmd *cobra.Command, args []string) {
	curseDBFile, err := os.Open(cursedb)
	if err != nil {
		log.Fatalf("error loading cursedb: %s", err)
	}
	jsonDB, err := curseforge.DecodeDB(curseDBFile)
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
		for _, addon := range jsonDB.Addons {
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

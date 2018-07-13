package cmd

import (
	"log"

	"github.com/coline-carle/zhevra/app"
	"github.com/spf13/cobra"
)

const url = "http://clientupdate-v6.cursecdn.com/feed/addons/1/v10/complete.json.bz2"

var filename string

var importdbCmd = &cobra.Command{
	Use:   "importdb",
	Short: "import a curse database for the skee of debugging",
	Run:   runImportdbCmd,
}

func runImportdbCmd(cmd *cobra.Command, args []string) {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("unexpected error decoding the database: %s", err)
	}

	defer a.Close()

	db, err := a.DownloadDatabase()
	a.ImportCurseDB(db)
}

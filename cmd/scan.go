package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/wow-sweetlie/zhevra/storage/sqlite"
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

	for _, f := range files {
		if f.IsDir() {
			fmt.Println(f.Name())
		}
	}
}

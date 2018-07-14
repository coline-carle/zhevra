package main

import (
	"log"
	"os"

	"github.com/coline-carle/zhevra/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("unexpected error decoding the database: %s", err)
	}

	defer app.Close()

	err = app.DetectWow()
	if err != nil {
		os.Exit(1)
	}

	err = app.ScanAddonsDirectories()
	if err != nil {
		log.Fatal(err)
	}

	err = app.UpdateCurseDatabase()
	if err != nil {
		os.Exit(1)
	}
}

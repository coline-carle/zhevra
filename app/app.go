package app

import (
	"fmt"
	"log"
	"time"

	"github.com/coline-carle/zhevra/addondir"
	"github.com/coline-carle/zhevra/storage"
)

// App is the context of the main app
type App struct {
	StorageLastMod *time.Time `json:"storage_lastmod,omitempty"`
	storage        storage.Storage
	wowDir         string
	directories    map[string]addondir.Hashmap
}

// Close the databaase save the state
func (app *App) Close() {
	app.storage.Close()
	err := app.WriteState()
	if err != nil {
		log.Fatal(err)
	}
}

// DetectWowDirectory detect wow directory
func (app *App) DetectWowDirectory() error {
	var err error
	app.wowDir, err = GuessWowDir()
	if err != nil {
		fmt.Println("failed to detect wow directory")
		fmt.Printf("error: %s", err)
		return err
	}
	fmt.Printf("wow directory found: %s\n", app.wowDir)
	return nil
}

// ScanAddonsDirectories scan all addon folders
func (app *App) ScanAddonsDirectories() error {
	var err error
	if len(app.wowDir) == 0 {
		err = app.DetectWowDirectory()
		if err != nil {
			return err
		}
	}
	app.directories, err = addondir.AddonDirectories(app.wowDir)
	log.Println(len(app.directories))
	for dir, hashes := range app.directories {
		fmt.Println(dir)
		for file, hash := range hashes {
			fmt.Printf("\t- %s [%x]\n", file, hash)
		}

	}

	return err
}

// NewApp init App struct, connect to storage and migrate
func NewApp() (*App, error) {
	storage, err := LoadDatabase()
	if err != nil {
		log.Fatal(err)
	}
	app, err := LoadState()
	if err != nil {
		log.Fatal(err)
	}
	app.storage = storage
	return app, nil
}

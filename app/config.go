package app

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"unicode"

	"github.com/coline-carle/zhevra/storage"
	"github.com/pkg/errors"
)

// AppName is the name of the app :)
const (
	AppName     = "zhevra"
	StateFile   = "state.json"
	AddonDBFile = "addons.sqlite"
)

// LoadDatabase open sqlite database and migrate it
func LoadDatabase() (storage.Storage, error) {
	path := filepath.Join(DataDir(), AddonDBFile)
	storage, err := storage.NewStorage(path)
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	err = storage.Migrate()
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	return storage, nil
}

// WriteState Save the state to the disk
func (app *App) WriteState() error {
	statePath := filepath.Join(DataDir(), StateFile)
	stateFile, err := os.Create(statePath)
	if err != nil {
		return err
	}
	defer stateFile.Close()
	return json.NewEncoder(stateFile).Encode(app)
}

// LoadState load the state of the application from the json file
func LoadState() (*App, error) {
	statePath := filepath.Join(DataDir(), StateFile)

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &App{}, nil
	}
	stateFile, err := os.Open(statePath)
	if err != nil {
		return nil, err
	}
	defer stateFile.Close()
	app := &App{}
	err = json.NewDecoder(stateFile).Decode(app)
	return app, err
}

// DataDir a function that return the app directory on all supported os
func DataDir() string {
	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	appNameLower := string(unicode.ToLower(rune(AppName[0]))) + AppName[1:]

	return filepath.Join(homeDir, "."+appNameLower)
}

func init() {
	// ensure the app datadir is created
	appDataDir := DataDir()
	if _, err := os.Stat(appDataDir); err != nil {
		err = os.Mkdir(appDataDir, 750)
		if err != nil {
			log.Fatalf("could not create config dir: %s\n", err)
		}
	}
}

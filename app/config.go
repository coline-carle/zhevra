package app

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"unicode"
)

// AppName is the name of the app :)
const (
	AppName     = "zhevra"
	AddonDBFile = "addons.sqlite"
)

// DatabasePath the full path to the sqlite file
func DatabasePath() string {
	return filepath.Join(DataDir(), AddonDBFile)
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

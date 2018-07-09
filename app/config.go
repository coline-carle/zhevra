package app

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
	"unicode"
)

// AppName is the name of the app :)
const (
	AppName     = "zhevra"
	AddonDBFile = "addons.sqlite"
)

// Config struct store the base configuration of the application
type Config struct {
	CurseDB CurseSnapshot
	WoWs    []WoWInfo
}

type WoWInfo struct {
	ID     int
	Folder string
	Addons []AddonInfo
}

type AddonInfo struct {
	id         int64
	version    string
	ReleasedAt time.Time
}

type CurseSnapshot struct {
	CreatedAt time.Time
}

// DatabasePath the full path to the sqlite file
func DatabasePath() string {
	return filepath.Join(DataDir(), addonDBFile)
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
	appNameLower := string(unicode.ToLower(rune(c.AppName[0]))) + c.AppName[1:]

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

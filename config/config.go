package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"unicode"
)

// Config struct store the base configuration of the application
type Config struct {
	AppName string
	dbFile  string
}

// DefaultConfig is the config that will be used by default
var DefaultConfig = &Config{
	AppName: "zhevra",
	dbFile:  "database.sqlite",
}

// DatabasePath the full path to the sqlite file
func (c *Config) DatabasePath() string {
	return filepath.Join(c.AppDataDir(), c.dbFile)
}

// AppDataDir a function that return the app directory on all supported os
func (c *Config) AppDataDir() string {
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
	appDataDir := DefaultConfig.AppDataDir()
	if _, err := os.Stat(appDataDir); err != nil {
		err = os.Mkdir(appDataDir, 750)
		if err != nil {
			log.Fatalf("could not create config dir: %s\n", err)
		}
	}
}

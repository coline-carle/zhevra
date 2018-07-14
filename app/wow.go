package app

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const versionFieldname = "Version!STRING:0"

// DetectWow Detect wow version and directory
func (app *App) DetectWow() error {
	var err error
	err = app.DetectWowDirectory()
	if err != nil {
		return err
	}
	return app.DetectWowVersion()
}

// DetectWowVersion detect the current versin of the game
func (app *App) DetectWowVersion() error {
	var err error
	if len(app.wowDir) == 0 {
		err = app.DetectWowDirectory()
		if err != nil {
			return err
		}
	}
	path := filepath.Join(app.wowDir, ".build.info")
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("can't open file: %s, for wow build detection\n", err)
		return err
	}

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = '|'
	lines, err := r.ReadAll()
	if len(lines) != 2 {
		fmt.Println("invalid .build.info file")
		return errors.New("invalid .build.info file")
	}
	line := lines[0]
	var i int
	var val string
	for i, val = range line {
		if val == versionFieldname {
			break
		}
	}
	if val != versionFieldname {
		fmt.Println("can'f find version filed in .build.info file")
		return errors.New("no version field")
	}

	app.wowVersion = lines[1][i]

	fmt.Printf("detected wow version: %s\n", app.wowVersion)

	return nil
}

package addondir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	ErrAddonsDirectoryNotFound = errors.New("addon directory can't be accessed")
)

// AddonDirectories return the map of addon directories with their hashmap
func AddonDirectories(wowDir string) (map[string]Hashmap, error) {
	addonsBaseDir := filepath.Join(wowDir, "Interface", "AddOns")
	files, err := ioutil.ReadDir(addonsBaseDir)
	if err != nil {
		return nil, ErrAddonsDirectoryNotFound
	}

	directories := make(map[string]Hashmap)

	for _, f := range files {
		if f.IsDir() {
			if IsAddonDir(addonsBaseDir, f.Name()) {
				directories[f.Name()], err = MD5All(filepath.Join(addonsBaseDir, f.Name()))
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return directories, nil
}

// IsAddonDir return true if the directory is a valid addon
// addonsDir the base directory of addons
// dir the addon directory to test
func IsAddonDir(addonsDir string, dir string) bool {
	tocFilename := fmt.Sprintf("%s.toc", dir)
	tocPath := filepath.Join(addonsDir, dir, tocFilename)
	_, err := os.Stat(tocPath)
	return !os.IsNotExist(err)
}

package addondir

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

const (
	wow6432RegKey    = `SOFTWARE\WOW6432Node\Blizzard Entertainment\World of Warcraft`
	wowRegKey        = `SOFTWARE\WOW6432Node\Blizzard Entertainment\World of Warcraft`
	pathKey          = "InstallPath"
	defaultInstall64 = `%PROGRAMFILES(x86)%\World of Warcraft`
	defaultInstall32 = `%PROGRAMFILES%\World of Warcraft`
)

// ErrWowNotDetected is the error returned when the wow folder can not be found
var ErrWowNotDetected = errors.New("world of warcraft not detected")

// GuessWowDir return wow directory name
func GuessWowDir() (string, error) {
	var path string
	var err error
	path, err = guessFromReg(wow6432RegKey)
	if err == nil && checkDir(path) {
		return path, nil
	}
	path, err = guessFromReg(wow6432RegKey)
	if err == nil && checkDir(path) {
		return path, nil
	}
	if checkDir(defaultInstall64) {
		return defaultInstall64, nil
	}
	if checkDir(defaultInstall32) {
		return defaultInstall32, nil
	}
	return "", errors.New("world of warcraft not detected")
}

func pathExists(path string, filename string) bool {
	fullpath := filepath.Join(path, filename)
	if _, err := os.Stat(fullpath); err == nil {
		return true
	}
	return false
}
func checkDir(dir string) bool {
	anyPath := []string{
		"Wow-64.exe",
		"Wow.exe",
	}
	for _, filename := range anyPath {
		if pathExists(dir, filename) {
			return true
		}
	}
	return false
}

func guessFromReg(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, key, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	path, _, err := k.GetStringValue(pathKey)
	if err != nil {
		return "", err
	}
	return path, nil
}

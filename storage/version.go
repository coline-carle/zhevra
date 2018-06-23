package storage

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	versionRegexp = regexp.MustCompile(`(\d{1,2})\.(\d{1,2}).(\d{1,2})`)
)

// VersionToInt convert a version in he format x.x.x to a number
func VersionToInt(version string) (int, error) {
	matches := versionRegexp.FindStringSubmatch(version)
	if len(matches) != 4 {
		return 0, errors.New("invalid version format")
	}
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, err
	}
	return major*0x100*0x100 + minor*0x100 + patch, nil
}

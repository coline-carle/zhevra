package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wow-sweetlie/zhevra/storage"
)

func TestFindByDirectory(t *testing.T) {
	si, err := NewStorage("test.sqlite")
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a database connection", err)
	}
	defer si.Close()
	err = si.Migrate()
	if err != nil {
		t.Errorf("An error '%s' was not while migrating", err)
	}
	addon := testUtilCreateAddon(1)
	addon.Releases = append(addon.Releases, testUtilCreateRelease(1, addon.ID))
	addon.Releases = append(addon.Releases, testUtilCreateRelease(2, addon.ID))
	defer tearDown(si, "curse_addon", "id = $1", addon.ID)
	defer tearDown(si, "curse_release", "id = $1", 1)
	defer tearDown(si, "curse_release", "id = $1", 2)
	defer tearDown(si, "curse_release_directory", "release_id = $1", 1)
	defer tearDown(si, "curse_release_directory", "release_id = $1", 2)
	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var addons []storage.CurseAddon
	err = si.Tx(func(tx *sql.Tx) error {
		addons, err = si.FindAddonsWithDirectoryName(tx, "directory-1-1")
		return err
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(addons) != 1 {
		t.Fatalf("unexpected number of addon fetched. E:1 Got: %d\n", len(addons))
	}
}

func TestCreateAddon(t *testing.T) {
	si, err := NewStorage("test.sqlite")
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a database connection", err)
	}
	defer si.Close()
	err = si.Migrate()
	if err != nil {
		t.Errorf("An error '%s' was not while migrating", err)
	}
	addon := testUtilCreateAddon(1)
	addon.Releases = append(addon.Releases, testUtilCreateRelease(1, addon.ID))
	addon.Releases = append(addon.Releases, testUtilCreateRelease(2, addon.ID))
	defer tearDown(si, "curse_addon", "id = $1", addon.ID)
	defer tearDown(si, "curse_release", "id = $1", 1)
	defer tearDown(si, "curse_release", "id = $1", 2)
	defer tearDown(si, "curse_release_directory", "release_id = $1", 1)
	defer tearDown(si, "curse_release_directory", "release_id = $1", 2)
	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var fetchedAddon storage.CurseAddon

	err = si.Tx(func(tx *sql.Tx) error {
		fetchedAddon, err = si.GetCurseAddon(tx, addon.ID)
		return err
	})
	if err != nil {

		t.Errorf("unexpected error: %s", err)
	}
	assert.Equal(t, addon, fetchedAddon)
}

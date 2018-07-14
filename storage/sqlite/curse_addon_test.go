package sqlite

import (
	"database/sql"
	"testing"

	"github.com/coline-carle/zhevra/storage/model"
	"github.com/stretchr/testify/assert"
)

func TestFindByDirectory(t *testing.T) {
	si, err := NewStorage(":memory:")
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

	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var addons []model.CurseAddon
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

func TestDeleteAllAddons(t *testing.T) {
	si, err := NewStorage(":memory:")
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

	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})

	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	err = si.Tx(func(tx *sql.Tx) error {
		return si.DeleteAllAddons(tx)
	})

	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	if !isTableEmpty(si, "curse_addon") {
		t.Errorf("Expected curse_addon to be empty\n")
	}
	if !isTableEmpty(si, "curse_release") {
		t.Errorf("Expected curse_release to be empty\n")
	}

	if !isTableEmpty(si, "curse_release_directory") {
		t.Errorf("Expected curse_release_directory to be empty\n")
	}
}

func TestCreateAddon(t *testing.T) {
	si, err := NewStorage(":memory:")
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
	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var fetchedAddon model.CurseAddon

	err = si.Tx(func(tx *sql.Tx) error {
		fetchedAddon, err = si.GetCurseAddon(tx, addon.ID)
		return err
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	assert.Equal(t, addon, fetchedAddon)
}

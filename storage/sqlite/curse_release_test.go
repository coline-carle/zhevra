package sqlite

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wow-sweetlie/zhevra/storage"
)

func TestCreateCurseRelease(t *testing.T) {
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
	release := testUtilCreateRelease(1, addon.ID)
	err = si.Tx(func(tx *sql.Tx) error {
		err2 := si.CreateCurseAddon(tx, addon)
		if err2 != nil {
			return err2
		}
		return si.CreateCurseRelease(tx, release)
	})
	defer tearDown(si, "curse_addon", "id = $1", addon.ID)
	defer tearDown(si, "curse_release", "id = $1", release.ID)
	defer tearDown(si, "curse_release_directory", "release_id = $1", release.ID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var fetchedRelease storage.CurseRelease

	err = si.Tx(func(tx *sql.Tx) error {
		fetchedRelease, err = si.GetCurseRelease(tx, release.ID)
		return err
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	assert.Equal(t, release, fetchedRelease)
}

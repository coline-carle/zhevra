package sqlite

import (
	"database/sql"
	"testing"
	"time"

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
	addon := &storage.CurseAddon{ID: 1}
	release := &storage.CurseRelease{
		ID:          1,
		Filename:    "Filename.zip",
		CreatedAt:   time.Date(2018, 9, 9, 0, 0, 0, 0, time.UTC),
		URL:         "https://url",
		GameVersion: "7.3.2",
		AddonID:     addon.ID,
	}
	err = si.Tx(func(tx *sql.Tx) error {
		err2 := si.CreateCurseAddon(tx, addon)
		if err2 != nil {
			return err2
		}
		return si.CreateCurseRelease(tx, release)
	})
	defer tearDown(si, "curse_addon", "id = $1", addon.ID)
	defer tearDown(si, "curse_release", "id = $1", addon.ID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var fetchedRelease *storage.CurseRelease

	err = si.Tx(func(tx *sql.Tx) error {
		fetchedRelease, err = si.GetCurseRelease(tx, release.ID)
		return err
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	assert.Equal(t, release, fetchedRelease)
}

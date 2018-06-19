package sqlite

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wow-sweetlie/zhevra/storage"
)

func tearDown(s *Storage, table string, where string, params ...interface{}) {
	_, err := s.DB.Exec(
		fmt.Sprintf("DELETE FROM \"%s\" WHERE %s", table, where),
		params...)
	if err != nil {
		panic(fmt.Sprintf("Problem tearing down %s data: %v", table, err))
	}
}

func testUtilCreateRelease(ID int64, addonID int64) storage.CurseRelease {
	return storage.CurseRelease{
		ID:          ID,
		Filename:    fmt.Sprintf("Filename-%d.zip", ID),
		CreatedAt:   time.Date(2018, 9, 9, 0, 0, 0, 0, time.UTC),
		URL:         fmt.Sprintf("https://url/filename-%d", ID),
		GameVersion: "7.3.2",
		AddonID:     addonID,
	}
}

func testUtilCreateAddon(ID int64) storage.CurseAddon {
	return storage.CurseAddon{
		ID:            ID,
		Name:          fmt.Sprintf("Name-%d", ID),
		Summary:       fmt.Sprintf("Summary-%d", ID),
		URL:           fmt.Sprintf("https://url/addon-%d", ID),
		DownloadCount: 33 * ID,
		Releases:      []storage.CurseRelease{},
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

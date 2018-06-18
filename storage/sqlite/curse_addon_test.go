package sqlite

import (
	"database/sql"
	"fmt"
	"testing"

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
	addon := &storage.CurseAddon{
		ID:            1,
		Name:          "Name",
		Summary:       "Summary",
		URL:           "https://url",
		DownloadCount: 2,
	}
	defer tearDown(si, "curse_addon", "id = $1", addon.ID)
	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	var fetchedAddon *storage.CurseAddon

	err = si.Tx(func(tx *sql.Tx) error {
		fetchedAddon, err = si.GetCurseAddon(tx, addon.ID)
		return err
	})
	if err != nil {

		t.Errorf("unexpected error: %s", err)
	}
	assert.Equal(t, addon, fetchedAddon)
}

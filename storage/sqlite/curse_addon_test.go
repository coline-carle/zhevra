package sqlite

import (
	"database/sql"
	"testing"

	"github.com/wow-sweetlie/zhevra/storage"
)

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
	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

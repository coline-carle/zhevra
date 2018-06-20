package sqlite

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
)

func TestSucesssfullTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	si := &Storage{DB: db}
	addon := testUtilCreateAddon(55)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO curse_addon").WithArgs(
		addon.ID,
		addon.Name,
		addon.URL,
		addon.Summary,
		addon.DownloadCount,
	).WillReturnResult(sqlmock.NewResult(addon.ID, 1))
	mock.ExpectCommit()

	err = si.Tx(func(tx *sql.Tx) error {
		return si.CreateCurseAddon(tx, addon)
	})
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s\n", err)
	}
}

func TestSucesssfailedTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	si := &Storage{DB: db}
	addon := testUtilCreateAddon(55)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO curse_addon").WithArgs(
		addon.ID,
		addon.Name,
		addon.URL,
		addon.Summary,
		addon.DownloadCount,
	).WillReturnResult(sqlmock.NewResult(addon.ID, 1))
	mock.ExpectExec("INSERT INTO curse_addon").WithArgs(
		addon.ID,
		addon.Name,
		addon.URL,
		addon.Summary,
		addon.DownloadCount,
	).WillReturnError(errors.New("unexpected error: failed to CreateCurseAddon"))
	mock.ExpectRollback()

	err = si.Tx(func(tx *sql.Tx) error {
		err2 := si.CreateCurseAddon(tx, addon)
		if err2 != nil {
			return err2
		}
		return si.CreateCurseAddon(tx, addon)
	})
	if err == nil {
		t.Errorf("this statement was supposed to failed")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s\n", err)
	}
}

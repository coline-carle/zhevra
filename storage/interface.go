package storage

import (
	"database/sql"

	"github.com/wow-sweetlie/zhevra/storage"
	"github.com/wow-sweetlie/zhevra/storage/sqlite"
)

// Storage is the interface to zhevra database
type Storage interface {
	// General
	Migrate() error
	Tx(fn func(*sql.Tx) error) error
	Close()

	// curse addons
	CreateCurseAddon(tx *sql.Tx, addon storage.CurseAddon) error
	FindAddonsWithDirectoryName(
		tx *sql.Tx,
		directory string) (
		[]storage.CurseAddon,
		error,
	)
	GetCurseAddon(tx *sql.Tx, id int64) (storage.CurseAddon, error)

	// curse releases
	GetCurseRelease(tx *sql.Tx, id int64) (storage.CurseRelease, error)
	FindCurseReleasesByAddonID(tx *sql.Tx, id int64) ([]storage.CurseRelease, error)
}

// NewStorage create the appropriate storage
// for a given url
func NewStorage(url string) error {
	return sqlite.NewStorage(url)
}

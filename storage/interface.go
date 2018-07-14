package storage

import (
	"database/sql"

	"github.com/coline-carle/zhevra/storage/model"
)

// Storage is the interface to zhevra database
type Storage interface {
	// General
	Migrate() error
	Tx(fn func(*sql.Tx) error) error
	Close() error

	// curse addons
	CreateCurseAddon(tx *sql.Tx, addon model.CurseAddon) error
	DeleteAllAddons(tx *sql.Tx) error
	FindAddonsWithDirectoryName(
		tx *sql.Tx,
		directory string) (
		[]model.CurseAddon,
		error,
	)
	GetCurseAddon(tx *sql.Tx, id int64) (model.CurseAddon, error)

	// curse releases
	GetCurseRelease(tx *sql.Tx, id int64) (model.CurseRelease, error)
}

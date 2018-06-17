package sqlite

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
	"github.com/wow-sweetlie/zhevra/storage/sqlite/migrations"

	// sqlite3 driver for migrations
	"github.com/golang-migrate/migrate/database/sqlite3"
	// go-sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// StorageSQLite define sorage type
type StorageSQLite struct {
	*sql.DB
}

func NewStorage(filename string) (*StorageSQLite, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}
	return &StorageSQLite{db}, nil
}

func (s *StorageSQLite) CreateAddon(tx *sql.Tx, addon *storage.Addon) (id int,
	err error) {
	err = tx.QueryRow(
		"INSERT INTO addon (name) VALUES(?)",
		addon.Name,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to CreateAddon tx")
	}
	return id, nil
}

func (s *StorageSQLite) CreateCurseAddon(tx *sql.Tx, addon *storage.Addon) error {
	_, err := tx.Exec(
		"INSERT INTO addon (name) VALUES(?)",
		addon.Name,
	)
	if err != nil {
		return errors.Wrap(err, "failed to CreateAddon tx")
	}
	return nil
}

// Migrate the database
func (s *StorageSQLite) Migrate() error {
	assets := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	source, err := bindata.WithInstance(assets)
	if err != nil {
		return err
	}
	database, err := sqlite3.WithInstance(s.DB, &sqlite3.Config{})
	m, err := migrate.NewWithInstance(
		"go-bindata",
		source,
		"sqlite3",
		database,
	)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}

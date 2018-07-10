package sqlite

import (
	"database/sql"

	"github.com/coline-carle/zhevra/storage/sqlite/migrations"
	"github.com/golang-migrate/migrate"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/pkg/errors"

	// sqlite3 driver for migrations
	"github.com/golang-migrate/migrate/database/sqlite3"
	// go-sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Storage define storage type
type Storage struct {
	*sql.DB
}

// NewStorage create a new storage using sqlite3 database
// filename name of the sqlite3 database file
func NewStorage(filename string) (*Storage, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}
	return &Storage{db}, nil
}

// Migrate the database
func (s *Storage) Migrate() error {
	assets := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	source, err := bindata.WithInstance(assets)
	if err != nil {
		return err
	}
	database, err := sqlite3.WithInstance(s.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance(
		"go-bindata",
		source,
		"sqlite3",
		database,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

// Tx wrapper for a new databse transaction
func (s *Storage) Tx(fn func(*sql.Tx) error) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	err = fn(tx)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, "transaction rollback failed")
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil

}

// Close close the database
func (s *Storage) Close() error {
	return s.DB.Close()
}

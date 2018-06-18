package sqlite

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage/sqlite/migrations"

	// sqlite3 driver for migrations
	"github.com/golang-migrate/migrate/database/sqlite3"
	// go-sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Storage define storage type
type Storage struct {
	*sql.DB
}

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

func (s *Storage) Tx(fn func(*sql.Tx) error) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	err = fn(tx)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			errors.Wrap(err2, "transaction rollback failed")
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit transaction")
	}
	return nil

}

func (s *Storage) Close() {
	s.DB.Close()
}

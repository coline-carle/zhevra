package sqlite

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/wow-sweetlie/zhevra/store/migrations"

	// sqlite3 driver for migrations
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	// go-sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// StorageSQLite define sorage type
type StorageSQLite struct {
	*sql.DB
}

func NewSqliteStorage(filename string) (*StorageSQLite, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return &StorageSQLite{db}, nil
}

func CreateAddon(tx *sql.Tx, addon *storage.Addon) {
}

// Migrate the database
func (s *StorageSQLite) Migrate() error {
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	source, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}
	database, err := postgres.WithInstance(s.DB, &postgres.Config{})
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

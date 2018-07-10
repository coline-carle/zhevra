package app

import (
	"github.com/coline-carle/zhevra/storage"
	"github.com/pkg/errors"
)

// App is the context of the main app
type App struct {
	storage storage.Storage
}

// NewApp init App struct, connect to storage and migrate
func NewApp() (*App, error) {
	storage, err := storage.NewStorage(DatabasePath())
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	err = storage.Migrate()
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	return &App{
		storage: storage,
	}, nil
}

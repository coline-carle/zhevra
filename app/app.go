package app

import (
	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
)

// App is the context of the main app
type App struct {
	config  Config
	storage *Storage
}

// NewApp init App struct, connect to storage and migrate
func NewApp(config *Config) (*App, err) {
	storage, err := storage.NewStorage(config.DatabasePath)
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	err = storage.Migrate()
	if err != nil {
		return nil, errors.Wrap(err, "NewApp")
	}
	return &App{
		config:  config,
		storage: storage,
	}, nil
}
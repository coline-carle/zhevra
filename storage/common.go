package storage

import "github.com/coline-carle/zhevra/storage/sqlite"

// NewStorage create the appropriate storage
// for a given url
func NewStorage(url string) (Storage, error) {
	return sqlite.NewStorage(url)
}

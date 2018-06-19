package sqlite

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
)

// CreateCurseAddon add a new addon from curse to the database
func (s *Storage) CreateCurseAddon(
	tx *sql.Tx, addon *storage.CurseAddon) error {
	_, err := tx.Exec(
		`INSERT INTO curse_addon (
			id,
			name,
			url,
			summary,
			downloadcount
		) VALUES(?, ?, ?, ?, ?)
		`,
		addon.ID,
		addon.Name,
		addon.URL,
		addon.Summary,
		addon.DownloadCount,
	)
	if err != nil {
		return errors.Wrap(err, "failed to CreateCurseAddon")
	}
	return nil
}

// GetCurseAddon fetch a curse addon from the databaze by ID
// id curse id of the addon
func (s *Storage) GetCurseAddon(
	tx *sql.Tx, id int64) (*storage.CurseAddon, error) {
	addon := &storage.CurseAddon{}
	err := tx.QueryRow(`
		SELECT
			id,
			name,
			url,
			summary,
			downloadcount
		FROM
			curse_addon
		WHERE
			id = $1
		`, id).Scan(
		&addon.ID,
		&addon.Name,
		&addon.URL,
		&addon.Summary,
		&addon.DownloadCount,
	)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			err = storage.ErrCurseAddonDoesNotExists
			return nil, errors.Wrapf(err, "id %d", id)
		}
		return nil, errors.Wrap(err, "GetCurseAddon failed")
	}
	return addon, nil
}

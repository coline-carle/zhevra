package sqlite

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
)

func (s *Storage) rowsToAddons(rows *sql.Rows) ([]storage.CurseAddon, error) {
	defer rows.Close()
	addons := []storage.CurseAddon{}
	for rows.Next() {
		addon := storage.CurseAddon{}
		err := rows.Scan(
			&addon.ID,
			&addon.Name,
			&addon.URL,
			&addon.Summary,
			&addon.DownloadCount,
		)
		if err != nil {
			return addons, errors.Wrap(err, "rowsToAddons")
		}
		addons = append(addons, addon)
	}
	return addons, nil
}

// CreateCurseAddon add a new addon from curse to the database
func (s *Storage) CreateCurseAddon(
	tx *sql.Tx, addon storage.CurseAddon) error {
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

	for _, release := range addon.Releases {
		release.AddonID = addon.ID
		err = s.CreateCurseRelease(tx, release)
		if err != nil {
			return errors.Wrap(err, "failed to create release from CreateCurseAddon")
		}
	}
	return nil
}

// FindAddonsWithDirectoryName find all addon for a given directory name
func (s *Storage) FindAddonsWithDirectoryName(
	tx *sql.Tx, directory string) ([]storage.CurseAddon, error) {
	addons := []storage.CurseAddon{}
	rows, err := tx.Query(`
		SELECT
			DISTINCT(curse_addon.id),
			curse_addon.name,
			curse_addon.url,
			curse_addon.summary,
			curse_addon.downloadcount
		FROM
			curse_addon
		INNER JOIN curse_release ON curse_addon.id = curse_release.addon_id
		INNER JOIN curse_release_directory ON
			curse_release_directory.release_id = curse_release.id
		WHERE
			curse_release_directory.directory = $1
		`, directory)
	if err != nil {
		return addons, errors.Wrap(err, "FindAddonsWithDirectoryName failed")
	}
	return s.rowsToAddons(rows)
}

// GetCurseAddon fetch a curse addon from the databaze by ID
// id curse id of the addon
func (s *Storage) GetCurseAddon(
	tx *sql.Tx, id int64) (storage.CurseAddon, error) {
	addon := storage.CurseAddon{}
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
			return addon, errors.Wrapf(err, "id %d", id)
		}
		return addon, errors.Wrap(err, "GetCurseAddon failed")
	}
	addon.Releases, err = s.FindCurseReleasesByAddonID(tx, id)
	return addon, err
}

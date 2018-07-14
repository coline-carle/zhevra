package sqlite

import (
	"database/sql"

	"github.com/coline-carle/zhevra/storage/model"
	"github.com/pkg/errors"
)

func (s *Storage) rowsToAddons(
	tx *sql.Tx, rows *sql.Rows) ([]model.CurseAddon, error) {
	defer rows.Close()
	addons := []model.CurseAddon{}
	for rows.Next() {
		addon := model.CurseAddon{}
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
		addon.Releases, err = s.FindCurseReleasesByAddonID(tx, addon.ID)
		addons = append(addons, addon)
	}
	return addons, nil
}

// DeleteAllAddons delete all addons from the database
func (s *Storage) DeleteAllAddons(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM curse_addon`)
	if err != nil {
		return errors.Wrap(err, "failed to DeleteAllAddons")
	}
	return nil
}

// CreateCurseAddon add a new addon from curse to the database
func (s *Storage) CreateCurseAddon(
	tx *sql.Tx, addon model.CurseAddon) error {
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
	tx *sql.Tx, directory string) ([]model.CurseAddon, error) {
	addons := []model.CurseAddon{}
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
	return s.rowsToAddons(tx, rows)
}

// GetCurseAddon fetch a curse addon from the databaze by ID
// id curse id of the addon
func (s *Storage) GetCurseAddon(
	tx *sql.Tx, id int64) (model.CurseAddon, error) {
	addon := model.CurseAddon{}
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
			err = model.ErrCurseAddonDoesNotExists
			return addon, errors.Wrapf(err, "id %d", id)
		}
		return addon, errors.Wrap(err, "GetCurseAddon failed")
	}
	addon.Releases, err = s.FindCurseReleasesByAddonID(tx, id)
	return addon, err
}

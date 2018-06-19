package sqlite

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
)

// CreateCurseRelease save new addon in the database
func (s *Storage) CreateCurseRelease(
	tx *sql.Tx, release storage.CurseRelease) error {
	_, err := tx.Exec(
		`INSERT INTO curse_release (
			id,
			filename,
			created_at,
			url,
			game_version,
			addon_id
		) VALUES(?, ?, ?, ?, ?, ?)
		`,
		release.ID,
		release.Filename,
		release.CreatedAt.Unix(),
		release.URL,
		release.GameVersion,
		release.AddonID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to CreateCurseRelease")
	}
	return nil
}

// GetCurseByAddonID return all release for a given addon ID
func (s *Storage) FindCurseReleasesByAddonID(
	tx *sql.Tx, id int64) ([]storage.CurseRelease, error) {
	releases := []storage.CurseRelease{}
	rows, err := tx.Query(`
		SELECT
			id,
			filename,
			created_at,
			url,
			game_version,
			addon_id
		FROM
			curse_release
		WHERE
			addon_id = $1
		`, id)
	if err != nil {
		return releases, errors.Wrap(err, "FindCurseReleasesByAddonID")
	}
	defer rows.Close()
	return rowToReleases(rows)
}

func rowToReleases(rows *sql.Rows) ([]storage.CurseRelease, error) {
	releases := []storage.CurseRelease{}
	var date int64
	for rows.Next() {
		release := storage.CurseRelease{}
		err := rows.Scan(
			&release.ID,
			&release.Filename,
			&date,
			&release.URL,
			&release.GameVersion,
			&release.AddonID,
		)
		if err != nil {
			return releases, errors.Wrap(err, "rowToReleases")
		}
		release.CreatedAt = time.Unix(date, 0).UTC()
		releases = append(releases, release)
	}
	return releases, nil
}

// GetCurseRelease fetch and addon by ID
func (s *Storage) GetCurseRelease(
	tx *sql.Tx, id int64) (storage.CurseRelease, error) {
	release := storage.CurseRelease{}
	var date int64
	err := tx.QueryRow(`
		SELECT
			id,
			filename,
			created_at,
			url,
			game_version,
			addon_id
		FROM
			curse_release
		WHERE
			id = $1
		`, id).Scan(
		&release.ID,
		&release.Filename,
		&date,
		&release.URL,
		&release.GameVersion,
		&release.AddonID,
	)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			err = storage.ErrCurseReleaseDoesNotExists
			return release, errors.Wrapf(err, "id %d", id)
		}
		return release, errors.Wrap(err, "GetCurseRelease failed")
	}
	release.CreatedAt = time.Unix(date, 0).UTC()
	return release, nil
}

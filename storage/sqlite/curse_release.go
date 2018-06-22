package sqlite

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/wow-sweetlie/zhevra/storage"
)

// CreateCurseReleaseDirectories create as many folder as we have
// it does not return an error if the folder already exist
// return a list of id for the coressponding folders
func (s *Storage) CreateCurseReleaseDirectories(
	tx *sql.Tx, release storage.CurseRelease) error {
	for _, directory := range release.Directories {
		_, err := tx.Exec(`
				INSERT INTO curse_release_directory
					(release_id, directory)
				VALUES(?, ?)`,
			release.ID, directory)
		if err != nil {
			return err
		}
	}
	return nil
}

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
			addon_id,
			is_alternate
		) VALUES(?, ?, ?, ?, ?, ?, ?)
		`,
		release.ID,
		release.Filename,
		release.CreatedAt.Unix(),
		release.URL,
		release.GameVersion,
		release.AddonID,
		release.IsAlternate,
	)
	if err != nil {
		return errors.Wrap(err, "failed to CreateCurseRelease")
	}
	return s.CreateCurseReleaseDirectories(tx, release)
}

// FindCurseReleaseDirectoriesByReleaseID return all directories for a
// given release
func (s *Storage) FindCurseReleaseDirectoriesByReleaseID(
	tx *sql.Tx, id int64) ([]string, error) {
	directories := []string{}
	rows, err := tx.Query(`
		SELECT
			directory
		FROM
			curse_release_directory
		WHERE
			release_id = $1
		`, id)
	if err != nil {
		return directories, errors.Wrap(err, "FindCurseReleaseDirecotries failed")
	}
	defer rows.Close()
	return rowsToStringSlice(rows)
}

// FindCurseReleasesByAddonID return all release for a given addon ID
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
			addon_id,
			is_alternate
		FROM
			curse_release
		WHERE
			addon_id = $1
		`, id)
	if err != nil {
		return releases, errors.Wrap(err, "FindCurseReleasesByAddonID")
	}
	defer rows.Close()
	return s.rowsToReleases(tx, rows)
}

func rowsToStringSlice(rows *sql.Rows) ([]string, error) {
	directories := []string{}
	var directory string
	for rows.Next() {
		err := rows.Scan(&directory)
		if err != nil {
			return directories, errors.Wrap(err, "rowsToStringSlice")
		}
		directories = append(directories, directory)
	}
	return directories, nil
}

func (s *Storage) rowsToReleases(
	tx *sql.Tx, rows *sql.Rows) ([]storage.CurseRelease, error) {
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
			&release.IsAlternate,
		)
		if err != nil {
			return releases, errors.Wrap(err, "rowsToReleases")
		}
		release.CreatedAt = time.Unix(date, 0).UTC()
		release.Directories, err = s.FindCurseReleaseDirectoriesByReleaseID(tx, release.ID)
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
			addon_id,
			is_alternate
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
		&release.IsAlternate,
	)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			err = storage.ErrCurseReleaseDoesNotExists
			return release, errors.Wrapf(err, "id %d", id)
		}
		return release, errors.Wrap(err, "GetCurseRelease failed")
	}
	release.Directories, err = s.FindCurseReleaseDirectoriesByReleaseID(tx, release.ID)
	release.CreatedAt = time.Unix(date, 0).UTC()
	return release, err
}

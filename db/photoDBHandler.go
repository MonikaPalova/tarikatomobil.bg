package db

import (
	"database/sql"
	"errors"
	. "github.com/MonikaPalova/tarikatomobil.bg/model"
)

type PhotoDBHandler struct {
	conn *sql.DB
}

func (pdb PhotoDBHandler) UploadPhoto(p *Photo) *DBError {
	insertQuery := `INSERT INTO PHOTOS (id, bytes) VALUES (?, ?)`
	stmt, err := pdb.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	_, err = stmt.Exec(p.ID, p.Base64Content)
	if err != nil {
		if isDuplicateEntryError(err) {
			return NewDBError(err, ErrConflict)
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (pdb PhotoDBHandler) GetPhotoByID(photoID string) (*Photo, *DBError) {
	row := pdb.conn.QueryRow("SELECT id, bytes FROM PHOTOS WHERE id = ?", photoID)
	var p Photo
	if err := row.Scan(&p.ID, &p.Base64Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewDBError(err, ErrNotFound)
		}
		return nil, NewDBError(err, ErrInternal)
	}
	return &p, nil
}

func (pdb PhotoDBHandler) DeletePhoto(photoID string) *DBError {
	stmt, err := pdb.conn.Prepare("DELETE FROM photos WHERE id = ?")
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	if _, err = stmt.Exec(photoID); err != nil {
		return NewDBError(err, ErrInternal)
	}
	return nil
}

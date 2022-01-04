package db

import (
	"database/sql"
	. "github.com/MonikaPalova/tarikatomobil.bg/model"
)

type PhotoDBHandler struct {
	conn *sql.DB
}

func (pdb PhotoDBHandler) UploadPhoto(p *Photo) error {
	insertQuery := `INSERT INTO PHOTOS (id, bytes, extension) VALUES (?, ?, ?)`
	stmt, err := pdb.conn.Prepare(insertQuery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(p.ID, p.Base64Content, p.Extension)
	return err
}

func (pdb PhotoDBHandler) GetPhotoByID(photoID string) (*Photo, error) {
	row := pdb.conn.QueryRow("SELECT * FROM PHOTOS WHERE id = ?", photoID)
	var p Photo
	if err := row.Scan(&p.ID, &p.Base64Content, &p.Extension); err != nil {
		return nil, err
	}
	return &p, nil
}

func (pdb PhotoDBHandler) DeletePhoto(photoID string) error {
	stmt, err := pdb.conn.Prepare("DELETE FROM photos WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(photoID)
	return err
}

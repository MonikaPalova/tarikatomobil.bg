package db

import (
	"database/sql"
	. "github.com/MonikaPalova/tarikatomobil.bg/model"
)

type PhotoDBHandler struct {
	conn *sql.DB
}

func (pdb PhotoDBHandler) UploadPhoto(p *Photo) {

}
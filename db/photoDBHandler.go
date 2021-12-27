package db

import "database/sql"

type PhotoDBHandler struct {
	Conn *sql.Conn
}

func (p PhotoDBHandler) UploadPhoto() {

}
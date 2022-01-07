package db

import (
	"database/sql"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

const (
	createTablesFile = "./sql/create_tables.sql"
	defaultPhotoFile = "./resources/unknown_photo.bin"
)

type Database struct {
	conn             *sql.DB
	UsersDBHandler   *UsersDBHandler
	PhotosDBHandler  *PhotoDBHandler
	SessionDBHandler *SessionDBHandler
}

func InitDB(user, password, dbName string) (*Database, error) {
	// Connect to DB
	connString := fmt.Sprintf("%s:%s@/%s?multiStatements=true&parseTime=true", user, password, dbName)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	// Create tables if needed
	var createTableQuery []byte
	if createTableQuery, err = ioutil.ReadFile(createTablesFile); err != nil {
		return nil, err
	}
	if _, err = conn.Exec(string(createTableQuery)); err != nil {
		return nil, err
	}

	// Upload the default photo if it is not uploaded yet
	defaultPhoto := model.Photo{
		ID: model.DefaultPhotoID,
	}
	var defaultPhotoBytes []byte
	if defaultPhotoBytes, err = ioutil.ReadFile(defaultPhotoFile); err != nil {
		return nil, fmt.Errorf("could not load default photo from %s: %s", defaultPhotoFile, err.Error())
	}
	defaultPhoto.Base64Content = string(defaultPhotoBytes)

	// Fill and return a Database struct
	db := Database{
		conn:             conn,
		UsersDBHandler:   &UsersDBHandler{conn: conn},
		PhotosDBHandler:  &PhotoDBHandler{conn: conn},
		SessionDBHandler: &SessionDBHandler{conn: conn},
	}

	dbErr := db.PhotosDBHandler.UploadPhoto(&defaultPhoto)
	if dbErr != nil && dbErr.ErrorType != ErrConflict {
		return nil, fmt.Errorf("could not upload the default photo to the database: %s", dbErr.Err.Error())
	}

	return &db, nil
}

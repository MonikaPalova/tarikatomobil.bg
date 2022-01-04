package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

const (
	createTablesFile = "./sql/create_tables.sql"
)

type Database struct {
	conn            *sql.DB
	UsersDBHandler  *UsersDBHandler
	PhotosDBHandler *PhotoDBHandler
}

func InitDB(user, password, dbName string) (*Database, error) {
	// Connect to DB
	connString := fmt.Sprintf("%s:%s@/%s?multiStatements=true", user, password, dbName)
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

	// Fill and return a Database struct
	db := Database{
		conn:            conn,
		UsersDBHandler:  &UsersDBHandler{conn: conn},
		PhotosDBHandler: &PhotoDBHandler{conn: conn},
	}

	return &db, nil
}

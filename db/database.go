package db

import (
	"database/sql"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

type Database struct {
	conn *sql.DB
}

func (d *Database) Connect(user, password, dbName string) error {
	var err error

	connString := fmt.Sprintf("%s:%s@/%s?multiStatements=true", user, password, dbName)
	d.conn, err = sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	return d.createTables()
}

func (d Database) createTables() error {
	createTableQuery, err := ioutil.ReadFile("./sql/create_tables.sql")
	if err != nil {
		return err
	}
	_, err = d.conn.Exec(string(createTableQuery))
	fmt.Println(string(createTableQuery))
	return err
}

func (d Database) GetUsers() ([]model.User, error) {
	rows, err := d.conn.Query("SELECT * FROM USERS")
	if err != nil {
		return nil, err
	}
	users := []model.User{}

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d Database) CreateUser(user model.User) error {
	stmt, err := d.conn.Prepare("INSERT INTO USERS (id, name) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.ID, user.Name)
	return err
}

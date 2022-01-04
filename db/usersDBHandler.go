package db

import (
	"database/sql"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/go-sql-driver/mysql"
)

type UsersDBHandler struct {
	conn *sql.DB
}

func (uh *UsersDBHandler) GetUsers() ([]model.User, *DBError) {
	rows, err := uh.conn.Query("SELECT * FROM USERS")
	if err != nil {
		return nil, NewDBError(err, ErrInternal)
	}
	users := []model.User{}

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, NewDBError(err, ErrInternal)
		}
		users = append(users, user)
	}
	return users, nil
}

func (uh *UsersDBHandler) CreateUser(user model.User) *DBError {
	insertQuery := `
		INSERT INTO USERS (id, name, password, email, phone_number, photo_id, times_passenger, times_driver) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)	`
	stmt, err := uh.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	photoID := &user.PhotoID
	if len(*photoID) == 0 {
		photoID = nil
	}
	_, err = stmt.Exec(user.ID, user.Name, user.Password, user.Email, user.PhoneNumber, photoID, user.TimesPassenger, user.TimesDriver)
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == mysqlDuplicateEntryCode {
			return NewDBError(err, ErrConflict)
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

package db

import (
	"database/sql"
	"errors"
	. "github.com/MonikaPalova/tarikatomobil.bg/model"
)

type UsersDBHandler struct {
	conn *sql.DB
}

func (uh *UsersDBHandler) GetUserByName(username string) (UserWithoutPass, *DBError) {
	row := uh.conn.QueryRow(`
		SELECT name, email, phone_number, photo_id, times_passenger, times_driver
		FROM USERS WHERE name = ?`, username)
	var u UserWithoutPass
	if err := row.Scan(&u.Name, &u.Email, &u.PhoneNumber, &u.PhotoID, &u.TimesPassenger, &u.TimesDriver); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, NewDBError(err, ErrNotFound)
		}
		return u, NewDBError(err, ErrInternal)
	}
	return u, nil
}

func (uh *UsersDBHandler) CreateUser(user User) *DBError {
	insertQuery := `
		INSERT INTO USERS (name, password, email, phone_number, photo_id, times_passenger, times_driver) 
		VALUES(?, ?, ?, ?, ?, ?, ?)	`
	stmt, err := uh.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	photoID := &user.PhotoID
	if len(*photoID) == 0 {
		photoID = nil
	}
	_, err = stmt.Exec(user.Name, user.Password, user.Email, user.PhoneNumber, photoID, user.TimesPassenger, user.TimesDriver)
	if err != nil {
		if isDuplicateEntryError(err) {
			return NewDBError(err, ErrConflict)
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"
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
			return u, NewDBError(err, ErrNotFound, fmt.Sprintf("user %s does not exist", username))
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

	_, err = stmt.Exec(user.Name, user.Password, user.Email, user.PhoneNumber, user.PhotoID, user.TimesPassenger, user.TimesDriver)
	if err != nil {
		if isDuplicateEntryError(err) {
			return NewDBError(err, ErrConflict, fmt.Sprintf("user %s already exists", user.Name))
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (uh *UsersDBHandler) UpdateUser(username string, userPatch UserPatch) *DBError {
	updateQuery := `UPDATE users SET users.password = IF(?='', users.password, ?), 
 								  users.photo_id = IF(?='', users.photo_id, ?),
                  				  users.email = IF(?='', users.email, ?),
                  				  users.phone_number = IF(?='', users.phone_number, ?) WHERE users.name = ?`
	stmt, err := uh.conn.Prepare(updateQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	_, err = stmt.Exec(
		userPatch.Password, userPatch.Password,
		userPatch.PhotoID, userPatch.PhotoID,
		userPatch.Email, userPatch.Email,
		userPatch.PhoneNumber, userPatch.PhoneNumber,
		username)
	if err != nil {
		if isForeignKeyError(err) { // If the photo with given ID does not exist
			return NewDBError(err, ErrConflict, fmt.Sprintf("a photo with ID %s does not exist", userPatch.PhotoID))
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

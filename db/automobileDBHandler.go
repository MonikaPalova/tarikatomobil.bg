package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type AutomobileDBHandler struct {
	conn *sql.DB
}

func (adb *AutomobileDBHandler) GetUserAutomobile(username string) (model.Automobile, *DBError) {
	row := adb.conn.QueryRow("SELECT reg_num, photo_id, comment, owner_name FROM automobiles WHERE owner_name = ?",
		username)
	var a model.Automobile
	if err := row.Scan(&a.RegNumber, &a.PhotoID, &a.Comment, &a.OwnerName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return a, NewDBError(err, ErrNotFound, fmt.Sprintf("user %s does not have an automobile", username))
		}
		return a, NewDBError(err, ErrInternal)
	}
	return a, nil
}

func (adb *AutomobileDBHandler) CreateAutomobile(a model.Automobile) *DBError {
	stmt, err := adb.conn.Prepare("INSERT INTO automobiles (reg_num, photo_id, comment, owner_name) VALUES(?,?,?,?)")
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	if _, err = stmt.Exec(a.RegNumber, a.PhotoID, a.Comment, a.OwnerName); err != nil {
		if isDuplicateEntryError(err) {
			return NewDBError(err, ErrConflict,
				fmt.Sprintf("automobile with number %s already exists or user %s already has an automobile",
					a.RegNumber, a.OwnerName))
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

func (adb *AutomobileDBHandler) UpdateAutomobile(patch model.AutomobilePatch, username string) *DBError {
	updateQuery := `UPDATE automobiles SET automobiles.photo_id = IF(?='', automobiles.photo_id, ?),
		       						   	   automobiles.comment = IF(?='', automobiles.comment, ?)
									   WHERE automobiles.owner_name = ?`

	stmt, err := adb.conn.Prepare(updateQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	_, err = stmt.Exec(patch.PhotoID, patch.PhotoID, patch.Comment, patch.Comment, username)
	if err != nil {
		if isForeignKeyError(err) { // If the photo with given ID does not exist
			return NewDBError(err, ErrConflict, fmt.Sprintf("a photo with ID %s does not exist", patch.PhotoID))
		}
		return NewDBError(err, ErrInternal)
	}
	return nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type SessionDBHandler struct {
	conn *sql.DB
}

func (sdb *SessionDBHandler) GetSessionByID(sessionID string) (*model.Session, *DBError) {
	row := sdb.conn.QueryRow("SELECT id, owner, expiration FROM sessions WHERE id = ?", sessionID)
	var session model.Session
	if err := row.Scan(&session.ID, &session.Owner, &session.Expiration); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewDBError(err, ErrNotFound, fmt.Sprintf("session with ID %s does not exist", sessionID))
		}
		return nil, NewDBError(err, ErrInternal)
	}
	return &session, nil
}

func (sdb *SessionDBHandler) InsertSession(session model.Session) *DBError {
	insertQuery := `INSERT INTO sessions (id, owner, expiration) VALUES (?, ?, ?)`
	stmt, err := sdb.conn.Prepare(insertQuery)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}

	_, err = stmt.Exec(session.ID, session.Owner, session.Expiration)
	if err != nil {
		return NewDBError(err, ErrInternal)
	}
	return nil
}

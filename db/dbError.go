package db

import "github.com/go-sql-driver/mysql"

const (
	mysqlDuplicateEntryCode   = 1062
	mysqlForeignKeyErrorCode1 = 1451
	mysqlForeignKeyErrorCode2 = 1452
)

type DBErrorType uint8

// Error types enum
const (
	ErrNotFound DBErrorType = iota // When an entity is not found
	ErrConflict                    // When a constraint is not satisfied
	ErrInternal                    // For driver and connection errors
)

type DBError struct {
	Err         error
	UserMessage string // The DB error exposes DB internals (it's logged) - this message shall be used as an API response
	ErrorType   DBErrorType
}

func NewDBError(err error, errorType DBErrorType, userMessage ...string) *DBError {
	dbErr := DBError{
		Err:         err,
		ErrorType:   errorType,
		UserMessage: "",
	}
	if len(userMessage) > 0 {
		dbErr.UserMessage = userMessage[0]
	}
	return &dbErr
}

func isDuplicateEntryError(err error) bool {
	return checkForSpecificError(err, mysqlDuplicateEntryCode)
}

func isForeignKeyError(err error) bool {
	return checkForSpecificError(err, mysqlForeignKeyErrorCode1) || checkForSpecificError(err, mysqlForeignKeyErrorCode2)
}

func checkForSpecificError(err error, errorCode uint16) bool {
	driverErr, ok := err.(*mysql.MySQLError)
	return ok && driverErr.Number == errorCode
}

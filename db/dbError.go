package db

import "github.com/go-sql-driver/mysql"

const (
	mysqlDuplicateEntryCode  = 1062
	mysqlForeignKeyErrorCode = 1452
)

type DBErrorType uint8

// Error types enum
const (
	ErrNotFound DBErrorType = iota // When an entity is not found
	ErrConflict                    // When a constraint is not satisfied
	ErrInternal                    // For driver and connection errors
)

type DBError struct {
	Err       error
	ErrorType DBErrorType
}

func NewDBError(err error, errorType DBErrorType) *DBError {
	dbErr := DBError{
		Err:       err,
		ErrorType: errorType,
	}
	return &dbErr
}

func isDuplicateEntryError(err error) bool {
	return checkForSpecificError(err, mysqlDuplicateEntryCode)
}

func isForeignKeyError(err error) bool {
	return checkForSpecificError(err, mysqlForeignKeyErrorCode)
}

func checkForSpecificError(err error, errorCode uint16) bool {
	if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == errorCode {
		return true
	}
	return false
}

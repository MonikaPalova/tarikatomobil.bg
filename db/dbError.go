package db

const (
	mysqlDuplicateEntryCode = 1062
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

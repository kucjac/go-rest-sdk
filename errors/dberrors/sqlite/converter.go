package sqlite

import (
	"database/sql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/mattn/go-sqlite3"
)

// SQLiteConverter is DBErrorConverter interface implementation
// for sqlite3 database.
//
type SQLiteConverter struct {
	errorMap map[interface{}]dbe.DBError
}

// Convert converts the provided error into *DBError type.
// It is method that implements DBErrorConverter Interface
func (r *SQLiteConverter) Convert(err error) *dbe.DBError {
	// Check if the error is of '*sqlite3.Error' type
	sqliteErr, ok := err.(sqlite3.Error)
	if !ok {
		// if not check sql errors
		if err == sql.ErrNoRows {
			return dbe.ErrNoResult.NewWithError(err)
		} else if err == sql.ErrTxDone {
			return dbe.ErrTxDone.NewWithError(err)
		}
		return dbe.ErrUnspecifiedError.NewWithError(err)
	}

	var dbError dbe.DBError
	// Check if Error.ExtendedCode is in recogniser
	dbError, ok = r.errorMap[sqliteErr.ExtendedCode]
	if ok {
		return dbError.NewWithError(err)
	}

	// otherwise check if Error.Code is in the recogniser
	dbError, ok = r.errorMap[sqliteErr.Code]
	if ok {
		return dbError.NewWithError(err)
	}

	// if no error is specified return Unspecified Error
	return dbe.ErrUnspecifiedError.NewWithError(err)
}

func New() *SQLiteConverter {
	return &SQLiteConverter{errorMap: defaultSQLiteErrorMap}
}

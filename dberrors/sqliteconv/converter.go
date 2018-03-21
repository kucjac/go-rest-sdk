package sqliteconv

import (
	"database/sql"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/mattn/go-sqlite3"
)

// SQLiteConverter is ErrorConverter interface implementation
// for sqlite3 database.
//
type SQLiteConverter struct {
	errorMap map[interface{}]dberrors.Error
}

// Convert converts the provided error into *Error type.
// It is method that implements ErrorConverter Interface
func (r *SQLiteConverter) Convert(err error) *dberrors.Error {
	// Check if the error is of '*sqlite3.Error' type
	sqliteErr, ok := err.(sqlite3.Error)
	if !ok {
		// if not check sql errors
		if err == sql.ErrNoRows {
			return dberrors.ErrNoResult.NewWithError(err)
		} else if err == sql.ErrTxDone {
			return dberrors.ErrTxDone.NewWithError(err)
		}
		return dberrors.ErrUnspecifiedError.NewWithError(err)
	}

	var dbError dberrors.Error
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
	return dberrors.ErrUnspecifiedError.NewWithError(err)
}

func New() *SQLiteConverter {
	return &SQLiteConverter{errorMap: defaultSQLiteErrorMap}
}

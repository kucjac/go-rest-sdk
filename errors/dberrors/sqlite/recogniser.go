package sqlite

import (
	"database/sql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/mattn/go-sqlite3"
)

// SQLiteRecogniser implements Recogniser interface
type SQLiteErrorRecogniser map[interface{}]*dbe.DBError

func (r SQLiteErrorRecogniser) Recognise(err error) error {
	// Check if the error is of '*sqlite3.Error' type
	sqliteErr, ok := err.(*sqlite3.Error)
	if !ok {
		// if not check sql errors
		if err == sql.ErrNoRows {
			return dbe.ErrNoResult
		} else if err == sql.ErrTxDone {
			return dbe.ErrTxDone
		}
		return err
	}

	var dbError *dbe.DBError

	// Check if Error.ExtendedCode is in recogniser
	dbError, ok = r[sqliteErr.ExtendedCode]
	if ok {
		return dbError
	}

	// otherwise check if Error.Code is in the recogniser
	dbError, ok = r[sqliteErr.Code]
	if ok {
		return dbError
	}

	// if no error is specified return Unspecified Error
	return dbe.ErrUnspecifiedError
}

package postgres

import (
	"database/sql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/lib/pq"
)

type PostgresErrorRecogniser map[interface{}]*dbe.DBError

// Recognise - implements DBErrorRecogniser
func (p PostgresErrorRecogniser) Recognise(err error) error {
	pgError, ok := err.(*pq.Error)
	if !ok {
		// The error may be of sql.ErrNoRows type
		if err == sql.ErrNoRows {
			return dbe.ErrNoResult
		} else if err == sql.ErrTxDone {
			return dbe.ErrTxDone
		}
		return err
	}

	var dbError *dbe.DBError
	// First check if recogniser has entire error code in it
	dbError, ok = p[pgError.Code]
	if ok {
		return dbError
	}

	// If the ErrorCode is not present, check the code class
	dbError, ok = p[pgError.Code.Class()]
	if ok {
		return dbError
	}

	// If the Error Class is not presen in the error map
	// return ErrDBNotMapped
	return dbe.ErrUnspecifiedError
}

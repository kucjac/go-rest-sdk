package postgres

import (
	"database/sql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/lib/pq"
)

// PostgresErrorConverter is an implementation of DBErrorConverter.
type PostgresErrorConverter map[interface{}]*dbe.DBError

// Convert converts the given error into *DBError.
// The method checks if given error is of known type, and then returns it.ty
// If an error is unknown it returns new 'dberrors.ErrUnspecifiedError'.
// At first converter checks if an error is of *pq.Error type.
// Having a postgres *pq.Error it checks if an ErrorCode is in the map,
// and returns it if true. Otherwise method checks if the ErrorClass exists in map.
// If it is present, new *DBError of given type is returned.
func (p PostgresErrorConverter) Convert(err error) (dbeErr *dbe.DBError) {
	pgError, ok := err.(*pq.Error)
	if !ok {
		// The error may be of sql.ErrNoRows type
		if err == sql.ErrNoRows {
			dbeErr = dbe.ErrNoResult.New()
			dbeErr.Message = err.Error()
			return
		} else if err == sql.ErrTxDone {
			dbeErr = dbe.ErrTxDone.New()
			dbeErr.Message = err.Error()
			return
		}
		dbeErr = dbe.ErrUnspecifiedError.New()
		dbeErr.Error()
		return
	}

	var (
		dbError      *dbe.DBError
		dbErrorProto dbe.DBError
	)

	// First check if recogniser has entire error code in it
	dbErrorProto, ok = p[pgError.Code]
	if ok {

		return dbErrorProto.NewWithMessage(pgError.Error())
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

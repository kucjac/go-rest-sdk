package postgres

import (
	"database/sql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/lib/pq"
)

// PostgresErrorConverter is an implementation of DBErrorConverter.
type PGConverter struct {
	errorMap map[interface{}]dbe.DBError
}

// Convert converts the given error into *DBError.
// The method checks if given error is of known type, and then returns it.ty
// If an error is unknown it returns new 'dberrors.ErrUnspecifiedError'.
// At first converter checks if an error is of *pq.Error type.
// Having a postgres *pq.Error it checks if an ErrorCode is in the map,
// and returns it if true. Otherwise method checks if the ErrorClass exists in map.
// If it is present, new *DBError of given type is returned.
func (p *PGConverter) Convert(err error) (dbeErr *dbe.DBError) {
	pgError, ok := err.(*pq.Error)
	if !ok {
		// The error may be of sql.ErrNoRows type
		if err == sql.ErrNoRows {
			return dbe.ErrNoResult.NewWithError(err)
		} else if err == sql.ErrTxDone {
			return dbe.ErrTxDone.NewWithError(err)
		}
		return dbe.ErrUnspecifiedError.NewWithError(err)

	}

	// DBError prototype
	var dbErrorProto dbe.DBError

	// First check if recogniser has entire error code in it
	dbErrorProto, ok = p.errorMap[pgError.Code]
	if ok {
		return dbErrorProto.NewWithError(err)
	}

	// If the ErrorCode is not present, check the code class
	dbErrorProto, ok = p.errorMap[pgError.Code.Class()]
	if ok {
		return dbErrorProto.NewWithError(err)
	}

	// If the Error Class is not presen in the error map
	// return ErrDBNotMapped
	return dbe.ErrUnspecifiedError.NewWithError(err)
}

// New creates new PGConverter
// It is already inited and ready to use.
func New() *PGConverter {
	return &PGConverter{errorMap: defaultPGErrorMap}
}

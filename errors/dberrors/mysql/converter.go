package mysql

import (
	"database/sql"
	msql "github.com/go-sql-driver/mysql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
)

// MySQLConverter is a DBErrorConverter interface implementation
// The Converter can convert provided error into *DBError with specific logic.
// Check the Convert method documentation for more information on how it distinguish given error
type MySQLConverter struct {
	// codeMap puts an error code or sqlstate into map and returns dberrors.DBError prototype
	codeMap map[interface{}]dbe.DBError

	// sqlStateMap is a helper map that recognises the sqlstate for given error code
	sqlStateMap map[uint16]string
}

// Convert converts provided 'err' error into *dbe.DBError type.
// With this method MySQLConverter implements DBErrorConverter interface.
// Convert distinguish  and convert specific error of types sql.Err*, msql.Err*,
// and *msql.MySQLError. If an error is of different type it returns new entity of
// dberrors.ErrUnspecifiedError
// If the error is of *msql.MySQLError type the method checks its code.
// If the code matches with internal code map it returns proper entity of *dbe.DBError.
// If the code does not exists in the code map, the method gets sqlstate for given code
// and checks if this sqlstate is in the code map.
// If the sqlstate does not exists in the code map, the first two numbers from the sqlstate
// are being checked in the codeMap as a 'sqlstate class'.
// If not found Convert returns new entity for dberrors.UnspecifiedError
func (m *MySQLConverter) Convert(err error) *dbe.DBError {
	// Check whether the given error is of *msql.MySQLError
	mySQLErr, ok := err.(*msql.MySQLError)
	if !ok {
		// Otherwise check if it sql.Err* or other errors from mysql package
		switch err {
		case msql.ErrInvalidConn, msql.ErrNoTLS, msql.ErrOldProtocol,
			msql.ErrMalformPkt:
			return dbe.ErrConnExc.NewWithError(err)
		case sql.ErrNoRows:
			return dbe.ErrNoResult.NewWithError(err)

		case sql.ErrTxDone:
			return dbe.ErrTxDone.NewWithError(err)

		default:
			return dbe.ErrUnspecifiedError.NewWithError(err)
		}
	}
	var dbErr dbe.DBError

	// Check if Error Number is in recogniser
	dbErr, ok = m.codeMap[mySQLErr.Number]
	if ok {
		// Return if found
		return dbErr.NewWithError(err)
	}

	// Otherwise check if given sqlstate is in the codeMap
	sqlState, ok := m.sqlStateMap[mySQLErr.Number]
	if !ok || len(sqlState) != 5 {
		return dbe.ErrUnspecifiedError.NewWithError(err)
	}
	dbErr, ok = m.codeMap[sqlState]
	if ok {
		return dbErr.NewWithError(err)
	}

	// First two letter from sqlState represents error class
	// Check if class is in error map
	sqlStateClass := sqlState[0:2]
	dbErr, ok = m.codeMap[sqlStateClass]
	if ok {
		return dbErr.NewWithError(err)
	}

	return dbe.ErrUnspecifiedError.NewWithError(err)
}

// New creates new already inited MySQLConverter
func New() *MySQLConverter {
	return &MySQLConverter{
		codeMap:     mysqlErrMap,
		sqlStateMap: codeSQLState,
	}
}

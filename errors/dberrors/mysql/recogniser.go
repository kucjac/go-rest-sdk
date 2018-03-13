package mysql

import (
	"database/sql"
	msql "github.com/go-sql-driver/mysql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
)

type MySQLErrorRecogniser struct {
	// codeMap puts an error code or sqlstate into map and returns *dbe.DBError
	codeMap map[interface{}]*dbe.DBError

	// sqlStateMap is a helper map that recognises the sqlstate from error code
	sqlStateMap map[uint16]string
}

func (m *MySQLErrorRecogniser) Recognise(err error) error {
	var dbErr *dbe.DBError
	// Check whether the given error is of *msql.MySQLError
	mySQLErr, ok := err.(*msql.MySQLError)
	if !ok {
		// Otherwise check if it sql.Err* or other errors from mysql package
		switch err {
		case msql.ErrInvalidConn, msql.ErrNoTLS, msql.ErrOldProtocol,
			msql.ErrMalformPkt, sql.ErrConnDone:
			dbErr = dbe.ErrConnExc
			dbErr.Message += " " + err.Error()
			return dbErr
		case sql.ErrNoRows:
			dbErr = dbe.ErrNoResult
			return dbErr
		case sql.ErrTxDone:
			dbErr = dbe.ErrTxDone
			return dbErr
		default:
			dbErr = dbe.ErrUnspecifiedError
			dbErr.Message += " " + err.Error()
			return dbErr
		}
	}

	// Check if Error Number is in recogniser
	dbErr, ok = m.codeMap[mySQLErr.Number]
	if ok {
		// Return if found
		return dbErr
	}

	// Otherwise check if given sqlstate is in the codeMap
	sqlState, ok := m.sqlStateMap[mySQLErr.Number]
	if !ok || len(sqlState) != 5 {
		return dbe.ErrUnspecifiedError
	}
	dbErr, ok = m.codeMap[sqlState]
	if ok {
		return dbErr
	}

	// First two letter from sqlState represents error class
	// Check if class is in error map
	sqlStateClass := sqlState[0:2]
	dbErr, ok = m.codeMap[sqlStateClass]
	if ok {
		return dbErr
	}

	return dbe.ErrUnspecifiedError
}

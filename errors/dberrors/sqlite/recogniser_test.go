package sqlite

import (
	"database/sql"
	"errors"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSQLiteRecogniser(t *testing.T) {
	Convey("Using a SQLite3 Error Recogniser", t, func() {

		sqliteRecogniser := SQLiteRecogniser

		Convey("Having a list of sqlite errors", func() {

			sqliteErrors := map[*sqlite3.Error]*dbe.DBError{
				{Code: sqlite3.ErrWarning}:  dbe.ErrWarning,
				{Code: sqlite3.ErrNotFound}: dbe.ErrNoResult,
				{Code: sqlite3.ErrCantOpen}: dbe.ErrConnExc,
				{Code: sqlite3.ErrNotADB}:   dbe.ErrConnExc,
				{Code: sqlite3.ErrMismatch}: dbe.ErrDataException,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintPrimaryKey}: dbe.ErrIntegrConstViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintFunction}: dbe.ErrIntegrConstViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintCheck}: dbe.ErrCheckViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintForeignKey}: dbe.ErrForeignKeyViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintUnique}: dbe.ErrUniqueViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintNotNull}: dbe.ErrNotNullViolation,
				{Code: sqlite3.ErrProtocol}: dbe.ErrInvalidTransState,
				{Code: sqlite3.ErrRange}:    dbe.ErrInvalidSyntax,
				{Code: sqlite3.ErrError}:    dbe.ErrInvalidSyntax,
				{Code: sqlite3.ErrAuth}:     dbe.ErrInvalidAuthorization,
				{Code: sqlite3.ErrPerm}:     dbe.ErrInsufficientPrivilege,
				{Code: sqlite3.ErrFull}:     dbe.ErrInsufficientResources,
				{Code: sqlite3.ErrTooBig}:   dbe.ErrProgramLimitExceeded,
				{Code: sqlite3.ErrNoLFS}:    dbe.ErrSystemError,
				{Code: sqlite3.ErrInternal}: dbe.ErrInternalError,
				{Code: sqlite3.ErrEmpty}:    dbe.ErrUnspecifiedError,
			}

			Convey("For given *sqlite.Error, specific database error should be returner.", func() {
				for sqliteErr, dbErr := range sqliteErrors {
					recognisedErr := sqliteRecogniser.Recognise(sqliteErr)
					So(recognisedErr, ShouldEqual, dbErr)
				}
			})
		})

		Convey("Having an error of type sql.Err*, error is converted into *DBError type.",
			func() {
				var err error
				Println("No rows")
				err = sql.ErrNoRows
				recognisedErr := sqliteRecogniser.Recognise(err)
				So(recognisedErr, ShouldEqual, dbe.ErrNoResult)

				Println("Tx done")
				err = sql.ErrTxDone
				recognisedErr = sqliteRecogniser.Recognise(err)
				So(recognisedErr, ShouldEqual, dbe.ErrTxDone)
			})

		Convey("Having an error of different type than *sqlite3.Error and sql.Err*", func() {
			err := errors.New("Unknown error type")
			recognisedErr := sqliteRecogniser.Recognise(err)
			So(recognisedErr, ShouldEqual, dbe.ErrUnspecifiedError)
		})
	})
}

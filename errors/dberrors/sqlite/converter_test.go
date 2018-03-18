package sqlite

import (
	"database/sql"
	"errors"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewConverter(t *testing.T) {
	Convey("Using New function creates non-empty SQLiteConverter.", t, func() {
		var converter *SQLiteConverter
		converter = New()

		So(converter, ShouldNotBeNil)

		So(len(converter.errorMap), ShouldBeGreaterThan, 0)

		Convey("The SQLiteConverter implements DBErrorConverter", func() {
			So(converter, ShouldImplement, (*dbe.DBErrorConverter)(nil))
		})
	})
}

func TestSQLiteRecogniser(t *testing.T) {
	Convey("Using a SQLite3 Error Converter", t, func() {

		var converter *SQLiteConverter = New()

		Convey("Having a list of sqlite errors", func() {

			sqliteErrors := map[*sqlite3.Error]dbe.DBError{
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
					recognisedErr := converter.Convert(sqliteErr)
					So(recognisedErr.Compare(dbErr), ShouldBeTrue)
				}
			})
		})

		Convey("Having an error of type sql.Err*, error is converted into *DBError type.", func() {
			var err error
			err = sql.ErrNoRows
			recognisedErr := converter.Convert(err)
			So(recognisedErr.Compare(dbe.ErrNoResult), ShouldBeTrue)

			err = sql.ErrTxDone
			recognisedErr = converter.Convert(err)
			So(recognisedErr.Compare(dbe.ErrTxDone), ShouldBeTrue)
		})

		Convey("Having an error of different type than *sqlite3.Error and sql.Err*", func() {
			err := errors.New("Unknown error type")
			recognisedErr := converter.Convert(err)
			So(recognisedErr.Compare(dbe.ErrUnspecifiedError), ShouldBeTrue)
		})
	})
}

package sqliteconv

import (
	"database/sql"
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
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

		Convey("The SQLiteConverter implements Converter", func() {
			So(converter, ShouldImplement, (*dberrors.Converter)(nil))
		})
	})
}

func TestSQLiteRecogniser(t *testing.T) {
	Convey("Using a SQLite3 Error Converter", t, func() {

		var converter *SQLiteConverter = New()

		Convey("Having a list of sqlite errors", func() {

			sqliteErrors := map[sqlite3.Error]dberrors.Error{
				{Code: sqlite3.ErrWarning}:  dberrors.ErrWarning,
				{Code: sqlite3.ErrNotFound}: dberrors.ErrNoResult,
				{Code: sqlite3.ErrCantOpen}: dberrors.ErrConnExc,
				{Code: sqlite3.ErrNotADB}:   dberrors.ErrConnExc,
				{Code: sqlite3.ErrMismatch}: dberrors.ErrDataException,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintPrimaryKey}: dberrors.ErrUniqueViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintFunction}: dberrors.ErrIntegrConstViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintCheck}: dberrors.ErrCheckViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintForeignKey}: dberrors.ErrForeignKeyViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintUnique}: dberrors.ErrUniqueViolation,
				{Code: sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintNotNull}: dberrors.ErrNotNullViolation,
				{Code: sqlite3.ErrProtocol}: dberrors.ErrInvalidTransState,
				{Code: sqlite3.ErrRange}:    dberrors.ErrInvalidSyntax,
				{Code: sqlite3.ErrError}:    dberrors.ErrInvalidSyntax,
				{Code: sqlite3.ErrAuth}:     dberrors.ErrInvalidAuthorization,
				{Code: sqlite3.ErrPerm}:     dberrors.ErrInsufficientPrivilege,
				{Code: sqlite3.ErrFull}:     dberrors.ErrInsufficientResources,
				{Code: sqlite3.ErrTooBig}:   dberrors.ErrProgramLimitExceeded,
				{Code: sqlite3.ErrNoLFS}:    dberrors.ErrSystemError,
				{Code: sqlite3.ErrInternal}: dberrors.ErrInternalError,
				{Code: sqlite3.ErrEmpty}:    dberrors.ErrUnspecifiedError,
			}

			Convey("For given *sqlite.Error, specific database error should be returner.", func() {
				for sqliteErr, dbErr := range sqliteErrors {
					recognisedErr := converter.Convert(sqliteErr)
					Println(sqliteErr)
					Println(recognisedErr)
					Println(dbErr)
					So(recognisedErr.Compare(dbErr), ShouldBeTrue)
				}
			})
		})

		Convey("Having an error of type sql.Err*, error is converted into *Error type.", func() {
			var err error
			err = sql.ErrNoRows
			recognisedErr := converter.Convert(err)
			So(recognisedErr.Compare(dberrors.ErrNoResult), ShouldBeTrue)

			err = sql.ErrTxDone
			recognisedErr = converter.Convert(err)
			So(recognisedErr.Compare(dberrors.ErrTxDone), ShouldBeTrue)
		})

		Convey("Having an error of different type than *sqlite3.Error and sql.Err*", func() {
			err := errors.New("Unknown error type")
			recognisedErr := converter.Convert(err)
			So(recognisedErr.Compare(dberrors.ErrUnspecifiedError), ShouldBeTrue)
		})
	})
}

package postgres

import (
	"database/sql"
	"errors"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewConverter(t *testing.T) {
	Convey("While using New function the converter is already inited.", t, func() {
		var converter *PGConverter
		converter = New()
		So(converter, ShouldNotBeNil)
		So(len(converter.errorMap), ShouldBeGreaterThan, 0)
	})
}

func TestPGRecogniser(t *testing.T) {
	Convey("Using Postgress Converter", t, func() {
		var converter *PGConverter = New()

		Convey("Having a list of typical postgres errors", func() {
			postgresErrors := map[*pq.Error]dbe.DBError{
				{Code: pq.ErrorCode("01000")}: dbe.ErrWarning,
				{Code: pq.ErrorCode("01007")}: dbe.ErrWarning,
				{Code: pq.ErrorCode("02000")}: dbe.ErrNoResult,
				{Code: pq.ErrorCode("P0002")}: dbe.ErrNoResult,
				{Code: pq.ErrorCode("08006")}: dbe.ErrConnExc,
				{Code: pq.ErrorCode("21000")}: dbe.ErrCardinalityViolation,
				{Code: pq.ErrorCode("22012")}: dbe.ErrDataException,
				{Code: pq.ErrorCode("23000")}: dbe.ErrIntegrConstViolation,
				{Code: pq.ErrorCode("23001")}: dbe.ErrRestrictViolation,
				{Code: pq.ErrorCode("23502")}: dbe.ErrNotNullViolation,
				{Code: pq.ErrorCode("23503")}: dbe.ErrForeignKeyViolation,
				{Code: pq.ErrorCode("23505")}: dbe.ErrUniqueViolation,
				{Code: pq.ErrorCode("23514")}: dbe.ErrCheckViolation,
				{Code: pq.ErrorCode("25001")}: dbe.ErrInvalidTransState,
				{Code: pq.ErrorCode("25004")}: dbe.ErrInvalidTransState,
				{Code: pq.ErrorCode("28000")}: dbe.ErrInvalidAuthorization,
				{Code: pq.ErrorCode("28P01")}: dbe.ErrInvalidPassword,
				{Code: pq.ErrorCode("2D000")}: dbe.ErrInvalidTransTerm,
				{Code: pq.ErrorCode("3F000")}: dbe.ErrInvalidSchemaName,
				{Code: pq.ErrorCode("40000")}: dbe.ErrTransRollback,
				{Code: pq.ErrorCode("42P06")}: dbe.ErrInvalidSyntax,
				{Code: pq.ErrorCode("42501")}: dbe.ErrInsufficientPrivilege,
				{Code: pq.ErrorCode("53100")}: dbe.ErrInsufficientResources,
				{Code: pq.ErrorCode("54011")}: dbe.ErrProgramLimitExceeded,
				{Code: pq.ErrorCode("58000")}: dbe.ErrSystemError,
				{Code: pq.ErrorCode("XX000")}: dbe.ErrInternalError,
				{Code: pq.ErrorCode("P0003")}: dbe.ErrUnspecifiedError,
			}

			Convey("For given postgres error, specific database error should return", func() {
				for pgErr, dbErr := range postgresErrors {
					convertedErr := converter.Convert(pgErr)
					So(convertedErr.Compare(dbErr), ShouldBeTrue)
				}
			})
		})

		Convey("When sql errors are returned, they are also converted into dberror", func() {
			errNoResults := converter.Convert(sql.ErrNoRows)
			So(errNoResults.Compare(dbe.ErrNoResult), ShouldBeTrue)

			errTxDone := converter.Convert(sql.ErrTxDone)
			So(errTxDone.Compare(dbe.ErrTxDone), ShouldBeTrue)
		})

		Convey("Having unknown error not of *pq.Error type forwards it", func() {

			fwdErr := converter.Convert(errors.New("Forwarded"))
			So(fwdErr.Compare(dbe.ErrUnspecifiedError), ShouldBeTrue)

		})
	})

}

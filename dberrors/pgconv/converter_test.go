package pgconv

import (
	"database/sql"
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
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
		So(converter, ShouldImplement, (*dberrors.Converter)(nil))
	})
}

func TestPGRecogniser(t *testing.T) {
	Convey("Using Postgress Converter", t, func() {
		var converter *PGConverter = New()

		Convey("Having a list of typical postgres errors", func() {
			postgresErrors := map[*pq.Error]dberrors.Error{
				{Code: pq.ErrorCode("01000")}: dberrors.ErrWarning,
				{Code: pq.ErrorCode("01007")}: dberrors.ErrWarning,
				{Code: pq.ErrorCode("02000")}: dberrors.ErrNoResult,
				{Code: pq.ErrorCode("P0002")}: dberrors.ErrNoResult,
				{Code: pq.ErrorCode("08006")}: dberrors.ErrConnExc,
				{Code: pq.ErrorCode("21000")}: dberrors.ErrCardinalityViolation,
				{Code: pq.ErrorCode("22012")}: dberrors.ErrDataException,
				{Code: pq.ErrorCode("23000")}: dberrors.ErrIntegrConstViolation,
				{Code: pq.ErrorCode("23001")}: dberrors.ErrRestrictViolation,
				{Code: pq.ErrorCode("23502")}: dberrors.ErrNotNullViolation,
				{Code: pq.ErrorCode("23503")}: dberrors.ErrForeignKeyViolation,
				{Code: pq.ErrorCode("23505")}: dberrors.ErrUniqueViolation,
				{Code: pq.ErrorCode("23514")}: dberrors.ErrCheckViolation,
				{Code: pq.ErrorCode("25001")}: dberrors.ErrInvalidTransState,
				{Code: pq.ErrorCode("25004")}: dberrors.ErrInvalidTransState,
				{Code: pq.ErrorCode("28000")}: dberrors.ErrInvalidAuthorization,
				{Code: pq.ErrorCode("28P01")}: dberrors.ErrInvalidPassword,
				{Code: pq.ErrorCode("2D000")}: dberrors.ErrInvalidTransTerm,
				{Code: pq.ErrorCode("3F000")}: dberrors.ErrInvalidSchemaName,
				{Code: pq.ErrorCode("40000")}: dberrors.ErrTransRollback,
				{Code: pq.ErrorCode("42P06")}: dberrors.ErrInvalidSyntax,
				{Code: pq.ErrorCode("42501")}: dberrors.ErrInsufficientPrivilege,
				{Code: pq.ErrorCode("53100")}: dberrors.ErrInsufficientResources,
				{Code: pq.ErrorCode("54011")}: dberrors.ErrProgramLimitExceeded,
				{Code: pq.ErrorCode("58000")}: dberrors.ErrSystemError,
				{Code: pq.ErrorCode("XX000")}: dberrors.ErrInternalError,
				{Code: pq.ErrorCode("P0003")}: dberrors.ErrUnspecifiedError,
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
			So(errNoResults.Compare(dberrors.ErrNoResult), ShouldBeTrue)

			errTxDone := converter.Convert(sql.ErrTxDone)
			So(errTxDone.Compare(dberrors.ErrTxDone), ShouldBeTrue)
		})

		Convey("Having unknown error not of *pq.Error type forwards it", func() {

			fwdErr := converter.Convert(errors.New("Forwarded"))
			So(fwdErr.Compare(dberrors.ErrUnspecifiedError), ShouldBeTrue)

		})
	})

}

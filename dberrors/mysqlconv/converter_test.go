package mysqlconv

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/kucjac/go-rest-sdk/dberrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNew(t *testing.T) {
	Convey("The 'New()' function creates new already inited '*MySQLConverter' entity", t, func() {
		var converter *MySQLConverter
		converter = New()

		So(converter, ShouldNotBeNil)

		So(len(converter.codeMap), ShouldBeGreaterThan, 0)
		So(len(converter.sqlStateMap), ShouldBeGreaterThan, 0)

		Convey("The *MySQLConverter implements dberrors.Converter interface", func() {
			So(converter, ShouldImplement, (*dberrors.Converter)(nil))
		})
	})
}

func TestMySQLRecogniser(t *testing.T) {
	Convey("Having MySQLConverter.", t, func() {
		var converter *MySQLConverter = New()

		Convey("Check if selected MySQL Errors would return for given dberrors Error", func() {
			errorMap := map[*mysql.MySQLError]dberrors.Error{
				{Number: 1022}: dberrors.ErrUniqueViolation,
				{Number: 1046}: dberrors.ErrInvalidCatalogName,
				{Number: 1048}: dberrors.ErrNotNullViolation,
				{Number: 1050}: dberrors.ErrInvalidSyntax,
				{Number: 1062}: dberrors.ErrUniqueViolation,
				{Number: 1114}: dberrors.ErrProgramLimitExceeded,
				{Number: 1118}: dberrors.ErrProgramLimitExceeded,
				{Number: 1129}: dberrors.ErrInternalError,
				{Number: 1130}: dberrors.ErrInvalidAuthorization,
				{Number: 1131}: dberrors.ErrInvalidAuthorization,
				{Number: 1132}: dberrors.ErrInvalidPassword,
				{Number: 1133}: dberrors.ErrInvalidPassword,
				{Number: 1169}: dberrors.ErrUniqueViolation,
				{Number: 1182}: dberrors.ErrTransRollback,
				{Number: 1216}: dberrors.ErrForeignKeyViolation,
				{Number: 1217}: dberrors.ErrForeignKeyViolation,
				{Number: 1227}: dberrors.ErrInsufficientPrivilege,
				{Number: 1251}: dberrors.ErrInvalidAuthorization,
				{Number: 1400}: dberrors.ErrInvalidTransState,
				{Number: 1401}: dberrors.ErrInternalError,
				{Number: 1451}: dberrors.ErrForeignKeyViolation,
				{Number: 1452}: dberrors.ErrForeignKeyViolation,
				{Number: 1557}: dberrors.ErrUniqueViolation,
				{Number: 1568}: dberrors.ErrUniqueViolation,
				{Number: 1698}: dberrors.ErrInvalidPassword,
				//Nil
				{Number: 1317}: dberrors.ErrUnspecifiedError,
				{Number: 1040}: dberrors.ErrConnExc,
				//Non mapped number
				{Number: 1000}: dberrors.ErrUnspecifiedError,
			}

			for msqlErr, dbErr := range errorMap {
				dbErrInMap := converter.Convert(msqlErr)

				So(dbErrInMap.Compare(dbErr), ShouldBeTrue)
			}
		})
		Convey("Having error of different type than *mysql.Error", func() {
			errorMap := map[error]dberrors.Error{
				sql.ErrNoRows:           dberrors.ErrNoResult,
				sql.ErrTxDone:           dberrors.ErrTxDone,
				mysql.ErrInvalidConn:    dberrors.ErrConnExc,
				mysql.ErrNoTLS:          dberrors.ErrConnExc,
				mysql.ErrMalformPkt:     dberrors.ErrConnExc,
				mysql.ErrOldProtocol:    dberrors.ErrConnExc,
				mysql.ErrNativePassword: dberrors.ErrUnspecifiedError,
			}

			for err, dbErr := range errorMap {
				dbErrInMap := converter.Convert(err)
				// Printf("%v: %v\n", err, dbErrInMap)
				So(dbErrInMap.Compare(dbErr), ShouldBeTrue)
			}
		})
	})
}

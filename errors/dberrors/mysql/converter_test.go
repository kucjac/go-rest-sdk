package mysql

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
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

		Convey("The *MySQLConverter implements dberrors.DBErrorConverter interface", func() {
			So(converter, ShouldImplement, (*dbe.DBErrorConverter)(nil))
		})
	})
}

func TestMySQLRecogniser(t *testing.T) {
	Convey("Having MySQLConverter.", t, func() {
		var converter *MySQLConverter = New()

		Convey("Check if selected MySQL Errors would return for given dberrors DBError", func() {
			errorMap := map[*mysql.MySQLError]dbe.DBError{
				{Number: 1022}: dbe.ErrUniqueViolation,
				{Number: 1046}: dbe.ErrInvalidCatalogName,
				{Number: 1048}: dbe.ErrNotNullViolation,
				{Number: 1050}: dbe.ErrInvalidSyntax,
				{Number: 1062}: dbe.ErrUniqueViolation,
				{Number: 1114}: dbe.ErrProgramLimitExceeded,
				{Number: 1118}: dbe.ErrProgramLimitExceeded,
				{Number: 1129}: dbe.ErrInternalError,
				{Number: 1130}: dbe.ErrInvalidAuthorization,
				{Number: 1131}: dbe.ErrInvalidAuthorization,
				{Number: 1132}: dbe.ErrInvalidPassword,
				{Number: 1133}: dbe.ErrInvalidPassword,
				{Number: 1169}: dbe.ErrUniqueViolation,
				{Number: 1182}: dbe.ErrTransRollback,
				{Number: 1216}: dbe.ErrForeignKeyViolation,
				{Number: 1217}: dbe.ErrForeignKeyViolation,
				{Number: 1227}: dbe.ErrInsufficientPrivilege,
				{Number: 1251}: dbe.ErrInvalidAuthorization,
				{Number: 1400}: dbe.ErrInvalidTransState,
				{Number: 1401}: dbe.ErrInternalError,
				{Number: 1451}: dbe.ErrForeignKeyViolation,
				{Number: 1452}: dbe.ErrForeignKeyViolation,
				{Number: 1557}: dbe.ErrUniqueViolation,
				{Number: 1568}: dbe.ErrUniqueViolation,
				{Number: 1698}: dbe.ErrInvalidPassword,
				//Nil
				{Number: 1317}: dbe.ErrUnspecifiedError,
				{Number: 1040}: dbe.ErrConnExc,
				//Non mapped number
				{Number: 1000}: dbe.ErrUnspecifiedError,
			}

			for msqlErr, dbErr := range errorMap {
				dbErrInMap := converter.Convert(msqlErr)

				So(dbErrInMap.Compare(dbErr), ShouldBeTrue)
			}
		})
		Convey("Having error of different type than *mysql.Error", func() {
			errorMap := map[error]dbe.DBError{
				sql.ErrNoRows:           dbe.ErrNoResult,
				sql.ErrTxDone:           dbe.ErrTxDone,
				mysql.ErrInvalidConn:    dbe.ErrConnExc,
				mysql.ErrNoTLS:          dbe.ErrConnExc,
				mysql.ErrMalformPkt:     dbe.ErrConnExc,
				mysql.ErrOldProtocol:    dbe.ErrConnExc,
				mysql.ErrNativePassword: dbe.ErrUnspecifiedError,
			}

			for err, dbErr := range errorMap {
				dbErrInMap := converter.Convert(err)
				// Printf("%v: %v\n", err, dbErrInMap)
				So(dbErrInMap.Compare(dbErr), ShouldBeTrue)
			}
		})
	})
}

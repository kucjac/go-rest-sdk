package gormrepo

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGORMConverterInit(t *testing.T) {

	Convey("Subject: Creating new GORMErrorConverter and initialize it with Init method", t, func() {

		Convey("Having *GORMErrorConverter entity and some *gorm.DB connections", func() {
			var errorConverter *GORMErrorConverter

			dbSqlite, _ := gorm.Open("sqlite3", "./tests.db")
			dbPostgres, _ := gorm.Open("postgres", "host=myhost port=myport")
			dbMySQL, _ := gorm.Open("mysql", "user:password@/dbname")
			dbMSSQL, _ := gorm.Open("mssql", "sqlserver://username:password@localhost:1433?database=dbname")

			var dbNil *gorm.DB

			gormSupported := []*gorm.DB{dbSqlite, dbPostgres, dbMySQL}

			Convey("While using Init method", func() {
				var err error

				Convey("If the dialect is supported, specific converter would be set", func() {
					for _, db := range gormSupported {
						errorConverter = new(GORMErrorConverter)
						err = errorConverter.Init(db)
						So(err, ShouldBeNil)
					}
				})

				Convey("If the dialect is unsupported an error would be returned", func() {
					errorConverter = new(GORMErrorConverter)
					err = errorConverter.Init(dbMSSQL)
					So(err, ShouldBeError)
				})

				Convey("If provided nil pointer an error would be thrown.", func() {
					errorConverter = new(GORMErrorConverter)
					err = errorConverter.Init(dbNil)
					So(err, ShouldBeError)
				})

			})

		})

	})

}

func TestGORMErrorConverterConvert(t *testing.T) {

	Convey("Subject: Converting an error into *DBError using GormErrorConverter method Convert", t, func() {

		Convey("Having inited GORMErrorConverter", func() {
			db, _ := gorm.Open("sqlite3", "./tests.db")

			errorConverter := new(GORMErrorConverter)
			err := errorConverter.Init(db)

			So(err, ShouldBeNil)

			Convey("Providing any error would result with *DBerror", func() {
				convertErrors := []error{gorm.ErrCantStartTransaction,
					gorm.ErrInvalidTransaction,
					gorm.ErrInvalidSQL,
					gorm.ErrUnaddressable,
					gorm.ErrRecordNotFound,
					dberrors.ErrCardinalityViolation.New(),
					dberrors.ErrWarning.New(),
					errors.New("Some error"),
					sql.ErrNoRows,
				}

				for _, err := range convertErrors {
					converted := errorConverter.Convert(err)
					So(converted, ShouldHaveSameTypeAs, &dberrors.DBError{})
				}
			})

		})

	})

}

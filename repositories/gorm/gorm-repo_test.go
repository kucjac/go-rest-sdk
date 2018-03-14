package gorm

import (
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitGORMRepository(t *testing.T) {
	Convey("Having a GORM pointer for sqlite3 database", t, func() {
		db, err := gorm.Open("sqlite3", "test.db")
		if err != nil {
			panic(err.Error())
		}

		Convey("And GORMRepository based on it", func() {
			repo := GORMRepository{}
			err = repo.Init(db)

			So(err, ShouldBeNil)
		})
	})
}

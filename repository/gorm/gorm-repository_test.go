package gormrepo

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var modelMigrated bool

func TestNewGORMRepository(t *testing.T) {

	Convey("Subject: Creating New GORMRepository with New() function", t, func() {

		Convey("Having *gorm.DB", func() {
			var err error
			var gormRepo *GORMRepository
			db, _ := gorm.Open("sqlite3", "./tests.db")
			defer db.Close()
			Convey("If the *gorm.DB is supported the function should Create new GORMRepository", func() {
				gormRepo, err = New(db)
				So(err, ShouldBeNil)
				So(gormRepo, ShouldNotBeNil)
				So(gormRepo.converter, ShouldNotBeNil)

			})

			db = nil
			Convey("The function should return error if nil *gorm.DB pointer  provided as an argument", func() {
				gormRepo, err := New(db)

				So(err, ShouldBeError)
				So(gormRepo, ShouldBeNil)
			})

			db, _ = gorm.Open("mssql", "sqlserver://username:password@localhost:1433?database=dbname")

			Convey("Should return error if *gorm.DB dialect is unknown", func() {
				gormRepo, err := New(db)

				So(err, ShouldBeError)
				So(gormRepo, ShouldBeNil)
			})
		})
	})
}

func TestGORMRepositoryCreate(t *testing.T) {

	Convey("Subject: Creating new model record in the database", t, func() {

		Convey("Using a Create method on GORMRepository with some model entity as an argument", func() {

			db, err := openGormSqlite()
			So(err, ShouldBeNil)

			defer db.Close()

			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			foo := Foo{Name: "Foo"}

			// bar := Bar{Name: "Bar"}

			// fooBar := Foo{Name: "FooBar", Bar: &bar}

			Convey("Should Create new model record in the database based on the gorm rules.", func() {

				So(db.NewRecord(foo), ShouldBeTrue)
				gormRepo.Create(foo)
				So(err, ShouldBeNil)

				So(db.NewRecord(foo), ShouldBeFalse)

				// So(db.NewRecord(fooBar), ShouldBeTrue)
				// err = gormRepo.Create(&fooBar)
				// So(err, ShouldBeNil)

				// So(db.NewRecord(fooBar), ShouldBeFalse)
			})

			Convey("Should throw *DBError if an error occurs during the creation process", nil)

		})

	})

}

func TestGORMRepositoryGet(t *testing.T) {

	Convey("Subject: Getting model record from the database", t, func() {

		Convey("Using a Get method on GORMRepository with some model as an argument.", func() {

			Convey("Should return appropiate model entity", nil)

			Convey("Should throw *DBError if an error occurs", nil)

		})

	})

}

func TestGORMRepositorySelect(t *testing.T) {

	Convey("Subject: List all records for given request model.", t, func() {

		Convey("Using a List method on GORMRepository with provided argument model", func() {

			Convey("Should list all records for given 'req' restrictions.", nil)

			Convey("Should not list records from sets disjoint for given 'req' restrictions", nil)

			Convey("Should return *DBError if an error occurs", nil)

		})

	})

}

func TestGORMRepositorySelectWithParams(t *testing.T) {

	Convey("Subject: List all records for given request model with provided *ListParameters", t, func() {

		Convey("Using a ListWithParams method on GORMRepository with provided model  and listParameters as an argument", func() {

			Convey("Should list all records for given 'req' restricions and listParameters", nil)

			Convey("Should not list records from sets disjoint for given restrictions", nil)

			Convey("Should return *DBError if an error occurs", nil)

		})

	})

}

func TestGORMRepositoryUpdate(t *testing.T) {

	Convey("Subject: Updates whole model entity so that any field would be replaced", t, func() {

		Convey("Using an Update method on GORMRepository with provided 'req' model entity.", func() {

			Convey("Should update whole records for given restrictions", nil)

			Convey("Should not update only provided fields", nil)

			Convey("Should return *DBError if an error occurs", nil)

		})

	})

}

func TestGORMRepositoryPatch(t *testing.T) {

	Convey("Subject: Updates the database record field from 'what' that matches given on 'where' restriction.", t, func() {

		Convey("Using a Patch method on GORMRepository with provided 'what' fields to change on 'where' restrictions .", func() {

			Convey("Should update all records 'what' fields that  match given 'where'", nil)

			Convey("Should not update all record for given 'where'", nil)

			Convey("Should return *DBError if an error occurs", nil)

		})

	})

}

func TestGORMRepositoryDelete(t *testing.T) {

	Convey("Subject: Deleting the database records that matches the 'req' object.", t, func() {

		Convey("Using a Delete method on GORMRepository for given 'req' restrictions", func() {

			Convey("Should delete all files that matches the 'req' restrictions", nil)

			Convey("Should not delete files for disjoint restrictions", nil)

			Convey("Should return *DBError if an error occurs", nil)

		})

	})

}

func openGormSqlite() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./tests.db")
	if err != nil {
		return nil, err
	}
	if !modelMigrated {
		err = migrateModels(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func migrateModels(db *gorm.DB) error {
	if db != nil {
		db.AutoMigrate(&Bar{}, &Foo{})
		modelMigrated = true
		Println("Migrated")
		return nil
	}
	return errors.New("Nil pointer provided")
}

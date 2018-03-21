package gormrepo

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/repository"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Foo struct {
	ID    uint
	Name  string
	Bar   *Bar `gorm:"foreignkey:BarID"`
	BarID uint `gorm:"index"`
}

type Bar struct {
	ID       uint
	Name     string
	Property int
}

type Foobar struct {
	ID   uint
	Name string `gorm:"type:varchar(5);unique"`
}

type NotInDB struct{}

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

			bar := Bar{Name: "Bar"}

			fooBar := Foo{Name: "FooBar", Bar: &bar}

			Convey("Should Create new model record in the database based on the gorm rules.", func() {

				So(db.NewRecord(foo), ShouldBeTrue)

				err = gormRepo.Create(&foo)
				So(err, ShouldBeNil)
				So(db.NewRecord(foo), ShouldBeFalse)

				So(db.NewRecord(bar), ShouldBeTrue)
				err = gormRepo.Create(&bar)
				So(err, ShouldBeNil)
				So(db.NewRecord(bar), ShouldBeFalse)

				So(db.NewRecord(fooBar), ShouldBeTrue)
				err = gormRepo.Create(&fooBar)
				So(err, ShouldBeNil)
				So(db.NewRecord(fooBar), ShouldBeFalse)

			})
			Convey(`Should throw *dberrors.Error if an error occurs 
				during the creation process`, func() {
				So(db.NewRecord(foo), ShouldBeTrue)
				err = gormRepo.Create(&foo)
				So(err, ShouldBeNil)

				err = gormRepo.Create(&foo)
				So(err, ShouldBeError)
				So(err, ShouldHaveSameTypeAs, &dberrors.Error{})

				dbErr := err.(*dberrors.Error)
				proto, err := dbErr.GetPrototype()
				So(err, ShouldBeNil)
				So(proto, ShouldResemble, dberrors.ErrUniqueViolation)
			})
		})

	})
}

func TestGORMRepositoryGet(t *testing.T) {

	Convey("Subject: Getting model record from the database", t, func() {

		Convey("Using a Get method on GORMRepository with some model as an argument.", func() {
			db, err := openGormSqlite()
			So(err, ShouldBeNil)
			defer db.Close()
			defer clearDB(db)

			Convey("Having some foo type Foo and bar type Bar objects", func() {
				var bar Bar = Bar{Name: "Some Bar"}

				db.Create(&bar)
				So(err, ShouldBeNil)

				var foo Foo = Foo{Name: "Some name", Bar: &bar}
				db.Create(&foo)

				gormRepo, _ := New(db)

				Convey("Should return appropiate model entity", func() {
					var req Foo = Foo{ID: foo.ID}
					getted, err := gormRepo.Get(&req)
					So(err, ShouldBeNil)

					// Getting Foo object as 'foo'

					getFoo, ok := getted.(*Foo)
					So(ok, ShouldBeTrue)

					So(getFoo.ID, ShouldEqual, foo.ID)
					So(getFoo.Name, ShouldEqual, foo.Name)

					// Getting Bar object as 'bar'

					getted, err = gormRepo.Get(&Bar{ID: bar.ID})
					So(err, ShouldBeNil)

					getBar, ok := getted.(*Bar)
					So(ok, ShouldBeTrue)

					So(*getBar, ShouldResemble, bar)

					Convey("Given object should equal the object inserted before", func() {
						getFoo.Bar = getBar
						So(*getFoo, ShouldResemble, foo)
					})

				})

				Convey("Should throw *Error if an error occurs", func() {

					res, err := gormRepo.Get(&Foo{ID: 99999})
					So(err, ShouldBeError)
					So(res, ShouldBeNil)
					So(err, ShouldHaveSameTypeAs, &dberrors.Error{})

				})
			})

		})

	})

}

func TestGORMRepositorySelect(t *testing.T) {

	Convey("Subject: List all records for given request model.", t, func() {

		Convey("Using a List method on GORMRepository with provided argument model", func() {

			db, err := openGormSqlite()
			So(err, ShouldBeNil)
			defer db.Close()
			defer clearDB(db)

			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			someBars := []*Bar{{Name: "First"}, {Name: "Second"}, {Name: "Third"}}

			for _, bar := range someBars {
				db.Create(bar)
			}

			defer Convey("Should list all records for given 'req' restrictions.", func() {
				res, err := gormRepo.List(&Bar{})
				So(err, ShouldBeNil)

				list, ok := res.([]*Bar)
				So(ok, ShouldBeTrue)

				for i := range someBars {
					So(list, ShouldContain, someBars[i])
				}

			})

			Convey("Should not list records from sets disjoint for given 'req' restrictions", func() {
				res, err := gormRepo.List(&Bar{Name: someBars[1].Name})
				So(err, ShouldBeNil)

				list, ok := res.([]*Bar)
				So(ok, ShouldBeTrue)

				So(list, ShouldContain, someBars[1])
				So(list, ShouldNotContain, someBars[0])
				So(list, ShouldNotContain, someBars[2])

			})

			Convey("Should return *Error if an error occurs", func() {

				res, err := gormRepo.List(&NotInDB{})
				So(err, ShouldBeError)
				So(res, ShouldBeNil)

				_, ok := err.(*dberrors.Error)
				So(ok, ShouldBeTrue)

			})

		})

	})

}

func TestGORMRepositorySelectWithParams(t *testing.T) {

	Convey(`Subject: List all records for given request model 
		with provided *ListParameters`, t, func() {

		Convey(`Using a ListWithParams method on GORMRepository with provided model  and 
			listParameters as an argument`, func() {

			// Open Sqlite connection
			db, err := openGormSqlite()
			So(err, ShouldBeNil)
			defer db.Close()
			defer clearDB(db)

			// Get GormRepo
			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			// Settle database
			var someBars []*Bar = seedBars(db)

			Convey(`Should list all records for given 'req' restricions 
				and listParameters`, func() {
				var params *repository.ListParameters

				Convey("Using Limit paramter should limit the result rows returned", func() {
					var limit int = 3
					params = &repository.ListParameters{Limit: limit}
					res, err := gormRepo.ListWithParams(&Bar{}, params)
					So(err, ShouldBeNil)

					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					So(len(list), ShouldBeLessThanOrEqualTo, limit)
					So(list, ShouldContain, someBars[0])
					So(list, ShouldContain, someBars[1])
					So(list, ShouldContain, someBars[2])

					Convey(`Should not list records from sets disjoint for given restrictions`, func() {
						So(list, ShouldNotContain, someBars[3])
						So(list, ShouldNotContain, someBars[4])
						So(list, ShouldNotContain, someBars[5])
					})
				})
				Convey(`Using limit and offset parameter should limit the 
					result rows return shifted by the offset value`, func() {
					var limit, offset int = 3, 2
					var bars []*Bar
					// db.Find(&bars, &Bar{})

					db.Offset(offset).Limit(limit).Find(&bars, &Bar{})

					params = &repository.ListParameters{Limit: limit, Offset: offset}

					res, err := gormRepo.ListWithParams(&Bar{}, params)
					So(err, ShouldBeNil)

					// List should be of []*Bar type
					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					Println(list)
					// Limit limts the returned rows
					So(len(list), ShouldBeLessThanOrEqualTo, limit)

					// Offset makes the results shifted by its value
					So(list, ShouldNotContain, someBars[0])
					So(list, ShouldNotContain, someBars[1])

					So(list, ShouldContain, someBars[2])
					So(list, ShouldContain, someBars[3])
					So(list, ShouldContain, someBars[4])
				})

				Convey(`Using the ids field for struct with primary key as int type
					should select only those entites with provided ids`, func() {
					var indices []int = []int{2, 4, 5}
					params = &repository.ListParameters{IDs: indices}

					res, err := gormRepo.ListWithParams(&Bar{}, params)
					So(err, ShouldBeNil)

					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					So(len(list), ShouldEqual, len(indices))
					for _, id := range indices {
						So(list, ShouldContain, someBars[id-1])
					}
				})

				Convey("Using the Order param orders the query", func() {
					var order string = "name desc"
					params = &repository.ListParameters{Order: order}

					res, err := gormRepo.ListWithParams(&Bar{}, params)
					So(err, ShouldBeNil)

					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					// The list should be sorted by name in a descending manner
					for i := 0; i < len(list)-1; i++ {
						So(list[i].Name[0], ShouldBeGreaterThanOrEqualTo, list[i+1].Name[0])
					}
				})
				Convey(`Using no params should 	return full List`, func() {
					res, err := gormRepo.ListWithParams(&Bar{}, &repository.ListParameters{})
					So(err, ShouldBeNil)

					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					So(len(list), ShouldEqual, len(someBars))
				})

				Convey("Using param nil pointer should return full list for given type", func() {
					res, err := gormRepo.ListWithParams(&Bar{}, params)
					So(err, ShouldBeNil)

					list, ok := res.([]*Bar)
					So(ok, ShouldBeTrue)

					So(len(list), ShouldEqual, len(someBars))
				})
			})

			Convey("Should return *Error if an error occurs", func() {
				res, err := gormRepo.ListWithParams(&NotInDB{},
					&repository.ListParameters{Limit: 10})
				So(err, ShouldBeError)
				So(res, ShouldBeNil)

				_, ok := err.(*dberrors.Error)
				So(ok, ShouldBeTrue)

			})

		})

	})

}

func TestGORMRepositoryUpdate(t *testing.T) {

	Convey("Subject: Updates whole model entity so that any field would be replaced", t, func() {
		db, err := openGormSqlite()
		So(err, ShouldBeNil)

		defer db.Close()
		defer clearDB(db)

		var bars []*Bar = seedBars(db)
		Convey(`Using an Update method on GORMRepository 
			with provided 'req' model entity.`, func() {

			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			Convey("Should update whole records for given restrictions", func() {
				Convey(`By providing only one property and ID, the given update method would 
					clear other properties`, func() {
					var bar, updated *Bar = bars[0], &Bar{Name: "FirstUpdated"}

					// In order to keep ID the updated entity must contain ID of the requested ID
					updated.ID = bar.ID
					err = gormRepo.Update(updated)
					So(err, ShouldBeNil)

					So(updated.ID, ShouldEqual, bar.ID)
					So(updated.Name, ShouldNotEqual, bar.Name)
					So(updated.Property, ShouldNotEqual, bar.Property)
					So(updated.Property, ShouldEqual, 0)
				})

				Convey(`By not providing ID - Primary Key in the 'req' object
				 the updated object would have a new ID.`, func() {
					var bar, updated *Bar = bars[1], &Bar{Name: "SecondNoID"}
					err = gormRepo.Update(updated)
					So(err, ShouldBeNil)

					So(updated.ID, ShouldNotEqual, bar.ID)
					So(updated.Name, ShouldNotEqual, bar.Name)
				})

				Convey(`By providing ID - Primary Key in the 'req' object
					and none other paramaters the query would result in an empty record.`, func() {
					var bar, updated *Bar = bars[2], &Bar{}
					updated.ID = bar.ID

					// This would search for records with id as bar.ID and property as
					//bars[3].property and update their fields so that Name = 'Changed'
					err = gormRepo.Update(updated)
					So(err, ShouldBeNil)

					So(updated.Name, ShouldEqual, "")
					So(updated.Property, ShouldEqual, 0)
				})

				Convey(`If provided ID is not in the table new object is 
					created with provided ID`, func() {
					err = gormRepo.Update(&Bar{ID: 12345, Name: "Name for new"})
					So(err, ShouldBeNil)

					var bar *Bar = &Bar{ID: 12345}
					err = db.Find(bar).Error
					So(err, ShouldBeNil)

					So(bar.Name, ShouldEqual, "Name for new")

				})
			})

			Convey("Should return *Error if an error occurs", func() {
				err = gormRepo.Update(&NotInDB{})
				So(err, ShouldBeError)

				So(err, ShouldHaveSameTypeAs, &dberrors.Error{})

			})

		})

	})

}

func TestGORMRepositoryPatch(t *testing.T) {

	Convey("Subject: Updates the database record field from 'what' that matches given on 'where' restriction.", t, func() {

		db, err := openGormSqlite()
		So(err, ShouldBeNil)

		defer db.Close()
		defer clearDB(db)

		var bars []*Bar = seedBars(db)

		Convey(`Using a Patch method on GORMRepository with provided 'what' fields 
			to change on 'where' restrictions.`, func() {
			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			Convey("Should update all 'what' records fields that  match given 'where'", func() {
				var what *Bar = &Bar{Name: "ChangedFirst", Property: 1234}
				err = gormRepo.Patch(what, &Bar{ID: bars[0].ID})
				So(err, ShouldBeNil)

				var getted *Bar = &Bar{ID: bars[0].ID}
				db.Find(getted)

				So(getted.Name, ShouldEqual, what.Name)
				So(getted.Property, ShouldEqual, what.Property)
			})

			Convey("Should not update all record for given 'where'", func() {
				var req *Bar = &Bar{Name: "Changed Second"}
				err = gormRepo.Patch(req, &Bar{ID: bars[1].ID})
				So(err, ShouldBeNil)

				var getted *Bar = &Bar{ID: bars[1].ID}
				db.Find(getted)

				Convey("Fields not included in 'req' should not change", func() {
					So(getted.ID, ShouldEqual, bars[1].ID)
					So(getted.Property, ShouldEqual, bars[1].Property)
				})

				Convey("Fields included in 'req' should change", func() {
					So(getted.Name, ShouldNotEqual, bars[1].Name)
				})

			})

			Convey("Should return *Error if an error occurs", func() {
				var foobar *Foobar = &Foobar{Name: "Name"}

				db.Create(foobar)
				var name string = "Long Name longer than possible"
				db.Create(&Foobar{Name: name})
				err = gormRepo.Patch(&Foobar{Name: name},
					Foobar{ID: foobar.ID},
				)
				So(err, ShouldBeError)
				So(err, ShouldHaveSameTypeAs, &dberrors.Error{})

				var bar *Bar = &Bar{Name: "Non existend id name"}
				err = gormRepo.Patch(bar, Bar{ID: 12345})
				So(err, ShouldBeError)

			})

		})

	})

}

func TestGORMRepositoryDelete(t *testing.T) {

	Convey("Subject: Deleting the database records that matches the 'req' object.", t, func() {

		db, err := openGormSqlite()
		So(err, ShouldBeNil)

		defer db.Close()
		defer clearDB(db)

		var bars []*Bar = seedBars(db)
		Convey("Using a Delete method on GORMRepository for given 'req' restrictions", func() {
			gormRepo, err := New(db)
			So(err, ShouldBeNil)

			Convey("Should delete all files that matches the 'req' restrictions", func() {
				var where []int = []int{1, 2, 3}

				err = gormRepo.Delete(&Bar{}, where)
				So(err, ShouldBeNil)

				var gettedBars []*Bar

				db.Find(&gettedBars)

				So(len(gettedBars), ShouldEqual, 3)
				So(len(gettedBars), ShouldNotEqual, len(bars))
				So(gettedBars[0], ShouldNotEqual, bars[0])

				err = gormRepo.Delete(&Bar{}, Bar{ID: 5})
				So(err, ShouldBeNil)

				var bar *Bar = &Bar{ID: 5}
				err = db.Find(bar).Error

				So(err, ShouldBeError)
				So(bar, ShouldNotResemble, bars[4])
			})

			Convey("Should return *Error if an error occurs", func() {
				err = gormRepo.Delete(&Foobar{}, Foobar{ID: 1234})
			})

		})

	})

}

func openGormSqlite() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./tests.db")
	if err != nil {
		return nil, err
	}
	err = migrateModels(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateModels(db *gorm.DB) error {
	if db != nil {
		db.AutoMigrate(&Bar{}, &Foo{}, &Foobar{})
		return nil
	}
	return errors.New("Nil pointer provided")
}

func clearDB(db *gorm.DB) {
	db.DropTableIfExists(&Bar{}, &Foo{}, &Foobar{})
}

func seedBars(db *gorm.DB) (bars []*Bar) {
	bars = []*Bar{
		{Name: "First", Property: 4141}, {Name: "Second", Property: 1213},
		{Name: "Third", Property: 21410}, {Name: "Fourth", Property: 111},
		{Name: "Fifth", Property: 15102}, {Name: "Sixth", Property: 12410}}

	for _, bar := range bars {
		db.Create(bar)
	}
	return bars
}

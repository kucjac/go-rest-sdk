package forms

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Foo struct {
	Bar string `json:"bar"`
}

type FooNess struct {
	BarNess uint `form:"custombarness"`
}

type IDSetterModel struct {
	ID uint64
}

func (m *IDSetterModel) SetID(id uint64) {
	m.ID = id
}

func TestMapForm(t *testing.T) {
	type mapTest struct {
		IntField          int
		Int8Field         int8
		Int16Field        int16
		Int32Field        int32
		Int64Field        int64
		UintField         uint
		Uint8Field        uint8 `testing:"uint8field"`
		Uint16Field       uint16
		Uint32Field       uint32
		Uint64Field       uint64
		Float32Field      float32
		Float64Field      float64
		StringField       string
		BoolField         bool
		FooField          Foo
		Unreadable        int `form:"-"`
		privateUnreadable int
		NotInQuery        int
		Timefield         time.Time `form:"timefield" time_format:"2006-01-02" time_location:"Poland"`
		IntSlice          []int     `testing:"intslicetest"`
		NoValueTime       time.Time `form:"novaluetime" time_format:"2006-01-02" time_utc:"1"`
		BlankTimeFormat   time.Time `form:"blanktimeformat" time_format:""`
		EmptyBool         bool
		Fooness           FooNess
	}

	Convey("Having a struct containing all possible basic fields", t, func() {

		mapTestObj := &mapTest{}

		Convey("And a form map containing key matching maptest field names", func() {
			correct := map[string][]string{
				"intfield":     {"-1"},
				"int8field":    {"127"},
				"int16field":   {"12301"},
				"int32field":   {"621021"},
				"int64field":   {"4300000000"},
				"uintfield":    {"3"},
				"uint8field":   {"255"},
				"uint16field":  {"65535"},
				"uint32field":  {"4294967295"},
				"uint64field":  {"1840000000000000000"},
				"float32field": {"3.2032"},
				"float64field": {"21431.21501021"},
				"boolfield":    {"true"},
				"stringfield":  {"samsing"},
				"bar":          {"asada"},
				"intslice":     {"1", "2"},
				"timefield":    {"1994-11-05"},
				"novaluetime":  {""},
				"emptybool":    {""},
			}
			Convey("Using policy that throws errors", func() {
				policy := &Policy{FailOnError: true, Tag: "form"}

				Convey("mapForm function match all fields within that struct", func() {
					err := mapForm(mapTestObj, correct, policy)

					So(err, ShouldBeNil)

					So(mapTestObj.IntField, ShouldEqual, -1)
					So(mapTestObj.Int8Field, ShouldEqual, 127)
					So(mapTestObj.Int16Field, ShouldEqual, 12301)
					So(mapTestObj.Int32Field, ShouldEqual, 621021)
					So(mapTestObj.Int64Field, ShouldEqual, 4300000000)
					So(mapTestObj.UintField, ShouldEqual, 3)
					So(mapTestObj.Uint8Field, ShouldEqual, 255)
					So(mapTestObj.Uint16Field, ShouldEqual, 65535)
					So(mapTestObj.Uint32Field, ShouldEqual, 4294967295)
					So(mapTestObj.Uint64Field, ShouldEqual, 1840000000000000000)
					So(mapTestObj.Float32Field, ShouldEqual, 3.2032)
					So(mapTestObj.Float64Field, ShouldEqual, 21431.21501021)
					So(mapTestObj.BoolField, ShouldEqual, true)
					So(mapTestObj.StringField, ShouldEqual, "samsing")
					So(mapTestObj.FooField.Bar, ShouldEqual, "asada")
					Print(mapTestObj.Timefield)

					So(mapTestObj.Timefield.Year(), ShouldEqual, 1994)
				})
			})
		})
		Convey("Using policy that does throw errors", func() {
			policy := &Policy{FailOnError: true, Tag: "testing", TaggedOnly: true}

			Convey("Having a map containg one correct and one incorrect value", func() {
				mapform := map[string][]string{
					"intslicetest": {"1", "incorrect type"},
				}

				Convey("Should not add any value if one of form arguments are incorrect", func() {
					mapTest2Obj := &mapTest{}
					err := mapForm(mapTest2Obj, mapform, policy)

					So(err, ShouldBeError)
					So(mapTest2Obj.IntSlice, ShouldNotContain, 1)
					So(mapTest2Obj.IntSlice, ShouldNotContain, "maciek")
				})
			})
			Convey("Having a map containg incorrect field value", func() {
				mapform2 := map[string][]string{
					"uint8field": {"-1"},
				}
				Convey("Should not bind incorrect value", func() {
					mapTestObj3 := &mapTest{}
					err := mapForm(mapTestObj3, mapform2, policy)

					So(err, ShouldBeError)
					Printf("Uint8Value: %v", mapTestObj3.Uint8Field)
					So(mapTestObj3.Uint8Field, ShouldNotEqual, -1)

				})
			})

			Convey("Having an incorrect value for nested struct", func() {
				mapform3 := map[string][]string{
					"barness": {"-1"},
				}
				Convey("Should throw error and not bind any value to it", func() {
					mapTestObj4 := &mapTest{}
					err := mapForm(mapTestObj4, mapform3, &Policy{FailOnError: true})

					So(mapTestObj4.Fooness.BarNess, ShouldBeZeroValue)
					So(err, ShouldBeError)

				})
			})
		})
	})
}

func TestBindQuery(t *testing.T) {
	Convey("Having a request containing query", t, func() {
		req := httptest.NewRequest("POST", "/foo?bar=content", nil)

		Convey("Using a policy with no tags", func() {
			policy := &Policy{TaggedOnly: false}

			Convey("Where the model is of type Foo", func() {
				fooModel := &Foo{}
				err := BindQuery(req, fooModel, policy)

				So(err, ShouldBeNil)
				So(fooModel.Bar, ShouldEqual, "content")
			})
		})

		Convey("Providing no policy, the default would be set", func() {
			var policy *Policy = nil

			Convey("With the model of type Foo", func() {
				fooModel := &Foo{}
				err := BindQuery(req, fooModel, policy)

				So(err, ShouldBeNil)
				So(fooModel.Bar, ShouldEqual, "content")
			})
		})
	})
	Convey("Having a request containing query", t, func() {
		req := httptest.NewRequest("POST", "/test?inttype=content", nil)

		Convey("Using a policy that fails on error", func() {
			policy := &Policy{FailOnError: true}

			Convey("Binding to model of type Test", func() {
				type Test struct {
					Inttype int
				}
				testObj := Test{}

				err := BindQuery(req, &testObj, policy)
				So(err, ShouldBeError)
				So(testObj.Inttype, ShouldBeZeroValue)
			})
		})
	})
}

func TestBindJSON(t *testing.T) {
	Convey("Having a request containing json type Body", t, func() {

		req := httptest.NewRequest("POST", "/jsontype", strings.NewReader(`{"Bar":"barcontent"}`))

		Convey("Using no policy", func() {
			var policy *Policy

			Convey("BindJSON decodes the body into a model of type Foo", func() {
				model := Foo{}
				err := BindJSON(req, &model, policy)
				So(err, ShouldBeNil)
				So(model.Bar, ShouldEqual, "barcontent")
			})
		})

	})

	Convey("Having a request with incorrect json body", t, func() {
		req := httptest.NewRequest("POST", "/jsonincorrect", strings.NewReader(`{"Bar":"nobrackets"`))

		Convey("Using a policy with FailOnError", func() {
			var policy *Policy = &Policy{FailOnError: true}

			Convey("Decoding the json request will return error", func() {
				model := Foo{}
				err := BindJSON(req, &model, policy)

				So(err, ShouldBeError)
				So(model.Bar, ShouldBeZeroValue)
			})
		})
	})
}

func TestSetBoolValue(t *testing.T) {
	Convey("Having a reflect.Value of type bool", t, func() {
		var BoolValue bool = true
		field := reflect.ValueOf(BoolValue)

		Convey("Providing incorrect boolean parse value", func() {
			val := "falsing"

			Convey("The error should occur", func() {
				err := setBoolField(val, field)
				So(err, ShouldBeError)
			})
		})
	})
}

func TestSetID(t *testing.T) {
	Convey("Having a id string value", t, func() {

		var correctID string = "32"
		var incorrectID string = "-32"

		Convey("And a model containing uint ID field", func() {
			type Model struct {
				ID uint
			}
			model := &Model{}

			Convey("Trying to set incorrectID to its ID field ", func() {
				err := SetID(model, incorrectID)

				So(err, ShouldBeError)
				So(model.ID, ShouldBeZeroValue)
			})

			Convey("Trying to set correctID to its ID field", func() {
				err := SetID(model, correctID)

				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 32)
			})
		})
		Convey("And a model that does not contain any ID field", func() {
			type Model struct {
				NotIDField string
				// private fields cannot be set - ignored
				id uint
			}

			model := &Model{}

			Convey("Trying to set any ID value", func() {

				err := SetID(model, correctID)

				So(err, ShouldBeError)
				So(err, ShouldEqual, ErrIncorrectModel)
			})
		})
		Convey("A model that implements IDSetter interface", func() {
			idSetterModel := &IDSetterModel{}

			Convey("ID is being set by idSetter method", func() {
				err := SetID(idSetterModel, correctID)

				So(err, ShouldBeNil)
				So(idSetterModel.ID, ShouldEqual, 32)

				Convey("But still it must be correct value", func() {
					err = SetID(idSetterModel, incorrectID)

					So(err, ShouldBeError)
				})
			})
		})
	})

}

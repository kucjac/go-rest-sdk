package restsdk

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

func TestGetType(t *testing.T) {
	Convey("Having an object of type *Foo", t, func() {
		obj := &Foo{Bar: "bar"}

		Convey("We obtain a type Foo", func() {
			fooType := getType(obj)

			So(fooType, ShouldEqual, reflect.TypeOf(Foo{}))
		})
	})
	Convey("Even if we have a type of ****Foo", t, func() {

		// Get ****Foo
		obj1 := &Foo{Bar: "bar"}
		obj2 := &obj1
		obj3 := &obj2
		obj := &obj3

		Convey("We again obtain a type of Foo", func() {
			fooType := getType(obj)

			So(reflect.New(fooType).Elem().Interface(), ShouldHaveSameTypeAs, Foo{})
		})
	})
}

func TestObjOfType(t *testing.T) {
	Convey("Having an object of type Foo", t, func() {
		obj := Foo{Bar: "bar"}

		Convey("Getting new object of the same type Foo", func() {
			newObj := ObjOfType(obj)
			So(reflect.TypeOf(obj), ShouldEqual, reflect.TypeOf(newObj))

			Convey("But the new object should be different", func() {
				So(newObj, ShouldNotEqual, obj)
			})
		})
	})

	Convey("Having a pointer to object of type Foo", t, func() {
		obj := &Foo{Bar: "bar"}

		Convey("By using ObjOfType() object of type Foo is obtained", func() {
			newObj := ObjOfType(obj)

			So(reflect.TypeOf(newObj), ShouldNotEqual, reflect.TypeOf(obj))
			So(reflect.TypeOf(newObj), ShouldEqual, reflect.TypeOf(Foo{}))
		})
	})
}

func TestObjOfPtrType(t *testing.T) {
	Convey("Having an object of type *Foo", t, func() {
		obj := new(Foo)
		obj.Bar = "bar"

		Convey("Getting pointer to new object of the same as 'obj' ", func() {
			newPtrObj := ObjOfPtrType(obj)

			So(reflect.TypeOf(newPtrObj), ShouldEqual, reflect.TypeOf(obj))

			Convey("But the pointer should not point to the same object", func() {

				So(newPtrObj, ShouldNotEqual, obj)
			})
		})
	})

	Convey("Having an object of type Foo", t, func() {
		obj := Foo{Bar: "bar"}

		Convey("By using ObjOfPtrType with 'object'", func() {
			newPtrObj := ObjOfPtrType(obj)

			Convey("The new object is of type *Foo", func() {
				So(reflect.TypeOf(newPtrObj), ShouldEqual, reflect.TypeOf(&Foo{}))

				Convey("But the new object do not point to the requested object", func() {
					So(newPtrObj, ShouldNotPointTo, &obj)
				})
			})
		})
	})
}

func TestSliceOfPtrType(t *testing.T) {
	Convey("Having requested object of type Foo", t, func() {
		obj := Foo{Bar: "bar"}

		Convey("Using it as an argument of SliceOfPtrType() requested object", func() {

			slice := SliceOfPtrType(obj)

			Convey("The result would be a slice of type *Foo", func() {
				typedSlice, ok := slice.([]*Foo)
				So(ok, ShouldBeTrue)

				Convey("The slice should be empty", func() {
					So(len(typedSlice), ShouldEqual, 0)
				})
			})
		})
	})
}

func TestStructName(t *testing.T) {
	Convey("Having an object of type Foo", t, func() {
		obj := Foo{Bar: "bar"}

		Convey("We get the struct name by using StructName", func() {
			sname := StructName(obj)

			So(sname, ShouldEqual, "Foo")
		})
	})
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
				policy := &FormPolicy{FailOnError: true, Tag: "form"}

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
		Convey("Using policy that does not throw errors", func() {
			policy := &FormPolicy{FailOnError: true, Tag: "testing", TaggedOnly: true}

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
		})
	})
}

func TestBindQuery(t *testing.T) {
	Convey("Having a request containing query", t, func() {
		req := httptest.NewRequest("POST", "/foo?bar=content", nil)

		Convey("Using a policy with no tags", func() {
			policy := &FormPolicy{TaggedOnly: false}

			Convey("Where the model is of type Foo", func() {
				fooModel := &Foo{}
				err := BindQuery(req, fooModel, policy)

				So(err, ShouldBeNil)
				So(fooModel.Bar, ShouldEqual, "content")
			})
		})

		Convey("Providing no policy, the default would be set", func() {
			var policy *FormPolicy = nil

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
			policy := &FormPolicy{FailOnError: true}

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
			var policy *FormPolicy

			Convey("BindJSON decodes the body into a model of type Foo", func() {
				model := Foo{}
				err := BindJSON(req, &model, policy)
				So(err, ShouldBeNil)
				So(model.Bar, ShouldEqual, "barcontent")
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

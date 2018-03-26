package refutils

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

type Foo struct {
	Bar string `json:"bar"`
}

type FooNess struct {
	BarNess uint `form:"custombarness"`
}

func TestObjsOfTheSameType(t *testing.T) {
	Convey("Subject: ObjsOfTheSameType compares the type of two structs", t, func() {
		Convey("Having two records of the same type", func() {
			a := Foo{}
			b := Foo{}
			Convey("The function should return true", func() {
				value := ObjsOfTheSameType(a, b)
				So(value, ShouldBeTrue)
			})
		})
		Convey("Having one object of *Foo type and second of Foo type", func() {
			a := &Foo{}
			b := Foo{}

			Convey("The function should return true", func() {
				value := ObjsOfTheSameType(a, b)
				So(value, ShouldBeTrue)
			})
		})

		Convey("Having objects of different types", func() {
			a := Foo{}
			b := FooNess{}

			Convey("The function should return false", func() {
				value := ObjsOfTheSameType(a, b)
				So(value, ShouldBeFalse)
			})
		})

		Convey("Having objects of different ptr types", func() {
			a := &Foo{}
			b := &FooNess{}

			Convey("The function should return false", func() {
				value := ObjsOfTheSameType(a, b)
				So(value, ShouldBeFalse)
			})
		})
	})
}

func TestGetType(t *testing.T) {
	Convey("Having an object of type *Foo", t, func() {
		obj := &Foo{Bar: "bar"}

		Convey("We obtain a type Foo", func() {
			fooType := GetType(obj)

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
			fooType := GetType(obj)

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

func TestPtrSliceOfPtrType(t *testing.T) {
	Convey("Subject: Create new *[]*Model for given interface{} 'Model'", t, func() {
		Convey("Having some record of type Foo", func() {
			foo := Foo{Bar: "bar"}

			Convey(`The PtrSliceOfPtrType on the 'Foo' entity,
			 should return non nil *[]*Foo`, func() {
				returned := PtrSliceOfPtrType(foo)

				So(returned, ShouldHaveSameTypeAs, &[]*Foo{})
				typed := returned.(*[]*Foo)
				So(typed, ShouldNotBeNil)
			})
		})
	})
}

func TestSliceOfType(t *testing.T) {
	Convey("Subject: Create new []Model for given interface 'model'", t, func() {
		Convey("Having some record of type Foo", func() {
			foo := Foo{Bar: "bar"}

			Convey(`The SliceOfType on given record should return []Foo slice.`, func() {
				returned := SliceOfType(foo)

				_, ok := returned.([]Foo)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestModelName(t *testing.T) {
	Convey(`Subject: Retrieve pluralized (for slices), lowercased model name
	 for any provided model`, t, func() {
		Convey("Having some record of type Foo", func() {
			fooRecord := Foo{Bar: "bar"}

			Convey("ModelName function should return lowercased 'foo'", func() {
				name := ModelName(fooRecord)

				So(name, ShouldEqual, "foo")
			})
		})
		Convey("Having some slice of type Foo", func() {
			fooSlice := []Foo{}

			Convey("ModelName should return pluralized and lowercased 'foos'", func() {
				name := ModelName(fooSlice)
				So(name, ShouldEqual, "foos")
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

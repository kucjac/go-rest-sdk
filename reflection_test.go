package restsdk

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

type Foo struct {
	Bar string
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

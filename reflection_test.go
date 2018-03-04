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

			So(reflect.New(fooType).Elem().Interface(), ShouldHaveSameTypeAs, Foo{})
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

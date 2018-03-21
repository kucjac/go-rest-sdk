package refutils

import (
	"reflect"
)

func ObjsOfTheSameType(first interface{}, second interface{}) bool {
	firstType := getType(first)
	secondType := getType(second)
	if firstType == secondType {
		return true
	}
	return false
}

// ObjOfType returns the empty object of type given in request
// For example if req is type Foo it returns empty object of type Foo
// The function deeply checks type of the request so that even if it is
// type *****Foo the result would be of type Foo
func ObjOfType(req interface{}) (res interface{}) {
	t := getType(req)
	res = reflect.New(t).Elem().Interface()
	return res
}

// ObjOfPtrType returns the object of pointer type given in request
// For example if req is type Foo it returns new object of type *Foo
// The function deeply checks type of the request so that even if it is
// type *****Foo the result would be of type Foo
func ObjOfPtrType(req interface{}) (res interface{}) {
	t := getType(req)
	res = reflect.New(t).Interface()
	return res
}

// SliceOfPtrType returns empty slice of a pointers of a type provided in request
// In example providing type Foo in request, the function returns res of type []*Foo
func SliceOfPtrType(req interface{}) (res interface{}) {
	t := getType(req)
	res = reflect.New(reflect.SliceOf(reflect.PtrTo(t))).Elem().Interface()
	return res
}

// PtrSliceOfPtrType returns a pointer to an empty slice of pointers of a type provided in request
// In example providing type Foo in request, the funciton returns *[]*Foo
func PtrSliceOfPtrType(req interface{}) (res interface{}) {
	t := getType(req)
	res = reflect.New(reflect.SliceOf(reflect.PtrTo(t))).Interface()
	return res
}

// SliceOfType returns empty slice of a type provided in request
// In example providing type Foo in request, the function returns res of type []Foo
func SliceOfType(req interface{}) (res interface{}) {
	t := getType(req)
	res = reflect.New(reflect.SliceOf(t)).Elem().Interface()
	return res
}

// StructName returns the name of the provided model
func StructName(model interface{}) string {
	t := getType(model)
	return t.Name()
}

// GetType returns the type of the 'req' source object.
// I.e. for *Foo, []Foo, []*Foo object the function returns the Type Foo
func GetType(req interface{}) reflect.Type {
	return getType(req)
}

func getType(req interface{}) reflect.Type {
	// Get Type of requested object
	t := reflect.TypeOf(req)

	// If the object is a pointer or a slice dereference it
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t
}

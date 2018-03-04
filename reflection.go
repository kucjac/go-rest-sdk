package restsdk

import (
	"reflect"
)

func getType(req interface{}) reflect.Type {
	// Get Type of requested object
	t := reflect.TypeOf(req)

	// If the object is a pointer and dereference it
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}
	return t
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

func StructName(req interface{}) string {
	t := getType(req)
	return t.Name()
}

// func BindQuery(req interface{}, query map[string][]string) error {
// 	if req == nil {
// 		return ErrNilPointerProvided
// 	}

// 	t := reflect.TypeOf(req)
// 	if t.Kind() != reflect.Ptr {
// 		return ErrPtrNotProvided
// 	}

// 	for i := 0; i < t.NumField(); i++ {
// 		tag := t.Field(i).Tag.Get("urlquery")
// 		if tag == "-" {
// 			continue
// 		} else if tag == "" {
// 			tag = strings.ToLower(t.Field(i).Name)
// 		}

// 		queryVal := query.Get(key)
// 		if queryVal != "" {

// 		}

// 	}

// }

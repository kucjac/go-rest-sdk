package restsdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FormPolicy struct {
	TaggedOnly  bool
	FailOnError bool
	Tag         string
}

var (
	DefaultFormPolicy = FormPolicy{
		TaggedOnly:  false,
		FailOnError: false,
		Tag:         "form",
	}
)

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

// StructName returns the name of the provided model
func StructName(model interface{}) string {
	t := getType(model)
	return t.Name()
}

// BindQuery binds the url Query
// for the given request to the provided model
func BindQuery(req *http.Request, model interface{}, policy *FormPolicy) error {
	if policy == nil {
		policy = &DefaultFormPolicy
	}
	values := req.URL.Query()
	err := mapForm(model, values, policy)
	if err != nil {
		return err
	}
	return nil
}

// BindJSON binds the request body and decode it into provided model
func BindJSON(req *http.Request, model interface{}, policy *FormPolicy) error {
	if policy == nil {
		policy = &DefaultFormPolicy
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(model)
	if policy.FailOnError && err != nil {
		return err
	}
	return nil
}

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

func mapForm(ptr interface{}, form map[string][]string, policy *FormPolicy) error {
	// Get type of pointer
	t := reflect.TypeOf(ptr).Elem()

	// Get value of pointer
	v := reflect.ValueOf(ptr).Elem()

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		sField := v.Field(i)

		// Check if field is settable - can be addresable or is public
		if !sField.CanSet() {
			continue
		}

		sFieldKind := sField.Kind()

		// Check if the field has a tag query
		formTag := tField.Tag.Get(policy.Tag)

		// If tag is set to '-' don't map values
		if formTag == "-" {
			continue
		}

		if formTag == "" {
			if sFieldKind == reflect.Struct {
				// mapQuery recursively if the field is a struct
				err := mapForm(sField.Addr().Interface(), form, policy)
				// check error only if the policy requirers it
				if policy.FailOnError {
					if err != nil {
						return err
					}
				}
				continue
			} else if !policy.TaggedOnly {
				// get the lowercased name of a field
				formTag = strings.ToLower(tField.Name)
			} else {
				// if the policy is tagged only continue
				continue
			}

		}

		// Check if the query contains the tag
		formValue, ok := form[formTag]
		if !ok {
			continue
		}

		elemNum := len(formValue)

		// The query value can conatin more than one value
		// If the field is a slice and the queryValue
		// has any values assign it if possible
		if sFieldKind == reflect.Slice && elemNum > 0 {
			sliceKind := sField.Type().Elem().Kind()
			fieldSlice := reflect.MakeSlice(sField.Type(), elemNum, elemNum)
			// Iterate over query elements and add to fieldSlice
			for j := 0; j < elemNum; j++ {
				// set with given value
				err := setFieldWithType(sliceKind, formValue[j], fieldSlice.Index(i))
				if policy.FailOnError && err != nil {
					return err
				}
			}
			// Set 'ptr' value for field of index 'i' with 'fieldSlice'
			v.Field(i).Set(fieldSlice)
		} else {
			// check if the query is of type time
			if _, isTime := sField.Interface().(time.Time); isTime {
				err := setTimeField(formValue[0], tField, sField)
				if policy.FailOnError && err != nil {
					return err
				}
			} else {
				err := setFieldWithType(tField.Type.Kind(), formValue[0], sField)
				if policy.FailOnError && err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// setFieldWithType sets given 'field' of 'fieldKind' with value 'val'.
// When the value is not of given Kind, the function throws error
func setFieldWithType(
	fieldKind reflect.Kind,
	val string,
	field reflect.Value,
) (err error) {
	switch fieldKind {
	case reflect.String:
		field.SetString(val)
	case reflect.Int:
		setIntField(val, field, 0)
	case reflect.Int8:
		setIntField(val, field, 8)
	case reflect.Int16:
		setIntField(val, field, 16)
	case reflect.Int32:
		setIntField(val, field, 32)
	case reflect.Int64:
		setIntField(val, field, 64)

	case reflect.Uint:
		setUintField(val, field, 0)
	case reflect.Uint8:
		setUintField(val, field, 8)
	case reflect.Uint16:
		setUintField(val, field, 16)
	case reflect.Uint32:
		setUintField(val, field, 32)
	case reflect.Uint64:
		setUintField(val, field, 64)

	case reflect.Float32:
		setFloatField(val, field, 32)
	case reflect.Float64:
		setFloatField(val, field, 64)

	case reflect.Bool:
		setBoolField(val, field)

	default:
		return ErrUnknownType
	}
	return nil
}

func setIntField(val string, field reflect.Value, bitSize int) (err error) {
	var intValue int64
	intValue, err = strconv.ParseInt(val, 10, bitSize)
	if err != nil {
		return err
	}
	// Set value if no error
	field.SetInt(intValue)
	return nil
}

func setUintField(val string, field reflect.Value, bitSize int) (err error) {
	var uintValue uint64

	// Parse unsigned int
	uintValue, err = strconv.ParseUint(val, 10, bitSize)

	if err != nil {
		return err
	}

	// Set uint
	field.SetUint(uintValue)
	return nil
}

func setFloatField(val string, field reflect.Value, bitSize int) (err error) {
	var floatValue float64

	// Parse float
	floatValue, err = strconv.ParseFloat(val, bitSize)
	if err != nil {
		return err
	}
	field.SetFloat(floatValue)
	return nil
}

func setBoolField(val string, field reflect.Value) (err error) {
	var boolValue bool
	// set default if empty
	if val == "" {
		val = "false"
	}
	boolValue, err = strconv.ParseBool(val)
	if err != nil {
		return err
	}
	field.SetBool(boolValue)
	return nil
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		return errors.New("Blank time format")
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

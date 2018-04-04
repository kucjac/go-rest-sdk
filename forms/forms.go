package forms

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUnknownType    = errors.New("Unknown data type")
	ErrIncorrectModel = errors.New("Given model do not have ID field. In order to set ID, it should implement IDSetter or contain field ID")
)

// IDSetter defines interface for data Models that allows
// Setting the provided models ID.
// Defines SetID() method.
type IDSetter interface {
	SetID(id uint64)
}

// BindQuery binds the url Query
// for the given request to the provided model
// The function mechanics is based on provided request URL Query
// as well as model and BindPolicy.
// The policy defines if the query binding should search only for
// the fields that contains tags, defines the tags and decides wether the
// function should return an error if any occurs during operation.
// If no policy is provided (or nil) then the function return quickly with nil error.
func BindQuery(req *http.Request, model interface{}, policy *BindPolicy) error {
	if policy == nil {
		return nil
	}
	values := req.URL.Query()
	err := mapForm(model, values, policy, policy.SearchDepthLevel)
	if err != nil {
		return err
	}
	return nil
}

// BindJSON binds the reads the provided request body
// and decode it into provided model.
// If an error occurred during model binding, error returns
func BindJSON(req *http.Request, model interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(model)
	if err != nil {
		return err
	}
	return nil
}

// SetID sets the ID of provided model.
// If model implements IDSetter interface it uses SetID method at first.
// Otherwise checks whether provided model contains 'ID' or 'Id' field
// And parses the 'id' argument
// Returns error if provided argument is not appropiate for given field
// 	or there is no ID field in the model
func SetID(model interface{}, id string) error {
	// Check if given model implements IDSetter interface
	if setter, ok := model.(IDSetter); ok {
		uintID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return err
		}
		setter.SetID(uintID)
		return nil
	}

	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		sField := v.Field(i)

		if !sField.CanSet() {
			continue
		}

		if strings.ToLower(tField.Name) == "id" {
			err := setFieldWithType(sField.Kind(), id, sField)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrIncorrectModel
}

func mapForm(
	model interface{},
	form map[string][]string,
	policy *BindPolicy,
	searchDepthLevel int,
) error {

	// Get type of pointer
	t := reflect.TypeOf(model).Elem()

	// Get value of pointer
	v := reflect.ValueOf(model).Elem()

	// iterate over model field
	for i := 0; i < t.NumField(); i++ {

		tField := t.Field(i)
		sField := v.Field(i)

		// isTime flags if the field is of time.Time type
		var isTime bool

		// Check if field is settable - can be addresable or is public
		// or if the field is of type Interface
		if !sField.CanSet() || sField.Kind() == reflect.Interface {
			continue
		}

		// Check if the field has a tag query
		fieldTag := tField.Tag.Get(policy.Tag)

		// If tag is set to '-' don't map values
		if fieldTag == "-" || (policy.TaggedOnly && fieldTag == "") {
			continue
		}

		// Init object if it is of ptr type.
		if sField.Kind() == reflect.Ptr {
			var initialize bool
			switch tField.Type.Elem().Kind() {
			case reflect.Ptr, reflect.Interface:
				continue
			case reflect.Struct:
				if sField.IsNil() && policy.SearchDepthLevel > 0 {
					initialize = true
				} else if policy.SearchDepthLevel <= 0 {
					continue
				}
			default:
				// if the sField is nil - create new item of type given struct type
				if sField.IsNil() {
					initialize = true
				}
			}
			if initialize {
				sField.Set(reflect.New(tField.Type.Elem()))
				sField = sField.Elem()
			}
		}

		if sField.Kind() == reflect.Struct {
			_, isTime = sField.Interface().(time.Time)
			if !isTime {
				if searchDepthLevel > 0 {
					// mapQuery recursively if the field is a struct
					err := mapForm(sField.Addr().Interface(), form, policy, searchDepthLevel-1)
					// check error only if the policy requirers it
					if err != nil && policy.FailOnError {
						return err
					}
				}
				continue
			}
		}

		if fieldTag == "" {
			fieldTag = strings.ToLower(tField.Name)
		}

		// Check if the query contains the tag
		formValue, ok := form[fieldTag]
		if !ok {
			continue
		}
		elemNum := len(formValue)

		// The query value can conatin more than one value
		// If the field is a slice and the queryValue
		// has any values assign it if possible
		if sField.Kind() == reflect.Slice {
			if elemNum <= 0 {
				continue
			}
			sliceKind := sField.Type().Elem().Kind()
			fieldSlice := reflect.MakeSlice(sField.Type(), elemNum, elemNum)

			// Check if the field is a slice of time.Time
			if sliceKind == reflect.Struct {

				_, isTime = fieldSlice.Index(0).Interface().(time.Time)
				if isTime {
					err := setSliceTimeField(formValue, tField, fieldSlice, policy.FailOnError)
					if err != nil && policy.FailOnError {
						return err
					}
				}
				continue
			}

			// Iterate over query elements and add to fieldSlice
			for i := 0; i < elemNum; i++ {
				// set with given value
				err := setFieldWithType(sliceKind, formValue[i], fieldSlice.Index(i))
				if err != nil && policy.FailOnError {
					return err
				}
			}

			// Set 'model' value for field of index 'i' with 'fieldSlice'
			v.Field(i).Set(fieldSlice)
		} else if isTime {
			err := setTimeField(formValue[0], tField, sField)
			if policy.FailOnError && err != nil {
				return err
			}
		} else {
			// check if the field is of type time
			err := setFieldWithType(tField.Type.Kind(), formValue[0], sField)
			if policy.FailOnError && err != nil {
				return err
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
		err = setIntField(val, field, 0)
	case reflect.Int8:
		err = setIntField(val, field, 8)
	case reflect.Int16:
		err = setIntField(val, field, 16)
	case reflect.Int32:
		err = setIntField(val, field, 32)
	case reflect.Int64:
		err = setIntField(val, field, 64)

	case reflect.Uint:
		err = setUintField(val, field, 0)
	case reflect.Uint8:
		err = setUintField(val, field, 8)
	case reflect.Uint16:
		err = setUintField(val, field, 16)
	case reflect.Uint32:
		err = setUintField(val, field, 32)
	case reflect.Uint64:
		err = setUintField(val, field, 64)

	case reflect.Float32:
		err = setFloatField(val, field, 32)
	case reflect.Float64:
		err = setFloatField(val, field, 64)

	case reflect.Bool:
		err = setBoolField(val, field)

	default:
		return ErrUnknownType
	}
	return err
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
	timeFormat, local, err := prepareTimeField(structField)
	if err != nil {
		return err
	}
	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	t, err := time.ParseInLocation(timeFormat, val, local)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setSliceTimeField(
	values []string,
	structField reflect.StructField,
	value reflect.Value,
	failOnError bool,
) error {
	timeFormat, local, err := prepareTimeField(structField)
	if err != nil {
		return err
	}

	for i := 0; i < len(values); i++ {
		if values[i] == "" {
			value.Index(i).Set(reflect.ValueOf(time.Time{}))
			continue
		}

		t, err := time.ParseInLocation(timeFormat, values[i], local)
		if err != nil && failOnError {
			return err
		}
		value.Index(i).Set(reflect.ValueOf(t))
	}
	return nil

}

func prepareTimeField(
	structField reflect.StructField,
) (timeFormat string, local *time.Location, err error) {
	timeFormat = structField.Tag.Get("time_format")
	if timeFormat == "" {
		err = errors.New("Blank time format")
		return
	}

	local = time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		local = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		local, err = time.LoadLocation(locTag)
		if err != nil {
			return
		}
	}
	return
}

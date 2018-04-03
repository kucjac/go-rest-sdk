package forms

import (
	"github.com/kucjac/go-rest-sdk/refutils"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ParamGetterFunc defines the function that retrieve the parameters
// from the specific third-party routing framework on the base
// of the provided parameterName string and req *http.Request
// if individual implementation needs more arguments push them into
// request's context.
type ParamGetterFunc func(paramName string, req *http.Request) (string, error)

// BindParams is a function that searches for fields that matches
// given url (routing) parameters and binds them to provided 'model'
// Uses ParamPolicy as a rules set. If no policy provided
// the function sets 'DefaultParamPolicy'.
// The model must be settable (addressable).
// By default BindParams searches only for the 'ID' field, and the default
// param tag is 'param'. Default ID parameter value should be used
// as a lowercased - model name.
// If the model implements IDSetter interface, the performance of this
// function should be better.
func BindParams(
	req *http.Request,
	model interface{},
	getParam ParamGetterFunc,
	policy *ParamPolicy,
) error {
	if policy == nil {
		return nil
	}
	return mapParam(model, getParam, req, policy, policy.SearchDepthLevel, "")
}

//mapParameters from the url
// rules:
//	- if policy idOnly iterate over the model and search for ID field and structures that are
// 		in the range of SearchDepthLevel
//	- By default the ID or 'Id' field is treated as an 'ID'. If the model contain different id
// 		field name, then set the field tag as an 'id'. if ID is not used as an ID - taggit with '-'
//	- if policy is TagOnly - iterate only over tagged fields.
//	- If a struct field has a tag. Then it goes as a parameter name for the model's id. By defult
//		model's id parameter is - model's name.
//	- if FailOnError is true - then any error that will occur during binding would be returned
//		otherwise the mapping would map the fields until the other rules allows it to continue.
// 	- if SearchDepthLevel is greater than 0, then the mapping function allows to check nested
//		the nested struct or ptr to struct fields recursively with decreased searchDepth level
// 		decreased by one.
//	- if the field is a slice, array or interface then it won't be taken into mapping.
func mapParam(
	model interface{},
	getParam ParamGetterFunc,
	req *http.Request,
	policy *ParamPolicy,
	searchDepthLevel int,
	paramName string,
) error {
	// idAlreadySet is control flag that defines if id was set
	var idAlreadySet bool

	// if paramName is an empty string - set it by default as Model struct Name.
	if paramName == "" {
		paramName = refutils.ModelName(model)
	}
	// Get the parameter from the getParam function
	modelID, err := getParam(paramName, req)
	if err != nil {
		return err
	}

	// if there is no
	if policy.IDOnly && modelID == "" {
		return ErrIncorrectModel
	}

	// If given model implements IDSetter, set ID  using SetID method
	// Returns if an error occurs or policy.IDOnly is true.
	if setter, ok := model.(IDSetter); ok {
		uintID, err := strconv.ParseUint(modelID, 10, 64)
		if err != nil {
			return err
		}
		setter.SetID(uintID)

		if policy.IDOnly && searchDepthLevel == 0 {
			return nil
		}
		idAlreadySet = true
	}

	// Get reflect.Type of model
	t := reflect.TypeOf(model).Elem()

	// Get reflect.Value of model
	v := reflect.ValueOf(model).Elem()

	// iterate  over model fields
	for i := 0; i < t.NumField(); i++ {

		tField := t.Field(i)
		sField := v.Field(i)

		// Check if field is settable - can be addresable or is public
		if !sField.CanSet() {
			continue
		}

		switch sField.Kind() {
		case reflect.Slice, reflect.Interface, reflect.Array:
			continue
		default:
		}

		// Check if the field has a tag query
		fieldTag := tField.Tag.Get(policy.Tag)

		// If tag is set to '-' don't map values
		if fieldTag == "-" || (fieldTag == "" && policy.TaggedOnly) {
			continue
		}

		// if sField is a Ptr check where it points to.
		if sField.Kind() == reflect.Ptr {
			switch tField.Type.Elem().Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Interface:
				continue
			default:
				// if the sField is nil - create new item of type given struct type
				if sField.IsNil() && policy.SearchDepthLevel > 0 {
					sField.Set(reflect.New(tField.Type.Elem()))
				} else if policy.SearchDepthLevel <= 0 {
					continue
				}
				sField = sField.Elem()
			}
		}

		// if the field is of Struct Type
		if sField.Kind() == reflect.Struct {
			// distinguish the
			// if it is time field set it as a time
			_, isTime := sField.Interface().(time.Time)
			if isTime {
				err = paramSetTime(sField, tField, getParam, req, fieldTag)
				if err != nil && policy.FailOnError {
					return err
				}
				continue
			}

			// search nested structs if and only if search depth level is greater than zero.
			if searchDepthLevel > 0 {
				// recursively seach provided struct with decreased search level
				err := mapParam(sField.Addr().Interface(),
					getParam, req, policy, searchDepthLevel-1, fieldTag)
				// return error if occurs and policy allows it
				if err != nil && policy.FailOnError {
					return err
				}
			}
			continue
		}

		var lowerCasedName string = strings.ToLower(tField.Name)

		// if given field is an id, which was not yet set
		if !idAlreadySet && modelID != "" && (lowerCasedName == "id" || fieldTag == "id") {

			err := setFieldWithType(sField.Kind(), modelID, sField)
			if err != nil {
				return err
			}

			// Stop if DeepSearch parameter is false
			if policy.IDOnly && policy.SearchDepthLevel == 0 {
				return nil
			}

			// if DeepSearch is true set idAlreadySet flag to true
			idAlreadySet = true
		} else if fieldTag == "" {
			fieldTag = lowerCasedName
		}

		// if IDOnly rule is true - go to the next field
		if policy.IDOnly {
			continue
		}

		// check the value of given paramtag in the ParamGetterFunc provided as an argument
		paramValue, err := getParam(fieldTag, req)
		// if an error occured check what is the param policy
		if err != nil && policy.FailOnError {
			return err

			// if policy allows errors continue to the next field
			continue
		}

		// if no value present continue with iteration
		if paramValue == "" {
			continue
		}

		// try to set the field with provided 'paramValue'
		err = setFieldWithType(sField.Kind(), paramValue, sField)

		// return an error if occurs and
		// if the policy doesn't allow to continue
		if err != nil && policy.FailOnError {
			return err

		} else if lowerCasedName == "id" {
			// if the correctly setted field was an id
			// set the 'idAlreadySet' flag to true
			idAlreadySet = true
		}
	}

	if !idAlreadySet {
		return ErrIncorrectModel
	}
	return nil
}

func paramSetTime(sField reflect.Value,
	tField reflect.StructField,
	getParam ParamGetterFunc,
	req *http.Request,
	fieldTag string,
) error {
	timeValue, err := getParam(fieldTag, req)
	if err != nil {
		return err
	}

	err = setTimeField(timeValue, tField, sField)
	if err != nil {
		return err
	}
	return nil
}

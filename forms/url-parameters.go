package forms

import (
	"github.com/kucjac/go-rest-sdk/refutils"
	"log"
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
		policy = &DefaultParamPolicy
	}
	return mapParam(model, getParam, req, policy)
}

//mapParameters from the url
func mapParam(
	model interface{},
	getParam ParamGetterFunc,
	req *http.Request,
	policy *ParamPolicy,
) error {

	modelName := refutils.ModelName(model)
	modelID, err := getParam(modelName, req)
	if err != nil {
		return err
	}
	log.Printf("\n\n%v", modelName)
	var idAlreadySet bool

	//If given model implements IDSetter, set ID  using SetID method
	if setter, ok := model.(IDSetter); ok {
		uintID, err := strconv.ParseUint(modelID, 10, 64)
		if err != nil {
			return err
		}
		setter.SetID(uintID)

		if !policy.DeepSearch {
			return nil
		}
		idAlreadySet = true
	}

	// Get type of pointer
	t := reflect.TypeOf(model).Elem()

	// Get value of pointer
	v := reflect.ValueOf(model).Elem()

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
		paramTag := tField.Tag.Get(policy.Tag)

		// If tag is set to '-' don't map values
		if paramTag == "-" {
			continue
		}

		if policy.DeepSearch {
			switch sField.Kind() {
			case reflect.Struct:
				// if it is time field set it as a time
				_, isTime := sField.Interface().(time.Time)
				if isTime {
					if (policy.TaggedOnly && paramTag == "") || paramTag == "" {
						continue
					}
					err = paramSetTime(sField, tField, getParam, req, paramTag)
					if policy.FailOnError && err != nil {
						return err
					}
					continue
				}
				// recursively seach provided struct
				if !policy.TaggedOnly || policy.TaggedOnly && paramTag != "" {
					err := mapParam(sField.Addr().Interface(), getParam, req, policy)
					if policy.FailOnError && err != nil {
						return err
					}

				}

				continue

			case reflect.Ptr:
				// sField.Set(reflect.New(tField.Type.Elem()))
				// the function allows only ptr of singular referenced elements
				switch tField.Type.Elem().Kind() {
				case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Interface:
					continue
				case reflect.Struct:
					sField.Set(reflect.New(tField.Type.Elem()))
					// Check if after dereferencing it is a of a type time.Time
					_, isTime := sField.Elem().Interface().(time.Time)
					if isTime {
						if (policy.TaggedOnly && paramTag == "") || paramTag == "" {
							continue
						}
						err = paramSetTime(sField.Elem(), tField, getParam, req, paramTag)
						if policy.FailOnError && err != nil {
							return err
						}
					} else if !isTime {
						err := mapParam(sField.Interface(), getParam, req, policy)
						if policy.FailOnError && err != nil {
							return err
						}
					}
					continue
				default:
					if policy.TaggedOnly && paramTag == "" {
						continue
					}
					sField.Set(reflect.New(tField.Type.Elem()))
					sField = sField.Elem()
				}

			default:
			}
		} else {
			// if not deep search the function is looking only for the 'ID' field or field with
			// param tagged 'id' that is of basic type
			log.Printf("SfieldKind: %v", sField.Kind())
			switch sField.Kind() {
			case reflect.Struct:
				continue
			case reflect.Ptr:
				switch tField.Type.Elem().Kind() {
				case reflect.Ptr, reflect.Slice,
					reflect.Array, reflect.Interface,
					reflect.Struct:
					continue
				default:
					sField.Set(reflect.New(tField.Type.Elem()))
				}
			default:
			}
		}

		var lowerCasedName string = strings.ToLower(tField.Name)
		log.Println(lowerCasedName)

		if paramTag == "" {
			// if the policy doesn't require tagged only fields
			// set the paramTag value as lowercased field.Name
			log.Println("preparing")
			if !policy.TaggedOnly {
				log.Println("Went through")
				// if given field is an id, which was not yet set
				if !idAlreadySet && lowerCasedName == "id" && modelID != "" {

					err := setFieldWithType(sField.Kind(), modelID, sField)
					if err != nil {
						return err
					}

					// Stop if DeepSearch parameter is false
					if !policy.DeepSearch {
						return nil
					}

					// if DeepSearch is true set idAlreadySet flag to true
					idAlreadySet = true
				} else {
					paramTag = lowerCasedName
				}
			} else {
				// if policy requires tag (TaggedOnly) continue with other fields
				continue
			}
		}

		// check the value of given paramtag in the ParamGetterFunc provided as an argument
		paramValue, err := getParam(paramTag, req)
		if err != nil {
			// if an error occured check what is the param policy
			if policy.FailOnError {
				return err
			}
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
		if err != nil {
			if policy.FailOnError {
				return err
			}
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
	paramTag string,
) error {
	timeValue, err := getParam(paramTag, req)
	if err != nil {
		return err
	}

	err = setTimeField(timeValue, tField, sField)
	if err != nil {
		return err
	}
	return nil
}

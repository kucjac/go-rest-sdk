package forms

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
	// "time"
)

func TestMapParams(t *testing.T) {
	Convey("Subject: Map parameters to the given model", t, func() {

		var err error

		req := httptest.NewRequest("GET", "/url/", nil)
		policy := DefaultParamPolicy.Copy()
		var valueMap map[string]string
		valueMap = make(map[string]string)

		Convey("If an error occurred during getting parameter with getParam func", func() {
			err = mapParam(&Model{}, getParamErrFunc(map[string]error{
				"model": errors.New("error")},
			), req, policy, policy.SearchDepthLevel, "")
			So(err, ShouldBeError)
		})

		Convey(`If policy is set to IDOnly and model's parameter is not found in with getParam function, an error should occur.`, func() {
			policy = &(*policy)
			policy.IDOnly = true
			err = mapParam(&Model{}, emptyGetParamFunc, req,
				policy, policy.SearchDepthLevel, "")
			So(err, ShouldBeError)
		})

		Convey(`Having some model that implements IDSetter interface`, func() {
			var model IDSetterModel
			Convey(`And the param value is of uint type, no error should occur and the ID field
					should be set with provded value`, func() {
				valueMap := map[string]string{"idsettermodel": "15001900"}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy,
					policy.SearchDepthLevel, "")
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 15001900)
			})

			Convey(`And the param is of incorrect value, an error should occur`, func() {
				valueMap := map[string]string{"idsettermodel": "-200"}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy, policy.SearchDepthLevel, "")
				So(err, ShouldBeError)
				So(model.ID, ShouldNotEqual, -200)
			})

			Convey(`if the policy is set tot IDOnly the function and the SearchDepthLevel is 
				is 0, the function return immediately.`, func() {
				policy := DefaultParamPolicy.Copy()
				policy.IDOnly = true
				// by default SearchDepthLevel is set to 0, setting just for testing visibility
				policy.SearchDepthLevel = 0
				valueMap := map[string]string{"idsettermodel": "123"}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy,
					policy.SearchDepthLevel, "")
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 123)
			})
		})

		Convey(`Having a model that contains params and private fields inner structs and 
			ptr type objects.`, func() {
			type ShortModel struct {
				ID    int
				Field string
			}

			type ShortModel2 struct {
				ID int
			}

			type ModelWithoutID struct {
				Field1 string
			}

			type SomeModel struct {
				ID             int
				Name           string
				BarID          int `param:"bar"`
				privateField   string
				NotIncluded    string `param:"-"`
				Short          ShortModel
				ShortWithParam ShortModel   `param:"short"`
				PtrModel       *ShortModel2 `param:"ptr-short"`
				Slice          []int
				SlicePtr       *[]int
				IntPtr         *int
			}

			var model SomeModel
			var err error
			req := httptest.NewRequest("GET", "/", nil)
			policy := DefaultParamPolicy.Copy()

			valueMap = map[string]string{
				"somemodel":  "23",
				"bar":        "44",
				"shortmodel": "55",
				"short":      "66",
				"ptr-short":  "77",
				"intptr":     "88",
			}

			Convey(`If SearchDepthLevel is set to 0, no nested models 
				would get parameters.`, func() {
				policy.SearchDepthLevel = 0
				model = SomeModel{}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "")
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 23)
				So(model.BarID, ShouldEqual, 44)
				So(model.Short.ID, ShouldNotEqual, 55)
				So(model.ShortWithParam.ID, ShouldNotEqual, 66)
				So(model.PtrModel, ShouldBeNil)
				So(*model.IntPtr, ShouldEqual, 88)
			})

			Convey(`If SearchDepthLevel is greater than 1, then nested models
				 would map parameters.`, func() {
				policy.SearchDepthLevel = 1
				model = SomeModel{}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "")
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 23)
				So(model.BarID, ShouldEqual, 44)
				So(model.Short.ID, ShouldEqual, 55)
				So(model.ShortWithParam.ID, ShouldEqual, 66)
				So(model.PtrModel, ShouldNotBeNil)
				So(model.PtrModel.ID, ShouldEqual, 77)

			})
		})
		Convey("If policy FailOnError is true", func() {
			policy.FailOnError = true
			Convey("Having struct that has incorrectly set time field (i.e. no time tags)", func() {
				type modelWithIncorrectTime struct {
					ID   int
					Time time.Time `param:"mytime"`
				}

				valueMap = map[string]string{
					"model":  "123",
					"mytime": "1213",
				}

				model := modelWithIncorrectTime{}
				policy.SearchDepthLevel = 1
				err = mapParam(&model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "model")
				So(err, ShouldBeError)
			})
			Convey("Having struct with nested structs that throws errors", func() {
				type modelWithNested struct {
					ID     int
					Nested struct{ Name string }
				}
				valueMap = map[string]string{
					"model":  "123",
					"nested": "12",
				}

				model := modelWithNested{}
				policy.SearchDepthLevel = 1
				err = mapParam(&model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "model")
				So(err, ShouldBeError)
			})
		})
		Convey(`If model does not implement IDSetter`, func() {
			type SomeModel struct {
				Name    string
				SomeInt uint `param:"errory"`
				ID      int
			}
			model := new(SomeModel)
			policy := DefaultParamPolicy.Copy()

			Convey("If provided incorrect parameter value for the model", func() {
				policy.IDOnly = true
				valueMap = map[string]string{
					"model": "incorrect",
				}
				err := mapParam(model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "model")
				So(err, ShouldBeError)
			})
			Convey("If policy is IDOnly with correct model value", func() {
				policy.IDOnly = true
				valueMap = map[string]string{
					"model": "123",
				}
				err := mapParam(model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "model")
				So(err, ShouldBeNil)
			})
			Convey("If provided parameter throws an error with getParam func", func() {
				policy.FailOnError = true
				errorMap := map[string]error{
					"name": errors.New("Some error"),
				}
				err := mapParam(model, getParamErrFunc(errorMap), req, policy,
					policy.SearchDepthLevel, "model")
				So(err, ShouldBeError)
			})
			Convey(`If provided non id parameter throws 
				an error with setting by field func`, func() {
				policy.FailOnError = true
				valueMap = map[string]string{
					"errory": "-1",
				}
				err := mapParam(model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "")
				So(err, ShouldBeError)
			})
			Convey(`If setting parameter is named as 'id'`, func() {
				valueMap = map[string]string{
					"id": "1234",
				}
				err := mapParam(model, getParamFuncWithValues(valueMap), req,
					policy, policy.SearchDepthLevel, "")
				So(err, ShouldBeNil)
			})

		})

	})
}

func TestBindParams(t *testing.T) {
	Convey("Subject: bind url parameters using BindParam function", t, func() {
		Convey(`Having some model, url parameter function and http request`, func() {
			req := httptest.NewRequest("GET", "/", nil)
			policy := DefaultParamPolicy.Copy()
			Convey(`If used a nonaddressable (not a pointer to) model,
					the function should Panic`, func() {
				model := ModelWithID{}
				So(func() { BindParams(req, model, emptyGetParamFunc, policy) }, ShouldPanic)
			})

			Convey(`If no params provided, or the ParamGetterFunc provided
					doesn't contain models id, an error should be thrown`, func() {
				model := &ModelWithID{}
				err := BindParams(req, model, emptyGetParamFunc, policy)
				So(err, ShouldBeError)
			})
			Convey(`If no policy provided, by default
					the function returns nil error`, func() {
				var policy *ParamPolicy = nil
				model := &ModelWithID{}
				err := BindParams(req, model, emptyGetParamFunc, policy)
				So(err, ShouldBeNil)
			})
			Convey(`If no paramGetterFunc provided the function would return error`, func() {
				policy := DefaultParamPolicy.Copy()
				model := &ModelWithID{}
				var getParam ParamGetterFunc
				err := BindParams(req, model, getParam, policy)
				So(err, ShouldBeError)
			})

		})
	})
}

func TestParamSetTime(t *testing.T) {
	sField := reflect.ValueOf(1)
	tField := reflect.StructField{}
	errorMap := map[string]error{"test": errors.New("Error")}
	err := paramSetTime(sField, tField, getParamErrFunc(errorMap), nil, "test")
	if err == nil {
		t.Error(err)
	}
}

func getParamFuncWithValues(
	paramValues map[string]string,
) ParamGetterFunc {
	return func(paramName string, req *http.Request) (string, error) {
		value := paramValues[paramName]
		return value, nil
	}
}

func emptyGetParamFunc(param string, req *http.Request) (string, error) {
	return "", nil
}

func getParamErrFunc(
	paramErr map[string]error,
) ParamGetterFunc {
	return func(paramName string, req *http.Request) (string, error) {
		err := paramErr[paramName]
		return "", err
	}
}

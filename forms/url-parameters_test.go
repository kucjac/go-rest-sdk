package forms

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestMapParams(t *testing.T) {
	Convey("Subject: Map parameters to the given model", t, func() {
		Convey("Having some model that implements IDSetter interface", func() {
			model := IDSetterModel{}
			req := httptest.NewRequest("GET", "/url", nil)
			policy := DefaultParamPolicy.New()

			err := mapParam(&model, emptyGetParamFunc, req, policy)
			So(err, ShouldBeError)

			valueMap := map[string]string{"idsettermodel": "15001900"}

			err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)

			So(err, ShouldBeNil)
			So(model.ID, ShouldEqual, 15001900)

			model = IDSetterModel{}

			policy.DeepSearch = true
			err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)

			So(err, ShouldBeNil)
			So(model.ID, ShouldEqual, 15001900)
		})

		Convey("Having a model that contains params and private fields", func() {
			model := ModelWithID{}
			req := httptest.NewRequest("GET", "/", nil)
			policy := DefaultParamPolicy.New()
			policy.TaggedOnly = true

			valueMap := map[string]string{"modelwithid": "1501", "bar": "1234"}
			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeError)

			policy.TaggedOnly = false
			policy.DeepSearch = false

			Convey("For ID value not set with different parameter", func() {
				valueMap = map[string]string{"id": "1230"}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeNil)
			})

			Convey("For ID Value setteble with 'modelwithid' param", func() {
				valueMap = map[string]string{"modelwithid": "1123"}
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 1123)

				model.ID = 0

				policy.DeepSearch = true
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeNil)
				So(model.ID, ShouldEqual, 1123)
			})
		})
		Convey("For ID value with different param than model name", func() {
			req := httptest.NewRequest("GET", "/", nil)
			model := ModelWithParam{}
			valueMap := map[string]string{"fieldorf": "1234"}
			policy := DefaultParamPolicy.New()
			policy.DeepSearch = true
			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeNil)
			So(model.ID, ShouldEqual, 1234)

			Convey("If a param contain incorrect value for given type ", func() {
				valueMap = map[string]string{"fieldorf": "maciek"}
				policy.FailOnError = true

				err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})
		})

		Convey(`Having some model and some errors occurs during 
			getting param with ParamGetterFunc`, func() {
			req := httptest.NewRequest("GET", "/", nil)
			model := ModelWithID{}
			errorMap := map[string]error{
				"modelwithid": errors.New("Some error"),
			}

			err := mapParam(&model, getParamErrFunc(errorMap), req, DefaultParamPolicy.New())
			So(err, ShouldBeError)

			errorMap = map[string]error{
				"fieldorf": errors.New("Some error."),
			}
			paramModel := ModelWithParam{}
			policy := DefaultParamPolicy.New()
			policy.DeepSearch = true
			err = mapParam(&paramModel, getParamErrFunc(errorMap), req, policy)
			So(err, ShouldBeError)

			policy.FailOnError = false
			err = mapParam(&paramModel, getParamErrFunc(errorMap), req, policy)
			So(err, ShouldBeError)
		})

		Convey("Having Default param policy and param containing model name", func() {
			model := ModelWithID{}
			req := httptest.NewRequest("GET", "/", nil)

			policy := DefaultParamPolicy.New()
			valueMap := map[string]string{"modelwithid": "1234"}

			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeNil)
			So(model.ID, ShouldEqual, 1234)

		})

		Convey(`Having a DeepSearched ParamPolicy on a model with 
			structs and a time field`, func() {
			type deepModel struct {
				ID                int
				Name              string
				ReferencedModel   ModelWithID
				Disallowed        int `param:"-"`
				private           int
				DateCreated       time.Time  `param:"date" time_format:"2006-01-02" time_utc:"1" time_location:"Asia/Chongqing"`
				PtrDate           *time.Time `param:"ptrdate" time_format:"2006-01-02"`
				SliceIsNotAllowed []int
				PtrModel          *ModelWithID `param:"ptrmodel"`
				Ptrint            *int         `param:"ptrint"`
				PtrSlice          *[]int
			}
			model := deepModel{}
			req := httptest.NewRequest("GET", "/", nil)
			policy := DefaultParamPolicy.New()
			policy.DeepSearch = true

			valueMap := map[string]string{
				"deepmodel":         "123",
				"modelwithid":       "456",
				"date":              "2017-01-20",
				"ptrdate":           "2017-04-30",
				"sliceisnotallowed": "1,2,3",
				"ptrmodel":          "1",
				"ptrint":            "1",
			}
			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeNil)
			So(model.ID, ShouldEqual, 123)
			So(model.ReferencedModel.ID, ShouldEqual, 456)
			year, month, day := model.DateCreated.Date()
			So(year, ShouldEqual, 2017)
			So(int(month), ShouldEqual, 01)
			So(day, ShouldEqual, 20)
			Println(model.PtrDate)
			So(model.PtrDate, ShouldNotBeNil)
			year, month, day = model.PtrDate.Date()
			So(*model.Ptrint, ShouldEqual, 1)

			Convey("If provided incorrect time field and FailOnError is set to true", func() {
				model := deepModel{}
				valueMap["date"] = "201701-20"
				policy.FailOnError = true
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})

			Convey("If provided incorrect param for nested struct", func() {
				model = deepModel{}
				valueMap["modelwithid"] = "maciej"
				policy := DefaultParamPolicy.New()
				policy.DeepSearch = true
				policy.FailOnError = true
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})

			Convey("If provided incorrect param for *nested struct", func() {
				model = deepModel{}
				valueMap["ptrdate"] = "2017093-1"
				policy.FailOnError = true
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})

			Convey("If tagged Only ", func() {
				policy := DefaultParamPolicy.New()
				policy.TaggedOnly = true
				policy.DeepSearch = true
				valueMap["ptrmodel"] = ""
				valueMap["ptrdate"] = ""
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})
		})
		Convey("Having a model with Tagged ID and ptr struct fields", func() {
			type modelTaggedId struct {
				ID      int `param:"id"`
				PtrDate *time.Time
				PtrInt  *int
			}
			policy := DefaultParamPolicy.New()
			policy.TaggedOnly = true
			policy.DeepSearch = true
			req := httptest.NewRequest("GET", "/", nil)
			model := modelTaggedId{}
			valueMap := map[string]string{"id": "123"}
			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeNil)
		})

		Convey("Having a model with ptr type elem and policy without deep search", func() {
			type modelWithPtrs struct {
				PtrStruct *[]int
				PtrInt    *int
				PtrModel  *ModelWithID
				Strct     ModelWithID
				ID        int
			}

			policy := DefaultParamPolicy.New()
			policy.TaggedOnly = false
			req := httptest.NewRequest("GET", "/", nil)
			model := modelWithPtrs{}
			valueMap := map[string]string{"modelwithptrs": "123"}
			err := mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
			So(err, ShouldBeNil)

			Convey("If fail on error", func() {
				policy.DeepSearch = true
				policy.FailOnError = true
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeError)
			})

			Convey("If deep search is true", func() {
				policy = DefaultParamPolicy.New()
				policy.DeepSearch = false
				err = mapParam(&model, getParamFuncWithValues(valueMap), req, policy)
				So(err, ShouldBeNil)
			})
		})

		Convey("Having model and a getParam func which fails on error", func() {
			type modelToTest struct {
				Name string
				ID   int
			}
			policy := DefaultParamPolicy.New()
			policy.DeepSearch = true
			policy.FailOnError = true
			errorMap := map[string]error{
				"name": errors.New("Error"),
			}
			req := httptest.NewRequest("GET", "/", nil)
			err := mapParam(&modelToTest{}, getParamErrFunc(errorMap), req, policy)
			So(err, ShouldBeError)
		})
	})
}

func TestBindParams(t *testing.T) {
	Convey("Subject: bind url parameters using BindParam function", t, func() {
		Convey(`Having some model, url parameter function and http request`, func() {
			req := httptest.NewRequest("GET", "/", nil)
			Convey(`If used a nonaddressable (not a pointer to) model, 
				the function should Panic`, func() {
				model := ModelWithID{}
				So(func() { BindParams(req, model, emptyGetParamFunc, DefaultParamPolicy.New()) }, ShouldPanic)
			})

			Convey(`If no params provided, or the ParamGetterFunc provided 
				doesn't contain models id, an error should be thrown`, func() {
				model := &ModelWithID{}
				err := BindParams(req, model, emptyGetParamFunc, DefaultParamPolicy.New())
				So(err, ShouldBeError)
			})
			Convey(`If no policy provided, by default 
				'DefaultParamPolicy', would be used`, func() {
				var policy *ParamPolicy = nil
				model := &ModelWithID{}
				BindParams(req, model, emptyGetParamFunc, policy)
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

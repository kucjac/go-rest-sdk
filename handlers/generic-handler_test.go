package handlers

import (
	"encoding/json"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/repository/mockrepo"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Model struct {
	ID   int
	Name string
}

func TestNewHandler(t *testing.T) {
	Convey(`Subject: New GenericHandler `, t, func() {
		errHandler := errhandler.New()
		Convey(`Having some repository and error handler 
			a GenericHandler should be created`, func() {
			repo := &mockrepo.MockRepository{}

			genericHandler, err := New(repo, errHandler, nil, nil)

			So(err, ShouldBeNil)
			So(genericHandler, ShouldNotBeNil)
			So(genericHandler.Log, ShouldNotBeNil)
		})

		Convey("If repository or error handler is nil, an error should be returned", func() {
			var repo repository.Repository

			genericHandler, err := New(repo, errHandler, nil, nil)

			So(err, ShouldBeError)
			So(genericHandler, ShouldBeNil)
		})

		Convey(`Having some genericHandler using New() method 
			creates and returns by callback a copy of given handler`, func() {
			repo := &mockrepo.MockRepository{}
			genericHandler, err := New(repo, errHandler, nil, nil)

			So(err, ShouldBeNil)
			createdHandler := genericHandler.New()

			So(createdHandler, ShouldNotBeNil)
			So(createdHandler, ShouldResemble, genericHandler)
			So(createdHandler, ShouldNotPointTo, genericHandler)

		})
	})
}

func TestCallbacks(t *testing.T) {
	Convey("Subject: Callback method for GenericHandler", t, func() {

		Convey("Having some generic Handler", func() {
			handler := &GenericHandler{}

			Convey(`By using WithQueryPolicy a policy is being saved in the handler and the handler is being returned`, func() {
				So(handler.QueryPolicy, ShouldBeNil)

				callbacked := handler.WithQueryPolicy(forms.DefaultBindPolicy.Copy())

				So(handler.QueryPolicy, ShouldNotBeNil)
				So(callbacked, ShouldPointTo, handler)
			})
			Convey(`WithParamPolicy method sets the ParamPolicy and callback the handler`, func() {
				So(handler.ParamPolicy, ShouldBeNil)

				callbacked := handler.WithParamPolicy(forms.DefaultParamPolicy.Copy())
				So(handler.ParamPolicy, ShouldNotBeNil)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`WithListParameters sets the ListParameters and 
				the flag IncludeListCount`, func() {

				So(handler.ListParams, ShouldBeNil)

				callbacked := handler.WithListParameters(
					&repository.ListParameters{Limit: 1},
				)

				So(handler.ListParams, ShouldNotBeNil)

				So(callbacked, ShouldPointTo, handler)
			})

			Convey("WithSelectCount sets the IncludeCount flag", func() {
				So(handler.IncludeListCount, ShouldBeFalse)
				callbacked := handler.WithSelectCount(true)
				So(handler.IncludeListCount, ShouldBeTrue)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`With URLParameters sets the UseURLParameters flag`, func() {
				So(handler.UseURLParams, ShouldBeZeroValue)

				callbacked := handler.WithURLParams(true)

				So(handler.UseURLParams, ShouldBeTrue)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`WithParamGetterFunc sets the ParamGetterFunc for given
				handler.`, func() {
				So(handler.GetParams, ShouldBeNil)

				mockParamGetter := func(param string, req *http.Request) (string, error) {
					return "", nil
				}

				callbacked := handler.WithParamGetterFunc(mockParamGetter)

				So(handler.GetParams, ShouldNotBeNil)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`WithResposneBody sets the response.Response for given handler`, func() {
				So(handler.ResponseBody, ShouldBeNil)

				callbacked := handler.WithResponseBody(&MockResponser{})
				So(handler.ResponseBody, ShouldResemble, &MockResponser{})
				So(callbacked, ShouldPointTo, handler)
			})

		})
	})
}

func TestCreateHandlerfunc(t *testing.T) {
	Convey(`Subject: Create Handlerfunc for GenericHandler`, t, func() {
		Convey("Having some http server and GenericHandler", func() {
			server := http.NewServeMux()
			repo := &mockrepo.MockRepository{}
			errHandler := errhandler.New()
			handler, err := New(repo, errHandler, nil, nil)
			So(err, ShouldBeNil)

			policy := forms.DefaultBindPolicy.Copy()
			policy.FailOnError = true
			server.Handle("/models", handler.Create(Model{}))

			type NestedModel struct {
				ID         int
				Nested     Model
				SomeString string `json:"someString"`
			}

			Convey("Having some request with incorrect json body", func() {
				req := httptest.NewRequest("POST", "/models", strings.NewReader(
					`"name": "nonclosed`,
				))
				rw := httptest.NewRecorder()

				server.ServeHTTP(rw, req)

				Convey(`A response with status 400 and containing restError with InvalidJSONDocument, should be responsed`, func() {
					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(rw.Code, ShouldEqual, 400)
					So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
				})
			})
			Convey(`Having a UseURLParameters flag`, func() {
				handler.WithURLParams(true)
				paramPolicy := forms.DefaultParamPolicy.Copy()
				handler.WithParamPolicy(paramPolicy)
				Convey(`If an error occurred during param biding 
					rest Internal Error would be jsoned.`, func() {

					server.Handle("/models/5/nested", handler.Create(NestedModel{}))

					req := httptest.NewRequest("POST", "/models/5/nested", strings.NewReader(
						`{"someString": "some info"}`))
					rw := httptest.NewRecorder()

					server.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(rw.Code, ShouldEqual, 500)
					So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
				})
				Convey("If no error occurred during param binding", func() {

					paramPolicy.SearchDepthLevel = 1

					handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{
						"model":       "6",
						"nestedmodel": "3",
					}))
					nested := &NestedModel{ID: 3, SomeString: "some info", Nested: Model{
						ID: 6,
					}}

					Convey("if no error occurred from database", func() {
						repo.On("Create", nested).Return(nil)
						handler.Repo = repo

						server.Handle("/models/6/nested", handler.Create(NestedModel{}))

						req := httptest.NewRequest("POST", "/models/6/nested", strings.NewReader(
							`{"someString": "some info"}`))
						rw := httptest.NewRecorder()

						server.ServeHTTP(rw, req)

						body, err := readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors, ShouldBeEmpty)
					})
					Convey("if an error occurred from database error", func() {
						repo.On("Create", nested).Return(dberrors.ErrCardinalityViolation.New())

						server.Handle("/models/6/nested", handler.Create(NestedModel{}))

						req := httptest.NewRequest("POST", "/models/6/nested", strings.NewReader(
							`{"someString": "some info"}`))
						rw := httptest.NewRecorder()

						server.ServeHTTP(rw, req)

						body, err := readBody(rw)
						So(err, ShouldBeNil)
						So(body.Errors, ShouldNotBeEmpty)
					})

				})

			})
		})
	})
}

func TestGetMethod(t *testing.T) {
	Convey("Subject: Get method for GenericHandler", t, func() {
		server := http.NewServeMux()
		repo := &mockrepo.MockRepository{}
		errHandler := errhandler.New()

		handler, err := New(repo, errHandler, nil, nil)
		So(err, ShouldBeNil)

		paramPolicy := forms.DefaultParamPolicy.Copy()
		handler.WithParamPolicy(paramPolicy)

		Convey(`Using URL parameters would bind them to model or throws errors`, func() {
			req := httptest.NewRequest("GET", "/models/1", nil)
			rw := httptest.NewRecorder()
			handler = handler.New()
			server.Handle("/models/1", handler.WithURLParams(true).Get(Model{}))

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)
			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})
		Convey("if correctly get from db", func() {
			repo.On("Get", &Model{ID: 1}).Return(&Model{ID: 1, Name: "This"}, nil)
			req := httptest.NewRequest("GET", "/models/1", nil)
			rw := httptest.NewRecorder()
			handler = handler.New()
			handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "1"}))
			handler.Repo = repo
			server.Handle("/models/1", handler.WithURLParams(true).Get(Model{}))

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)
			So(body.Errors, ShouldBeEmpty)

			model, ok := body.Content["model"].(map[string]interface{})
			So(ok, ShouldBeTrue)
			So(model["ID"], ShouldEqual, 1)
			So(model["Name"], ShouldEqual, "This")
		})

		Convey("If an error occured while getting from db", func() {
			repo.On("Get", &Model{ID: 1}).Return(nil, dberrors.ErrNoResult.New())
			req := httptest.NewRequest("GET", "/models/1", nil)
			rw := httptest.NewRecorder()
			handler = handler.New()
			handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "1"}))
			handler.Repo = repo
			server.Handle("/models/1", handler.WithURLParams(true).Get(Model{}))

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)

			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrResourceNotFound), ShouldBeTrue)
		})
	})
}

func TestListMethod(t *testing.T) {
	Convey("Subject: List method for GenericHandler", t, func() {
		server := http.NewServeMux()
		repo := &mockrepo.MockRepository{}
		errHandler := errhandler.New()

		handler, err := New(repo, errHandler, nil, nil)
		So(err, ShouldBeNil)
		Convey("Having an error on query", func() {
			policy := forms.DefaultBindPolicy.Copy()
			policy.FailOnError = true

			handler.
				WithQueryPolicy(policy).
				WithListParameters(&repository.ListParameters{})

			server.Handle("/models", handler.List(Model{}))

			req := httptest.NewRequest("GET", "/models?id=string", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)

			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInvalidQueryParameter), ShouldBeTrue)

			req = httptest.NewRequest("GET", "/models?limit=string", nil)
			rw = httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err = readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInvalidQueryParameter), ShouldBeTrue)
		})

		Convey("With url parameter", func() {
			paramPolicy := forms.DefaultParamPolicy.Copy()
			paramPolicy.FailOnError = true

			handler.WithParamPolicy(paramPolicy).WithURLParams(true)

			server.Handle("/models/5/nested", handler.List(Model{}))
			req := httptest.NewRequest("GET", "/models/5/nested", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})

		Convey("if listParameters are taken from query", func() {
			handler.WithQueryPolicy(forms.DefaultBindPolicy.Copy())
			listParam := &repository.ListParameters{Limit: 5}
			handler.WithListParameters(listParam).WithSelectCount(true)

			repo.On("ListWithParams", &Model{}, listParam).Return([]*Model{
				{ID: 1, Name: "String"},
				{ID: 2, Name: "Ss"},
			}, nil)
			repo.On("Count", Model{}).Return(2, nil)
			server.Handle("/models", handler.List(Model{}))

			req := httptest.NewRequest("GET", "/models", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldBeEmpty)
			So(body.Content, ShouldNotBeEmpty)

		})
		Convey("Getting a list with no count", func() {
			repo.On("List", &Model{}).Return([]*Model{
				{ID: 1, Name: "First"},
				{ID: 2, Name: "Second"},
			}, nil)

			server.Handle("/models", handler.List(Model{}))

			req := httptest.NewRequest("GET", "/models", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldBeEmpty)
			So(body.Content, ShouldNotBeEmpty)
		})

		Convey("If DBError would be returned, it is returned as resterror", func() {
			repo.On("List", &Model{}).Return([]*Model{}, dberrors.ErrInternalError.New())

			server.Handle("/models", handler.List(Model{}))

			req := httptest.NewRequest("GET", "/models", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})
		Convey("If error occurred during count db", func() {
			repo.On("List", &Model{}).Return([]*Model{}, nil)
			repo.On("Count", Model{}).Return(0, dberrors.ErrInternalError.New())

			handler.WithSelectCount(true)

			server.Handle("/models/listcount", handler.List(Model{}))
			req := httptest.NewRequest("GET", "/models/listcount", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})
	})
}

func TestUpdateMethod(t *testing.T) {
	Convey("Subject: Update method for GenericHandler", t, func() {
		server := http.NewServeMux()
		repo := &mockrepo.MockRepository{}
		errHandler := errhandler.New()

		handler, err := New(repo, errHandler, nil, nil)
		So(err, ShouldBeNil)

		Convey("If error occurred during binding json", func() {
			req := httptest.NewRequest("PUT", "/models/123", strings.NewReader(
				`{"something": "incorrect}`))
			rw := httptest.NewRecorder()

			server.Handle("/models/123", handler.Update(Model{}))

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
		})

		Convey("if an error occurred during param binding", func() {
			paramPolicy := forms.DefaultParamPolicy.Copy()
			handler.WithParamPolicy(paramPolicy)
			// no ParamGetterFunc

			server.Handle("/models/1234", handler.WithURLParams(true).Update(Model{}))
			req := httptest.NewRequest("PUT", "/models/1234",
				strings.NewReader(`{"name": "my"}`),
			)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})

		Convey("If everything is going fine, an object should be returned", func() {
			repo.On("Update", &Model{ID: 1234, Name: "my"}).Return(nil)
			handler.WithParamPolicy(forms.DefaultParamPolicy.Copy())
			handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "1234"}))
			server.Handle("/models/1234", handler.WithURLParams(true).Update(Model{}))
			req := httptest.NewRequest("PUT", "/models/1234",
				strings.NewReader(`{"name": "my"}`),
			)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)

			So(err, ShouldBeNil)
			So(body.Errors, ShouldBeEmpty)
			So(body.Content, ShouldNotBeEmpty)
		})

		Convey("If an error occurred during updating to database", func() {
			repo.On("Update", &Model{ID: 1234, Name: "my"}).
				Return(dberrors.ErrUniqueViolation.New())

			handler.WithParamPolicy(forms.DefaultParamPolicy.Copy())
			handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "1234"}))
			server.Handle("/models/1234", handler.WithURLParams(true).Update(Model{}))

			req := httptest.NewRequest("PUT", "/models/1234",
				strings.NewReader(`{"name": "my"}`),
			)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)

			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrResourceAlreadyExists), ShouldBeTrue)
		})
	})
}

func TestPatchMethod(t *testing.T) {
	Convey("Subject: Patch method for GenericHandler", t, func() {
		server := http.NewServeMux()
		repo := &mockrepo.MockRepository{}
		errHandler := errhandler.New()

		handler, err := New(repo, errHandler, nil, nil)
		So(err, ShouldBeNil)

		Convey("If an error occured during binding URL Parameters", func() {
			paramPolicy := forms.DefaultParamPolicy.Copy()
			// no ParamGetterFunc

			server.Handle("/model/badurl", handler.
				WithURLParams(true).
				WithParamPolicy(paramPolicy).
				Patch(Model{}))

			req := httptest.NewRequest("PATCH", "/model/badurl", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, err := readBody(rw)
			So(err, ShouldBeNil)

			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})
		Convey("If an error occurred during json binding", func() {
			server.Handle("/models/123", handler.Patch(Model{}))

			req := httptest.NewRequest("PATCH", "/models/123",
				strings.NewReader(`{"name": "my name"`))
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, _ := readBody(rw)
			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
		})

		Convey("If correctly taken values form url and json", func() {
			paramPolicy := forms.DefaultParamPolicy.Copy()

			handler.WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "123"}))
			server.Handle("/models/123", handler.
				WithURLParams(true).
				WithParamPolicy(paramPolicy).
				Patch(Model{}))

			req := httptest.NewRequest("PATCH", "/models/123", strings.NewReader(
				`{"name": ""}`))
			rw := httptest.NewRecorder()

			Convey("If an error occurred during patching on db", func() {
				repo.On("Patch", &Model{}, &Model{ID: 123}).
					Return(dberrors.ErrCheckViolation.New())
				server.ServeHTTP(rw, req)
				body, _ := readBody(rw)
				So(body.Errors, ShouldNotBeEmpty)
				So(body.Errors[0].Compare(resterrors.ErrInvalidInput), ShouldBeTrue)
			})
			Convey("If no error occurred during patching, but while getting", func() {
				repo.On("Patch", &Model{}, &Model{ID: 123}).Return(nil)
				repo.On("Get", &Model{ID: 123}).Return(nil, dberrors.ErrInternalError.New())

				server.ServeHTTP(rw, req)
				body, _ := readBody(rw)
				So(body.Errors, ShouldNotBeEmpty)
				So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
			})
			Convey("If no error occurred than the patched object is being responded", func() {
				repo.On("Patch", &Model{Name: "Some Name"}, &Model{ID: 123}).Return(nil)
				responded := &Model{ID: 123, Name: "Some Name"}
				repo.On("Get", &Model{ID: 123}).Return(responded, nil)

				req = httptest.NewRequest("PATCH", "/models/123",
					strings.NewReader(`{"name":"Some Name"}`))

				server.ServeHTTP(rw, req)

				body, _ := readBody(rw)
				So(body.Errors, ShouldBeEmpty)
				So(body.Content, ShouldNotBeEmpty)

				model, ok := body.Content["model"].(map[string]interface{})
				So(ok, ShouldBeTrue)

				So(model["ID"], ShouldEqual, responded.ID)
				So(model["Name"], ShouldEqual, responded.Name)
			})

		})
	})
}

func TestDeleteMethod(t *testing.T) {
	Convey("Subject: Delete method for GenericHandler", t, func() {
		server := http.NewServeMux()
		repo := &mockrepo.MockRepository{}
		errHandler := errhandler.New()

		handler, err := New(repo, errHandler, nil, nil)
		So(err, ShouldBeNil)

		Convey(`if an error occurred during param binding 
			a rest error should be returned`, func() {
			server.Handle("/delete/invalidurl", handler.
				WithParamPolicy(forms.DefaultParamPolicy.Copy()).
				WithURLParams(true).Delete(Model{}),
			)
			req := httptest.NewRequest("DELETE", "/delete/invalidurl", nil)
			rw := httptest.NewRecorder()

			server.ServeHTTP(rw, req)

			body, _ := readBody(rw)
			So(body.Errors, ShouldNotBeEmpty)
			So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
		})

		Convey("If no error occurred during param binding", func() {
			server.Handle("/models/123", handler.
				WithURLParams(true).
				WithParamGetterFunc(getParamFuncWithValues(map[string]string{"model": "123"})).
				WithParamPolicy(forms.DefaultParamPolicy.Copy()).Delete(Model{}))
			req := httptest.NewRequest("DELETE", "/models/123", nil)
			rw := httptest.NewRecorder()

			Convey("But an error occurred during deleting from repository", func() {
				repo.On("Delete", &Model{}, &Model{ID: 123}).Return(dberrors.ErrNoResult.New())

				server.ServeHTTP(rw, req)
				body, _ := readBody(rw)

				So(body.Errors, ShouldNotBeEmpty)
				So(body.Errors[0].Compare(resterrors.ErrResourceNotFound), ShouldBeTrue)
			})

			Convey("if everytihng is ok, an item should be deleted from repository", func() {
				repo.On("Delete", &Model{}, &Model{ID: 123}).Return(nil)

				server.ServeHTTP(rw, req)
				body, _ := readBody(rw)

				So(body.Errors, ShouldBeEmpty)
				So(rw.Code, ShouldEqual, 200)
			})

		})
	})
}

func TestJSONMethod(t *testing.T) {
	Convey("Subject: JSON method for GenericHandler", t, func() {
		Convey("Having some GenericHandler, request, response and some mux", func() {
			handler, err := New(&mockrepo.MockRepository{}, errhandler.New(), nil, nil)
			So(err, ShouldBeNil)

			mux := http.NewServeMux()

			jsonHandlerFunc := func(status int, body response.Responser) http.HandlerFunc {
				return func(rw http.ResponseWriter, req *http.Request) {
					handler.JSON(rw, req, status, body)
				}
			}

			mux.Handle("/JSONNoStatus", jsonHandlerFunc(0, &response.DefaultBody{}))
			mux.Handle("/JSONErrMarshal", jsonHandlerFunc(123, &MockResponser{}))

			Convey(`Having a request on '/JSONNoStatus' 
				that sets status by default to 200`, func() {
				req := httptest.NewRequest("GET", "/JSONNoStatus", nil)
				rw := httptest.NewRecorder()

				mux.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, 200)
			})

			Convey(`Having a request that marshals incorrect response.Response,
				and marshal default response.DefaultBody`, func() {
				req := httptest.NewRequest("GET", "/JSONErrMarshal", nil)
				rw := httptest.NewRecorder()

				mux.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, http.StatusInternalServerError)

				body, err := readBody(rw)
				So(err, ShouldBeNil)
				So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
			})
		})
	})
}

func TestHandleDBError(t *testing.T) {
	Convey("Subject: private functions", t, func() {
		handler, _ := New(&mockrepo.MockRepository{}, errhandler.New(), nil, nil)
		Convey("Providing unknown error to the handleDBError", func() {
			mydbErr := &dberrors.Error{ID: 1234}

			req := httptest.NewRequest("GET", "/", nil)
			rw := httptest.NewRecorder()

			handler.handleDBError(rw, req, mydbErr)
		})

		Convey("Providing  StatusResponser", func() {
			body := &response.DetailedBody{}
			handler.ResponseBody = body
			So(body, ShouldImplement, (*response.StatusResponser)(nil))
			Convey("To getResponseBodyErr", func() {
				handler.getResponseBodyErr(123)
			})
			Convey("To getResponseBodyContent", func() {
				handler.getResponseBodyContent(123)
			})
		})
	})
}

func erroredSetIDFunc(req *http.Request, model interface{}) error {
	return ErrIncorrectModel
}

func monerroredSetIDFunc(req *http.Request, model interface{}) error {
	return nil
}

type MockResponser struct {
	Channel chan (int) `json:"invalid"`
}

func (m *MockResponser) AddContent(content ...interface{}) {}

func (m *MockResponser) AddErrors(errors ...*resterrors.Error) {}

func (m *MockResponser) WithErrors(errors ...*resterrors.Error) response.Responser {
	return m
}

func (m *MockResponser) WithContent(content ...interface{}) response.Responser {
	return m
}

func (m *MockResponser) New() response.Responser {
	return m
}

func (m *MockResponser) NewErrored() response.Responser {
	return m
}

func readBody(rw *httptest.ResponseRecorder) (body *response.DefaultBody, err error) {
	rsp, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		return nil, err
	}
	body = new(response.DefaultBody)
	err = json.Unmarshal(rsp, &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getParamFuncWithValues(
	paramValues map[string]string,
) forms.ParamGetterFunc {
	return func(paramName string, req *http.Request) (string, error) {
		value := paramValues[paramName]
		return value, nil
	}
}

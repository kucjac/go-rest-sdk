package ginhandler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/repository/mocks"
	"github.com/kucjac/go-rest-sdk/response"
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type Model struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ModelWithForm struct {
	Model
	Age int `json:"age" form:"age"`
}

type IncorrectModel struct {
	Name string
}

func TestNew(t *testing.T) {

	Convey("Subject: New JSONHandler", t, func() {
		var err error
		var gjh1, gjh2, gjh3 *JSONHandler

		Convey(`Having some repository that implements repository.GenericRepository 
			and error handler the correct json handler should be created`, func() {
			repo := &mocks.MockRepository{}
			errHandler := errhandler.New()
			gjh1, err = New(repo, errHandler, nil)

			So(err, ShouldBeNil)
			So(gjh1, ShouldNotBeNil)
			gjh1 = nil
		})

		Convey("If either repository or error handler is not provided an error would be returned", func() {
			repo := &mocks.MockRepository{}
			var nilErrHandler *errhandler.ErrorHandler = nil

			var nilRepo repository.GenericRepository = nil
			errHandler := errhandler.New()

			gjh1, err = New(repo, nilErrHandler, nil)
			So(err, ShouldBeError)
			So(gjh1, ShouldBeNil)

			gjh2, err = New(nilRepo, errHandler, nil)
			So(err, ShouldBeError)
			So(gjh2, ShouldBeNil)

			gjh3, err = New(nilRepo, nilErrHandler, nil)
			So(err, ShouldBeError)
			So(gjh3, ShouldBeNil)

			gjh1, gjh2, gjh3 = nil, nil, nil
		})

		Convey("The responseBody argument may be nil, in such a case the default *DefaultBody would be used.", func() {
			repo := &mocks.MockRepository{}
			errHandler := errhandler.New()
			var body response.Responser

			gjh1, err = New(repo, errHandler, body)
			So(err, ShouldBeNil)

			So(gjh1.responseBody, ShouldHaveSameTypeAs, &response.DefaultBody{})

		})
	})
}

func TestWithResponseBody(t *testing.T) {
	Convey("Subject: Setting responseBody with callback method WithResponseBody", t, func() {
		Convey("Having some JSONHandler with default responseBody", func() {
			handler, _ := New(&mocks.MockRepository{}, errhandler.New(), nil)

			So(handler.responseBody, ShouldResemble, &response.DefaultBody{})
			Convey(`Using WithResponseBody() method sets the responseBody as in the argument and returns given handler as callback`, func() {
				callbacked := handler.WithResponseBody(&response.DetailedBody{})

				So(callbacked, ShouldEqual, handler)
				So(callbacked.responseBody, ShouldResemble, &response.DetailedBody{})
			})
		})
	})
}

func TestCreateHandlerfunc(t *testing.T) {
	Convey("Subject: Create gin.Handlerfunc", t, func() {
		var err error
		var gjh *JSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var policy *forms.Policy
		var body *response.DefaultBody
		var req *http.Request
		var rw *httptest.ResponseRecorder
		gin.SetMode(gin.TestMode)

		Convey("Having an error handler and some gin router", func() {

			router = gin.New()

			errHandler = errhandler.New()

			repo := &mocks.MockRepository{}
			policy = &forms.Policy{FailOnError: true}

			gjh, err = New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			gjh.jsonPolicy = policy

			router.POST("/model", gjh.Create(Model{}))

			Convey("Handling POST request with incorrect json form", func() {
				incorrectBody := strings.NewReader(`{"id":"stringID", "name":"IncorrectID"}`)

				req = httptest.NewRequest("POST", "/model", incorrectBody)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)
				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldNotBeEmpty)
				So(body.Content, ShouldBeEmpty)

				incorrectJSON := strings.NewReader(`{"id":1`)

				req = httptest.NewRequest("POST", "/model", incorrectJSON)
				rw = httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(rw.Code, ShouldEqual, 400)
				So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)

				So(body.Content, ShouldBeEmpty)
			})

			Convey(`Handling POST request of correct JSON form
				where some clientside error occurs while connecting repository`, func() {
				var dberr *dberrors.Error = dberrors.ErrUniqueViolation.New()
				repo.On("Create", &Model{Name: "Duplicated"}).Return(dberr)

				duplicatedModel := strings.NewReader(`{"name": "Duplicated"}`)

				req = httptest.NewRequest("POST", "/model", duplicatedModel)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err = readBody(rw)
				So(err, ShouldBeNil)

				resterr, err := errHandler.Handle(dberr)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldContain, resterr)
			})

			Convey(`Handling POST request of correct JSON form
				when some unknown error occurs while using repository`, func() {

				var dberr *dberrors.Error = &dberrors.Error{Title: "Unknown error"}
				repo.On("Create", &Model{Name: "Bad error"}).Return(dberr)

				req = httptest.NewRequest("POST", "/model",
					strings.NewReader(`{"name": "Bad error"}`))
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
				So(body.Content, ShouldBeEmpty)
			})

			Convey(`Handling POST request of correct JSON form
				when some internal dberror occurs while using repository`, func() {
				var dberr *dberrors.Error = dberrors.ErrInsufficientResources.New()
				repo.On("Create", &Model{Name: "Internal Error"}).Return(dberr)

				req = httptest.NewRequest("POST", "/model",
					strings.NewReader(`{"name": "Internal Error"}`))
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
				So(body.Content, ShouldBeEmpty)
			})

			FocusConvey(`Handling POST request of correct JSON with succesful creation`, func() {
				correct := &Model{Name: "Correct Model"}
				// the repo would add id = 1
				repo.On("Create", correct).Return(nil).
					Run(func(args mock.Arguments) { args[0].(*Model).ID = 1 })

				req = httptest.NewRequest("POST", "/model",
					strings.NewReader(`{"name": "Correct Model"}`))
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldBeEmpty)

				modelFields, ok := body.Content["model"].(map[string]interface{})
				So(ok, ShouldBeTrue)

				So(modelFields["name"], ShouldEqual, correct.Name)
				So(modelFields["id"], ShouldEqual, 1)
			})
		})
	})
}

func TestGetHandlerfunc(t *testing.T) {

	Convey("Subject: Get gin.Handlerfunc", t, func() {
		var err error
		var gjh *JSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine

		Convey("Having a ginJSONHandler and some gin router", func() {
			router = gin.New()
			repo := &mocks.MockRepository{}
			errHandler = errhandler.New()

			gjh, err = New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			Convey("If parameter is named differently the handler would return response with http 500 error.", func() {

				router.GET("/incorrect/:model_id", gjh.Get(Model{}))

				req := httptest.NewRequest("GET", "/incorrect/1", strings.NewReader(`{"name": "Generic Model"}`))
				rw := httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
			})

			Convey(`If the model doesn't have ID field or doesn't implement
				IDSetter interface the method would return response with 500 error.`, func() {

				router.GET("/correct/:incorrectmodel", gjh.Get(IncorrectModel{}))
				req := httptest.NewRequest("GET", "/correct/1", strings.NewReader(`{"name": "No id model"}`))
				rw := httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())

			})

			Convey("While getting from repository an error is handled by errorHandler", func() {

				// Internal Error - not recognised
				repo.On("Get", &Model{ID: 50}).Return(nil,
					&dberrors.Error{Title: "Should not exist",
						Message: "I am not recognised by "},
				)

				// Internal Error - recognised
				repo.On("Get", &Model{ID: 51}).Return(nil, dberrors.ErrInternalError.New())

				// Client side error
				repo.On("Get", &Model{ID: 52}).Return(nil, dberrors.ErrNoResult.New())

				So(err, ShouldBeNil)

				router.GET("/errhandler/:model", gjh.Get(Model{}))

				Convey(`The handler may return non response error 'error'
					that means some kind of internal server error`, func() {

					req := httptest.NewRequest("GET", "/errhandler/50", nil)
					rw := httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
				})

				Convey("It can also handle dbError as InternalError", func() {
					req := httptest.NewRequest("GET", "/errhandler/51", nil)
					rw := httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
					So(body.Content, ShouldBeEmpty)
				})

				Convey("Or to the client side '400' errors", func() {
					req := httptest.NewRequest("GET", "/errhandler/52", nil)
					rw := httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldNotContain, resterrors.ErrInternalError.New())
					So(body.Errors, ShouldContain, resterrors.ErrResourceNotFound.New())
					So(body.Content, ShouldBeEmpty)
				})
			})

			Convey("If the given id entity exists the handler func would return it in the body->content->modelname", func() {

				mockedModel := &Model{ID: 1, Name: "First"}
				repo.On("Get", &Model{ID: 1}).Return(mockedModel, nil)

				So(err, ShouldBeNil)

				router.GET("/correct/:model", gjh.Get(Model{}))

				req := httptest.NewRequest("GET", "/correct/1", nil)
				rw := httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldBeEmpty)

				respModel, ok := body.Content["model"]
				So(ok, ShouldBeTrue)

				modelMap, ok := respModel.(map[string]interface{})
				So(ok, ShouldBeTrue)

				So(modelMap["name"], ShouldEqual, mockedModel.Name)
				So(modelMap["id"], ShouldEqual, mockedModel.ID)
			})

		})

	})

}

func TestListHandlerfunc(t *testing.T) {
	Convey("Subject: List gin.Handlerfunc", t, func() {
		var err error
		var gjh *JSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var repo *mocks.MockRepository = &mocks.MockRepository{}
		var req *http.Request
		var rw *httptest.ResponseRecorder

		Convey("Having some JSONHandler and some gin.Router", func() {

			errHandler = errhandler.New()
			gjh, err = New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			policy := &forms.Policy{FailOnError: true, TaggedOnly: true, Tag: "form"}

			router = gin.New()
			router.GET("/models", gjh.New().WithQueryPolicy(policy).List(ModelWithForm{}))

			Convey(`Handling the GET method with query policy`, func() {

				req = httptest.NewRequest("GET", "/models?age=fine", nil)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors[0].Compare(resterrors.ErrInvalidQueryParameter), ShouldBeTrue)
				So(body.Content, ShouldBeEmpty)
			})

			router.GET("/parametrized/models",
				gjh.New().WithParamPolicy(policy).List(Model{}))

			Convey(`Handling the GET method request on '/parametrized/models',
				with parameters policy in JSONHandler and incorrect
				parameters in url Query`, func() {

				req = httptest.NewRequest("GET", "/parametrized/models?ids=incorrect,3,4", nil)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors[0].Compare(resterrors.ErrInvalidQueryParameter), ShouldBeTrue)
				So(body.Content, ShouldBeNil)
			})

			Convey(`Handling the GET method request on '/parametrized/models'
				with parameters policy and correct url query parameters`, func() {
				var mockModels []*Model = []*Model{
					{ID: 1, Name: "First"},
					{ID: 2, Name: "Second"},
					{ID: 3, Name: "Third"},
				}
				repo.On("ListWithParams", &Model{}, &repository.ListParameters{IDs: []int{1, 2, 3}}).Return(mockModels, nil)
				req = httptest.NewRequest("GET", "/parametrized/models?ids=1&ids=2&ids=3", nil)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Errors, ShouldBeEmpty)

				var models []*Model

				modelJson, err := json.Marshal(body.Content["models"])
				So(err, ShouldBeNil)
				err = json.Unmarshal(modelJson, &models)
				So(err, ShouldBeNil)

				So(models, ShouldResemble, mockModels)
			})

			router.GET("/errored/models", gjh.WithQueryPolicy(&forms.Policy{TaggedOnly: false}).List(Model{}))
			Convey(`Having GET method request on '/errored/models'`, func() {
				Convey("When some clientside dberror occurs", func() {
					dbErr := dberrors.ErrNoResult.New()
					repo.On("List", &Model{Name: "Marcin"}).Return(nil, dbErr)

					restErr, _ := errHandler.Handle(dbErr)

					req = httptest.NewRequest("GET", "/errored/models?name=Marcin", nil)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors[0].Compare(*restErr), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})
				Convey("When some Internal error occurs", func() {
					dbErr := dberrors.ErrInternalError.New()
					repo.On("List", &Model{Name: "Some Internal Error"}).Return(nil, dbErr)

					restErr, _ := errHandler.Handle(dbErr)

					req = httptest.NewRequest("GET", "/errored/models?name="+url.QueryEscape("Some Internal Error"), nil)

					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldNotBeEmpty)
					So(body.Errors[0].Compare(*restErr), ShouldBeTrue)
					So(body.Content, ShouldBeNil)
				})
				Convey("When some unknown error occurs while handling error", func() {
					repo.On("List", &Model{Name: "Unknown error"}).Return(nil, &dberrors.Error{
						Message: "Unspecified and unknown error"})

					req = httptest.NewRequest("GET",
						"/errored/models?name="+url.QueryEscape("Unknown error"),
						nil,
					)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldNotBeEmpty)
					So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})

			})
		})

	})
}

func TestUpdateHandlerfunc(t *testing.T) {
	Convey("Subject: JSONHandler Update method", t, func() {
		var err error
		var gjh *JSONHandler
		var errHandler *errhandler.ErrorHandler = errhandler.New()
		var router *gin.Engine
		var repo *mocks.MockRepository = &mocks.MockRepository{}
		var req *http.Request
		var rw *httptest.ResponseRecorder
		var policy *forms.Policy = &forms.Policy{FailOnError: true}

		Convey("Having some gin.Router and JSONHandler", func() {
			router = gin.New()

			gjh, err = New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			Convey("Handling the request with method PUT on '/models/1'", func() {

				Convey(`If invalid id parameter is provided in the router url,
					an internal error would be returned `, func() {
					router.PUT("/incorrectid/:incorrectmodel_id", gjh.Update(Model{}))

					req = httptest.NewRequest("PUT", "/incorrectid/1", nil)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})

				Convey(`If an incorrect JSON body request was provided
					then the client side restError should be send`, func() {
					router.PUT("/correctid/:model", gjh.WithJSONPolicy(policy).Update(Model{}))

					jsonBody := strings.NewReader(`{"name" = "`)
					req = httptest.NewRequest("PUT", "/correctid/2", jsonBody)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})

				Convey(`If an incorrect model for setting id was provided, an internal error would be sent`, func() {

					router.PUT("/correctid/:incorrectmodel", gjh.Update(IncorrectModel{}))

					req = httptest.NewRequest("PUT", "/correctid/1",
						strings.NewReader(`{"name":"incorrect model"}`))
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})

				Convey("If the repository.Update method returns dberrors.Error", func() {
					router.PUT("/errored/:model", gjh.Update(Model{}))
					Convey(`If unrecognized dberrors.Error was provided an resterrors.ErrInternal should be returned`, func() {
						var dbErr *dberrors.Error = &dberrors.Error{Message: "Some unknown internal error occured."}
						repo.On("Update", &Model{ID: 3, Name: "Piesek"}).Return(dbErr)

						req = httptest.NewRequest("PUT", "/errored/3",
							strings.NewReader(`{"name": "Piesek"}`))
						rw = httptest.NewRecorder()

						router.ServeHTTP(rw, req)

						body, err := readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
						So(body.Content, ShouldBeEmpty)
					})

					Convey(`If an error was recognised as non Internal
						response should contain 400 code`, func() {

						dbErr := dberrors.ErrIntegrConstViolation.New()
						repo.On("Update", &Model{ID: -1, Name: "SomeName"}).Return(dbErr)

						restErr, _ := errHandler.Handle(dbErr)

						req = httptest.NewRequest("PUT", "/errored/-1",
							strings.NewReader(`{"name": "SomeName"}`))
						rw = httptest.NewRecorder()

						router.ServeHTTP(rw, req)

						body, err := readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors[0].Compare(*restErr), ShouldBeTrue)
						So(body.Content, ShouldBeEmpty)
					})
					Convey("If recognized Internal error was returned", func() {

						dbErr := dberrors.ErrSystemError.New()
						repo.On("Update", &Model{ID: 66, Name: "Very funny name"}).Return(dbErr)

						restErr, _ := errHandler.Handle(dbErr)

						req = httptest.NewRequest("PUT", "/errored/66",
							strings.NewReader(`{"name": "Very funny name"}`))
						rw = httptest.NewRecorder()

						router.ServeHTTP(rw, req)

						body, err := readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors[0].Compare(*restErr), ShouldBeTrue)
						So(body.Content, ShouldBeNil)
					})
				})
				Convey(`If correct request was made
					the updated item body would be returned`, func() {

					router.PUT("/correct/:model", gjh.New().WithJSONPolicy(policy).Update(Model{}))
					mockModel := &Model{ID: 6543, Name: "Correct Update"}
					repo.On("Update", mockModel).Return(nil)

					req = httptest.NewRequest("PUT", "/correct/6543",
						strings.NewReader(`{"name":"Correct Update"}`))

					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldBeEmpty)
					content, ok := body.Content["model"]
					So(ok, ShouldBeTrue)

					contentJson, err := json.Marshal(content)
					So(err, ShouldBeNil)

					var model *Model
					err = json.Unmarshal(contentJson, &model)
					So(err, ShouldBeNil)
					So(model, ShouldResemble, mockModel)
				})
			})
		})
	})
}

func TestPatchHandlerfunc(t *testing.T) {
	Convey("Subject: Patch Handlefunc for *JSONHandler", t, func() {

		var err error
		var handler *JSONHandler
		var errHandler *errhandler.ErrorHandler = errhandler.New()
		var router *gin.Engine
		var repo *mocks.MockRepository = &mocks.MockRepository{}
		var req *http.Request
		var rw *httptest.ResponseRecorder
		var policy *forms.Policy = &forms.Policy{FailOnError: true}
		var body *response.DefaultBody

		Convey("Having some *gin.Engine with *ginhandler.JSONHandler.", func() {
			handler, err = New(repo, errHandler, &response.DefaultBody{})
			So(err, ShouldBeNil)

			router = gin.New()

			router.PATCH("/incorrectparam/:incorrectparam_name", handler.Patch(Model{}))

			Convey(`Having a request with PATCH method on '/incorrectparam/1' 
				path that is handled by handler.Patch method with 'Model'`, func() {
				req = httptest.NewRequest("PATCH",
					"/incorrectparam/:incorrectparam_name",
					nil)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				Convey("Then the response should contain internal error", func() {
					body, err = readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
				})
			})
			Convey("If the handler binds restrict JSONPolicy", func() {
				handler.WithJSONPolicy(policy)

				router.PATCH("/incorrectjson/:model", handler.Patch(Model{}))
				Convey(`Having a request with PATCH method on '/incorrectjson/1',
					that contain incorrect json body...`, func() {

					req = httptest.NewRequest("PATCH", "/incorrectjson/1",
						strings.NewReader(`{"name": 123`),
					)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					Convey("should response with InvalidJSONDocument error", func() {
						body, err = readBody(rw)
						So(err, ShouldBeNil)
						Println(body.Errors[0])
						So(body.Errors[0].Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
					})
				})
			})
			Convey(`If the provided Model doesn't have ID field or doesn't implement
				IDSetter interface`, func() {

				router.PATCH("/incorrectmodel/:incorrectmodel", handler.Patch(IncorrectModel{}))

				Convey("Having a request with PATCH method on '/incorrectmodel/1'", func() {

					req = httptest.NewRequest("PATCH", "/incorrectmodel/1",
						strings.NewReader(`{"name": "Incorrect Model Name"}`))

					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					Convey("The Internal error should be resposned", func() {
						body, err = readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
					})
				})
			})
			Convey("If dberrors.Error occured during processing repository", func() {

				router.PATCH("/errors/:model", handler.Patch(Model{}))

				Convey(`Having request with PATCH method on '/errors/1',
				 when an error is unrecognised by the converter`, func() {
					unknownDBError := &dberrors.Error{Message: "Unknown Error"}

					repo.On("Patch", &Model{Name: "First"}, &Model{ID: 1}).
						Return(unknownDBError)
					req = httptest.NewRequest("PATCH",
						"/errors/1",
						strings.NewReader(`{"name": "First"}`),
					)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					Convey("An internal error should be responsed", func() {
						body, err = readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
					})
				})

				Convey(`Having a request with PATCH method on '/errors/2,
					when a client side error occured using a repository`, func() {
					dbErr := dberrors.ErrNoResult.New()
					restErr, err := errHandler.Handle(dbErr)
					So(err, ShouldBeNil)

					repo.On("Patch", &Model{Name: "Second"}, &Model{ID: 2}).Return(dbErr)

					req = httptest.NewRequest("PATCH", "/errors/2",
						strings.NewReader(`{"name": "Second"}`),
					)
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					Convey("A client side error should be responsed (with code '4xx')", func() {

						body, err = readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors[0].Compare(*restErr), ShouldBeTrue)
						So(rw.Code, ShouldEqual, 400)
					})
				})
				Convey(`Having a request with PATCH method on '/errors/3', 
					when an internal error occured during processing a repository`, func() {
					dbErr := dberrors.ErrInternalError.New()
					repo.On("Patch", &Model{Name: "Third"}, &Model{ID: 3}).Return(dbErr)

					req = httptest.NewRequest("PATCH", "/errors/3",
						strings.NewReader(`{"name": "Third"}`))
					rw = httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					Convey("An Internal Server error should be responsed", func() {
						body, err = readBody(rw)
						So(err, ShouldBeNil)

						So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
						So(rw.Code, ShouldEqual, 500)
					})
				})
			})
			Convey(`Having a succesful request with PATCH method on '/models/1`, func() {

				router.PATCH("/models/:model", handler.Patch(Model{}))
				repo.On("Patch", &Model{Name: "Success"}, &Model{ID: 1}).Return(nil)

				req = httptest.NewRequest("PATCH", "/models/1",
					strings.NewReader(`{"name": "Success"}`))
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				Convey("Should response with succesful Model in the content", func() {

					body, err = readBody(rw)
					So(err, ShouldBeNil)

					So(body.Errors, ShouldBeEmpty)
					So(rw.Code, ShouldEqual, 200)
				})

			})

		})

	})
}

func TestGetResponseBody(t *testing.T) {
	Convey("Subject: Test getResponseBody method for *JSONHandler", t, func() {

		repo := &mocks.MockRepository{}
		errHandler := errhandler.New()
		Convey(`Having a JSONHandler with responseBody that 
			doesn't implement StatusResponser`, func() {
			handler, _ := New(repo, errHandler, &response.DefaultBody{})

			Convey(`Using get body with error would return with 
				new *DefaultBody with given error`, func() {
				restErr := resterrors.ErrInvalidURI.New()
				body := handler.getResponseBodyErr(1234, restErr)

				resbody, ok := body.(*response.DefaultBody)
				So(ok, ShouldBeTrue)

				So(resbody.Errors, ShouldContain, restErr)
			})
		})
		Convey(`Having a JSONHandler with responseBody that implements StatusResponser`, func() {
			handler, _ := New(repo, errHandler, &response.DetailedBody{})
			Convey("Using getResponseBodyErr would set status and add error", func() {
				restErr := resterrors.ErrInvalidURI.New()
				body := handler.getResponseBodyErr(1234, restErr)

				stResBody, ok := body.(*response.DetailedBody)
				So(ok, ShouldBeTrue)

				So(stResBody.Errors, ShouldContain, restErr)
				So(stResBody.HttpStatus, ShouldEqual, 1234)
			})

			Convey("Using getResponseBodyCon method would set status and add content", func() {
				content := Model{ID: 1, Name: "Name"}
				body := handler.getResponseBodyCon(132, content)

				stResBody, ok := body.(*response.DetailedBody)
				So(ok, ShouldBeTrue)

				So(stResBody.Content["model"], ShouldResemble, content)
				So(stResBody.HttpStatus, ShouldEqual, 132)
			})

		})
	})
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

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

	Convey("Subject: New GinJSONHandler", t, func() {
		var err error
		var gjh1, gjh2, gjh3 *GinJSONHandler

		Convey(`Having some repository that implements repository.GenericRepository 
			and error handler the correct json handler should be created`, func() {
			repo := &mocks.MockRepository{}
			errHandler := errhandler.New()
			gjh1, err = New(repo, errHandler)

			So(err, ShouldBeNil)
			So(gjh1, ShouldNotBeNil)
			gjh1 = nil
		})

		Convey("If either repository or error handler is not provided an error would be returned", func() {
			repo := &mocks.MockRepository{}
			var nilErrHandler *errhandler.ErrorHandler = nil

			var nilRepo repository.GenericRepository = nil
			errHandler := errhandler.New()

			gjh1, err = New(repo, nilErrHandler)
			So(err, ShouldBeError)
			So(gjh1, ShouldBeNil)

			gjh2, err = New(nilRepo, errHandler)
			So(err, ShouldBeError)
			So(gjh2, ShouldBeNil)

			gjh3, err = New(nilRepo, nilErrHandler)
			So(err, ShouldBeError)
			So(gjh3, ShouldBeNil)

			gjh1, gjh2, gjh3 = nil, nil, nil
		})

		Convey("The policy argument may be nil, in such a case the default policy would be used.", func() {
			repo := &mocks.MockRepository{}
			errHandler := errhandler.New()
			var policy *forms.Policy

			gjh1, err = New(repo, errHandler)
			gjh1.JSONPolicy = policy
			So(err, ShouldBeNil)

		})
	})
}

func TestCreate(t *testing.T) {
	Convey("Subject: Create gin.Handlerfunc", t, func() {
		var err error
		var gjh *GinJSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var policy *forms.Policy
		var body *response.Body
		var req *http.Request
		var rw *httptest.ResponseRecorder
		gin.SetMode(gin.TestMode)

		Convey("Having an error handler and some gin router", func() {

			router = gin.New()

			errHandler = errhandler.New()

			repo := &mocks.MockRepository{}
			policy = &forms.Policy{FailOnError: true}

			gjh, err = New(repo, errHandler)
			So(err, ShouldBeNil)

			gjh.JSONPolicy = policy

			router.POST("/model", gjh.Create(Model{}))

			Convey("Handling POST request with incorrect json form", func() {
				incorrectBody := strings.NewReader(`{"id":"stringID", "name":"IncorrectID"}`)

				req = httptest.NewRequest("POST", "/model", incorrectBody)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)
				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 400)
				// So(body.Errors, ShouldContain, resterrors.ErrInvalidJSONDocument.New())
				So(body.Content, ShouldBeEmpty)

				incorrectJSON := strings.NewReader(`{"id":1`)

				req = httptest.NewRequest("POST", "/model", incorrectJSON)
				rw = httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				body, err = readBody(rw)
				So(err, ShouldBeNil)

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 400)
				for _, restErr := range body.Errors {
					So(restErr.Compare(resterrors.ErrInvalidJSONDocument), ShouldBeTrue)
				}
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

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 400)

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

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 500)
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

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 500)
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

				So(body.Status, ShouldEqual, response.StatusOk)
				So(body.HttpStatus, ShouldEqual, 200)
				So(body.Errors, ShouldBeEmpty)

				modelFields, ok := body.Content["model"].(map[string]interface{})
				So(ok, ShouldBeTrue)

				So(modelFields["name"], ShouldEqual, correct.Name)
				So(modelFields["id"], ShouldEqual, 1)
			})
		})
	})
}

func TestGet(t *testing.T) {

	Convey("Subject: Get gin.Handlerfunc", t, func() {
		var err error
		var gjh *GinJSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine

		Convey("Having a ginJSONHandler and some gin router", func() {
			router = gin.New()
			repo := &mocks.MockRepository{}
			errHandler = errhandler.New()
			gjh, err = New(repo, errHandler)

			Convey("If parameter is named differently the handler would return response with http 500 error.", func() {

				So(err, ShouldBeNil)
				router.GET("/incorrect/:model_id", gjh.Get(Model{}))

				req := httptest.NewRequest("GET", "/incorrect/1", strings.NewReader(`{"name": "Generic Model"}`))
				rw := httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body := response.Body{}
				rsp, err := ioutil.ReadAll(rw.Body)
				So(err, ShouldBeNil)

				err = json.Unmarshal(rsp, &body)
				So(err, ShouldBeNil)

				So(body.Status, ShouldEqual, response.StatusError)
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

				So(body.Status, ShouldEqual, response.StatusError)
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

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
				})

				Convey("It can also handle dbError as InternalError", func() {
					req := httptest.NewRequest("GET", "/errhandler/51", nil)
					rw := httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.HttpStatus, ShouldEqual, 500)
					So(body.Errors, ShouldContain, resterrors.ErrInternalError.New())
					So(body.Content, ShouldBeEmpty)
				})

				Convey("Or to the client side '400' errors", func() {
					req := httptest.NewRequest("GET", "/errhandler/52", nil)
					rw := httptest.NewRecorder()

					router.ServeHTTP(rw, req)

					body, err := readBody(rw)
					So(err, ShouldBeNil)

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.HttpStatus, ShouldEqual, 400)

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

				So(body.Status, ShouldEqual, response.StatusOk)
				So(body.Errors, ShouldBeEmpty)

				Println(body)
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

func TestList(t *testing.T) {
	Convey("Subject: List gin.Handlerfunc", t, func() {
		var err error
		var gjh *GinJSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var repo *mocks.MockRepository = &mocks.MockRepository{}
		var req *http.Request
		var rw *httptest.ResponseRecorder

		Convey("Having some GinJSONHandler and some gin.Router", func() {

			errHandler = errhandler.New()
			gjh, err = New(repo, errHandler)
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

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 400)
				So(body.Errors[0].Compare(resterrors.ErrInvalidQueryParameter), ShouldBeTrue)
				So(body.Content, ShouldBeEmpty)
			})

			router.GET("/parametrized/models",
				gjh.New().WithParamPolicy(policy).List(Model{}))

			Convey(`Handling the GET method request on '/parametrized/models', 
				with parameters policy in GinJSONHandler and incorrect 
				parameters in url Query`, func() {

				req = httptest.NewRequest("GET", "/parametrized/models?ids=incorrect,3,4", nil)
				rw = httptest.NewRecorder()

				router.ServeHTTP(rw, req)

				body, err := readBody(rw)
				So(err, ShouldBeNil)

				So(body.Status, ShouldEqual, response.StatusError)
				So(body.HttpStatus, ShouldEqual, 400)
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

				So(body.Status, ShouldEqual, response.StatusOk)
				So(body.HttpStatus, ShouldEqual, 200)
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

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.HttpStatus, ShouldEqual, 400)
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

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.HttpStatus, ShouldEqual, 500)
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

					So(body.Status, ShouldEqual, response.StatusError)
					So(body.HttpStatus, ShouldEqual, 500)
					So(body.Errors, ShouldNotBeEmpty)
					So(body.Errors[0].Compare(resterrors.ErrInternalError), ShouldBeTrue)
					So(body.Content, ShouldBeEmpty)
				})

			})
		})

	})
}

func readBody(rw *httptest.ResponseRecorder) (body *response.Body, err error) {
	rsp, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		return nil, err
	}
	body = new(response.Body)
	err = json.Unmarshal(rsp, &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

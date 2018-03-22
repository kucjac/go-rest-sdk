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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Model struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

		Convey("The policy argument may be nil, in such a case the default policy would be used.", func() {
			repo := &mocks.MockRepository{}
			errHandler := errhandler.New()
			var policy *forms.FormPolicy

			gjh1, err = New(repo, errHandler, policy)
			So(err, ShouldBeNil)

			So(gjh1.formPolicy, ShouldNotBeNil)
			So(*gjh1.formPolicy, ShouldResemble, forms.DefaultFormPolicy)
		})
	})
}

func TestCreate(t *testing.T) {
	Convey("Subject: Create gin.Handlerfunc", t, func() {
		var err error
		var gjh *GinJSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var policy *forms.FormPolicy
		var body *response.Body
		var req *http.Request
		var rw *httptest.ResponseRecorder
		gin.SetMode(gin.TestMode)

		Convey("Having an error handler and some gin router", func() {

			router = gin.New()

			errHandler = errhandler.New()

			repo := &mocks.MockRepository{}
			policy = &forms.FormPolicy{FailOnError: true}

			gjh, err = New(repo, errHandler, policy)
			So(err, ShouldBeNil)

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

			Convey(``, nil)
		})
	})
}

func TestGet(t *testing.T) {

	Convey("Subject: Get gin.Handlerfunc", t, func() {
		var err error
		var gjh *GinJSONHandler
		var errHandler *errhandler.ErrorHandler
		var router *gin.Engine
		var policy *forms.FormPolicy

		Convey("Having a ginJSONHandler and some gin router", func() {
			router = gin.New()
			repo := &mocks.MockRepository{}
			errHandler = errhandler.New()

			Convey("If parameter is named differently the handler would return response with http 500 error.", func() {
				type GenericModel struct {
					ID   int
					Name string
				}
				gjh, err = New(repo, errHandler, policy)
				So(err, ShouldBeNil)

				router.GET("/incorrect/:model", gjh.Get(GenericModel{}))

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

				mockrepo := &mocks.MockRepository{}

				gjh, err = New(mockrepo, errHandler, policy)
				So(err, ShouldBeNil)

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
				mockRepo := &mocks.MockRepository{}

				// Internal Error - not recognised
				mockRepo.On("Get", &Model{ID: 50}).Return(nil,
					&dberrors.Error{Title: "Should not exist",
						Message: "I am not recognised by "},
				)

				// Internal Error - recognised
				mockRepo.On("Get", &Model{ID: 51}).Return(nil, dberrors.ErrInternalError.New())

				// Client side error
				mockRepo.On("Get", &Model{ID: 52}).Return(nil, dberrors.ErrNoResult.New())

				gjh, err := New(mockRepo, errHandler, policy)
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

				mockRepo := &mocks.MockRepository{}
				mockedModel := &Model{ID: 1, Name: "First"}
				mockRepo.On("Get", &Model{ID: 1}).Return(mockedModel, nil)

				gjh, err = New(mockRepo, errHandler, policy)
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

package handlers

import (
	"encoding/json"
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

			genericHandler, err := New(repo, errHandler, nil)

			So(err, ShouldBeNil)
			So(genericHandler, ShouldNotBeNil)
		})

		Convey("If repository or error handler is nil, an error should be returned", func() {
			var repo repository.Repository

			genericHandler, err := New(repo, errHandler, nil)

			So(err, ShouldBeError)
			So(genericHandler, ShouldBeNil)
		})

		Convey(`Having some genericHandler using New() method 
			creates and returns by callback a copy of given handler`, func() {
			repo := &mockrepo.MockRepository{}
			genericHandler, err := New(repo, errHandler, nil)

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
				So(handler.queryPolicy, ShouldBeNil)

				callbacked := handler.WithQueryPolicy(&forms.DefaultPolicy)

				So(handler.queryPolicy, ShouldNotBeNil)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`By using WithJSONPolicy a policy is being saved in the handler,
				and the handler itself is being returned`, func() {
				So(handler.jsonPolicy, ShouldBeNil)

				policy := &forms.DefaultPolicy
				callbacked := handler.WithJSONPolicy(policy)

				So(handler.jsonPolicy, ShouldEqual, policy)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`By using WithListPolicy a listPolicy is being saved in the handler,
				and the handler itself is being returned by callback`, func() {

				So(handler.listPolicy, ShouldBeNil)

				listPolicy := &forms.DefaultListPolicy

				callbacked := handler.WithListPolicy(listPolicy)

				So(handler.listPolicy, ShouldEqual, listPolicy)
				So(callbacked, ShouldPointTo, handler)
			})

			Convey(`By using WithSetIDFunc a setIDFunc is being saved in the handler,
				and the handler itself is being returned by callback`, func() {

				So(handler.idSetFunc, ShouldBeNil)

				customSetIDFunc := func(req *http.Request, model interface{}) error {
					return nil
				}
				callbacked := handler.WithSetIDFunc(customSetIDFunc)

				So(handler.idSetFunc, ShouldNotBeNil)
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
			handler, err := New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			server.Handle("/models", handler.WithJSONPolicy(&forms.Policy{FailOnError: true}).
				Create(Model{}))

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
		})
	})
}

func TestJSONMethod(t *testing.T) {
	Convey("Subject: JSON method for GenericHandler", t, func() {
		Convey("Having some GenericHandler, request, response and some mux", func() {
			handler, err := New(&mockrepo.MockRepository{}, errhandler.New(), nil)
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

func erroredSetIDFunc(req *http.Request, model interface{}) error {
	return ErrIncorrectModel
}

func monerroredSetIDFunc(req *http.Request, model interface{}) error {
	return nil
}

type MockResponser struct {
	Channel chan (int) `json:"invalid"`
}

func (m *MockResponser) AddContent(content ...interface{}) {

}

func (m *MockResponser) AddErrors(errors ...*resterrors.Error) {

}

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

package gorillamux

import (
	"github.com/gorilla/mux"
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/repository/mockrepo"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Model struct {
	ID int
}

type InvalidModel struct {
	Name string
}

func TestNew(t *testing.T) {
	Convey("Subject: New GenericHandler with Gorilla Mux IDSetFunction", t, func() {
		Convey("Having non nil repository and errHandler, should create GenericHandler", func() {
			repo := &mockrepo.MockRepository{}
			errHandler := errhandler.New()
			genericHandler, err := New(repo, errHandler, nil)
			So(err, ShouldBeNil)
			So(genericHandler, ShouldNotBeNil)
		})
		Convey("If any of repository or errhandler is nil, New should return error", func() {
			var repo repository.Repository
			errHandler := errhandler.New()

			gHandler, err := New(repo, errHandler, nil)
			So(err, ShouldBeError)
			So(gHandler, ShouldBeNil)

		})
	})
}

func TestGorillaMuxSetIDFunc(t *testing.T) {
	Convey("Subject: SetIDFunc for gorilla/mux router ", t, func() {
		Convey("Having some gorilla mux", func() {

			handleGet := func(model interface{}) http.HandlerFunc {
				return func(rw http.ResponseWriter, req *http.Request) {
					err := GorillaMuxSetIDFunc(req, model)
					if err != nil {
						http.Error(rw, "Internal Error", http.StatusInternalServerError)
						return
					}
				}
			}

			m := mux.NewRouter()
			m.Handle("/models/{model}", handleGet(&Model{}))
			m.Handle("/invalidmodels/{invalidmodel}", handleGet(&InvalidModel{}))
			m.Handle("/invalidpath/{modelingsa}", handleGet(&Model{}))

			Convey("If everything is correct no errors should be responsed", func() {
				req := httptest.NewRequest("GET", "/models/1", nil)
				rw := httptest.NewRecorder()

				m.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, http.StatusOK)
			})

			Convey(`If a model that doesn't have possibility to set id is provided 
				an error should be responded`, func() {
				req := httptest.NewRequest("GET", "/invalidmodels/3", nil)
				rw := httptest.NewRecorder()

				m.ServeHTTP(rw, req)
				So(rw.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey(`If invalid model name was set within routing path an 
				error should be responded`, func() {
				req := httptest.NewRequest("GET", "/invalidpath/4", nil)
				rw := httptest.NewRecorder()

				m.ServeHTTP(rw, req)
				So(rw.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

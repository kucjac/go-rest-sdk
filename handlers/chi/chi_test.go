package chihandler

import (
	"github.com/go-chi/chi"
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

func TestCheckSetIDFunc(t *testing.T) {
	Convey("Subject go-chi SetIDFunc for GenericHandler", t, func() {
		Convey("Having chi Mux, with some routes", func() {
			m := chi.NewMux()

			handleGet := func(model interface{}) http.HandlerFunc {
				return func(rw http.ResponseWriter, req *http.Request) {
					err := ChiSetIDFunc(req, model)
					if err != nil {
						http.Error(rw, "Internal Error", http.StatusInternalServerError)
						return
					}
				}
			}

			m.Get("/models/{model}", handleGet(&Model{}))
			m.Get("/incorrectmodel/{invalidmodel}", handleGet(&InvalidModel{}))
			m.Get("/invalidpath/{modeling}", handleGet(&Model{}))
			Println(m.Routes())

			Convey(`If a correct path and correct model is provided 
				no error should be returned.`, func() {
				req := httptest.NewRequest("GET", "/models/1", nil)
				rw := httptest.NewRecorder()
				m.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, 200)
			})

			Convey(`If an incorrect model is provided (i.e. No ID Field)
			 then an internal error should be returned`, func() {
				req := httptest.NewRequest("GET", "/incorrectmodel/1", nil)
				rw := httptest.NewRecorder()
				m.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey(`If an model name is provided in the routing path,
				then an error should be responsed`, func() {
				req := httptest.NewRequest("GET", "/invalidpath/3", nil)
				rw := httptest.NewRecorder()

				m.ServeHTTP(rw, req)

				So(rw.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestNew(t *testing.T) {
	Convey("Subject: New go-chi based GenericHandler", t, func() {
		Convey(`Having some repository, errorHandler a new GenericHandler 
			shoule be created`, func() {
			repo := &mockrepo.MockRepository{}
			errHandler := errhandler.New()

			genericHandler, err := New(repo, errHandler, nil)
			So(err, ShouldBeNil)

			So(genericHandler, ShouldNotBeNil)
		})

		Convey(`If no repo or errhandler would be provided to New() function,
			then an error would be returned instead`, func() {
			var repo repository.Repository
			errHandler := errhandler.New()

			genericHandler, err := New(repo, errHandler, nil)
			So(err, ShouldBeError)
			So(genericHandler, ShouldBeNil)
		})
	})
}

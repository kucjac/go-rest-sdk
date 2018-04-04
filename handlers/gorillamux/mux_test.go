package gorillamux

import (
	"github.com/kucjac/go-rest-sdk/errhandler"
	"github.com/kucjac/go-rest-sdk/repository"
	"github.com/kucjac/go-rest-sdk/repository/mockrepo"
	. "github.com/smartystreets/goconvey/convey"
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
			genericHandler, err := New(repo, errHandler, nil, nil)
			So(err, ShouldBeNil)
			So(genericHandler, ShouldNotBeNil)
		})
		Convey("If any of repository or errhandler is nil, New should return error", func() {
			var repo repository.Repository
			errHandler := errhandler.New()

			gHandler, err := New(repo, errHandler, nil, nil)
			So(err, ShouldBeError)
			So(gHandler, ShouldBeNil)

		})
	})
}

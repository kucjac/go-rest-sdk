package chihandler

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
	Convey("Subject: New go-chi based GenericHandler", t, func() {
		Convey(`Having some repository, errorHandler a new GenericHandler 
			shoule be created`, func() {
			repo := &mockrepo.MockRepository{}
			errHandler := errhandler.New()

			genericHandler, err := New(repo, errHandler, nil, nil)
			So(err, ShouldBeNil)

			So(genericHandler, ShouldNotBeNil)
		})

		Convey(`If no repo or errhandler would be provided to New() function,
			then an error would be returned instead`, func() {
			var repo repository.Repository
			errHandler := errhandler.New()

			genericHandler, err := New(repo, errHandler, nil, nil)
			So(err, ShouldBeError)
			So(genericHandler, ShouldBeNil)
		})
	})
}

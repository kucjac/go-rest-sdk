package response

import (
	"errors"
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	type TestStruct struct {
		ID   uint
		Name string
	}
	Convey("While having some struct", t, func() {
		s := &TestStruct{ID: 1, Name: "Test"}

		Convey("The correct response would be created with it as a result", func() {
			response := New()
			response.AddContent("test content", s)

			Convey("The response result will containt that struct", func() {
				So(response.Status, ShouldEqual, StatusOk)
				So(response.HttpStatus, ShouldEqual, 200)
				So(response.Content, ShouldContainKey, "test content")
				So(response.Content["test content"], ShouldEqual, s)
			})
		})
	})
}

func TestNewWithError(t *testing.T) {
	Convey("While processing something, an error occured", t, func() {
		err := &resterrors.Error{Title: "Some error"}
		Convey("The error is of http type 400 - Bad Request", func() {
			httpStatus := http.StatusBadRequest

			Convey("Prepared response contain that error and status", func() {
				response := NewWithError(httpStatus, err)

				So(response.HttpStatus, ShouldEqual, 400)
				So(response.Errors, ShouldContain, err)
				So(response.Status, ShouldEqual, StatusError)
			})
		})
	})
}

func TestBodyAddErrors(t *testing.T) {
	Convey("Having some Response", t, func() {
		res := NewWithError(400)

		Convey("There occured some errors", func() {
			err1 := &resterrors.Error{Title: "Some error 1"}
			err2 := &resterrors.Error{Title: "Some error 2"}
			Convey("Adding them to response", func() {
				res.AddErrors(err1, err2)

				So(res.Errors, ShouldContain, err1)
				So(res.Errors, ShouldContain, err2)
			})
		})

	})
}

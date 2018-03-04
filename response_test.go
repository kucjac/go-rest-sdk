package restsdk

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestResponseWithOk(t *testing.T) {
	type TestStruct struct {
		ID   uint
		Name string
	}
	Convey("While having some struct", t, func() {
		t := &TestStruct{ID: 1, Name: "Test"}

		Convey("The correct response would be created with it as a result", func() {
			response := ResponseWithOk()
			response.AddResult("testResult", t)

			Convey("The response result will containt that struct", func() {
				So(response.Status, ShouldEqual, StatusOk)
				So(response.HttpStatus, ShouldEqual, 200)
				So(response.Result, ShouldContainKey, "testResult")
				So(response.Result["testResult"], ShouldEqual, t)
			})
		})
	})
}

func TestResponseWithError(t *testing.T) {
	Convey("While processing something, an error occured", t, func() {
		err := &ResponseError{
			ErrorCategory: ErrorCategory{
				Code:  "101",
				Title: "Test error",
			},
			ID:     "12",
			Detail: "Some detail",
			Status: "400",
		}
		Convey("The error is of http type 400 - Bad Request", func() {
			httpStatus := http.StatusBadRequest

			Convey("Prepared response contain that error and status", func() {
				response := ResponseWithError(httpStatus, err)

				So(response.HttpStatus, ShouldEqual, 400)
				So(response.Errors, ShouldContain, err)
				So(response.Status, ShouldEqual, StatusError)
			})
		})
	})
}

func TestResponseAddErrors(t *testing.T) {
	Convey("Having some Response", t, func() {
		res := ResponseWithError(400)

		Convey("There occured some errors", func() {
			err1 := errors.New("Some error 1")
			err2 := errors.New("Some error 2")
			Convey("Adding them to response", func() {
				res.AddErrors(err1, err2)

				So(res.Errors, ShouldContain, err1)
				So(res.Errors, ShouldContain, err2)
			})
		})

	})
}

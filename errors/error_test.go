package resterrors

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResponseErrorAddLink(t *testing.T) {
	Convey("Given an Error Category", t, func() {
		var testError ResponseError = ResponseError{
			ID:     "806x",
			Detail: "Test Detail",
		}

		Convey("We add link for the following urlbase: http://host.com/errors", func() {
			var urlBase string = "http://host.com/errors"
			err := testError.AddLink(urlBase)

			Convey(`Then there should be no error and 
				the link should contain urlBase with code`, func() {
				So(err, ShouldBeNil)
				So(
					testError.Links.About,
					ShouldEqual,
					fmt.Sprintf("%s/%s", urlBase, testError.Code))
			})
		})

		testError.Links = nil

		Convey("Now we add the link where the last sign is backslash '/'", func() {
			var urlBase string = "http://host.com/errors2/"
			err := testError.AddLink(urlBase)

			Convey("Then there should be no error", func() {
				So(err, ShouldBeNil)

				Convey("And the value should contain only one backslash at the end", func() {
					So(testError.Links.About,
						ShouldEqual,
						fmt.Sprintf("%s%s", urlBase, testError.Code))
				})
			})

		})

		testError.Links = nil
		Convey("While having incorrect url", func() {
			var urlBase string = "http://192.168.0.%31/"
			err := testError.AddLink(urlBase)

			Convey("There should be an error while parsing the link", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestResponseErrorErrorMethod(t *testing.T) {
	Convey("Having a ResponseError", t, func() {
		rerr := &ResponseError{
			ErrorCategory: ErrorCategory{
				Code:  "8132",
				Title: "The title",
			},
		}
		Convey("The Error method should be combination of code and title", func() {
			errValue := rerr.Error()
			So(errValue, ShouldContainSubstring, rerr.Code)
			So(errValue, ShouldContainSubstring, rerr.Title)
		})
	})
}

func TestMarshalingResponseError(t *testing.T) {
	Convey("Having a Response Error", t, func() {
		resErr := &ResponseError{
			ErrorCategory: ErrorCategory{
				Code:  "123",
				Title: "The Title",
			},
			ID:     "1231",
			Status: "400",
			Detail: "Detailed info",
		}

		Convey("While marshaling it to json", func() {
			resErrJsoned, err := json.Marshal(resErr)

			So(err, ShouldBeNil)
			Convey(`The json should not contain category object, 
				instead it should contain a combination of all error and category values`,
				func() {
					resErrJsonString := string(resErrJsoned)
					So(resErrJsonString, ShouldContainSubstring, "\"code\":\"123\"")
					So(resErrJsonString, ShouldContainSubstring, "\"title\":\"The Title\"")
					So(resErrJsonString, ShouldNotContainSubstring, "links")

					So(resErrJsonString, ShouldContainSubstring, "\"id\":\"1231\"")
					So(resErrJsonString, ShouldContainSubstring, "\"status\":\"400\"")
					So(resErrJsonString, ShouldContainSubstring, "\"detail\":\"Detailed info\"")

				})
		})
	})
}

func TestUmarshalingResponseError(t *testing.T) {
	Convey("Having a json Response Error", t, func() {
		jsonError := `{"id":"123","status":"404","detail":"Detailed info",
		"code":"12","title":"The title","links":{"about":"someurl/to/error/12"}}`
		Convey("The unmarshaling into 'ResponseError'", func() {
			var resErr *ResponseError
			err := json.Unmarshal([]byte(jsonError), &resErr)
			So(err, ShouldBeNil)

			Convey(`Should result in unmarshaling part of json into ErrorCategory
				and the rest into response error`, func() {
				So(resErr.Code, ShouldEqual, "12")
				So(resErr.Title, ShouldEqual, "The title")
				So(resErr.Links.About, ShouldEqual, "someurl/to/error/12")
				So(resErr.ID, ShouldEqual, "123")
				So(resErr.Status, ShouldEqual, "404")
				So(resErr.Detail, ShouldEqual, "Detailed info")
			})
		})
	})

	Convey("Having a json ResponseError with incorrect types", t, func() {
		jsonError := `{"id":123, "status":"403"}`
		Convey("Unmarshaling it into ResponseError", func() {
			var resErr *ResponseError
			err := json.Unmarshal([]byte(jsonError), &resErr)
			Convey("It should produce unmarshal error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestResponseErrorWithCategory(t *testing.T) {
	incorrectIDCategory := ErrorCategory{
		Code:  "2052y",
		Title: "Incorrect ID in the query URL",
	}
	Convey("While an error occurred", t, func() {
		var occuredErr = errors.New("Error that occurred")

		Convey("And we categorize it as some category", func() {
			category := incorrectIDCategory

			Convey("By using the 'ResponseErrorWithCategory' we get 'ResponseError'", func() {
				err := ResponseErrorWithCategory(occuredErr, category)

				So(err, ShouldBeError)
				So(err.Code, ShouldEqual, incorrectIDCategory.Code)
				So(err.Title, ShouldEqual, incorrectIDCategory.Title)
				So(err.Detail, ShouldEqual, occuredErr.Error())
			})
		})
	})
}

func TestErrorCategoryStringer(t *testing.T) {
	Convey("Having given ErrorCategory", t, func() {
		category := &ErrorCategory{Code: "2xxxcode", Title: "Description of the problem"}

		Convey("The String method should contain code and title", func() {
			categoryString := category.String()
			So(categoryString, ShouldContainSubstring, category.Code)
			So(categoryString, ShouldContainSubstring, category.Title)
		})
	})
}

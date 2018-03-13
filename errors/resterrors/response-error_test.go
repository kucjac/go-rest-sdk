package resterrors

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResponseErrorAddLink(t *testing.T) {
	Convey("Given a Response Error.", t, func() {
		var testError ResponseError = ResponseError{
			ID:     "806x",
			Detail: &Detail{Title: "Test Detail"},
		}

		Convey("We add link for the following urlbase: http://host.com/errors", func() {
			var urlBase string = "http://host.com/errors"
			err := testError.AddLink(urlBase)

			Convey(`Then there should be no error and 
				the link should contain urlBase with code`, func() {
				So(err, ShouldBeNil)
				So(testError.Links.About,
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

func TestResponseErrorExtendDetail(t *testing.T) {
	Convey("Having a response error with initial detail", t, func() {
		err := &ResponseError{Detail: &Detail{Title: "Detail"}}

		Convey("With usage of AddDetialInfo method the Detail would be extended by argument", func() {
			var info string = "Added info."
			err.AddDetailInfo(info)

			So(err.Detail.Title, ShouldEqual, "Detail")
			So(err.Detail.Info, ShouldContain, info)
			So(err.Detail.Title, ShouldNotEqual, "DetailExtend")
			So(err.Detail, ShouldResemble, &Detail{Title: "Detail", Info: []string{info}})
		})

		errWithEmptyDetail := &ResponseError{}

		Convey("When the detail is empty, extending it just adds the value", func() {
			So(errWithEmptyDetail.Detail, ShouldBeNil)

			errWithEmptyDetail.AddDetailInfo("Extended")

			So(errWithEmptyDetail.Detail, ShouldNotBeNil)
			So(errWithEmptyDetail.Detail.Info, ShouldContain, "Extended")
		})
	})

}

func TestResponseErrorErrorMethod(t *testing.T) {
	Convey("Having a ResponseError", t, func() {
		rerr := &ResponseError{
			Code: "CODE8132",
			ID:   "ID123",
		}
		Convey("The Error method should be combination of code and title", func() {
			errValue := rerr.Error()
			So(errValue, ShouldContainSubstring, rerr.Code)
			So(errValue, ShouldContainSubstring, rerr.ID)
		})
	})
}

func TestMarshalingResponseError(t *testing.T) {
	Convey("Having a Response Error", t, func() {
		resErr := &ResponseError{
			Code:   "123",
			Title:  "The Title",
			ID:     "1231",
			Status: "400",
			Detail: &Detail{Title: "Detailed info"},
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
					So(resErrJsonString, ShouldContainSubstring, "\"title\":\"Detailed info\"")

				})
		})
	})
}

func TestUmarshalingResponseError(t *testing.T) {
	Convey("Having a json Response Error", t, func() {
		jsonError := `{"id":"123","status":"404","detail": {"title": "Detailed info", "info": ["Info info"]},
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
				So(resErr.Detail, ShouldResemble, &Detail{Title: "Detailed info", Info: []string{"Info info"}})
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

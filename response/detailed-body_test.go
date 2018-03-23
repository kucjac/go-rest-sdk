package response

import (
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Foo struct {
	ID   uint
	Name string
}

func TestDetailedNew(t *testing.T) {

	Convey("Having some DetailedBody entity", t, func() {
		body := &DetailedBody{}
		Convey("While having some struct", func() {

			Convey("The correct response would be created with inited Content", func() {
				response := body.New().(*DetailedBody)

				Convey(`The response result will containt that struct, if no status provided or the status is not of type 'int' the HttpStaus would be set to 
					200 by default`, func() {

					So(response.Status, ShouldEqual, StatusOk)
					So(response.HttpStatus, ShouldEqual, 200)
					So(response.Content, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestDetailedNewErrored(t *testing.T) {
	Convey("Subject: NewErrored() method for *DetailedBody", t, func() {

		Convey("Having some Detailed body", func() {
			body := &DetailedBody{}

			Convey(`By using NewErrored() method new *DetailedBody would  be created.
				By default the Status field is set to StatusError,
				And HttpStatus is set to 500`, func() {

				response, ok := body.NewErrored().(*DetailedBody)
				So(ok, ShouldBeTrue)

				So(response.Status, ShouldEqual, StatusError)
				So(response.HttpStatus, ShouldEqual, 500)
			})
		})

	})
}

func TestDetailedAddErrors(t *testing.T) {
	Convey("Having some Errored DetailedBody", t, func() {
		body := &DetailedBody{Status: StatusError}

		Convey("There occured some errors", func() {
			err1 := &resterrors.Error{Title: "Error 1"}
			err2 := &resterrors.Error{Title: "Error 2"}
			Convey("Adding them to response", func() {
				body.AddErrors(err1, err2)
				So(body.Errors, ShouldContain, err1)
				So(body.Errors, ShouldContain, err2)
			})
		})

	})
}

func TestDetailedWithErrors(t *testing.T) {
	Convey("Subject: WithErrors() method for *DetailedBody", t, func() {
		Convey("Having some *DetailedBody", func() {
			body := &DetailedBody{}

			Convey(`By using WithErrors method provided 
				errors would be appended into Errors field
				and the callback of *DetailedBody would be returned`, func() {
				err1 := &resterrors.Error{Title: "Error 1"}
				err2 := &resterrors.Error{Title: "Error 2"}
				theSameBody := body.WithErrors(err1, err2)

				So(theSameBody, ShouldEqual, body)
				So(body.Errors, ShouldContain, err1)
				So(body.Errors, ShouldContain, err2)
			})
		})
	})
}

func TestDetailedAddContent(t *testing.T) {
	Convey("Having some Detailed Body", t, func() {
		body := (&DetailedBody{}).New().(*DetailedBody)

		Convey("Adding to it's Content some testing Model", func() {
			model := Foo{ID: 1, Name: "Model"}
			models := []*Foo{{ID: 2, Name: "Second"}, {ID: 2, Name: "Third"}}
			body.AddContent(model, models)

			Convey("Then the body content should contain this testing model", func() {
				So(body.Content["foo"], ShouldResemble, model)
				So(body.Content["foos"], ShouldResemble, models)
			})

			Convey(`If having custom basic type and not wanting key like 'int' for int type,
				just create wrapper struct`, func() {
				type customInt int
				custom := customInt(5)
				body.AddContent(custom)

				So(body.Content["customint"], ShouldResemble, custom)
			})
		})
	})
}

func TestDetailedWithContent(t *testing.T) {
	Convey("Subject: WithContent() method for *DetailedBody", t, func() {

		Convey("Having some *DetailedBody with inited Content field", func() {
			body := &DetailedBody{Content: make(map[string]interface{})}

			Convey(`By using WithContent method the provided 'content' models 
				would be added to the Content field with specific key names 
				and the method would return 'body' entity itself after processing`, func() {
				someModel := Foo{ID: 1, Name: "First"}
				fewModels := []Foo{{ID: 2, Name: "Second"}, {ID: 3, Name: "Third"}}

				newBody := body.WithContent(someModel, fewModels)

				So(newBody, ShouldEqual, body)
				So(body.Content["foo"], ShouldResemble, someModel)
				So(body.Content["foos"], ShouldResemble, fewModels)
			})
		})
	})
}

func TestDetailedWithStatus(t *testing.T) {
	Convey("Subject: WithStatus method for *DetailedBody", t, func() {
		Convey("Having some *DetailedBody entity", func() {
			body := &DetailedBody{HttpStatus: 1234}

			Convey("By using WithStatus method the HttpStatus may be changed", func() {

				Convey("If 'status' argument is of type 'int'", func() {
					var status int = 300
					theSameBody := body.WithStatus(status)
					So(theSameBody, ShouldEqual, body)
					So(body.HttpStatus, ShouldEqual, status)
				})
				Convey("Should not changed if the status is of different type", func() {
					theSameBody := body.WithStatus("500")

					So(theSameBody, ShouldEqual, body)
					So(body.HttpStatus, ShouldNotEqual, 500)
				})
			})
		})
	})
}

func TestDetailedImplements(t *testing.T) {
	Convey("Subject: *DetailedBody implements Responser", t, func() {
		Convey("Having some *DetailedBody, it would implement Responser interface", func() {
			detailedBody := &DetailedBody{}

			So(detailedBody, ShouldImplement, (*Responser)(nil))
		})
	})
	Convey("Subject *DetailedBody implements StatusResponser", t, func() {
		Convey("Having some *DetailedBody, it would implement StatusResponser interface", func() {
			detailedBody := &DetailedBody{}

			So(detailedBody, ShouldImplement, (*StatusResponser)(nil))
		})
	})
}

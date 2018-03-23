package response

import (
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDefaultNew(t *testing.T) {
	Convey("Subject: New method for *DefaultBody", t, func() {

		Convey("Having some non-nil *DefaultBody", func() {
			defaultBody := &DefaultBody{}

			Convey(`By using New method creates new *DefaultBody entity 
				with non empty Content field`, func() {
				newBody, ok := defaultBody.New().(*DefaultBody)
				So(ok, ShouldBeTrue)

				So(newBody, ShouldHaveSameTypeAs, defaultBody)
				// Pointer comparing
				So(newBody, ShouldNotEqual, defaultBody)
				So(newBody.Content, ShouldNotBeNil)

			})
		})
	})
}

func TestDefaultNewErrored(t *testing.T) {
	Convey("Subject: NewErrored method for *DefaultBody", t, func() {

		Convey("Having some non-nil *DefaultBody", func() {
			defaultBody := &DefaultBody{}

			Convey(`NewErrored() method creates new *DefaultBody entity, 
				status argument is not used in this implementation. 
				Using NewErrored doesn't initiate Content field`, func() {

				newErrored, ok := defaultBody.NewErrored().(*DefaultBody)
				So(ok, ShouldBeTrue)

				So(newErrored, ShouldHaveSameTypeAs, defaultBody)
				So(newErrored, ShouldNotEqual, defaultBody)
				So(newErrored.Content, ShouldBeNil)
			})
		})
	})
}

func TestDefaultAddContent(t *testing.T) {
	Convey("Subject: AddContent() method for *DefaultBody", t, func() {
		Convey("Having some defaultBody with inited Content", func() {
			body := &DefaultBody{Content: make(map[string]interface{})}

			Convey(`Using AddContent() with some 'model', 
				would result in adding the model to the Content
				with lowercased name (pluralized if slice)`, func() {

				model := Foo{ID: 1, Name: "Some Name"}
				models := []Foo{{ID: 2, Name: "Second"}, {ID: 3, Name: "Third"}}

				body.AddContent(model)
				body.AddContent(models)

				So(body.Content["foo"], ShouldResemble, model)
				So(body.Content["foos"], ShouldResemble, models)
			})
		})
	})
}

func TestDefaultWithContent(t *testing.T) {
	Convey("Subject: WithContent() method for *DefaultBody", t, func() {
		Convey("Having some *DefaultBody with inited Content field", func() {
			body := &DefaultBody{Content: make(map[string]interface{})}

			Convey(`Using WithContent would add provided models to the Content map,
				with keys as lowercased (pluralized if slice provided) struct names.`, func() {

				model := Foo{ID: 1, Name: "First"}
				models := []Foo{{ID: 2, Name: "Second"}, {ID: 3, Name: "Third"}}

				theSameBody := body.WithContent(model, models)
				So(theSameBody, ShouldEqual, body)
				So(body.Content["foo"], ShouldResemble, model)
				So(body.Content["foos"], ShouldResemble, models)
			})
		})
	})
}

func TestDefaultAddError(t *testing.T) {
	Convey("Subject: AddErrors() method for *DefaultBody", t, func() {

		Convey("Having some *DefaultBody", func() {
			body := &DefaultBody{}

			Convey(`Using AddErrors() method would append given errors 
				to the Errors field for DefaultBody`, func() {
				err1 := &resterrors.Error{Title: "Error 1"}
				err2 := &resterrors.Error{Title: "Error 2"}

				body.AddErrors(err1, err2)

				So(body.Errors, ShouldContain, err1)
				So(body.Errors, ShouldContain, err2)
			})
		})
	})
}

func TestWithErrors(t *testing.T) {
	Convey("Subject: WithErros() method for *DefaultBody", t, func() {

		Convey("Having some *DefaultBody", func() {
			body := &DefaultBody{}

			Convey(`Using WithErrors() method would append given errors
				to the Error field for *DefualtBody and then returns itself
				as a callback function`, func() {
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

func TestDefaultImplementsResponser(t *testing.T) {
	Convey("Subject: *DefaultBody implements Responser interface", t, func() {
		Convey("Having some *DefaultBody it should implement Responser interface", func() {
			body := &DefaultBody{}

			So(body, ShouldImplement, (*Responser)(nil))
		})
	})
}

package repositories

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestListParametersContainParameters(t *testing.T) {
	Convey("Having a filled ListParameters", t, func() {

		Convey("With Only Limit should be true", func() {
			listParameters := ListParameters{}
			listParameters.Limit = 5
			ok := listParameters.ContainsParameters()

			So(ok, ShouldBeTrue)
		})

		Convey("With Only Offset set should also be true", func() {
			listParameters := ListParameters{}
			listParameters.Offset = 10
			ok := listParameters.ContainsParameters()

			So(ok, ShouldBeTrue)
		})

		Convey("Or only with order parameter", func() {
			listParameters := ListParameters{}
			listParameters.Order = "ID"
			ok := listParameters.ContainsParameters()

			So(ok, ShouldBeTrue)
		})

		Convey("So does any combination of list parameters", func() {
			listParameters := ListParameters{Limit: 10, Offset: 20}
			ok := listParameters.ContainsParameters()

			So(ok, ShouldBeTrue)
		})
	})
	Convey("But Having a nonfilled List parameters returns false", t, func() {
		listParameters := ListParameters{}
		ok := listParameters.ContainsParameters()

		So(ok, ShouldBeFalse)
	})
}

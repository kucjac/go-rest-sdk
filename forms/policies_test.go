package forms

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewPolicies(t *testing.T) {
	Convey("Subject: Creating New polices", t, func() {
		Convey("Having a DefaultBindPolicy, and creating a new copy", func() {
			policy := DefaultBindPolicy.Copy()
			So(*policy, ShouldResemble, DefaultBindPolicy)
		})

		Convey("Having ParamPolicy and creating a copy of it using.Copy() method", func() {
			policy := DefaultParamPolicy.Copy()
			So(*policy, ShouldResemble, DefaultParamPolicy)
			So(policy, ShouldNotEqual, &DefaultParamPolicy)
		})
	})
}

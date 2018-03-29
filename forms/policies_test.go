package forms

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewPolicies(t *testing.T) {
	Convey("Subject: Creating New polices", t, func() {
		Convey("Having a DefaultPolicy, and creating a new copy", func() {
			policy := DefaultPolicy.New()
			So(*policy, ShouldResemble, DefaultPolicy)
		})
		Convey("Having ListPolicy and creating a copy of it using New() method", func() {
			policy := DefaultListPolicy.New()
			So(*policy, ShouldResemble, DefaultListPolicy)
		})
		Convey("Having JSONPolicy and creating a copy of it using New() method", func() {
			policy := DefaultJSONPolicy.New()
			So(*policy, ShouldResemble, DefaultJSONPolicy)
		})
		Convey("Having ParamPolicy and creating a copy of it using New() method", func() {
			policy := DefaultParamPolicy.New()
			So(*policy, ShouldResemble, DefaultParamPolicy)
			So(policy, ShouldNotEqual, &DefaultParamPolicy)
		})
	})
}

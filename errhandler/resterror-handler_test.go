package errhandler

import (
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNew(t *testing.T) {
	Convey("Creating new error handler.", t, func() {
		errorHandler := New()

		Convey("The newly created handler use defaultErrorMap by default", func() {
			So(errorHandler.dbToRest, ShouldResemble, defaultErrorMap)

		})
	})
}

func TestLoadCustomErrorMap(t *testing.T) {
	Convey("While having an Error Handler", t, func() {
		errorHandler := New()

		Convey("And a prepared custom error map with a custom resterror", func() {
			customError := &resterrors.Error{Code: "C123", Title: "Custom rest error"}

			customMap := map[dberrors.Error]*resterrors.Error{
				dberrors.ErrUnspecifiedError: customError,
			}

			Convey("For given Error some Error should be returned.", func() {
				someError := dberrors.ErrUnspecifiedError.New()
				prevRestErr, err := errorHandler.HandleError(someError)
				So(err, ShouldBeNil)
				So(prevRestErr, ShouldNotBeNil)

				FocusConvey("While loading custom error map containing given error pair", func() {
					errorHandler.LoadCustomErrorMap(customMap)

					FocusConvey("Given Error would return different Error", func() {

						afterRestErr, err := errorHandler.HandleError(someError)
						So(err, ShouldBeNil)
						So(afterRestErr, ShouldNotResemble, prevRestErr)
						So(afterRestErr, ShouldResemble, customError)
					})
				})

			})

		})
	})
}

func TestUpdateErrorMapEntry(t *testing.T) {
	Convey("Having a ErrorHandler containing default error map", t, func() {
		errorHandler := New()

		So(errorHandler.dbToRest, ShouldResemble, defaultErrorMap)

		Convey("Getting a *Error for given Error", func() {
			someErrorProto := dberrors.ErrCheckViolation
			someError := someErrorProto.New()

			prevRestErr, err := errorHandler.HandleError(someError)
			So(err, ShouldBeNil)

			Convey("While using UpdateErrorMapEntry method on the errorHandler.", func() {
				customError := &resterrors.Error{ID: "1234", Title: "My custom Error"}

				errorHandler.UpdateErrorMapEntry(someErrorProto, customError)

				FocusConvey(`Handling again given Error would result 
					with different *Error entity`, func() {
					afterRestErr, err := errorHandler.HandleError(someError)

					So(err, ShouldBeNil)
					So(afterRestErr, ShouldNotResemble, prevRestErr)
					So(afterRestErr, ShouldResemble, customError)
				})

			})
		})
	})
}

func TestHandleError(t *testing.T) {
	Convey("Having a ErrorHandler with default error map", t, func() {
		errorHandler := New()

		Convey("And a *Error based on the basic Error prototype", func() {
			someError := dberrors.ErrUniqueViolation.New()

			Convey(`Then handling given *Error would result 
				with some *Error entity`, func() {
				restError, err := errorHandler.HandleError(someError)

				So(err, ShouldBeNil)
				So(restError, ShouldHaveSameTypeAs, &resterrors.Error{})
				So(restError, ShouldNotBeNil)
			})

		})

		Convey("If the *Error is not based on basic Error prototype", func() {
			someCustomError := &dberrors.Error{ID: uint(240), Message: "Some error message"}
			Convey("Then handling this error would result with nil *Error and throwing an internal error.", func() {
				restError, err := errorHandler.HandleError(someCustomError)

				So(err, ShouldBeError)
				So(restError, ShouldBeNil)
			})
		})

		Convey(`If we set a non default error map, 
			that may not contain every Error entry as a key`, func() {
			someErrorProto := dberrors.ErrSystemError
			customErrorMap := map[dberrors.Error]*resterrors.Error{
				someErrorProto: {ID: "1921", Title: "Some Error"},
			}
			errorHandler.LoadCustomErrorMap(customErrorMap)

			Convey(`Then handling a *Error based on the basic Error prototype that is not in 
				the custom error map, would throw an internal error 
				and a nil *Error.`, func() {
				restError, err := errorHandler.HandleError(someErrorProto.New())

				So(err, ShouldBeNil)
				So(restError, ShouldHaveSameTypeAs, &resterrors.Error{})

				someDBFromProto := dberrors.ErrInvalidSyntax.New()
				otherError, err := errorHandler.HandleError(someDBFromProto)

				So(err, ShouldBeError)
				So(otherError, ShouldBeNil)

			})

		})

	})

}

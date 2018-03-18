package errors

import (
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/kucjac/go-rest-sdk/errors/resterrors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewErrorHandler(t *testing.T) {
	Convey("Creating new error handler.", t, func() {
		errorHandler := NewErrorHandler()

		Convey("The newly created handler use defaultErrorMap by default", func() {
			So(errorHandler.dbToRest, ShouldResemble, defaultErrorMap)

		})
	})
}

func TestLoadCustomErrorMap(t *testing.T) {
	Convey("While having an Error Handler", t, func() {
		errorHandler := NewErrorHandler()

		Convey("And a prepared custom error map with a custom resterror", func() {
			customError := &resterrors.RestError{Code: "C123", Title: "Custom rest error"}

			customMap := map[dberrors.DBError]*resterrors.RestError{
				dberrors.ErrUnspecifiedError: customError,
			}

			Convey("For given DBError some RestError should be returned.", func() {
				someDBError := dberrors.ErrUnspecifiedError.New()
				prevRestErr, err := errorHandler.HandleDBError(someDBError)
				So(err, ShouldBeNil)
				So(prevRestErr, ShouldNotBeNil)

				FocusConvey("While loading custom error map containing given error pair", func() {
					errorHandler.LoadCustomErrorMap(customMap)

					FocusConvey("Given DBError would return different RestError", func() {

						afterRestErr, err := errorHandler.HandleDBError(someDBError)
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
	Convey("Having a RestErrorHandler containing default error map", t, func() {
		errorHandler := NewErrorHandler()

		So(errorHandler.dbToRest, ShouldResemble, defaultErrorMap)

		Convey("Getting a *RestError for given DBError", func() {
			someDBErrorProto := dberrors.ErrCheckViolation
			someDBError := someDBErrorProto.New()

			prevRestErr, err := errorHandler.HandleDBError(someDBError)
			So(err, ShouldBeNil)

			Convey("While using UpdateErrorMapEntry method on the errorHandler.", func() {
				customRestError := &resterrors.RestError{ID: "1234", Title: "My custom RestError"}

				errorHandler.UpdateErrorMapEntry(someDBErrorProto, customRestError)

				FocusConvey(`Handling again given DBError would result 
					with different *RestError entity`, func() {
					afterRestErr, err := errorHandler.HandleDBError(someDBError)

					So(err, ShouldBeNil)
					So(afterRestErr, ShouldNotResemble, prevRestErr)
					So(afterRestErr, ShouldResemble, customRestError)
				})

			})
		})
	})
}

func TestHandleDBError(t *testing.T) {
	Convey("Having a RestErrorHandler with default error map", t, func() {
		errorHandler := NewErrorHandler()

		Convey("And a *DBError based on the basic DBError prototype", func() {
			someDBError := dberrors.ErrUniqueViolation.New()

			Convey(`Then handling given *DBError would result 
				with some *RestError entity`, func() {
				restError, err := errorHandler.HandleDBError(someDBError)

				So(err, ShouldBeNil)
				So(restError, ShouldHaveSameTypeAs, &resterrors.RestError{})
				So(restError, ShouldNotBeNil)
			})

		})

		Convey("If the *DBError is not based on basic DBError prototype", func() {
			someCustomDBError := &dberrors.DBError{ID: uint(240), Message: "Some error message"}
			Convey("Then handling this error would result with nil *RestError and throwing an internal error.", func() {
				restError, err := errorHandler.HandleDBError(someCustomDBError)

				So(err, ShouldBeError)
				So(restError, ShouldBeNil)
			})
		})

		Convey(`If we set a non default error map, 
			that may not contain every DBError entry as a key`, func() {
			someDBErrorProto := dberrors.ErrSystemError
			customErrorMap := map[dberrors.DBError]*resterrors.RestError{
				someDBErrorProto: {ID: "1921", Title: "Some Error"},
			}
			errorHandler.LoadCustomErrorMap(customErrorMap)

			Convey(`Then handling a *DBError based on the basic DBError prototype that is not in 
				the custom error map, would throw an internal error 
				and a nil *RestError.`, func() {
				restError, err := errorHandler.HandleDBError(someDBErrorProto.New())

				So(err, ShouldBeNil)
				So(restError, ShouldHaveSameTypeAs, &resterrors.RestError{})

				someDBFromProto := dberrors.ErrInvalidSyntax.New()
				otherRestError, err := errorHandler.HandleDBError(someDBFromProto)

				So(err, ShouldBeError)
				So(otherRestError, ShouldBeNil)

			})

		})

	})

}

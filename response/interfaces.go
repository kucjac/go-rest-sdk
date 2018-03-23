package response

import (
	"github.com/kucjac/go-rest-sdk/resterrors"
)

// Responser is an interface that is used for response bodies.
// It implements both ContentAdder and ErrorAdder as well as defines
// two additional methods:
// - New() that creates new Responser
// - NewErrored(status interface{}) - creates new Responser that is specified
//		for errored response. I.e. implementation differs when an error occured
//		In addition status is provided as an argument, where some implementations might used it.
type Responser interface {
	// AddContent and WithContent adds the content to the Responser implementation
	// Both methods should do it with the same rules, but WithContent() acts like
	// callback function that after processing returns itself
	AddContent(content ...interface{})
	WithContent(content ...interface{}) Responser

	// AddErrors and WithErrors adds the errors to the given Responser
	// Both method should do it with the same rules, but WithErrors() should act like a
	//callback function that returns itself after processing
	AddErrors(errors ...*resterrors.Error)
	WithErrors(errors ...*resterrors.Error) Responser

	// New() creates a new Responser entity
	// the status argument may be not used in implementations
	New() Responser

	// NewErrored() creates new Resposner that is defined for errored response.
	// This enables some implementations to differ when some error occured
	NewErrored() Responser
}

// StatusResponser is an interface that inherit Responser interface
// In addition it contains WithStatus method that sets the status for given StatusResponser
type StatusResponser interface {
	Responser
	//WithStatus is a callback function that sets the status for given Responser
	WithStatus(status interface{}) StatusResponser
}

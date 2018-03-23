package response

import (
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/resterrors"
)

// DetailedBody - basic REST API response structure
// Created on purpose of easily managable and
// It implements StatusResponser interface
type DetailedBody struct {
	// Status is an easy to check variable with only two possible values:
	// - 'ok'
	// - 'error'
	// If the status is 'ok' the response recipient can proceed to result variable
	Status Status `json:"status"`

	// HttpCode is a http status code applicable to this problem
	// While using multiple ResponseErrors with http status
	// it is a good practice to set this value as a leading Http Status
	// i.e.:
	//		- 200 - correct status
	//		- 400 - multiple client error - 4xx
	//		- 500 - multiple API server error - 5xx
	HttpStatus int `json:"httpCode,omitempty"`

	// Errors - list of errors that occurred
	// The server MAY choose to stop processing as soon as a problem is encountered, or it
	// MAY continue processing and encounter multiple problems.
	Errors []*resterrors.Error `json:"errors,omitempty"`

	// Content contains all response results
	// Composed as a map[string]interface{}
	// Every result should have it's own tag
	// i.e. "user" : User{1} - user object
	// 		"users" : []User{1,2} - list (plural)
	Content map[string]interface{} `json:"result,omitempty"`
}

// AddContent adds a content to the Detailed body Content
// The key for the Content is set as provided 'content' struct Name - lowercased
// I.e. type Model struct would use key 'model'
// But Slice of models []Model or []*Model would use pluralized name - 'models'
// For basic types like 'int' or 'string' use wrapper struct so that the name would
// be as proided i.e.: having some Limit variable of type int, by wrapping it as
// type Limit int and insert content as Limit(limitValue) would result storing
// the Limit content with key 'limit'
// Implements ContentAdder interface
func (d *DetailedBody) AddContent(content ...interface{}) {
	d.addContent(content...)
}

// WithContent adds provided content to the DetailedBody.Content field
// The rules are the same as with AddContent() method.
// In addition the method may be used as callback function returning itself
// after processing
func (d *DetailedBody) WithContent(content ...interface{}) Responser {
	d.addContent(content...)
	return d
}

// AddErrors adds errors to the Errors field within the *DetailedBody
func (d *DetailedBody) AddErrors(errors ...*resterrors.Error) {
	d.addErrors(errors...)
}

// WithErrors adds errors to the *DetailedBody.
// Acts like AddErrors() but in addition the method may be used as a callback
// that returns itself after processing
func (d *DetailedBody) WithErrors(errors ...*resterrors.Error) Responser {
	d.addErrors(errors...)
	return d
}

// New creates new response DetailedBody with positive status.
// The function initialize the DetailedBody with:
//	- Status: StatusOK
// 	- HttpStatus: 200
//	- empty inited 'Content'
func (d *DetailedBody) New() Responser {
	response := &DetailedBody{
		Status:     StatusOk,
		HttpStatus: 200,
		Content:    make(map[string]interface{}),
	}
	return response
}

// NewErrored prepares DetailedBody with Status: 'StatusError'
// By default HttpStatus is set to 500
func (d *DetailedBody) NewErrored() Responser {
	response := &DetailedBody{
		Status:     StatusError,
		HttpStatus: 500,
	}
	return response
}

// WithStatus if the provided status is of type int, the method
// sets the httpStatus with the provided in the argument
// and returns itself as a callback function
func (d *DetailedBody) WithStatus(status interface{}) StatusResponser {
	intStatus, ok := status.(int)
	if ok {
		d.HttpStatus = intStatus
	}
	return d
}

func (d *DetailedBody) addContent(contents ...interface{}) {
	for _, content := range contents {
		d.Content[refutils.ModelName(content)] = content
	}
}

func (d *DetailedBody) addErrors(errors ...*resterrors.Error) {
	d.Errors = append(d.Errors, errors...)
}

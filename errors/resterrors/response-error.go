package resterrors

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrNilPointerProvided = errors.New("Provided nil pointer")
	ErrPtrNotProvided     = errors.New("Provided arugment is raw struct type. Provide a pointer to struct")
	ErrUnknownType        = errors.New("Unknown type of the field")
)

// ResponseErrorLink is an object that contains
// link that leads to further details about this particular occurrence of the problem.
type ErrorLink struct {
	About string `json:"about"`
}

// ResponseErrorCategory is application specific error category
type ErrorCategory struct {
	// Code is an application-specific code, expressed as code
	Code string `json:"code,omitempty"`

	// Title is a short human-readable summary of the problem. SHOULD NOT change from occurrence to
	// occurrence of the problem
	Title string `json:"title,omitempty"`
}

// String implements Stringer interface
func (c *ErrorCategory) String() string {
	return fmt.Sprintf("%s: %s", c.Code, c.Title)

}

// ResponseError represents full JSON-API error. It's easier to
type ResponseError struct {
	// ErrorCategory is inherited and contains error category variables
	ErrorCategory

	// ID is a unique identifier for this particular occurence of the problem
	ID string `json:"id,omitempty"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Detail is a human-readable explanation of the problem that SHOULD describe specific
	// occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// Links contains the the link that leads to further details about this particular occurrence
	// of the problem
	Links *ErrorLink `json:"links,omitempty"`

	// Err keeps the internal error message for logging purpose
	err error
}

// ResponseErrorWithCategory prepares response error using 'category' argument.
func ResponseErrorWithCategory(err error, category ErrorCategory) *ResponseError {
	return &ResponseError{ErrorCategory: category, err: err}
}

// AddLink adds the link to the Error Category.
// Parameters:
// - urlBase - string representing the url link to the error category
// The method checks if the urlBase is a correct url and then
// appends the error category code to the urlBase
func (r *ResponseError) AddLink(urlBase string) error {
	// if the url ends with '/', trim it
	if last := len(urlBase) - 1; last >= 0 && urlBase[last] == '/' {
		urlBase = urlBase[:last]
	}
	url, err := url.Parse(urlBase)
	if err != nil {
		return err
	}
	r.Links = &ErrorLink{About: fmt.Sprintf("%s/%s", url.String(), r.Code)}
	return nil
}

func (r *ResponseError) ExtendDetail(moreInfo string) {
	if len(r.Detail) != 0 {
		last := r.Detail[len(r.Detail)-1:]

		if last == "." {
			r.Detail += " " + moreInfo
		} else if last == " " {
			r.Detail += moreInfo
		} else {
			r.Detail += ". " + moreInfo
		}
	} else {
		r.Detail = moreInfo
	}

}

// Error implements error interface
func (r *ResponseError) Error() string {
	return fmt.Sprintf("%s-%s: %s", r.Code, r.ID, r.err.Error())
}

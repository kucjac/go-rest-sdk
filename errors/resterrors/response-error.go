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

// Detail contains human readable detailed information about the specific problem
// If more specific information is available, it may be stored in the 'Info' field.
type Detail struct {
	Title string   `json:"title,omitempty"`
	Info  []string `json:"info,omitempty"`
}

// ResponseError represents full JSON-API error. It's easier to
type ResponseError struct {
	// Code is an application-specific code, expressed as code
	Code string `json:"code,omitempty"`

	// Title is a short human-readable summary of the problem. SHOULD NOT change from occurrence to
	// occurrence of the problem
	Title string `json:"title,omitempty"`

	// ID is a unique identifier for this particular occurence of the problem
	ID string `json:"id,omitempty"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Detail is a human-readable explanation of the problem that SHOULD describe specific
	// occurrence of the problem.
	Detail *Detail `json:"detail,omitempty"`

	// Links contains the the link that leads to further details about this particular occurrence
	// of the problem
	Links *ErrorLink `json:"links,omitempty"`
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

func (r *ResponseError) AddDetailInfo(moreInfo string) {
	if r.Detail == nil {
		r.Detail = &Detail{}
	}
	r.Detail.Info = append(r.Detail.Info, moreInfo)
}

// Error implements error interface
func (r *ResponseError) Error() string {
	return fmt.Sprintf("%s-%s", r.Code, r.ID)
}

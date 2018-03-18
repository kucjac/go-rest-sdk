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

// RestErrorLink is an object that contains
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

// Copy creates a copy of the Detail entity
func (d *Detail) Copy() *Detail {
	return d.copy()
}

func (d *Detail) copy() *Detail {
	var detail *Detail
	detail = &Detail{Title: d.Title}
	detail.Info = append(detail.Info, d.Info...)
	return detail
}

// RestError represents full JSON-API error. It's easier to
type RestError struct {
	// ID is a unique identifier for this particular occurence of the problem
	ID string `json:"id,omitempty"`

	// Links contains the the link that leads to further details about this particular occurrence
	// of the problem
	Links *ErrorLink `json:"links,omitempty"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Code is an application-specific code, expressed as code
	Code string `json:"code,omitempty"`

	// Title is a short human-readable summary of the problem. SHOULD NOT change from occurrence to
	// occurrence of the problem
	Title string `json:"title,omitempty"`

	// Detail is a human-readable explanation of the problem that SHOULD describe specific
	// occurrence of the problem.
	Detail *Detail `json:"detail,omitempty"`
}

// New creates new *RestError entity that is a copy of given RestError prototype.
func (r RestError) New() *RestError {
	return &RestError{Code: r.Code, Title: r.Title, Status: r.Status, Detail: r.Detail.copy()}
}

// AddLink adds the link to the Error Category.
// Parameters:
// - urlBase - string representing the url link to the error category
// The method checks if the urlBase is a correct url and then
// appends the error category code to the urlBase
func (r *RestError) AddLink(urlBase string) error {
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

// AddDetailInfo appends the provided 'infos' argument to the given RestError's Detail field.
// If the Detail field is nil the new Detail entity would be created.
func (r *RestError) AddDetailInfo(infos ...string) {
	if r.Detail == nil {
		r.Detail = &Detail{}
	}
	r.Detail.Info = append(r.Detail.Info, infos...)
}

// Compare compares the given RestError entity with an RestError prototype 'err'
// If both error and prototype has the same code the method returns 'true'.
func (r *RestError) Compare(err RestError) bool {
	if r.Code != err.Code {
		return false
	}
	return true
}

// Error implements error interface
func (r *RestError) Error() string {
	return fmt.Sprintf("%s-%s", r.Code, r.ID)
}

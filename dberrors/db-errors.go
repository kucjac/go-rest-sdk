package dberrors

import (
	"errors"
	"fmt"
)

// Converter is an interface that defines the form
// of converting third-party database errors into uniform
// dberrors.Error. The returned errors should be based on the
// prototypes provided in this package.
type Converter interface {
	Convert(err error) *Error
}

// Error is a unified Database Error.
//
// This package contain error prototypes with name starting with Err...
// On their base recogniser should create new errors.
// In order to compare the error entity with prototype use the 'Compare' method.
type Error struct {
	ID      uint
	Title   string
	Message string
}

// Compare - checks if the error is of the same type as given in the argument
//
// Error variables given in the package doesn't have details.
// Every *Error has its own Message. By comparing the error with
// Variables of type Error in the package the result will always be false
// This method allows to check if the error has the same ID as the error provided
// as an argument
func (d *Error) Compare(err Error) bool {
	if d.ID == err.ID {
		return true
	}
	return false
}

// GetPrototype returns the Error prototype on which the
// given database *Error entity was built.
func (d *Error) GetPrototype() (proto Error, err error) {
	var ok bool
	proto, ok = prototypeMap[d.ID]
	if !ok {
		return proto, errors.New("ID field not found or unrecognisable")
	}
	return proto, nil
}

// Error implements error interface
func (d *Error) Error() string {
	return fmt.Sprintf("%s: %s", d.Title, d.Message)
}

// New creates a copy of the given *Error
// Only ID and Title fields are copied to new Error
// the Message should be unique for given situation.
// Used on prototypes to create Prototype based Errors
func (d Error) New() *Error {
	return d.new()
}

// NewWithMessage creates new *Error copy of the Error with additional message.
// Used on prototypes to create new prototype based Error containing
// a situation specific message based on the given argument.
func (d Error) NewWithMessage(message string) (err *Error) {
	err = d.new()
	err.Message = message
	return
}

// NewWithError creates new Error copy based on the Error with a message.
// The message is an Error value from 'err' argument.
// Used on prototypes to create new prototype based Error containing
// a situation specific message based on provided error.
func (d Error) NewWithError(err error) (dbError *Error) {
	dbError = d.new()
	dbError.Message = err.Error()
	return
}

func (d Error) new() *Error {
	return &Error{ID: d.ID, Title: d.Title}
}

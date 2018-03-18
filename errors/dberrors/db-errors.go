package dberrors

import (
	"errors"
	"fmt"
)

// DBError is a unified Database Error.
//
// This package contain error prototypes with name starting with Err...
// On their base recogniser should create new errors.
// In order to compare the error entity with prototype use the 'Compare' method.
type DBError struct {
	ID      uint
	Title   string
	Message string
}

// Compare - checks if the error is of the same type as given in the argument
//
// DBError variables given in the package doesn't have details.
// Every *DBError has its own Message. By comparing the error with
// Variables of type DBError in the package the result will always be false
// This method allows to check if the error has the same ID as the error provided
// as an argument
func (d *DBError) Compare(err DBError) bool {
	if d.ID == err.ID {
		return true
	}
	return false
}

// GetPrototype returns the DBError prototype on which the
// 'd' *DBError entity was built.
func (d *DBError) GetPrototype() (DBError, error) {
	proto, ok := prototypeMap[d.ID]
	if !ok {
		return proto, errors.New("ID field not found or unrecognisable")
	}
	return proto, nil
}

// Error implements error interface
func (d *DBError) Error() string {
	return fmt.Sprintf("%s: %s", d.Title, d.Message)
}

// New creates new *DBError copy of the DBError
func (d DBError) New() *DBError {
	return d.new()
}

// NewWithMessage creates new *DBError copy of the DBError with additional message.
func (d DBError) NewWithMessage(message string) (dbError *DBError) {
	dbError = d.new()
	dbError.Message = message
	return
}

// NewWithError creates new DBError copy based on the DBError with a message.
// The message is an Error value from 'err' argument.
func (d DBError) NewWithError(err error) (dbError *DBError) {
	dbError = d.new()
	dbError.Message = err.Error()
	return
}

func (d DBError) new() *DBError {
	return &DBError{ID: d.ID, Title: d.Title}
}

// DBErrorConverter is an interface that converts errors into *DBError
type DBErrorConverter interface {
	Convert(err error) *DBError
}

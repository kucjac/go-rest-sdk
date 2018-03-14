package dberrors

import (
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

// Error implements error interface
func (d *DBError) Error() string {
	return fmt.Sprintf("%s: %s", d.Title, d.Message)
}

// New creates new *DBError copy of the DBError
func (d DBError) New() *DBError {
	return &DBError{ID: d.ID, Title: d.Title}
}

// NewWithMessage creates new *DBError copy of the DBError with additional message.
func (d DBError) NewWithMessage(message string) *DBError {
	return &DBError{ID: d.ID, Title: d.Title, Message: message}
}

var (

	// Warnings
	ErrWarning = DBError{ID: 1, Title: "Warning"}

	// ErrNoResult used as a replacement for ErrNoRows - for non-sql databases
	ErrNoResult = DBError{ID: 2, Title: "No Result"}

	// Connection Exception
	ErrConnExc = DBError{ID: 3, Title: "Connection exception"}

	ErrCardinalityViolation = DBError{ID: 4, Title: "Cardinality violation"}

	// Data Exception
	ErrDataException = DBError{ID: 5, Title: "Data Exception"}

	// Integrity Violation
	ErrIntegrConstViolation = DBError{ID: 6, Title: "Integrity constraint violation"}
	ErrRestrictViolation    = DBError{ID: 7, Title: "Restrict violation"}
	ErrNotNullViolation     = DBError{ID: 8, Title: "Not null violation"}
	ErrForeignKeyViolation  = DBError{ID: 9, Title: "Foreign-Key violation"}
	ErrUniqueViolation      = DBError{ID: 10, Title: "Unique violation"}
	ErrCheckViolation       = DBError{ID: 11, Title: "Check violation"}

	// Transactions
	ErrInvalidTransState = DBError{ID: 12, Title: "Invalid transaction state"}
	ErrInvalidTransTerm  = DBError{ID: 13, Title: "Invalid transaction termination"}
	ErrTransRollback     = DBError{ID: 14, Title: "Transaction Rollback"}

	// TxDone is an equivalent of sql.ErrTxDone error from sql package
	ErrTxDone = DBError{ID: 15, Title: "Transaction done"}

	// Invalid Authorization
	ErrInvalidAuthorization = DBError{ID: 16, Title: "Invalid Authorization Specification"}
	ErrInvalidPassword      = DBError{ID: 17, Title: "Invalid password"}

	// Invalid Schema Name
	ErrInvalidSchemaName = DBError{ID: 18, Title: "Invalid Schema Name"}

	// Invalid Catalog Name
	ErrInvalidCatalogName = DBError{ID: 19, Title: "Invalid Catalog Name"}

	// Syntax Error
	ErrInvalidSyntax         = DBError{ID: 20, Title: "Syntax Error"}
	ErrInsufficientPrivilege = DBError{ID: 21, Title: "Insufficient Privilege"}

	// Insufficient Resources
	ErrInsufficientResources = DBError{ID: 22, Title: "Insufficient Resources"}

	// Program Limit Exceeded
	ErrProgramLimitExceeded = DBError{ID: 23, Title: "Program Limit Exceeded"}

	// System Error
	ErrSystemError = DBError{ID: 24, Title: "System error"}

	// Internal Error
	ErrInternalError = DBError{ID: 25, Title: "Internal error"}

	// Unspecified Error - all other errors not included in this division
	ErrUnspecifiedError = DBError{ID: 26, Title: "Unspecified error"}
)

// DBErrorConverter is an interface that converts errors into *DBError
type DBErrorConverter interface {
	Convert(err error) *DBError
}

package dberrors

import (
	"fmt"
)

type DBError struct {
	ID      string
	Title   string
	Message string
}

func (d *DBError) Error() string {
	return fmt.Sprintf("%s: %s", d.Title, d.Message)
}

var (

	// Warnings
	ErrWarning = &DBError{Title: "Warning"}

	// ErrNoResult used as a replacement for ErrNoRows - for non-sql databases
	ErrNoResult = &DBError{Title: "No Result"}

	// Connection Exception
	ErrConnExc = &DBError{Title: "Connection exception"}

	// Data Exception
	ErrDataException = &DBError{Title: "Data Exception"}

	// Integrity Violation
	ErrIntegrConstViolation = &DBError{Title: "Integrity constraint violation"}
	ErrRestrictViolation    = &DBError{Title: "Restrict violation"}
	ErrNotNullViolation     = &DBError{Title: "Not null violation"}
	ErrForeignKeyViolation  = &DBError{Title: "Foreign-Key violation"}
	ErrUniqueViolation      = &DBError{Title: "Unique violation"}
	ErrCheckViolation       = &DBError{Title: "Check violation"}

	// Transactions
	ErrInvalidTransState = &DBError{Title: "Invalid transaction state"}
	ErrInvalidTransTerm  = &DBError{Title: "Invalid transaction termination"}
	ErrTransRollback     = &DBError{Title: "Transaction Rollback"}

	// TxDone is an equivalent of sql.ErrTxDone error from sql package
	ErrTxDone = &DBError{Title: "Transaction done"}

	// Invalid Authorization
	ErrInvalidAuthorization = &DBError{Title: "Invalid Authorization Specification"}
	ErrInvalidPassword      = &DBError{Title: "Invalid password"}

	// Invalid Schema Name
	ErrInvalidSchemaName = &DBError{Title: "Invalid Schema Name"}

	// Syntax Error
	ErrInvalidSyntax         = &DBError{Title: "Syntax Error"}
	ErrInsufficientPrivilege = &DBError{Title: "Insufficient Privilege"}

	// Insufficient Resources
	ErrInsufficientResources = &DBError{Title: "Insufficient Resources"}

	// Program Limit Exceeded
	ErrProgramLimitExceeded = &DBError{Title: "Program Limit Exceeded"}

	// System Error
	ErrSystemError = &DBError{Title: "System error"}

	// Internal Error
	ErrInternalError = &DBError{Title: "Internal error"}

	// Unspecified Error - all other errors not included in this division
	ErrUnspecifiedError = &DBError{Title: "Unspecified error"}
)

type DBErrorRecogniser interface {
	Recognise(err error) error
}

package dberrors

var (
	// Warnings
	ErrWarning = Error{ID: 1, Title: "Warning"}

	// ErrNoResult used as a replacement for ErrNoRows - for non-sql databases
	ErrNoResult = Error{ID: 2, Title: "No Result"}

	// Connection Exception
	ErrConnExc = Error{ID: 3, Title: "Connection exception"}

	ErrCardinalityViolation = Error{ID: 4, Title: "Cardinality violation"}

	// Data Exception
	ErrDataException = Error{ID: 5, Title: "Data Exception"}

	// Integrity Violation
	ErrIntegrConstViolation = Error{ID: 6, Title: "Integrity constraint violation"}
	ErrRestrictViolation    = Error{ID: 7, Title: "Restrict violation"}
	ErrNotNullViolation     = Error{ID: 8, Title: "Not null violation"}
	ErrForeignKeyViolation  = Error{ID: 9, Title: "Foreign-Key violation"}
	ErrUniqueViolation      = Error{ID: 10, Title: "Unique violation"}
	ErrCheckViolation       = Error{ID: 11, Title: "Check violation"}

	// Transactions
	ErrInvalidTransState = Error{ID: 12, Title: "Invalid transaction state"}
	ErrInvalidTransTerm  = Error{ID: 13, Title: "Invalid transaction termination"}
	ErrTransRollback     = Error{ID: 14, Title: "Transaction Rollback"}

	// TxDone is an equivalent of sql.ErrTxDone error from sql package
	ErrTxDone = Error{ID: 15, Title: "Transaction done"}

	// Invalid Authorization
	ErrInvalidAuthorization = Error{ID: 16, Title: "Invalid Authorization Specification"}
	ErrInvalidPassword      = Error{ID: 17, Title: "Invalid password"}

	// Invalid Schema Name
	ErrInvalidSchemaName = Error{ID: 18, Title: "Invalid Schema Name"}

	// Invalid Catalog Name
	ErrInvalidCatalogName = Error{ID: 19, Title: "Invalid Catalog Name"}

	// Syntax Error
	ErrInvalidSyntax         = Error{ID: 20, Title: "Syntax Error"}
	ErrInsufficientPrivilege = Error{ID: 21, Title: "Insufficient Privilege"}

	// Insufficient Resources
	ErrInsufficientResources = Error{ID: 22, Title: "Insufficient Resources"}

	// Program Limit Exceeded
	ErrProgramLimitExceeded = Error{ID: 23, Title: "Program Limit Exceeded"}

	// System Error
	ErrSystemError = Error{ID: 24, Title: "System error"}

	// Internal Error
	ErrInternalError = Error{ID: 25, Title: "Internal error"}

	// Unspecified Error - all other errors not included in this division
	ErrUnspecifiedError = Error{ID: 26, Title: "Unspecified error"}
)

var prototypeMap = map[uint]Error{
	uint(1):  ErrWarning,
	uint(2):  ErrNoResult,
	uint(3):  ErrConnExc,
	uint(4):  ErrCardinalityViolation,
	uint(5):  ErrDataException,
	uint(6):  ErrIntegrConstViolation,
	uint(7):  ErrRestrictViolation,
	uint(8):  ErrNotNullViolation,
	uint(9):  ErrForeignKeyViolation,
	uint(10): ErrUniqueViolation,
	uint(11): ErrCheckViolation,
	uint(12): ErrInvalidTransState,
	uint(13): ErrInvalidTransTerm,
	uint(14): ErrTransRollback,
	uint(15): ErrTxDone,
	uint(16): ErrInvalidAuthorization,
	uint(17): ErrInvalidPassword,
	uint(18): ErrInvalidSchemaName,
	uint(19): ErrInvalidCatalogName,
	uint(20): ErrInvalidSyntax,
	uint(21): ErrInsufficientPrivilege,
	uint(22): ErrInsufficientResources,
	uint(23): ErrProgramLimitExceeded,
	uint(24): ErrSystemError,
	uint(25): ErrInternalError,
	uint(26): ErrUnspecifiedError,
}

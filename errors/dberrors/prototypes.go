package dberrors

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

var prototypeMap = map[uint]DBError{
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

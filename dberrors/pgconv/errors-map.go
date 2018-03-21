package pgconv

import (
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/lib/pq"
)

var defaultPGErrorMap = map[interface{}]dberrors.Error{

	// Class 01 - Warnings
	pq.ErrorClass("01"): dberrors.ErrWarning,

	// Class 02 - No data
	pq.ErrorClass("02"):   dberrors.ErrNoResult,
	pq.ErrorCode("P0002"): dberrors.ErrNoResult,

	// Class 08 - Connection Exception
	pq.ErrorClass("08"): dberrors.ErrConnExc,

	// Class 21 - Cardinality Violation
	pq.ErrorClass("21"): dberrors.ErrCardinalityViolation,

	// Class 22 Data Exception
	pq.ErrorClass("22"): dberrors.ErrDataException,

	// Class 23 Integrity Violation errors
	pq.ErrorClass("23"):   dberrors.ErrIntegrConstViolation,
	pq.ErrorCode("23000"): dberrors.ErrIntegrConstViolation,
	pq.ErrorCode("23001"): dberrors.ErrRestrictViolation,
	pq.ErrorCode("23502"): dberrors.ErrNotNullViolation,
	pq.ErrorCode("23503"): dberrors.ErrForeignKeyViolation,
	pq.ErrorCode("23505"): dberrors.ErrUniqueViolation,
	pq.ErrorCode("23514"): dberrors.ErrCheckViolation,

	// Class 25 Invalid Transaction State
	pq.ErrorClass("25"): dberrors.ErrInvalidTransState,

	// Class 28 Invalid Authorization Specification
	pq.ErrorCode("28000"): dberrors.ErrInvalidAuthorization,
	pq.ErrorCode("28P01"): dberrors.ErrInvalidPassword,

	// Class 2D Invalid Transaction Termination
	pq.ErrorCode("2D000"): dberrors.ErrInvalidTransTerm,

	// Class 3F Invalid Schema Name
	pq.ErrorCode("3F000"): dberrors.ErrInvalidSchemaName,

	// Class 40 - Transaciton Rollback
	pq.ErrorClass("40"): dberrors.ErrTransRollback,

	// Class 42 - Invalid Syntax
	pq.ErrorClass("42"):   dberrors.ErrInvalidSyntax,
	pq.ErrorCode("42501"): dberrors.ErrInsufficientPrivilege,

	// Class 53 - Insufficient Resources
	pq.ErrorClass("53"): dberrors.ErrInsufficientResources,

	// Class 54 - Program Limit Exceeded
	pq.ErrorClass("54"): dberrors.ErrProgramLimitExceeded,

	// Class 58 - System Errors
	pq.ErrorClass("58"): dberrors.ErrSystemError,

	// Class XX - Internal Error
	pq.ErrorClass("XX"): dberrors.ErrInternalError,
}

package postgres

import (
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/lib/pq"
)

var PGRecogniser PostgresErrorRecogniser = PostgresErrorRecogniser{

	// Class 01 - Warnings
	pq.ErrorClass("01"): dbe.ErrWarning,

	// Class 02 - No data
	pq.ErrorClass("02"):   dbe.ErrNoResult,
	pq.ErrorCode("P0002"): dbe.ErrNoResult,

	// Class 08 - Connection Exception
	pq.ErrorClass("08"): dbe.ErrConnExc,

	// Class 21 - Cardinality Violation
	pq.ErrorClass("21"): dbe.ErrCardinalityViolation,

	// Class 22 Data Exception
	pq.ErrorClass("22"): dbe.ErrDataException,

	// Class 23 Integrity Violation errors
	pq.ErrorClass("23"):   dbe.ErrIntegrConstViolation,
	pq.ErrorCode("23000"): dbe.ErrIntegrConstViolation,
	pq.ErrorCode("23001"): dbe.ErrRestrictViolation,
	pq.ErrorCode("23502"): dbe.ErrNotNullViolation,
	pq.ErrorCode("23503"): dbe.ErrForeignKeyViolation,
	pq.ErrorCode("23505"): dbe.ErrUniqueViolation,
	pq.ErrorCode("23514"): dbe.ErrCheckViolation,

	// Class 25 Invalid Transaction State
	pq.ErrorClass("25"): dbe.ErrInvalidTransState,

	// Class 28 Invalid Authorization Specification
	pq.ErrorCode("28000"): dbe.ErrInvalidAuthorization,
	pq.ErrorCode("28P01"): dbe.ErrInvalidPassword,

	// Class 2D Invalid Transaction Termination
	pq.ErrorCode("2D000"): dbe.ErrInvalidTransTerm,

	// Class 3F Invalid Schema Name
	pq.ErrorCode("3F000"): dbe.ErrInvalidSchemaName,

	// Class 40 - Transaciton Rollback
	pq.ErrorClass("40"): dbe.ErrTransRollback,

	// Class 42 - Invalid Syntax
	pq.ErrorClass("42"):   dbe.ErrInvalidSyntax,
	pq.ErrorCode("42501"): dbe.ErrInsufficientPrivilege,

	// Class 53 - Insufficient Resources
	pq.ErrorClass("53"): dbe.ErrInsufficientResources,

	// Class 54 - Program Limit Exceeded
	pq.ErrorClass("54"): dbe.ErrProgramLimitExceeded,

	// Class 58 - System Errors
	pq.ErrorClass("58"): dbe.ErrSystemError,

	// Class XX - Internal Error
	pq.ErrorClass("XX"): dbe.ErrInternalError,
}

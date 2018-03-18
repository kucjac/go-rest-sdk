package sqlite

import (
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/mattn/go-sqlite3"
)

var defaultSQLiteErrorMap map[interface{}]dbe.DBError = map[interface{}]dbe.DBError{
	sqlite3.ErrWarning: dbe.ErrWarning,

	sqlite3.ErrNotFound: dbe.ErrNoResult,

	sqlite3.ErrCantOpen: dbe.ErrConnExc,
	sqlite3.ErrNotADB:   dbe.ErrConnExc,

	sqlite3.ErrMismatch: dbe.ErrDataException,

	sqlite3.ErrConstraint:           dbe.ErrIntegrConstViolation,
	sqlite3.ErrConstraintCheck:      dbe.ErrCheckViolation,
	sqlite3.ErrConstraintForeignKey: dbe.ErrForeignKeyViolation,
	sqlite3.ErrConstraintUnique:     dbe.ErrUniqueViolation,
	sqlite3.ErrConstraintNotNull:    dbe.ErrNotNullViolation,

	sqlite3.ErrProtocol: dbe.ErrInvalidTransState,

	sqlite3.ErrRange: dbe.ErrInvalidSyntax,
	sqlite3.ErrError: dbe.ErrInvalidSyntax,

	sqlite3.ErrAuth: dbe.ErrInvalidAuthorization,

	sqlite3.ErrPerm: dbe.ErrInsufficientPrivilege,

	sqlite3.ErrFull: dbe.ErrInsufficientResources,

	sqlite3.ErrTooBig: dbe.ErrProgramLimitExceeded,

	sqlite3.ErrNoLFS: dbe.ErrSystemError,

	sqlite3.ErrInternal: dbe.ErrInternalError,
}

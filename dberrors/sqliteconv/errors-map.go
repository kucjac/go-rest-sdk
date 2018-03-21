package sqliteconv

import (
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/mattn/go-sqlite3"
)

var defaultSQLiteErrorMap map[interface{}]dberrors.Error = map[interface{}]dberrors.Error{
	sqlite3.ErrWarning: dberrors.ErrWarning,

	sqlite3.ErrNotFound: dberrors.ErrNoResult,

	sqlite3.ErrCantOpen: dberrors.ErrConnExc,
	sqlite3.ErrNotADB:   dberrors.ErrConnExc,

	sqlite3.ErrMismatch: dberrors.ErrDataException,

	sqlite3.ErrConstraint:           dberrors.ErrIntegrConstViolation,
	sqlite3.ErrConstraintCheck:      dberrors.ErrCheckViolation,
	sqlite3.ErrConstraintForeignKey: dberrors.ErrForeignKeyViolation,
	sqlite3.ErrConstraintUnique:     dberrors.ErrUniqueViolation,
	sqlite3.ErrConstraintNotNull:    dberrors.ErrNotNullViolation,
	sqlite3.ErrConstraintPrimaryKey: dberrors.ErrUniqueViolation,

	sqlite3.ErrProtocol: dberrors.ErrInvalidTransState,

	sqlite3.ErrRange: dberrors.ErrInvalidSyntax,
	sqlite3.ErrError: dberrors.ErrInvalidSyntax,

	sqlite3.ErrAuth: dberrors.ErrInvalidAuthorization,

	sqlite3.ErrPerm: dberrors.ErrInsufficientPrivilege,

	sqlite3.ErrFull: dberrors.ErrInsufficientResources,

	sqlite3.ErrTooBig: dberrors.ErrProgramLimitExceeded,

	sqlite3.ErrNoLFS: dberrors.ErrSystemError,

	sqlite3.ErrInternal: dberrors.ErrInternalError,
}

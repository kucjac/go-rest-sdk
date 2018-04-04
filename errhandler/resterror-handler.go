package errhandler

import (
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/resterrors"
	"sync"
)

// DefaultErrorMap contain default mapping of dberrors.Error prototype into
// resterrors.Error. It is used by default by 'ErrorHandler' if created using New() function.
var DefaultErrorMap map[dberrors.Error]resterrors.Error = map[dberrors.Error]resterrors.Error{
	dberrors.ErrNoResult:              resterrors.ErrResourceNotFound,
	dberrors.ErrConnExc:               resterrors.ErrInternalError,
	dberrors.ErrCardinalityViolation:  resterrors.ErrInternalError,
	dberrors.ErrDataException:         resterrors.ErrInvalidInput,
	dberrors.ErrIntegrConstViolation:  resterrors.ErrInvalidInput,
	dberrors.ErrRestrictViolation:     resterrors.ErrInvalidInput,
	dberrors.ErrNotNullViolation:      resterrors.ErrInvalidInput,
	dberrors.ErrForeignKeyViolation:   resterrors.ErrInvalidInput,
	dberrors.ErrUniqueViolation:       resterrors.ErrResourceAlreadyExists,
	dberrors.ErrCheckViolation:        resterrors.ErrInvalidInput,
	dberrors.ErrInvalidTransState:     resterrors.ErrInternalError,
	dberrors.ErrInvalidTransTerm:      resterrors.ErrInternalError,
	dberrors.ErrTransRollback:         resterrors.ErrInternalError,
	dberrors.ErrTxDone:                resterrors.ErrInternalError,
	dberrors.ErrInvalidAuthorization:  resterrors.ErrInsufficientAccPerm,
	dberrors.ErrInvalidPassword:       resterrors.ErrInternalError,
	dberrors.ErrInvalidSchemaName:     resterrors.ErrInternalError,
	dberrors.ErrInvalidSyntax:         resterrors.ErrInternalError,
	dberrors.ErrInsufficientPrivilege: resterrors.ErrInsufficientAccPerm,
	dberrors.ErrInsufficientResources: resterrors.ErrInternalError,
	dberrors.ErrProgramLimitExceeded:  resterrors.ErrInternalError,
	dberrors.ErrSystemError:           resterrors.ErrInternalError,
	dberrors.ErrInternalError:         resterrors.ErrInternalError,
	dberrors.ErrUnspecifiedError:      resterrors.ErrInternalError,
}

// ErrorHandler defines the database dberrors.Error one-to-one mapping
// into resterrors.Error. The default error mapping is defined
// in package variable 'DefaultErrorMap'.
//
type ErrorHandler struct {
	dbToRest map[dberrors.Error]resterrors.Error
	sync.RWMutex
}

// NewErrorHandler creates new error handler with already inited ErrorMap
func New() *ErrorHandler {
	return &ErrorHandler{dbToRest: DefaultErrorMap}
}

// Handle enables dberrors.Error handling so that proper resterrors.Error is returned.
// It returns resterror.Error if given database error exists in the private error mapping.
// If provided dberror doesn't have prototype or no mapping exists for given dberrors.Error an
// application 'error' would be returned.
// Thread safety by using RWMutex.RLock
func (r *ErrorHandler) Handle(dberr *dberrors.Error) (*resterrors.Error, error) {
	// Get the prototype for given dberr
	dbProto, err := dberr.GetPrototype()
	if err != nil {
		return nil, err
	}

	// Get Rest
	r.RLock()
	restProto, ok := r.dbToRest[dbProto]
	r.RUnlock()
	if !ok {
		err = errors.New("Given database error is unrecognised by the handler")
		return nil, err
	}

	// // Create new entity
	resterr := restProto.New()
	return resterr, nil
}

// LoadCustomErrorMap enables replacement of the ErrorHandler default error map.
// This operation is thread safe - with RWMutex.Lock
func (r *ErrorHandler) LoadCustomErrorMap(errorMap map[dberrors.Error]resterrors.Error) {
	r.Lock()
	r.dbToRest = errorMap
	r.Unlock()
}

// UpdateErrorMapEntry changes single entry in the Error Handler error map.
// This operation is thread safe - with RWMutex.Lock
func (r *ErrorHandler) UpdateErrorEntry(
	dberr dberrors.Error,
	resterr resterrors.Error,
) {
	r.Lock()
	r.dbToRest[dberr] = resterr
	r.Unlock()
}

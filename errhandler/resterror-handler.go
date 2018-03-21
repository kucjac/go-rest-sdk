package errhandler

import (
	"errors"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/resterrors"
)

var defaultErrorMap map[dberrors.Error]*resterrors.Error = map[dberrors.Error]*resterrors.Error{
	dberrors.ErrWarning:               nil,
	dberrors.ErrNoResult:              resterrors.ErrResourceNotFound.New(),
	dberrors.ErrConnExc:               resterrors.ErrInternalError.New(),
	dberrors.ErrCardinalityViolation:  resterrors.ErrInternalError.New(),
	dberrors.ErrDataException:         resterrors.ErrInvalidInput.New(),
	dberrors.ErrIntegrConstViolation:  resterrors.ErrInvalidInput.New(),
	dberrors.ErrRestrictViolation:     resterrors.ErrInvalidInput.New(),
	dberrors.ErrNotNullViolation:      resterrors.ErrInvalidInput.New(),
	dberrors.ErrForeignKeyViolation:   resterrors.ErrInvalidInput.New(),
	dberrors.ErrUniqueViolation:       resterrors.ErrInvalidInput.New(),
	dberrors.ErrCheckViolation:        resterrors.ErrInvalidInput.New(),
	dberrors.ErrInvalidTransState:     resterrors.ErrInternalError.New(),
	dberrors.ErrInvalidTransTerm:      resterrors.ErrInternalError.New(),
	dberrors.ErrTransRollback:         resterrors.ErrInternalError.New(),
	dberrors.ErrTxDone:                resterrors.ErrInternalError.New(),
	dberrors.ErrInvalidAuthorization:  resterrors.ErrInsufficientAccPerm.New(),
	dberrors.ErrInvalidPassword:       resterrors.ErrInternalError.New(),
	dberrors.ErrInvalidSchemaName:     resterrors.ErrInternalError.New(),
	dberrors.ErrInvalidSyntax:         resterrors.ErrInternalError.New(),
	dberrors.ErrInsufficientPrivilege: resterrors.ErrInsufficientAccPerm.New(),
	dberrors.ErrInsufficientResources: resterrors.ErrInternalError.New(),
	dberrors.ErrProgramLimitExceeded:  resterrors.ErrInternalError.New(),
	dberrors.ErrSystemError:           resterrors.ErrInternalError.New(),
	dberrors.ErrInternalError:         resterrors.ErrInternalError.New(),
	dberrors.ErrUnspecifiedError:      resterrors.ErrInternalError.New(),
}

// ErrorHandler is a handler that
type ErrorHandler struct {
	dbToRest map[dberrors.Error]*resterrors.Error
}

// NewErrorHandler
func New() *ErrorHandler {
	return &ErrorHandler{dbToRest: defaultErrorMap}
}

// Handle
func (r *ErrorHandler) Handle(dberr *dberrors.Error,
) (resterr *resterrors.Error, err error) {
	var proto dberrors.Error
	var ok bool

	// Get the prototype for given dberr
	proto, err = dberr.GetPrototype()
	if err != nil {
		return nil, err
	}

	// Get Rest
	resterr, ok = r.dbToRest[proto]
	if !ok {
		err = errors.New("Given database error is unrecognised by the handler")
		return nil, err
	}

	return resterr, nil
}

// LoadCustomErrorMap
func (r *ErrorHandler) LoadCustomErrorMap(errorMap map[dberrors.Error]*resterrors.Error,
) {
	r.dbToRest = errorMap
}

// UpdateErrorMapEntry
func (r *ErrorHandler) UpdateErrorEntry(dberr dberrors.Error,
	resterr *resterrors.Error) {
	r.dbToRest[dberr] = resterr
}

# errhandler
Package errhandler handles dberrors.Error conversion into proper resterrors.Error.

By providing `ErrorHandler` struct the package allows conversion provided database `*dberrors.Error`
into proper `*resterror.Error`. This automates the database error handling process.

```go 
// Database error handler
type ErrorHandler struct {
	// contains error map of type map[dberrors.Error]*resterrors.Error	
}

// Handle handles provided *dberrors.Error into mapped *resterrors.Error
func (r *ErrorHandler) Handle(*dberrors.Error) *resterrors.Error {
	...
	return properRestError
}
```

In order to create new error handler for custom application use `New()` function

```go
...
// Having some dberror
var dbError *dberrors.Error
dbError = dberrors.ErrInternalError.New()

// If there is need of handling dberrors into resterrors
// create and use new ErrorHandler with New() function.
// This initialize the ErrorHandler by loading 
// default error map - 'DefaultErrorMap' from errhandler package
customErrorHandler := errhandler.New()

// if no error occured properRestError (of type *resterrors.Error) should be returned 
properRestError, err := customErrorHandler.Handle(dbError)
if err != nil {
	// some kind of internal error may return i.e. resterrors.ErrInternalError.New()
	properRestError = resterrors.ErrInternalError.New()
}
// use given rest error
...
```

If `DefaultErrorMap` doesn't satisfy custom application needs, whole error mapping might get replaced
by using `LoadCustomErrorMap`.
```go
// LoadCustomErrorMap replaces the default error mapping with the one provided in argument
func(r *ErrorHandler) LoadCustomErrorMap(errorMap map[dberrors.Error]resterrors.Error) {
	// set default error map to provided 'errorMap'
}

var dbError *dberrors.Error = dberrors.ErrCheckViolation.New()

// restErr should be based on resterrors.ErrInvalidInput prototype
restErr, err := customErrorHandler.Handle(dbError)
if err != nil {
	...
}

// Define customErrorMap
var customErrorMap map[dberrors.Error]resterrors.Error = map[dberrors.Error]resterrors.Error{
	...
	dberrors.ErrCheckViolation: resterrors.ErrOutOfRangeInput,
	...
}

// Let's replace error mapping
customErrorHandler.LoadCustomErrorMap(customErrorMap)

// After loading error map, resterror.Error of the type defined in the 'customErrorMap'
// should be returned
restErrAfterReplacemnt, err := customErrorHandler.Handle(dbError)
if err != nil {
	...
}

// ok value should be true
ok := restErrAfterReplacement.Compare(resterrors.ErrOutOfRangeInput)
```

If there is no need to replace every entry int the map use `UpdateErrorEntry` method.

```go
// UpdateErrorEntry replaces default error mapping for given dberror.Error, resterrors.Error prototypes
func(r *ErrorHandler) UpdateErrorEntry(dbError dberrors.Error, restErr resterrors.Error){
	// replce given 'dbError' in the mapped value to 'restErr'
}

// Having some dberror
var dbError *dberrors.Error = dberrors.ErrNotNullViolation.New()

// ErrorHandler should return by default ErrInvalidInput resterrors.Error
restError, err := customErrorHandler.Handle(dbError)
if err != nil {
	...
}

// ok boolean for comparing resterrors
var ok bool 

// The value of ok should be true
ok = restError.Compare(resterrors.ErrInvalidInput)

// Provided that the application needs different error mapping for that entry
// UpdateErrorEntry method replaces the default mapping for given arguments.
customErrorHandler.UpdateErrorEntry(dberrors.ErrNotNullViolation, resterror.ErrInternalError)

// Now handling the provided dbError should return new ErrInternalError entity
restErrorAfterUpdate, err := customErrorHandler.Handle(dbError)
if err != nil{
	...
}


// Now the value of the 'ok' bool should be false
ok = restErrorAfterUpdate.Compare(resterrors.ErrInvalidInput) 


// by comparing to ErrInternalError as in the UpdateErrorEntry method
// the returned value should be true
ok = restErrorAfterUpdate.Compare(resterrors.ErrInternalError) 
```


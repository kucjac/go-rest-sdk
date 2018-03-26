# dberrors
Package **dberrors** contains unified database Errors converted from third-party package errors.

In order to create the RESTful API that is indepenedent of the database type, the database
errors must be converted into single form.
This package defines database Errors with some Prototypes that are the most common error categories.

```go 
// Unified database error
type Error struct {
	ID 	uint
	Title 	string
	Message string
}

// database error prototypes
var (
	...
	ErrIntegrityConstViolation = Error{
		ID: 	6, 
		Title: 	"Integrity constraint violation",
	}
	...
)
```

In order to maintaing uniform form of the error converting, every database driver should 
implement the 'Converter' interface with it's custom error converter.

```go
// Converter interface
type Converter interface {
	Convert(err error) *dberrors.Error
}

...

type CustomConverter struct {
	// Contains some kind of error mapper 
	// from custom errors into dberrors.ErrPrototype
}

func (c *CustomConverter) Convert(err error) *dberrors.Error {
	// get dberrors.ErrPrototype from the mapping and
	// create new *dberrors.Error on it's base
	customErr := err.(CustomErrorType)
	proto, ok := mapping[customErr]
	if !ok {
		return dberrors.ErrUnspecifiedError.New()
	}
	return proto.New()
}
...
```

Good examples on how to write error converters are the converters defined in the __*mysqlconv*__, __*pgconv*__, __*sqliteconv*__ or __*gormconv*__ packages.

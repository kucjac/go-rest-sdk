package gormconv

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/dberrors/mysqlconv"
	"github.com/kucjac/go-rest-sdk/dberrors/pgconv"
	"github.com/kucjac/go-rest-sdk/dberrors/sqliteconv"
)

// GORMConverter defines error converter for multiple databases drivers
// used by the 'gorm' package.
// Implements 'Converter' interface.
type GORMConverter struct {
	converter dberrors.Converter
}

// New creates new *GORMConverter.
// On the base of the *gorm.DB argument it recognise given gorm.Dialect
// on the base of the dialect the function recognise the appropiate error converter.
// returns error if the nil pointer provided or unsupported db.Dialect
func New(db *gorm.DB) (conv *GORMConverter, err error) {
	conv = &GORMConverter{}
	err = conv.initialize(db)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// Convert implements dberrors.Converter
// Converts provided argument error into *dberrors.Error type
func (g *GORMConverter) Convert(err error) *dberrors.Error {
	switch err {
	case gorm.ErrCantStartTransaction, gorm.ErrInvalidTransaction:
		return dberrors.ErrInvalidTransState.NewWithError(err)
	case gorm.ErrInvalidSQL:
		return dberrors.ErrInvalidSyntax.NewWithError(err)
	case gorm.ErrUnaddressable:
		return dberrors.ErrUnspecifiedError.NewWithError(err)
	case gorm.ErrRecordNotFound:
		return dberrors.ErrNoResult.NewWithError(err)
	}
	// If error is not of gorm type
	// use db recogniser
	return g.converter.Convert(err)
}

// initialize provides initialization process of the GORMConverter
// returns error if nil pointer provided or unsupported database dialect.
func (g *GORMConverter) initialize(db *gorm.DB) error {
	if db == nil {
		return errors.New("Nil pointer provided")
	}
	dialect := db.Dialect()
	switch dialect.GetName() {
	case "postgres":
		g.converter = pgconv.New()
	case "mysql":
		g.converter = mysqlconv.New()
	case "sqlite3":
		g.converter = sqliteconv.New()
	default:
		return errors.New("Unsupported database dialect.")
	}
	return nil
}

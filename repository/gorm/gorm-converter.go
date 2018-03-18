package gorm

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/mysql"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/postgres"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/sqlite"
)

type GORMErrorConverter struct {
	converter dberrors.DBErrorConverter
}

func (g *GORMErrorConverter) Init(db *gorm.DB) error {
	dialect := db.Dialect()
	switch dialect.GetName() {
	case "postgres":
		g.converter = postgres.New()
	case "mysql":
		g.converter = mysql.New()
	case "sqlite3":
		g.converter = sqlite.New()
	default:
		return errors.New("Unsupported database dialect.")
	}
	return nil
}

func (g *GORMErrorConverter) Convert(err error) *dberrors.DBError {
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

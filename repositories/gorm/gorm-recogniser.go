package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/mysql"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/postgres"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/sqlite"
)

type GORMErrorRecogniser struct {
	dbRecogniser dberrors.DBErrorRecogniser
}

func (g *GORMErrorRecogniser) Init(db *gorm.DB) error {
	dialect := g.db.Dialect()
	switch dialect.GetName() {
	case "postgres":
		g.dbRecogniser = postgres.PGRecogniser
	case "mysql":
		g.dbRecogniser = mysql.MySQLRecogniser
	case "sqlite3":
		g.dbRecogniser = sqlite.SQLiteRecogniser
	default:
		return errors.New("Unsupported database dialect.")
	}
	return nil
}

func (g *GORMErrorRecogniser) Recognise(err error) error {
	switch err {
	case gorm.ErrCantStartTransaction, gorm.ErrInvalidTransaction:
		return dberrors.ErrInvalidTransState
	case gorm.ErrInvalidSQL:
		return dberrors.ErrInvalidSyntax
	case gorm.ErrUnaddressable:
		return err
	case gorm.ErrRecordNotFound:
		return dberrors.ErrNoResult
	}
	// If error is not of gorm type
	// use db recogniser
	return g.dbRecogniser.Recognise(err)
}

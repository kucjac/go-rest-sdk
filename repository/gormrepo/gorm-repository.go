package gormrepo

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/dberrors/gormconv"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"reflect"
)

// GORMRepository is an implementation of Repository interface for 'jinzhu/gorm' package.
// It is composed of *gorm.DB and dberrors.Converter.
type GORMRepository struct {
	db        *gorm.DB
	converter *gormconv.GORMConverter
}

func New(db *gorm.DB) (*GORMRepository, error) {
	gormRepo := &GORMRepository{}
	err := gormRepo.initialize(db)
	if err != nil {
		return nil, err
	}
	return gormRepo, nil

}

func (g *GORMRepository) initialize(db *gorm.DB) (err error) {
	if db == nil {
		err = errors.New("Nil pointer as an argument provided.")
		return
	}
	g.db = db

	// Get Error converter
	g.converter, err = gormconv.New(db)
	if err != nil {
		return err
	}
	return nil
}

func (g *GORMRepository) Create(req interface{}) *dberrors.Error {
	if err := g.db.Create(req).Error; err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Get(req interface{}) (res interface{}, dberr *dberrors.Error) {
	res = refutils.ObjOfPtrType(req)
	if err := g.db.First(res, req).Error; err != nil {
		return nil, g.converter.Convert(err)
	}
	return res, nil
}

func (g *GORMRepository) List(
	req interface{},
) (res interface{}, dberr *dberrors.Error) {
	// Get Slice of pointer type 'req'
	res = refutils.PtrSliceOfPtrType(req)

	// List objects provided with arguments probided in request
	if err := g.db.Find(res, req).Error; err != nil {
		return nil, g.converter.Convert(err)
	}

	return reflect.ValueOf(res).Elem().Interface(), nil
}

func (g *GORMRepository) ListWithParams(
	req interface{}, params *repository.ListParameters,
) (res interface{}, dberr *dberrors.Error) {
	if params == nil {
		return g.List(req)
	}

	if params.Offset == 0 && params.Limit == 0 && params.Order == "" && len(params.IDs) == 0 {
		return g.List(req)
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	// Get Slice of pointer type 'req'
	res = refutils.PtrSliceOfPtrType(req)

	var err error
	if len(params.IDs) > 0 {
		err = g.db.
			Offset(params.Offset).
			Limit(params.Limit).
			Order(params.Order).
			Where(params.IDs).
			Find(res, req).
			Error
	} else {
		err = g.db.
			Offset(params.Offset).
			Limit(params.Limit).
			Order(params.Order).
			Find(res, req).
			Error
	}
	if err != nil {
		dberr = g.converter.Convert(err)
		return nil, dberr
	}
	return reflect.ValueOf(res).Elem().Interface(), nil
}

func (g *GORMRepository) Count(req interface{}) (count int, dberr *dberrors.Error) {
	err := g.db.Model(req).Count(&count).Error
	if err != nil {
		return count, g.converter.Convert(err)
	}
	return count, nil
}

func (g *GORMRepository) Update(req interface{}) (dberr *dberrors.Error) {
	err := g.db.Save(req).Error
	if err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Patch(req, where interface{}) (dberr *dberrors.Error) {
	var db *gorm.DB
	db = g.db.Model(req).Where(where).Update(req)
	err := db.Error
	rows := db.RowsAffected
	if rows == 0 && err == nil {
		err = dberrors.ErrNoResult.NewWithMessage("No rows affected")
	}
	if err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Delete(req, where interface{}) *dberrors.Error {
	var db *gorm.DB
	db = g.db.Where(where).Delete(req)
	err := db.Error
	rows := db.RowsAffected
	if rows == 0 && err == nil {
		err = dberrors.ErrNoResult.NewWithMessage("No rows affected")
	}
	if err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

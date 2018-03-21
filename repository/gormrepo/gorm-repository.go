package gormrepo

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/dberrors/gormconv"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/repository"
	"reflect"
)

type GORMRepository struct {
	db        *gorm.DB
	converter dberrors.Converter
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

func (g *GORMRepository) Create(req interface{}) (err error) {
	if err = g.db.Create(req).Error; err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Get(req interface{}) (res interface{}, err error) {
	res = forms.ObjOfPtrType(req)
	if err = g.db.First(res, req).Error; err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}
	return res, nil
}

func (g *GORMRepository) List(
	req interface{},
) (res interface{}, err error) {
	// Get Slice of pointer type 'req'
	res = forms.PtrSliceOfPtrType(req)

	// List objects provided with arguments probided in request
	if err = g.db.Find(res, req).Error; err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}

	return reflect.ValueOf(res).Elem().Interface(), nil
}

func (g *GORMRepository) ListWithParams(
	req interface{}, params *repository.ListParameters,
) (res interface{}, err error) {
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
	res = forms.PtrSliceOfPtrType(req)

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
		err = g.converter.Convert(err)
		return nil, err
	}
	return reflect.ValueOf(res).Elem().Interface(), nil
}

func (g *GORMRepository) Update(req interface{}) (err error) {
	err = g.db.Save(req).Error
	if err != nil {
		err = g.converter.Convert(err)
		return err
	}
	return nil
}

func (g *GORMRepository) Patch(req, where interface{}) (err error) {
	model := forms.ObjOfPtrType(req)
	var db *gorm.DB
	db = g.db.Model(model).Where(where).Update(req)
	err = db.Error
	rows := db.RowsAffected
	if rows == 0 && err == nil {
		err = dberrors.ErrNoResult.NewWithMessage("No rows affected")
	}
	if err != nil {
		err = g.converter.Convert(err)
		return err
	}
	return nil
}

func (g *GORMRepository) Delete(req, where interface{}) (err error) {
	var db *gorm.DB
	db = g.db.Where(where).Delete(req)
	err = db.Error
	rows := db.RowsAffected
	if rows == 0 && err == nil {
		err = dberrors.ErrNoResult.NewWithMessage("No rows affected")
	}
	if err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

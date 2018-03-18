package gorm

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/postgres"
	"github.com/kucjac/go-rest-sdk/errors/dberrors/sqlite"
	"github.com/kucjac/go-rest-sdk/repositories"
	"reflect"
)

type GORMRepository struct {
	db        *gorm.DB
	converter dberrors.DBErrorConverter
}

// Initialize the gorm repository
func (g *GORMRepository) Init(db interface{}) (err error) {
	if db == nil {
		err = errors.New("Nil pointer as an argument provided.")
		return
	}
	conn, ok := db.(*gorm.DB)
	if !ok {
		err = errors.New(fmt.Sprintf("Incorrect argument type: %v", reflect.TypeOf(db)))
		return
	}
	g.db = conn

	// Initialize GORM Error Converter
	gormConverter := new(GORMErrorConverter)
	gormConverter.Init(db)

	// Assign GORM Error Converter as a repository converter
	g.converter = gormConverter

	if err != nil {
		return err
	}
	return nil
}

func (g *GORMRepository) Create(req interface{}) (err error) {
	if err = g.db.Create(&req).Error; err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Get(req interface{}) (res interface{}, err error) {
	res = restsdk.ObjOfPtrType(req)
	if err = g.db.First(&res, req).Error; err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}
	return res, nil
}

func (g *GORMRepository) List(
	req interface{},
) (res interface{}, err error) {
	// Get Slice of pointer type 'req'
	res = restsdk.SliceOfPtrType(req)

	// List objects provided with arguments probided in request
	if err = g.db.Find(&res, req).Error; err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}
	return res, nil
}

func (g *GORMRepository) ListWithParams(
	req interface{}, params *restsdk.ListParameters,
) (res interface{}, err error) {
	// Get Slice of pointer type 'req'
	res = restsdk.SliceOfPtrType(req)

	err = g.db.
		Offset(params.Offset).
		Limit(params.Limit).
		Order(params.Order).
		Find(&res, req).
		Error
	if err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}
	return res, nil
}

func (g *GORMRepository) Update(req interface{}) (err error) {
	err = g.db.Save(req).Error
	if err != nil {
		err = g.converter.Convert(err)
		return err
	}
	return nil
}

func (g *GORMRepository) Patch(req interface{}) (err error) {
	err = g.db.Update(&req).Error
	if err != nil {
		err = g.converter.Convert(err)
		return err
	}
	return nil
}

func (g *GORMRepository) Delete(req interface{}) (err error) {
	err = g.db.Delete(&req).Error
	if err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

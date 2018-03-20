package gormrepo

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/kucjac/go-rest-sdk/errors/dberrors"
	"github.com/kucjac/go-rest-sdk/forms"
	"github.com/kucjac/go-rest-sdk/repository"
)

type Foo struct {
	ID    uint
	Name  string
	Bar   *Bar
	BarID uint
}

type Bar struct {
	ID   uint
	Name string
}

type GORMRepository struct {
	db        *gorm.DB
	converter dberrors.DBErrorConverter
}

func New(db *gorm.DB) (*GORMRepository, error) {
	gormRepo := &GORMRepository{}
	err := gormRepo.init(db)
	if err != nil {
		return nil, err
	}
	return gormRepo, nil

}

func (g *GORMRepository) init(db *gorm.DB) (err error) {
	if db == nil {
		err = errors.New("Nil pointer as an argument provided.")
		return
	}
	g.db = db

	// Initialize GORM Error Converter
	gormConverter := new(GORMErrorConverter)

	err = gormConverter.Init(db)
	if err != nil {
		return err
	}

	// Assign GORM Error Converter as a repository converter
	g.converter = gormConverter

	return nil
}

func (g *GORMRepository) Create(req interface{}) (err error) {

	g.db.Debug().Find(&req)

	if err = g.db.Create(&req).Error; err != nil {
		return g.converter.Convert(err)
	}
	return nil
}

func (g *GORMRepository) Get(req interface{}) (res interface{}, err error) {
	res = forms.ObjOfPtrType(req)
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
	res = forms.SliceOfPtrType(req)

	// List objects provided with arguments probided in request
	if err = g.db.Find(&res, req).Error; err != nil {
		err = g.converter.Convert(err)
		return nil, err
	}
	return res, nil
}

func (g *GORMRepository) ListWithParams(
	req interface{}, params *repository.ListParameters,
) (res interface{}, err error) {
	// Get Slice of pointer type 'req'
	res = forms.SliceOfPtrType(req)

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

func (g *GORMRepository) Patch(what interface{}, where interface{}) (err error) {

	err = g.db.Update(&what).Select(&where).Error
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

package mgorepo

import (
	"github.com/jinzhu/inflection"
	"github.com/kucjac/go-rest-sdk"
	"gopkg.in/mgo.v2"
)

type MGORepository struct {
	db     *mgo.Session
	dbname string
}

func (m *MGORepository) Create(req interface{}) (err error) {
	err = m.collection(req).Insert(req)
	if err != nil {
		return err
	}
	return nil
}

func (m *MGORepository) Get(req interface{}) (res interface{}, err error) {
	return
}

func (m *MGORepository) List(req interface{}) (res []interface{}, err error) {
	return
}

func (m *MGORepository) ListPaginated(
	req interface{},
	limit, offset int) (res []interface{}, err error) {
	return
}

func (m *MGORepository) Update(req interface{}) (err error) {
	return
}

func (m *MGORepository) Patch(req interface{}) (err error) {
	return
}

func (m *MGORepository) Delete(req interface{}) (err error) {
	return
}

func (m *MGORepository) collection(req interface{}) *mgo.Collection {
	collection := inflection.Plural(restsdk.StructName(req))
	return m.db.DB(m.dbname).C(collection)
}

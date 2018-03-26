package mgorepo

import (
	"errors"
	"github.com/jinzhu/inflection"
	"github.com/kucjac/go-rest-sdk/dberrors"
	"github.com/kucjac/go-rest-sdk/refutils"
	"github.com/kucjac/go-rest-sdk/repository"
	"gopkg.in/mgo.v2"
)

type MGORepository struct {
	session   *mgo.Session
	dbname    string
	converter dberrors.Converter
}

// New creates new MGORepository
func New(session *mgo.Session, dbName string) (repo *MGORepository, err error) {
	if session == nil {
		return nil, errors.New("Nil pointer provided")
	}
	repo = &MGORepository{session: session, dbname: dbName}
	return repo, nil

}

func (m *MGORepository) Create(req interface{}) (dbErr *dberrors.Error) {
	err := m.collection(req).Insert(req)
	if err != nil {
		return m.converter.Convert(err)
	}
	return nil
}

func (m *MGORepository) Get(req interface{}) (res interface{}, dbErr *dberrors.Error) {
	/**

	TO DO

	*/

	return
}

func (m *MGORepository) List(
	req interface{},
) (res []interface{}, dbErr *dberrors.Error) {
	/**

	TO DO

	*/

	return
}

func (m *MGORepository) ListWithParams(
	req interface{}, listParameters *repository.ListParameters,
) (res []interface{}, dbErr *dberrors.Error) {
	/**

	TO DO

	*/

	return
}

func (m *MGORepository) Count(req interface{}) (count int, err *dberrors.Error) {
	/**

	TO DO

	*/

	return
}

func (m *MGORepository) Update(req interface{}) (dbErr *dberrors.Error) {
	/**

	TO DO

	*/

	return
}

func (m *MGORepository) Patch(req, where interface{}) (dbErr *dberrors.Error) {
	/**

	TO DO

	*/
	return
}

func (m *MGORepository) Delete(req, where interface{}) (dbErr *dberrors.Error) {
	/**

	TO DO

	*/
	return
}

func (m *MGORepository) collection(req interface{}) *mgo.Collection {
	collection := inflection.Plural(refutils.StructName(req))
	return m.session.DB(m.dbname).C(collection)
}

package repository

import (
	"github.com/kucjac/go-rest-sdk/dberrors"
)

type GenericRepository interface {
	// Create a new entry for specified 'req' object.
	Create(req interface{}) (err *dberrors.Error)

	// For specific 'req' object the Get method returns first result that match a query
	// where all fields matched with a 'req' would be present.
	// I.e. for:
	// type Foo struct {
	// 	ID int
	//	Name string
	//  Age int
	// }
	// if 'req' would contain ID=1 the record with id=1 would be returned
	// if 'req' fields Age=20 and Name='some name' the first record that match
	// Age=20 and Name='some name' would be returned
	Get(req interface{}) (res interface{}, err *dberrors.Error)

	// List search and list all objects of given type
	// 'req' specify what field values should be present in the result.
	// i.e. if 'req' is a of type Foo struct {Name string} then, providing
	// 'req' &Foo{Name: "Some name"} would list all entries containing field "Name='Some name'"
	List(req interface{}) (res interface{}, err *dberrors.Error)

	// List search and list all objects of given type
	// It extends the List method by providing query parameters.
	// Using list parameters allows i.e. paginate the results or specify order of the query
	// By providing 'IDs' field all entries with given ID's would be queried
	ListWithParams(req interface{}, params *ListParameters) (res interface{}, err *dberrors.Error)

	// Update replaces the whole object with given in argument
	// If no primary key provided in the 'req' or given 'primary key' is not found
	// Update creates new entity
	Update(req interface{}) (err *dberrors.Error)

	// Patch updates only selected fields in the 'req' object
	// selected from 'where' object
	// if where is nil object then all records for given table would be patched.
	Patch(req, where interface{}) (err *dberrors.Error)

	// Delete given records from database provided by 'req' object
	// where describes which entries should be deleted.
	// if where is nil all entries for given model should be deleted.
	Delete(req, where interface{}) (err *dberrors.Error)
}

// List Parameters contains fields common for queries
type ListParameters struct {
	IDs    []int
	Limit  int    `form:"page_size"`
	Offset int    `form:"page"`
	Order  string `form:"order"`
}

// ContainsParameters checks whether given 'ListParameters'
// Has any parameters of non-zero value
func (l ListParameters) ContainsParameters() bool {
	if l.Limit != 0 || l.Offset != 0 || l.Order != "" || len(l.IDs) != 0 {
		return true
	}
	return false
}

package repository

// List Parameters contains fields common for queries
type ListParameters struct {
	IDs    []int
	Limit  int `form:"page_size"`
	Offset int `form:"page"`
	Order  string
}

func (l ListParameters) ContainsParameters() bool {
	if l.Limit != 0 || l.Offset != 0 || l.Order != "" {
		return true
	}
	return false
}

type GenericRepository interface {
	// Create or add new entry
	Create(req interface{}) (err error)

	// Get single object
	Get(req interface{}) (res interface{}, err error)

	// List search and list all objects of given type
	List(req interface{}) (res interface{}, err error)

	// List search and list all objects of given type
	ListWithParams(req interface{}, params *ListParameters) (res interface{}, err error)

	// Update replaces the whole object with given in argument
	Update(req interface{}) (err error)

	// Patch updates only selected fields in the 'what' object
	// selected from 'where' object
	Patch(what interface{}, where interface{}) (err error)

	// Delete given object
	Delete(req interface{}) (err error)
}

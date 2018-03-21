package response

// Body - basic REST API response structure
// Created on purpose of easily managable and
type Body struct {
	// Status is an easy to check variable with only two possible values:
	// - 'ok'
	// - 'error'
	// If the status is 'ok' the response recipient can proceed to result variable
	Status Status `json:"status"`

	// HttpCode is a http status code applicable to this problem
	// While using multiple ResponseErrors with http status
	// it is a good practice to set this value as a leading Http Status
	// i.e.:
	//		- 200 - correct status
	//		- 400 - multiple client error - 4xx
	//		- 500 - multiple API server error - 5xx
	HttpStatus int `json:"httpCode,omitempty"`

	// Errors - list of errors that occurred
	// The server MAY choose to stop processing as soon as a problem is encountered, or it
	// MAY continue processing and encounter multiple problems.
	Errors []error `json:"errors,omitempty"`

	// Content contains all response results
	// Composed as a map[string]interface{}
	// Every result should have it's own tag
	// i.e. "user" : User{1} - user object
	// 		"users" : []User{1,2} - list (plural)
	Content map[string]interface{} `json:"result,omitempty"`
}

// AddContent adds a content for a 'key' string to the response Body
func (r *Body) AddContent(key string, result interface{}) {
	r.Content[key] = result
}

// AddErrors adds an error for the given response body.
func (r *Body) AddErrors(err ...error) {
	r.Errors = append(r.Errors, err...)
}

// New creates new response Body with positive status.
// The function initialize the Body with Status: StatusOK
// HttpStatus: 200 and empty 'Content'
func New() *Body {
	response := &Body{
		Status:     StatusOk,
		HttpStatus: 200,
		Content:    make(map[string]interface{}),
	}
	return response
}

// NewWithError prepares Body with Status: 'StatusError'
// The function takes httpStatus as a first argument. Whereas it can take any int
// it is not a good practice. Multiple errors are allow as the remaining arguments.
func NewWithError(httpStatus int, errors ...error) *Body {
	response := &Body{
		Status:     StatusError,
		HttpStatus: httpStatus,
	}
	response.Errors = append(response.Errors, errors...)
	return response
}

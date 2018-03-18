package response

package restsdk

// APIResponse - basic REST API response structure
// Created on purpose of easily managable and
type APIResponse struct {
	// Status is an easy to check variable with only two possible values:
	// - 'ok'
	// - 'error'
	// If the status is 'ok' the response recipient can proceed to result variable
	Status ResponseStatus `json:"status"`

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

	// Result is a response main content
	// Composed as a map[string]interface{}
	// Every result should have it's own tag
	// i.e. "user" : User{1} - user object
	// 		"users" : []User{1,2} - list (plural)
	Result map[string]interface{} `json:"result,omitempty"`
}

// AddResult adds a result with a 'key' string to the Response
// The result is saved then as a key:value in the Response.
func (r *APIResponse) AddResult(key string, result interface{}) {
	r.Result[key] = result
}

// AddErrors appends errors to the given response
func (r *APIResponse) AddErrors(err ...error) {
	r.Errors = append(r.Errors, err...)
}

// ResponseWithOk prepares APIResponse with Status: 'StatusOk'
// The response has already set httpStatus to 'OK' - 200.
func ResponseWithOk() *APIResponse {
	response := &APIResponse{
		Status:     StatusOk,
		HttpStatus: 200,
		Result:     make(map[string]interface{}),
	}
	return response
}

// ResponseWithError prepares APIResponse with Status: 'StatusError'
// The function first param is the provided httpStatus i.e. BadRequest - 400.
// The rest arguments are multiple errors.
func ResponseWithError(httpStatus int, errors ...error) *APIResponse {
	response := &APIResponse{
		Status:     StatusError,
		HttpStatus: httpStatus,
	}
	response.Errors = append(response.Errors, errors...)
	return response
}

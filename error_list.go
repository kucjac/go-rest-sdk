package restsdk

var (
	CatBadRequest = ErrorCategory{Code: "10c", Title: "Bad request"}
)

var (
	ErrInvalidJSONRequest = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1101",
		Detail:        "Provided request contains invalid json body.",
	}
)

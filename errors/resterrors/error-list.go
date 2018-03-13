package resterrors

// Categories
var (
	ErrHeadersNotSupported = &ResponseError{
		Code: "BRQ001", Title: "Condition Headers Not Supported",
		Status: "400",
	}

	ErrInvalidAuthenticationInfo = &ResponseError{
		Code: "BRQ002", Title: "Invalid Authentication Info",
		Detail: &Detail{Title: `The authentication information was not provided in the correct format. 
			Verify the value of Authorization header.`},
		Status: "400",
	}

	ErrInvalidHeaderValue = &ResponseError{
		Code: "BRQ003", Title: "Invalid Header Value",
		Detail: &Detail{Title: "The value provided in one of the HTTP headers was not in the correct format."},
		Status: "400",
	}

	ErrInvalidInput = &ResponseError{
		Code: "BRQ004", Title: "Invalid Input",
		Detail: &Detail{Title: "One of the request inputs is not valid."},
		Status: "400",
	}

	ErrInvalidQueryParameter = &ResponseError{
		Code: "BRQ005", Title: "Invalid Query Parameter Value",
		Detail: &Detail{Title: "An invalid value was specified for one of the query parameters in the request URI"},
		Status: "400",
	}

	ErrInvalidResourceName = &ResponseError{
		Code: "BRQ006", Title: "Invalid Resource Name",
		Detail: &Detail{Title: "The specified resource name contains invalid characters"},
		Status: "400",
	}

	ErrInvalidURI = &ResponseError{
		Code: "BRQ007", Title: "Invalid URI",
		Detail: &Detail{Title: "The requested URI does not represent any resource on the server"},
		Status: "400",
	}

	ErrInvalidJSONDocument = &ResponseError{
		Code: "BRQ008", Title: "Invalid JSON Document",
		Detail: &Detail{Title: "The specified JSON is not syntatically valid."},
		Status: "400",
	}

	ErrInvalidJSONNodeValue = &ResponseError{
		Code: "BRQ009", Title: "Invalid JSON Node Value",
		Detail: &Detail{Title: "The value provided for one of the JSON nodes in the requested body was not in the correct format"},
		Status: "400",
	}

	ErrMD5Mismatch = &ResponseError{
		Code: "BRQ010", Title: "MD5 Mismatch",
		Detail: &Detail{Title: "The MD5 value specified in the request did not match the MD5 value calculated by the server"},
		Status: "400",
	}

	ErrMetadataTooLarge = &ResponseError{
		Code: "BRQ011", Title: "Metadata Too Large",
		Detail: &Detail{Title: "The size of the specified metada exceeds the maximum size permitted"},
		Status: "400",
	}

	ErrMissingRequiredQueryParam = &ResponseError{
		Code: "BRQ012", Title: "Missing Required Query Parameter",
		Detail: &Detail{Title: "A required query parameter was not specified for this request"},
		Status: "400",
	}

	ErrMissingRequiredHeader = &ResponseError{
		Code: "BRQ013", Title: "Missing Required Header",
		Detail: &Detail{Title: "A required HTTP header was not specified"},
		Status: "400",
	}

	ErrMissingRequiredJSONNode = &ResponseError{
		Code: "BRQ014", Title: "Missing Required JSON Node",
		Detail: &Detail{Title: "A required JSON node was not specified in the request body"},
		Status: "400",
	}

	ErrOutOfRangeInput = &ResponseError{
		Code: "BRQ015", Title: "Request Input Out Of Range",
		Detail: &Detail{Title: "One of the request inputs is out of range"},
		Status: "400",
	}

	ErrOutOfRangeQueryParameterValue = &ResponseError{
		Code: "BRQ016", Title: "Query parameter value out of range",
		Detail: &Detail{Title: "A query parameter specified in the request URI is outside the permissible range"},
		Status: "400",
	}

	ErrUnsupportedHeader = &ResponseError{
		Code: "BRQ017", Title: "Unsupported header",
		Detail: &Detail{Title: "One of the HTTP headers specified in the request is not supported"},
		Status: "400",
	}

	ErrUnsupportedJSONField = &ResponseError{
		Code: "BRQ018", Title: "Unsupported JSON field.",
		Detail: &Detail{Title: "One of the JSON fields specified in the request body is not supported."},
		Status: "400",
	}

	ErrUnsupportedQueryParameter = &ResponseError{
		Code: "BRQ019", Title: "Unsupported query parameter.",
		Detail: &Detail{Title: "One of the query parameters in the request URI is not supported"},
		Status: "400",
	}

	// STATUS 403
	ErrAccountDisabled = &ResponseError{
		Code: "ATH001", Title: "Accound Disabled",
		Detail: &Detail{Title: "The specified account is disabled"},
		Status: "403",
	}

	ErrAuthenticationFailed = &ResponseError{
		Code: "ATH002", Title: "Authentication Failed",
		Detail: &Detail{Title: `Server failed to authenticate the request. Make sure the value of 
		Authorization header is formed correctly including the signature`},
		Status: "403",
	}

	ErrInsufficientAccPerm = &ResponseError{
		Code: "ATH003", Title: "Insufficient account permissions",
		Detail: &Detail{Title: "The account being accessed does not have sufficient permissions to execute this operation."},
		Status: "403",
	}

	// STATUS 404
	ErrResourceNotFound = &ResponseError{
		Code: "NTF001", Title: "Resource not found.",
		Detail: &Detail{Title: "The specified resource does not exists."},
		Status: "404",
	}

	// STATUS 405
	ErrMethodNotAllowed = &ResponseError{
		Code: "BRQ020", Title: "Unsupported http verb",
		Detail: &Detail{Title: "The resource doesn't support the specified HTTP verb."},
		Status: "405",
	}

	// STATUS 409

	ErrAccountAlreadyExists = &ResponseError{
		Code: "CON001", Title: "Account Already Exists",
		Detail: &Detail{Title: "The Specified account already exists"},
		Status: "409",
	}

	ErrResourceAlreadyExists = &ResponseError{
		Code: "CON002", Title: "Resource already exists.",
		Detail: &Detail{Title: "The specified resource already exists."},
		Status: "409",
	}

	// STATUS 413

	ErrRequestBodyTooLarge = &ResponseError{
		Code: "BRQ017", Title: "Request body too large.",
		Detail: &Detail{Title: "The size of the request body exceeds the maximum size permitted"},
		Status: "413",
	}

	// STATUS 500

	ErrInternalError = &ResponseError{
		Code: "INT001", Title: "Internal Server Error",
		Detail: &Detail{Title: "The server encountered an internal error. Please retry the request."},
		Status: "500",
	}

	ErrOperatinTimedOut = &ResponseError{
		Code: "INT002", Title: "Operation Timed Out",
		Detail: &Detail{Title: "The operation could not be completed within the permitted time"},
		Status: "500",
	}

	// STATUS 503

	ErrServerBusy1 = &ResponseError{
		Code: "UNV001", Title: "Server busy",
		Detail: &Detail{Title: "The server is currently unable to receive requests. Please retry your request."},
		Status: "503",
	}
	ErrServerBusy2 = &ResponseError{
		Code: "UNV002", Title: "Server busy",
		Detail: &Detail{Title: "Operations per second is over the account limit"},
		Status: "503",
	}
)

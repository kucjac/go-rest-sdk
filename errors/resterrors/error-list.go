package resterrors

// Contains common REST API Errors
//
// The following error list is based on the Microsoft Azure 'Common REST API Error Codes',
// Published at: https://docs.microsoft.com/en-us/rest/api/storageservices/common-rest-api-error-codes

// List of REST API errors
var (
	// STATUS 400 - CODE: 'BRQXXX'
	ErrHeadersNotSupported = &ResponseError{
		Code: "BRQ001", Title: "Condition headers not supported",
		Status: "400",
	}

	ErrInvalidAuthenticationInfo = &ResponseError{
		Code: "BRQ002", Title: "Invalid authentication info",
		Detail: &Detail{Title: `The authentication information was not provided in the correct format. 
			Verify the value of Authorization header.`},
		Status: "400",
	}

	ErrInvalidHeaderValue = &ResponseError{
		Code: "BRQ003", Title: "Invalid header value",
		Detail: &Detail{Title: "The value provided in one of the HTTP headers was not in the correct format."},
		Status: "400",
	}

	ErrInvalidInput = &ResponseError{
		Code: "BRQ004", Title: "Invalid input",
		Detail: &Detail{Title: "One of the request inputs is not valid."},
		Status: "400",
	}

	ErrInvalidQueryParameter = &ResponseError{
		Code: "BRQ005", Title: "Invalid query parameter value",
		Detail: &Detail{Title: "An invalid value was specified for one of the query parameters in the request URI"},
		Status: "400",
	}

	ErrInvalidResourceName = &ResponseError{
		Code: "BRQ006", Title: "Invalid resource name",
		Detail: &Detail{Title: "The specified resource name contains invalid characters"},
		Status: "400",
	}

	ErrInvalidURI = &ResponseError{
		Code: "BRQ007", Title: "Invalid URI",
		Detail: &Detail{Title: "The requested URI does not represent any resource on the server"},
		Status: "400",
	}

	ErrInvalidJSONDocument = &ResponseError{
		Code: "BRQ008", Title: "Invalid JSON document",
		Detail: &Detail{Title: "The specified JSON is not syntatically valid."},
		Status: "400",
	}

	ErrInvalidJSONFieldValue = &ResponseError{
		Code: "BRQ009", Title: "Invalid JSON field value",
		Detail: &Detail{Title: "The value provided for one of the JSON fields in the requested body was not in the correct format"},
		Status: "400",
	}

	ErrMD5Mismatch = &ResponseError{
		Code: "BRQ010", Title: "MD5 mismatch",
		Detail: &Detail{Title: "The MD5 value specified in the request did not match the MD5 value calculated by the server"},
		Status: "400",
	}

	ErrMetadataTooLarge = &ResponseError{
		Code: "BRQ011", Title: "Metadata too large",
		Detail: &Detail{Title: "The size of the specified metada exceeds the maximum size permitted"},
		Status: "400",
	}

	ErrMissingRequiredQueryParam = &ResponseError{
		Code: "BRQ012", Title: "Missing required query parameter",
		Detail: &Detail{Title: "A required query parameter was not specified for this request"},
		Status: "400",
	}

	ErrMissingRequiredHeader = &ResponseError{
		Code: "BRQ013", Title: "Missing required header",
		Detail: &Detail{Title: "A required HTTP header was not specified"},
		Status: "400",
	}

	ErrMissingRequiredJSONField = &ResponseError{
		Code: "BRQ014", Title: "Missing required JSON field",
		Detail: &Detail{Title: "A required JSON field was not specified in the request body"},
		Status: "400",
	}

	ErrOutOfRangeInput = &ResponseError{
		Code: "BRQ015", Title: "Request input out of range",
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

	// STATUS 403, CODE: 'AUTHXX'
	ErrAccountDisabled = &ResponseError{
		Code: "AUTH01", Title: "Accound disabled",
		Detail: &Detail{Title: "The specified account is disabled"},
		Status: "403",
	}

	ErrAuthenticationFailed = &ResponseError{
		Code: "AUTH02", Title: "Authentication failed",
		Detail: &Detail{Title: `Server failed to authenticate the request. Make sure the value of 
		Authorization header is formed correctly including the signature`},
		Status: "403",
	}

	ErrInsufficientAccPerm = &ResponseError{
		Code: "AUTH03", Title: "Insufficient account permissions",
		Detail: &Detail{Title: "The account being accessed does not have sufficient permissions to execute this operation."},
		Status: "403",
	}
	ErrAuthInvalidCredentials = &ResponseError{
		Code: "AUTH04", Title: "Invalid credentials",
		Detail: &Detail{Title: "Access is denied due to invalid credentials."},
		Status: "403",
	}

	// STATUS 404, CODE: 'NTFXXX'
	ErrResourceNotFound = &ResponseError{
		Code: "NTF001", Title: "Resource not found.",
		Detail: &Detail{Title: "The specified resource does not exists."},
		Status: "404",
	}

	// STATUS 405, CODE: "MNAXXX"
	ErrMethodNotAllowed = &ResponseError{
		Code: "MNA001", Title: "Unsupported http verb",
		Detail: &Detail{Title: "The resource doesn't support the specified HTTP verb."},
		Status: "405",
	}

	// STATUS 409, CODE: "CON001"
	ErrAccountAlreadyExists = &ResponseError{
		Code: "CON001", Title: "Account already exists",
		Detail: &Detail{Title: "The Specified account already exists"},
		Status: "409",
	}

	ErrResourceAlreadyExists = &ResponseError{
		Code: "CON002", Title: "Resource already exists.",
		Detail: &Detail{Title: "The specified resource already exists."},
		Status: "409",
	}

	// STATUS 413, CODE: 'RTLXXX'
	ErrRequestBodyTooLarge = &ResponseError{
		Code: "RTL001", Title: "Request body too large.",
		Detail: &Detail{Title: "The size of the request body exceeds the maximum size permitted"},
		Status: "413",
	}

	// STATUS 500, CODE: 'INTXXX'
	ErrInternalError = &ResponseError{
		Code: "INT001", Title: "Internal server error",
		Detail: &Detail{Title: "The server encountered an internal error. Please retry the request."},
		Status: "500",
	}

	ErrOperatinTimedOut = &ResponseError{
		Code: "INT002", Title: "Operation timed out",
		Detail: &Detail{Title: "The operation could not be completed within the permitted time"},
		Status: "500",
	}

	// STATUS 503, CODE: 'UNAVXX'
	ErrServerBusy1 = &ResponseError{
		Code: "UNAV01", Title: "Server busy",
		Detail: &Detail{Title: "The server is currently unable to receive requests. Please retry your request."},
		Status: "503",
	}
	ErrServerBusy2 = &ResponseError{
		Code: "UNAV02", Title: "Server busy",
		Detail: &Detail{Title: "Operations per second is over the account limit"},
		Status: "503",
	}
)

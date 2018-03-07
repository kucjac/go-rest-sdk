package resterrors

//Categories
var (
	//
	// Client side errors:
	//
	// - Integrity Constraint Violation #23
	CatIntegrityViolation = ErrorCategory{Code: "23", Title: "Integrity Violation"}
	// - Incorrect data type #22
	CatInvalidDataType = ErrorCategory{Code: "22", Title: "Invalid data type"}
	// - Unauthorized #28
	CatUnauthorized = ErrorCategory{Code: "28", Title: "Unauthorized"}
	// - Not Found - #20
	CatNotFound = ErrorCategory{Code: "20", Title: "Not Found"}
	// - Bad Request - incorrect syntax of the request
	CatBadRequest = ErrorCategory{Code: "C1", Title: "Bad Request"}
	// - Invalid Parameters - mostly correct sytnax but invalid parameters provided in request
	CatInvalidParameters = ErrorCategory{Code: "C2", Title: "Invalid parameters"}
	// - Too Many Requests
	CatTooManyRequests = ErrorCategory{Code: "C3", Title: "Too many requests"}

	//
	// Server side errors:
	// - Insufficient resources # 53
	CatInsufficientResources = ErrorCategory{Code: "53", Title: " Insufficient Resources"}
	// - System errors #58
	CatSystemError = ErrorCategory{Code: "58", Title: "System Error"}
	// - Syntax errors DB #42
	CatSyntaxError = ErrorCategory{Code: "42", Title: "Syntax Error"}
	// - Invalid SQL Statement #26
	CatInvSQLStmt = ErrorCategory{Code: "26", Title: "Invalid SQL Statement"}
	// - Invalid Transaction State #25 #2D #0B
	CatInvalidTransaction = ErrorCategory{Code: "25", Title: "Invalid Transaction"}
	// - Cardinality Violation #21
	CatCardinalityViolation = ErrorCategory{Code: "21", Title: "Cardinatlity Violation"}
	// - Invalid role #22, #0L
	CatInvalidRole = ErrorCategory{Code: "22", Title: "Invalid Role"}
	// - Connection Error
	CatConnectionError = ErrorCategory{Code: "C4", Title: "Connection Error"}
	// - Internal Errors
	CatInternalError = ErrorCategory{Code: "0I", Title: "Internal Error"}
)

var ClientErrorCodes map[string]bool = map[string]bool{
	"23": true,
	"22": true,
	"28": true,
	"20": true,
	"C1": true,
	"C2": true,
	"C3": true,
}

var (
	//#23
	ErrDuplicatedValue = ResponseError{
		ErrorCategory: CatIntegrityViolation,
		ID:            "0001",
		Detail:        "Provided entity already exists.",
		Status:        "409",
	}
	ErrForeignKeyViolation = ResponseError{
		ErrorCategory: CatIntegrityViolation,
		ID:            "0002",
		Detail:        "Requested data violates foreign table constraint.",
		Status:        "409",
	}
	ErrCheckViolation = ResponseError{
		ErrorCategory: CatIntegrityViolation,
		ID:            "0003",
		Detail:        "Requested data violates model schema.",
		Status:        "409",
	}
	//#22
	ErrInvalidDataType = ResponseError{
		ErrorCategory: CatInvalidDataType,
		ID:            "D001",
		Detail:        "Provided data is of bad type.",
		Status:        "422",
	}
	ErrInvalidJSONParameters = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P001",
		Detail:        "Provided incorrect json data parameters.",
		Status:        "422",
	}
	ErrInvalidQueryParameters = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P002",
		Detail:        "Query contains incorrect parameters.",
		Status:        "422",
	}
	ErrRequieresQueryParameters = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P003",
		Detail:        "Query does not contain required parameter.",
		Status:        "422",
	}
	ErrRequiresJSONField = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P004",
		Detail:        "JSON input data does not contain required fields.",
		Status:        "422",
	}
	ErrInvalidFormParameters = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P005",
		Detail:        "Form contains incorrect parameters.",
		Status:        "422",
	}
	ErrRequiresFormParameters = ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P005",
		Detail:        "Input form does not contain required parameters.",
		Status:        "422",
	}

	//20
	ErrEntityNotFound = ResponseError{
		ErrorCategory: CatNotFound,
		ID:            "NF01",
		Detail:        "The server cannot find requested resource",
		Status:        "404",
	}

	//C1
	ErrInvalidJSONRequest = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1101",
		Detail:        "The server cannot understand the request due to invalid JSON syntax.",
	}
	ErrInvalidURLQuerySyntax = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1102",
		Detail:        "The server cannot understand the request due to invalid query syntax.",
	}

	ErrBadFormSyntax = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1103",
		Detail:        "The server cannot understand the request due to invalid form syntax.",
	}

	ErrBadRequestNoID = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1104",
		Detail:        "No id provided in URL. This endpoints requires a valid unsigned int id to be provided.",
		Status:        "400",
	}

	//28 - Unauthorized
	ErrInvalidCredentails = ResponseError{
		Status:        "401",
		ErrorCategory: CatUnauthorized,
		ID:            "2001",
		Detail:        "Provided credentials are invalid",
	}
	ErrUnauthorizedAccess = ResponseError{
		Status:        "403",
		ErrorCategory: CatUnauthorized,
		ID:            "2002",
		Detail:        "Unauthorized Access",
	}

	//0I
	ErrInternalServerError = ResponseError{
		Status:        "500",
		ErrorCategory: CatInternalError,
		ID:            "F001",
		Detail:        "Internal Server Error",
	}
)

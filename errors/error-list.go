package resterrors

//Categories
var (

	/**

	Client side errors:

	*/

	// - Not Found - #20 - Entity not found.
	//		- Provided entity not found - i.e.
	//		- If no rows are
	CatNotFound = ErrorCategory{Code: "20", Title: "Not Found"}

	// - Integrity Constraint Violation #23 - used with errors:
	//		- Unique violation - already exists error
	//		- Not null violation - the following field cannot be null
	//		- Check violation - database type check error
	//		- Foreign-Key violation - database foreign key violation - i.e. does not exists
	CatDBIntegrityViolation = ErrorCategory{Code: "D0", Title: "Integrity Violation"}

	// - Unauthorized #28 - This includes errors such as:
	//		- Provided authentication credentials are incorrect
	//		- Unauthorized for given restricted content
	CatUnauthorized = ErrorCategory{Code: "A0", Title: "Unauthorized"}

	// - Bad Request - incorrect syntax of the request i.e.:
	//		- Cannot read the body of the request
	//		- Cannot unmarshal body to json/xml
	//		- invalid request url form
	CatBadRequest = ErrorCategory{Code: "R0", Title: "Bad Request"}

	// - Too Many Requests - for api load balancers, dispatchers to reduce the request count per
	//		the unit of time
	CatTooManyRequests = ErrorCategory{Code: "R1", Title: "Too many requests"}

	// - Invalid Parameters - mostly correct sytnax but invalid parameters provided in request
	//		used while validating given model. Should be used with:
	//		- with validators
	//		- binding to models - i.e. provided query contain incorrect value types
	CatInvalidParameters = ErrorCategory{Code: "P0", Title: "Invalid parameters"}

	/**

	Server side errors:

	*/

	//
	// All errors Title is - 'Internal Server Error'
	//
	// The detail would be shown if the API is

	// - Cardinality Violation #21 i.e.:
	// 		- return more rows than supposed in subquery
	CatDBCardinalityViolation = ErrorCategory{Code: "D1", Title: "Cardinatlity Violation"}

	// - Incorrect data type #22 - database error category, should be handled by database driver.
	//		In this category there it is also good to distinguish:
	//			- Divide by zero error
	//			- Null value not allowed - the validator didn't handled the entity correctly
	CatDBInvalidDataType = ErrorCategory{Code: "D2", Title: "Invalid data type"}

	// - Invalid SQL Statement #26
	CatDBInvSQLStmt = ErrorCategory{Code: "D3", Title: "Invalid SQL Statement"}

	// - Other Database error
	CatDBOthers = ErrorCategory{Code: "DO", Title: "Other Database Errors"}

	// - #0T Invalid Transaction State  #2D #0B
	CatInvalidTransaction = ErrorCategory{Code: "T0", Title: "Invalid Transaction"}

	// - Insufficient resources # 53
	//
	CatInsufficientResources = ErrorCategory{Code: "53", Title: " Insufficient Resources"}
	// - System errors #58
	CatSystemError = ErrorCategory{Code: "S0", Title: "System Error"}
	// - Syntax errors DB #42
	CatSyntaxError = ErrorCategory{Code: "S1", Title: "Syntax Error"}

	// - Connection Error
	CatConnectionError = ErrorCategory{Code: "C4", Title: "Connection Error"}
	// - Internal Errors
	CatInternalError = ErrorCategory{Code: "I0", Title: "Internal Error"}
)

var ClientErrorCodes map[string]bool = map[string]bool{
	"20": true,
	"D0": true,
	"A0": true,
	"R0": true,
	"R1": true,
	"P0": true,
}

var (
	// #23 ErrUniqueViolation
	ErrUniqueViolation = &ResponseError{
		ErrorCategory: CatDBIntegrityViolation,
		ID:            "0001",
		Detail:        "Provided entity already exists.",
		Status:        "409",
	}

	// ErrDBForeignKeyViolation - Foreign-Key Violtion
	// 	Foreign Key Violation Should be treated as invalid request or confilict
	ErrDBForeignKeyViolation = &ResponseError{
		ErrorCategory: CatDBIntegrityViolation,
		ID:            "0002",
		Detail:        "Referenced entity does not exists.",
		Status:        "409",
	}

	// ErrDBCheckViolation - database field checks are violated
	ErrDBCheckViolation = &ResponseError{
		ErrorCategory: CatDBIntegrityViolation,
		ID:            "0003",
		Detail:        "Requested data violates model schema.",
		Status:        "422",
	}

	ErrDBNotNullViolation = &ResponseError{
		ErrorCategory: CatDBIntegrityViolation,
		ID:            "0004",
		Detail:        "Requested data violates model schema.",
		Status:        "422",
	}

	// ErrDBInvalidDataType - #22
	// Trying to insert/upload data of bad type for given field
	//
	ErrDBInvalidDataType = &ResponseError{
		ErrorCategory: CatDBInvalidDataType,
		ID:            "D001",
		Detail:        "Provided data is of bad type.",
		Status:        "422",
	}

	ErrDBUnauthorized = &ResponseError{
		ErrorCategory: CatUnauthorized,
		ID:            "0001",
		Detail:        "Privilege not granted",
		Status:        "409",
	}

	ErrDBInvalidTransactionState = &ResponseError{
		ErrorCategory: CatInvalidTransaction,
		ID:            "00S1",
		Detail:        "Invalid Transaction State",
		Status:        "500",
	}

	ErrDBInvalidTransactionTermination = &ResponseError{
		ErrorCategory: CatInvalidTransaction,
		ID: "00S1",
		Detail: 
	}

	// ErrDBNotMapped is a
	ErrDBNotMapped = &ResponseError{
		ErrorCategory: CatInternalError,
		ID:            "I1",
		Detail:        "Internal Server Error",
		Status:        "500",
	}

	// ErrInvalidJSONParameters - error occurs while binding json form to model
	ErrInvalidJSONParameters = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "JS01",
		Detail:        "Provided incorrect json data parameters.",
		Status:        "422",
	}

	// ErrRequiresJSONField - error occurs while binding json form to
	ErrRequiresJSONField = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "JS04",
		Detail:        "JSON input data does not contain required fields.",
		Status:        "422",
	}

	// ErrInvalidQueryParameters - error occurs while the query is bound incorrectly to model
	ErrInvalidQueryParameters = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P002",
		Detail:        "Query contains incorrect parameters.",
		Status:        "422",
	}

	// ErrRequiresQueryParameters - error occurs while binding query to model and the query
	// does not containt required fields
	ErrRequieresQueryParameters = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P003",
		Detail:        "Query does not contain required parameter.",
		Status:        "422",
	}

	ErrInvalidFormParameters = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P005",
		Detail:        "Form contains incorrect parameters.",
		Status:        "422",
	}
	ErrRequiresFormParameters = &ResponseError{
		ErrorCategory: CatInvalidParameters,
		ID:            "P005",
		Detail:        "Input form does not contain required parameters.",
		Status:        "422",
	}

	//20
	ErrEntityNotFound = &ResponseError{
		ErrorCategory: CatNotFound,
		ID:            "NF01",
		Detail:        "The server cannot find requested resource",
		Status:        "404",
	}

	//C1
	ErrInvalidJSONRequest = &ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1101",
		Detail:        "The server cannot understand the request due to invalid JSON syntax.",
	}
	ErrInvalidURLQuerySyntax = &ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1102",
		Detail:        "The server cannot understand the request due to invalid query syntax.",
	}

	ErrBadFormSyntax = &ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1103",
		Detail:        "The server cannot understand the request due to invalid form syntax.",
	}

	ErrBadRequestNoID = &ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1104",
		Detail:        "No id provided in URL. This endpoints requires a valid unsigned int id to be provided.",
		Status:        "400",
	}

	//28 - Unauthorized

	ErrInvalidCredentails = &ResponseError{
		Status:        "401",
		ErrorCategory: CatUnauthorized,
		ID:            "2001",
		Detail:        "Provided credentials are invalid",
	}
	ErrUnauthorizedAccess = &ResponseError{
		Status:        "403",
		ErrorCategory: CatUnauthorized,
		ID:            "2002",
		Detail:        "Unauthorized Access",
	}

	//0I
	ErrInternalServerError = &ResponseError{
		Status:        "500",
		ErrorCategory: CatInternalError,
		ID:            "F001",
		Detail:        "Internal Server Error",
	}
)

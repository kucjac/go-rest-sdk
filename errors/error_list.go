package resterrors

var (
	CatBadRequest = ErrorCategory{Code: "10c", Title: "Bad request"}

	//SQL Errors
	CatDisconnect          = ErrorCategory{Code: "01002", Title: "Disconnect error"}
	CatDataTruncated       = ErrorCategory{Code: "01004", Title: "Data Truncated"}
	CatPrivilegeNotRevoced = ErrorCategory{Code: "01006", Title: "Privilege not revoked"}
	CatInvalidConnAtr      = ErrorCategory{Code: "01S00", Title: "Invalid Connection String Attribute"}
	CatErrInRow            = ErrorCategory{Code: "01S01", Title: "Error in row"}
	CatNoRowsUpdated       = ErrorCategory{Code: "01S03", Title: "No rows updated or deleted"}
	CatUpdatedMoreThanOne  = ErrorCategory{Code: "01S04", Title: "Updated more than one row"}
	CatWrongParamNo        = ErrorCategory{Code: "07001", Title: "Wrong number of parameters"}
)

var (
	ErrInvalidJSONRequest = ResponseError{
		ErrorCategory: CatBadRequest,
		ID:            "1101",
		Detail:        "Provided request contains invalid json body.",
	}
)

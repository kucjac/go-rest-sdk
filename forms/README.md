# forms
Package forms contains structures and functions for binding data models to multiple form types.

REST API data models are being set to multiple form kinds over the lifetime
of an application.

This package enables binding queries, json forms and url params (third-party router/mux libraries)
to the models of unknown type and unknown fields during the runtime of an application.

Few functions were fetched from *github.com/gin-gonic/gin/binding* where some were enhanced with binding policies.

The policies sets the rules for the binding functions mechanics.

The package contains three basic binding functions:
```go
	func BindQuery(req *http.Request, model interface{}, policy *Policy) error {
		// match given request.URL.Query() to the provided 
		// model fields		
	}

	func BindJSON(req *http.Request, model interface{}, policy *Policy) error {
		// decode given request Body from JSON type and 
		// match the fields to provided model.
	}

	func BindParams(req *http.Request, model interface{}, policy *ParamPolicy) error {
		// search given url/route parameters and match them
		// with proper model fields.
	}
```
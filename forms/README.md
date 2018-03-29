# forms
Package forms contains structures and functions for binding data models to multiple form types.

REST API data models are being set to multiple form kinds over the lifetime
of an application.

This package enables binding queries, json forms and url params (third-party router/mux libraries)
to the models of unknown type and unknown fields during the runtime of an application.

Few functions were fetched from *github.com/gin-gonic/gin/binding* where some were enhanced with binding policies.

The policies sets the rules for the binding functions mechanics.

### Binding functions:
The package contains three basic binding functions:
```go
// BindQuery binds the url.Query() to the provided model.
func BindQuery(req *http.Request, model interface{}, policy *Policy) error {
	// match given request.URL.Query() to the provided 
	// model fields		
}

// BindJSON binds the request.Body to the provided model.
func BindJSON(req *http.Request, model interface{}, policy *Policy) error {
	// decode given request Body from JSON type and 
	// match the fields to provided model.
}

// BindParams binds the route parameters to the provided model.
func BindParams(req *http.Request, model interface{}, getParam ParamGetterFunc, policy *ParamPolicy) error {	
	// search given url/route parameters and match them
	// with proper model fields.
}
```

### Policy
In order to customize the mechanics of these functinos the package provide three different policy types:
The basic policy structure is used for BindQuery and BindJSON as well as the root for the other policies.

```go
// Policy is a set of rules used during the process
// of model binding
type Policy struct {
	TaggedOnly  bool
	FailOnError bool
	Tag         string
}
```
The ListPolicy is based on the 'Policy' enhanced by the parameters used by list	handlers functions.
```go
// ListPolicy is a set of rules used during the process of model
// binding, enhanced with the 'List-parameters' for the list handler function.
type ListPolicy struct {
	Policy
	DefaultLimit int
	WithCount    bool
}
```

The ParamPolicy is based on the 'Policy'. It is used in BindParams function. It enhances the root policy with the 
possibility of deep search - bind multiple url parameters.
```go
// ParamPolicy is a set of rules used during the process of
// routing/ url params.
// Enhances the Policy with DeepSearch field. This field defines if the
// Param binding function should check every model's field.
type ParamPolicy struct {
	Policy
	DeepSearch bool
}
```

### Custom Policy:
If the mechanics of the binding functions are different than expected, then the behavior may be changed by providing custom 
policies.
```go

// Let's specify custom policy
var MyCustomPolicy *Policy

// It is a good practice to create it by copying the DefaultPolicy
// and then customize it
MyCustomPolicy = DefaultPolicy.New()

// Suppose the binding functions should look for the 'mycustomtag'
MyCustomPolicy.Tag = 'mycustomtag'

// If the field doesn't contain 'mycustomtag' binding function would omit it, unless it is a struct.
MyCustomPolicy.TaggedOnly = true
```

### ParamGetterFunc:
In order to `adapt` any url parameters mechanics for the third-party routers, the package defines the 'ParamGetterFunc'. 
It is an adapter function that for provided parameter name and http request it should return a parameter value or 
an error if something went wrong.
```go
// ParamGetterFunc defines the adaption function that retrieve the parameters
// from the specific third-party routing framework on the base
// of the provided parameterName string and req *http.Request
// if individual implementation needs more arguments push them into
// request's context.
type ParamGetterFunc func(paramName string, req *http.Request) (string, error)
```

Writing the custom ParamGetterFunc should be an easy task. If the implementation requires more 
arguments in order to get the parameter, Set the request context value with it.

### Custom ParamGetterFunc Example:
```go
// For example purpose lets suppose that the router is from 'mypkg' package, that needs multiple arguments 
// in its parameter function implementation.
// In order to adapt that function let's write MyCustomParamGetterFunc.

// Having some handler function that uses 'mypkg' addtional arguments that are necessary to get parameters.
func SomeHandlerFunc(rw http.ResponseWriter, req *http.Request){
	...
	ctx := req.Context()
	ctx = context.WithValue(ctx, "mypkgArgument", mypkg.Arguments)
	req.WithContext(ctx)	
	...
	err = BindParam(req, model, MyCustomParamGetterFunc, nil)
	...
}

// some custom error when no 'mypkg.Argument' value in the context.
var ErrNoArgument = errors.New("No mypkg.Argument in context")

// MyCustomParamGetterFunc is an adaption function that implements ParamGetterFunc for the imaginary 'mypkg' package.
func MyCustomParamGetterFunc(paramName string, req *http.Request) (string ,error){
	// Supposed that 'mypkg' requires the mypkg.Argument in order to get parameters from url
	// Get the value from the context
	argument, ok  := req.Context().Value("mypkgArgument").(mypkg.Argument)
	if !ok {
		return "", ErrNoArgument
	}
	// Get the parameters with an argumnt
	paramValue := mypkg.Parameters(paramName, argument)
	return paramValue
}
```

### BindQuery Example:
BindQuery binds the url queries to the provided model. 
It should be used inside handler function, so that the url would be matched to provided data model.
```go
// Let's have model that use custom query tags
type MyModel struct {
	Name string `myquerytag:"firstname"`
	Age int `myquerytag:"age"`
}

func SomeHandlerFunc(rw http.ResponseWriter, req *http.Request) {
	...
	var model MyModel

	var policy *Policy = forms.DefaultPolicy.New()
	policy.Tag = "myquerytag"
	policy.FailOnError = true

	err = BindQuery(req, &model, policy)
	if err != nil {
		// handle the error
	}
	...
}

// Having some basic mux with some route set to 'SomeHandlerFunc'
func main(){
	mux :=http.NewServeMux()
	mux.HandleFunc("/custom/route", SomeHandlerFunc)

	req := httptest.NewRequest("GET", "/custom/route?firstname=John&age=50")
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	// Then the model within the SomeHandlerFunc should be bound to the query
	// and its Name field should equal: 'John' while the Age field should equal: 50
}
```









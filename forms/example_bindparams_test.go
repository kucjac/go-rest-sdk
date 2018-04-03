package forms

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

type Model struct {
	ID   int
	Name string
	Bar  Bar
}

func (m Model) String() string {
	return fmt.Sprintf("{ID: %v, Name: %v, Bar: %v}", m.ID, m.Name, m.Bar)
}

type Bar struct {
	ID   int
	Date time.Time `param:"date" time_format:"2006-02-01"`
}

func (b Bar) String() string {
	return fmt.Sprintf("{ID: %v, Date: %4d-%02d-%02d}",
		b.ID, b.Date.Year(), b.Date.Day(), b.Date.Month())
}

// ChiParamGetterFunc is custom ParamGetterFunc
// This is just an example of wrapping third-party ParamGetting functions.
func ChiParamGetterFunc(paramName string, req *http.Request) (string, error) {
	// take the param from the request
	paramValue := chi.URLParam(req, paramName)
	return paramValue, nil
}

// ExampleBindParams is an example of bind params using third party routing - in this example
// it is 'github.com/go-chi/chi' package
func ExampleBindParams() {

	// In this example chi.Router would be used
	mux := chi.NewMux()

	handleFooBarDate := func(rw http.ResponseWriter, req *http.Request) {
		// Match the model with the routing url params
		model := &Model{}

		// set the policy to default - with DeepSearch
		policy := DefaultParamPolicy.Copy()
		policy.SearchDepthLevel = 1
		policy.FailOnError = true

		err := BindParams(req, model, ChiParamGetterFunc, policy)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Bind Parameter errors: %v", err), 500)
			return
		}

		httpResponse, err := json.Marshal(model)
		if err != nil {
			http.Error(rw, "Cannot marshal the model", 500)
			return
		}
		rw.Write(httpResponse)
	}

	// Let our mux handle url with multiple parameters
	mux.Get("/models/{model}/bars/{bar}/date/{date}", handleFooBarDate)

	req := httptest.NewRequest("GET", "/models/15/bars/55/date/2016-20-03", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	fmt.Printf("The router responded with: %v HTTP code\n", rw.Code)

	var BoundModel Model
	body, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		fmt.Println("Error occured when reading body of ResponseWriter")
		return
	}
	err = json.Unmarshal(body, &BoundModel)
	if err != nil {
		fmt.Println(string(body))
		return
	}

	fmt.Printf("Model: %v", BoundModel)
	// Output:
	// The router responded with: 200 HTTP code
	// Model: {ID: 15, Name: , Bar: {ID: 55, Date: 2016-20-03}}
}

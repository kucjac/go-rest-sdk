go-rest-sdk
=============
package that contains basic RESTful API structures.

[![Build Status](https://travis-ci.org/kucjac/go-rest-sdk.svg?branch=master)](https://travis-ci.org/kucjac/go-rest-sdk)
[![Coverage Status](https://coveralls.io/repos/github/kucjac/rest-response/badge.svg?branch=master)](https://coveralls.io/github/kucjac/rest-response?branch=master)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/kucjac/rest-response)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/kucjac/rest-response/master/LICENSE)

It is mostly used to unify the response value for REST Api.

### Features:

- Status easy to discover - the response contains the field "status" which is prepared to serve only as:
	- `"ok"` 
	- `"error"`
- Multiple Error support - the response supports multiple errors during processing. In example while processing the request on API, there occured minor error that does not break the context of the request. If in the same context another error was thrown, it may be saved with its description and http.Status. These errors may become helpful for the developer that uses this API.
- Well described Response Errors - ```ResponseError``` struct based on the JSON API error, contain multiple fields that makes life easier for developer who uses an API
- Support to Categorize API errors - ```ErrorCategory``` may be used as an error category to the API documentation.
- Multiple result support - the results (main content) of the response may be a multiple 'key':'result' value 

### Installation:

```go get -u github.com/kucjac/rest-response```


### Examples:

The tests are written using [`go-convey`](https://github.com/smartystreets/goconvey) package. Reading all tests should show how this package works.

*Inspired by:*
 - [Medium@Shazow](https://medium.com/@shazow/how-i-design-json-api-responses-71900f00f2db)
 - [JSON-API#errors](http://jsonapi.org/format/#errors)

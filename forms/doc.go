/*
package forms contain structures and functions for model bindings

REST API data models are being set to multiple form kinds over the lifetime
of an application.

This package enables binding queries, json forms and url params (third-party router/mux libraries)
to the models of unknown type and unknown fields during the runtime of an application.

Few functions were fetched from github.com/gin-gonic/gin/binding and enhanced with policies.
The policies sets the rules for the binding functions mechanics.

Binding functions:
	BindQuery	- used to set the query parameters into model
	BindJSON	- binds json form into provided model
	BindParams	- binds the url routing parameters to the given model

There are three types of the polices:
	Policy		- this is the basic policy structure. Used for BindQuery and BindJSON
	ListPolicy	- based on the 'Policy' enhancing it by the parameters used to list
				records.
	ParamPolicy	- based on the 'Policy' used in BindParams function. Enhances the policy with the
				possibility of deep search - bind params different than main object id.

The basic Policy contains three basic rules:
	TaggedOnly	- if set to true binding works only on correctly tagged fields
	FailOnError	- if set to true binding functions returns error if occurs. There is an exception
				for the 'BindParams' function - if the object main id was not set, even if this
				is set to false the function would return an error.
	Tag			- defines the 'tag' that would be used for the binding function using this policy

The ListPolicy enhances the policy with fields:
	DefaultLimit	- used for setting default limit of the recoreds return with the list handler 				functions.
	WithCount		- specified if the list handler function should include count of the whole 				collection.

The ParamPolicy enhances Policy with:
	DeepSearch	- if set to true the BindParam function searches for any matching field - param
			pair.


There are few default policies for different use purpose:
	DefaultPolicy 		- default policy
	DefaultJSONPolicy 	- default policy for binding JSON
	DefaultListPolicy	- default policy containing List-Parameters
	DefaultParamPolicy	- default policy for binding parameters.

By default they are used in their purposed functions.
In order to use a copy of these, just use New() method and a new copy would be returned.

*/
package forms

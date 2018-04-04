/*
package forms contain structures and functions for model bindings

REST API data models are being set to multiple form kinds over the lifetime
of an application.

This package enables binding queries, json forms and url params (third-party router/mux libraries)
to the models of unknown type and unknown fields during the runtime of an application.

The binding functions were written so that they use policies.
The policies sets the rules for the binding functions mechanics.

Few functions were fetched from github.com/gin-gonic/gin/binding.


Binding functions:
	BindQuery	- used to set the query parameters into model
	BindJSON	- binds json form into provided model
	BindParams	- binds the url routing parameters to the given model

There are Two types of the polices:
	BindPolicy	- this is the basic policy structure. Used for BindQuery and BindJSON
	ParamPolicy	- based on the 'Policy' used in BindParams function. Enhances the policy with the
				possibility of deep search - bind params different than main object id.

The basic Policy contains three basic rules:
	TaggedOnly	- if set to true binding works only on correctly tagged fields
	FailOnError	- if set to true binding functions returns error if occurs. There is an exception
				for the 'BindParams' function - if the object main id was not set, even if this
				is set to false the function would return an error.
	Tag			- defines the 'tag' that would be used for the binding function using this policy
	SearchDepthLevel - is the depth of the search for nested structs

The ParamPolicy enhances Policy with:
	IDOnly	- if set to true the BindParam function searches only for the ID field pair (with also
			nested structs if searchdepthlevel is greater than 0).


There are few default policies for different use purpose:
	DefaultPolicy 		- default policy
	DefaultParamPolicy	- default policy for binding parameters.

In order to use a copy of these, use Copy() method and a new copy would be returned.

*/
package forms

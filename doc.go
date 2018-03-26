/*
Package sdk is a pluggable SDK for creating RESTful API's in Golang.

Developed with an idea of creating RESTful API's using the most popular Golang tools,
not forcing the use of specific one.

The SDK removes the complexity of using different web frameworks, databases and data models
in a process of creating RESTful API's.

Almost all components are based on the interfaces. That enables combining multiple
independent third-party tools. This solution allows to easily develop
components either based on the 'go-rest-sdk' prepared tools or on custom implementations.

The package is divided into eight main components:
	dberrors 	# unifies the database errors. Defines the 'Converter' interface and database Errors prototypes
	errhandler	# handles is a mapping of database errors into resterrors. Defines 'ErrorHandler'
			that Handles provided 'dberrors.Error' and maps into 'resterrors.Error'
	forms		# enables binding provided model to different form types.
	handlers	# joins 'go-rest-sdk' packages to create model, web framework and database
			repository independent RESTful handlers.
	refutils	# contains reflect encapsulations useful for other subpackages
	repository	# defines database and models repositories. Defines 'Repository' interface.
	response	# contains body for the RESTful API responses. Defines 'Responser' and
			'StatusResponser' interfaces.
	resterrors	# defines RESTful response 'Errors', and their prototypes.

*/
package sdk

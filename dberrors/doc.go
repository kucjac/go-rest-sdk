/*
	Package dberrors contains unified database Errors converted from third-party package errors.

	In order to create the RESTful API that is indepenedent of the database type, the database
	errors must be converted into single form.
	This package defines database Errors with some Prototypes that are the most common error categories.

	In order to maintaing uniform form of the error converting, every database driver should
	implement the 'Converter' interface.
*/
package dberrors

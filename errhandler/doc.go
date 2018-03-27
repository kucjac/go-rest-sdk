/*
Package errhandler handles dberrors.Error conversion into resterrors.Error

In order to automate handling database errors into proper resterror.Errors,
errhandler package provide 'ErrorHandler' structure. It contains map where
for every database dberrors.Error key have its own resterror.Error value.

By default 'ErrorHandler' use 'DefaultErrorMap' as a container for conversion.
If there is a need of changing single mapping the 'UpdateErrorEntry' method should be used.
If more entries or whole map should be changed the 'LoadCustomErrorMap' method may be used.

*/
package errhandler

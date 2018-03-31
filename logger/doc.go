/*
package logger contains generic-logger, basic-logger and logging interfaces.

Logging is very important Rest API feature.
In order not to extort any specific logging package, a generic-logger has been created.
GenericLogger is a wrapper around third-party loggers which must implement of three
specified logging-interfaces.
This solution allows to use prepared function that requieres logging
using your favorite logger.

There is also BasicLogger logger that implements 'LeveledLogger' interface
that is may be used as a default logger.
It is very simple and lightweight implementation of leveled logger.

*/

package logger

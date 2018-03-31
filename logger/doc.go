/*
package logger contains LoggerWrapper, BasicLogger and logging interfaces.

In order not to extort any specific logging package, a LoggerWrapper has been created.
LoggerWrapper wraps around third-party loggers that implement one of the logging-interfaces:
	# StdLogger - standard library logger interface
	# LeveledLogger - basic leveled logger interface
	# ShortLeveledLogger - basic leveled logger interfaces with shortened method names
	# ExtendedLeveledLogger - a fully leveled logger interface

This solution allows to use ExtendedLeveledLogger interface methods for most of the third-party
logging packages.

There is also BasicLogger logger that implements 'LeveledLogger' interface.
It is very simple and lightweight implementation of leveled logger.
*/

package logger

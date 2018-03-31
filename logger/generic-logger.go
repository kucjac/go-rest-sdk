package logger

import (
	"errors"
	"fmt"
)

// GenericLogger is wrapper around any third-party logger that implements any of
// the following interfaces:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// For loggers that implements only StdLogger, GenericLogger tries to virtualize
// ExtendedLeveledLogger behaviour. It adds level name and adds it before logging message,
// without full levels implementation.
// If a logger implements LeveledLogger that doesn't have specific log line '****ln()' methods,
// it uses default non 'ln' functions - i.e. instead 'Infoln' uses 'Info'.
type GenericLogger struct {
	logger        interface{}
	currentLogger int
}

// NewGenericLogger creates a GenericLogger wrapper over provided 'logger' argument
// By default the function checks if provided logger implements logging interfaces
// in a following hierarchy:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// if logger doesn't implement an interface it tries to check the next in hierarchy.
// If it doesn't implement any of known logging interfaces the function returns error.
func NewGenericLogger(logger interface{}) (*GenericLogger, error) {
	return newGenericLogger(logger)
}

// MustGetGenericLogger creates a GenericLogger wrapper over provided 'logger' argument.
// By default the function checks if provided logger implements logging interfaces
// in a following hierarchy:
//	# ExtendedLeveledLogger
//	# ShortLeveledLogger
//	# LeveledLogger
//	# StdLogger
// if logger doesn't implement an interface it tries to check the next in hierarchy.
// If it doesn't implement any of known logging interfaces the function panics.
func MustGetGenericLogger(logger interface{}) *GenericLogger {
	generic, err := newGenericLogger(logger)
	if err != nil {
		panic(err)
	}
	return generic
}

func newGenericLogger(logger interface{}) (*GenericLogger, error) {
	generic := &GenericLogger{}
	var err error

	if l, ok := logger.(ExtendedLeveledLogger); ok {
		generic.logger = l
		generic.currentLogger = 4
		return generic, nil
	}

	if l, ok := logger.(ShortLeveledLogger); ok {
		generic.logger = l
		generic.currentLogger = 3
		return generic, nil
	}

	if l, ok := logger.(LeveledLogger); ok {
		generic.logger = l
		generic.currentLogger = 2
		return generic, nil
	}

	if l, ok := logger.(StdLogger); ok {
		generic.logger = l
		generic.currentLogger = 1
		return generic, nil
	}

	err = errors.New("Provided logger doesn't implement any known interfaces")
	return nil, err
}

// Print calls Output to print the GenericLogger.
// Arguments are handled in the manner of log.Print for StdLogger and
// Extended LeveledLogger as well as log.Info for LeveledLogger
func (c *GenericLogger) Print(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Print(args...)
	default:
	}
}

// Printf calls formatted Output to print the GenericLogger.
// Arguments are handled in the manner of log.Printf for StdLogger and
// Extended LeveledLogger as well as log.Infof for LeveledLogger
func (c *GenericLogger) Printf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Infof(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Infof(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Printf(format, args...)
	default:
	}
}

// Println calls Output to print the GenericLogger.
// Arguments are handled in the manner of log.Println for StdLogger and
// Extended LeveledLogger as well as log.Info for LeveledLogger
func (c *GenericLogger) Println(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Println(args...)
	default:

	}
}

// Debug calls Output to print the GenericLogger with DEBUG level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Debug for ExtendedLeveledLogger and LeveledLogger.
func (c *GenericLogger) Debug(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debug(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debug(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debug(args...)
	default:
	}
}

// Debugf calls formatted Output to print the GenericLogger with DEBUG level.
// Arguments are handled in the manner of log.Printf for StdLogger,
// log.Debugf for ExtendedLeveledLogger and LeveledLogger.
func (c *GenericLogger) Debugf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debugf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debugf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debugf(format, args...)
	default:
	}
}

// Debugln calls Output to print the GenericLogger with DEBUG level.
// Arguments are handled in the manner of log.Println for StdLogger,
// log.Debugln for ExtendedLeveledLogger and log.Debug for LeveledLogger .
func (c *GenericLogger) Debugln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(DEBUG, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Debug(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Debug(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Debugln(args...)
	default:
	}

}

// Debug calls Output to print the GenericLogger with INFO level.
// Arguments are handled in the manner of log.Print for StdLogger,
// log.Info for ExtendedLeveledLogger and LeveledLogger.
func (c *GenericLogger) Info(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Info(args...)
	default:
	}

}

func (c *GenericLogger) Infof(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Infof(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Infof(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Infof(format, args...)
	default:
	}
}

func (c *GenericLogger) Infoln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(INFO, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Info(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Info(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Infoln(args...)
	default:
	}
}

func (c *GenericLogger) Warning(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warning(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warn(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warning(args...)
	default:
	}
}

func (c *GenericLogger) Warningf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warningf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warnf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warningf(format, args...)
	default:
	}
}

func (c *GenericLogger) Warningln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(WARNING, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Warning(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Warn(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Warningln(args...)
	default:
	}
}

func (c *GenericLogger) Error(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, nil, args...)
		log.Print(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Error(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Error(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Error(args...)
	default:
	}
}

func (c *GenericLogger) Errorf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, &format, args...)
		log.Printf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Errorf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Errorf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Errorf(format, args...)
	default:
	}
}

func (c *GenericLogger) Errorln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(ERROR, nil, args...)
		log.Println(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Error(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Error(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Errorln(args...)
	default:
	}
}

func (c *GenericLogger) Fatal(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Fatal(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatal(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatal(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatal(args...)
	default:
	}
}

func (c *GenericLogger) Fatalf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, &format, args...)
		log.Fatalf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatalf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatalf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatalf(format, args...)
	default:
	}
}

func (c *GenericLogger) Fatalln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Fatalln(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Fatal(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Fatal(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Fatalln(args...)
	default:
	}
}

func (c *GenericLogger) Panic(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Panic(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panic(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panic(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panic(args...)
	default:
	}
}

func (c *GenericLogger) Panicf(format string, args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, &format, args...)
		log.Panicf(format, args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panicf(format, args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panicf(format, args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panicf(format, args...)
	default:
	}
}

func (c *GenericLogger) Panicln(args ...interface{}) {
	switch c.currentLogger {
	case 1:
		log := c.logger.(StdLogger)
		args = buildLeveled(CRITICAL, nil, args...)
		log.Panicln(args...)
	case 2:
		log := c.logger.(LeveledLogger)
		log.Panic(args...)
	case 3:
		log := c.logger.(ShortLeveledLogger)
		log.Panic(args...)
	case 4:
		log := c.logger.(ExtendedLeveledLogger)
		log.Panicln(args...)
	default:
	}
}

func buildLeveled(level Level, format *string, args ...interface{}) (leveled []interface{}) {
	if format == nil {
		leveled = append(leveled, fmt.Sprintf("%s: ", level))
		leveled = append(leveled, args...)
	} else {
		leveled = append(leveled, args...)
		msg := fmt.Sprintf("%s: %s", level, *format)
		*format = msg
	}
	return leveled
}

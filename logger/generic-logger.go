package logger

import (
	"errors"
)

type GenericLogger struct {
	stdLogger        StdLogger
	leveledLogger    LeveledLogger
	extLeveledLogger ExtendedLeveledLogger
	currentLogger    int
}

func NewGenericLogger(logger interface{}) (*GenericLogger, error) {
	generic := &GenericLogger{}
	var err error
	switch l := logger.(type) {
	case StdLogger:
		generic.stdLogger = l
		generic.currentLogger = 1
	case LeveledLogger:
		generic.leveledLogger = l
		generic.currentLogger = 2
	case ExtendedLeveledLogger:
		generic.extLeveledLogger = l
		generic.currentLogger = 3
	default:
		generic = nil
		err = errors.New("Provided logger doesn't implement any known interfaces")
	}
	return generic, err
}

func (c *GenericLogger) Print(args ...interface{}) {

}

func (c *GenericLogger) Printf(format string, args ...interface{}) {

}

func (c *GenericLogger) Println(args ...interface{}) {

}

func (c *GenericLogger) Debug(args ...interface{}) {

}

func (c *GenericLogger) Debugf(format string, args ...interface{}) {

}

func (c *GenericLogger) Debugln(args ...interface{}) {

}

func (c *GenericLogger) Info(args ...interface{}) {

}

func (c *GenericLogger) Infof(format string, args ...interface{}) {

}

func (c *GenericLogger) Infoln(args ...interface{}) {

}

func (c *GenericLogger) Warning(args ...interface{}) {

}

func (c *GenericLogger) Warningf(format string, args ...interface{}) {

}

func (c *GenericLogger) Warningln(args ...interface{}) {

}

func (c *GenericLogger) Error(args ...interface{}) {

}

func (c *GenericLogger) Errorf(format string, args ...interface{}) {

}

func (c *GenericLogger) Errorln(args ...interface{}) {

}

func (c *GenericLogger) Fatal(args ...interface{}) {

}

func (c *GenericLogger) Fatalf(format string, args ...interface{}) {

}

func (c *GenericLogger) Fatalln(args ...interface{}) {

}

func (c *GenericLogger) Panic(args ...interface{}) {

}

func (c *GenericLogger) Panicf(format string, args ...interface{}) {

}

func (c *GenericLogger) Panicln(args ...interface{}) {

}

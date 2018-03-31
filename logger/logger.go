package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
)

var logSequenceID uint64

func init() {
	logSequenceID = 0
}

/**

Levels

*/
// Level defines a logging level used in BasicLogger
type Level int

// Following levels are supported in BasicLogger
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	CRITICAL
	PRINT
)

var level_names = []string{
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"CRITICAL",
	"INFO",
}

func (l Level) String() string {
	return level_names[l]
}

/**

Message

*/

// Message is a basic logging record structure used in BasicLogger
type Message struct {
	id      uint64
	level   Level
	fmt     *string
	message *string
	args    []interface{}
}

// Message prepares the string message based on the format and args private fields
// of the message
func (m *Message) Message() string {
	return m.getMessage()
}

func (m *Message) getMessage() string {
	if m.message == nil {
		var msg string
		if m.fmt == nil {
			//println etc.
			msg = fmt.Sprint(m.args...)
		} else {
			msg = fmt.Sprintf(*m.fmt, m.args...)
		}
		m.message = &msg
	}
	return *m.message
}

// String returns string that concantates:
// id hash - 4 digits|time formatted in RFC339|level|message
func (m *Message) String() string {
	msg := fmt.Sprintf("%s|%04x: %s", m.level, m.id, m.getMessage())
	return msg
}

/**

BasicLogger

*/

// BasicLogger is simple leveled logger that implements LeveledLogger interface.
// It uses 5 basic log levels:
//	# DEBUG
//	# INFO
//	# WARNING
//	# ERROR
//	# CRITICAL
// By default DEBUG level is used. It may be reset using SetLevel() method.
// It allows to filter the logs by given level.
// I.e. Having BasicLogger with level Set to WARNING, then there would be
// no DEBUG and INFO logs (the hierarchy goes up only).
type BasicLogger struct {
	stdLogger *log.Logger
	level     Level
}

// NewBasicLogger creates new BasicLogger that shares common sequence id.
// By default it uses DEBUG level. It can be changed later using SetLevel() method.
// BasicLogger uses standard library *log.Logger for logging purpose.
// The arguments used in this function are described in log.New() method.
func NewBasicLogger(out io.Writer, prefix string, flags int) *BasicLogger {
	logger := &BasicLogger{
		stdLogger: log.New(out, prefix, flags),
		level:     DEBUG,
	}
	return logger
}

// SetLevel sets the level of logging for given Logger.
func (l *BasicLogger) SetLevel(level Level) {
	l.level = level
}

func (l *BasicLogger) Debug(args ...interface{}) {
	l.log(DEBUG, nil, args...)
}

func (l *BasicLogger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, args...)
}

func (l *BasicLogger) Info(args ...interface{}) {
	l.log(INFO, nil, args...)
}

func (l *BasicLogger) Infof(format string, args ...interface{}) {
	l.log(INFO, &format, args...)
}

func (l *BasicLogger) Print(args ...interface{}) {
	l.log(PRINT, nil, args...)
}

func (l *BasicLogger) Printf(format string, args ...interface{}) {
	l.log(PRINT, &format, args...)
}

func (l *BasicLogger) Warning(args ...interface{}) {
	l.log(WARNING, nil, args...)
}

func (l *BasicLogger) Warningf(format string, args ...interface{}) {
	l.log(WARNING, &format, args...)
}

func (l *BasicLogger) Error(args ...interface{}) {
	l.log(ERROR, nil, args...)
}

func (l *BasicLogger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, &format, args...)
}

func (l *BasicLogger) Fatal(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	os.Exit(1)
}

func (l *BasicLogger) Fatalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	os.Exit(1)
}

func (l *BasicLogger) Panic(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	panic(fmt.Sprint(args...))
}

func (l *BasicLogger) Panicf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	panic(fmt.Sprintf(format, args...))
}

/**

PRIVATE

*/

func (l *BasicLogger) log(level Level, format *string, args ...interface{}) {
	if !l.isLevelEnabled(level) {
		return
	}
	msg := &Message{
		id:    atomic.AddUint64(&logSequenceID, 1),
		level: level,
		fmt:   format,
		args:  args,
	}
	l.stdLogger.Output(2, msg.String())
}

func (l *BasicLogger) isLevelEnabled(level Level) bool {
	return level >= l.level
}

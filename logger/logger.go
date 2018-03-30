package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"
)

// StdLogger is the logger interface for standard log library
type StdLogger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

// LeveledLogger is a logger that uses basic logging leveles
type LeveledLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

// ExtendedLogger adds distinction between Leveled methods that starts new or not
// i.e.: 'Debugln' and 'Debug'.
// It also adds all Print's methods.
type ExtendedLeveledLogger interface {
	LeveledLogger

	Print(args ...interface{})
	Printf(args ...interface{})
	Println(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type Level int

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

var sequenceID uint64

func init() {
	sequenceID = 0
}

func (l Level) String() string {
	return level_names[l]
}

type Message struct {
	id      uint64
	time    time.Time
	level   Level
	fmt     *string
	message *string
	args    []interface{}
}

func (m *Message) Message() string {
	return m.getMessage()
}

func (m *Message) getMessage() string {
	if m.message == nil {
		var msg string
		if m.fmt == nil {
			//println etc.
			msg = fmt.Sprintln(m.args...)
		} else {
			msg = fmt.Sprintf(*m.fmt, m.args...)
		}
		m.message = &msg
	}
	return *m.message
}

func (m *Message) String() string {
	msg := fmt.Sprintf("%04x|%s|%s|%s", m.id, m.time, m.level, m.getMessage())
	return msg
}

// BasicLogger is simple leveled logger that implements ExtendedLeveledLogger
// It uses 5 basic log levels:
//	- DEBUG
//	- INFO
//	- WARNING
//	- ERROR
//	- CRITICAL
// By default DEBUG level is set. It may be set using SetLevel() method.
// This allows to control what output goes to the logs
type BasicLogger struct {
	stdLogger *log.Logger
	level     Level
}

// NewBasicLogger creates new BasicLogger that shares common sequence id.
// By default it uses DEBUG level. It can be changed later using SetLevel() method.
// The provided arguments creates standard *log.Logger. In this way you can set whatever output
// that is possible
func NewBasicLogger(out io.Writer, prefix string, flag int) *BasicLogger {
	logger := &BasicLogger{
		stdLogger: log.New(out, prefix, flag),
		level:     DEBUG,
	}
	return logger
}

// SetLevel sets the level of logging for given Logger.
func (l *BasicLogger) SetLevel(level Level) {
	l.level = level
}

func (l *BasicLogger) log(level Level, format *string, ln bool, args ...interface{}) {
	if !l.isLevelEnabled(level) || level == PRINT {
		return
	}
	msg := &Message{
		id:    atomic.AddUint64(&sequenceID, 1),
		time:  time.Now(),
		level: level,
		fmt:   format,
		args:  args,
	}
	l.stdLogger.Output(2, msg.String())
}

func (l *BasicLogger) isLevelEnabled(level Level) bool {
	return level >= l.level
}

func (l *BasicLogger) Debug(args ...interface{}) {
	l.log(DEBUG, nil, false, args...)
}

func (l *BasicLogger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, false, args...)
}

func (l *BasicLogger) Debugln(args ...interface{}) {
	l.log(DEBUG, nil, true, args...)
}

func (l *BasicLogger) Info(args ...interface{}) {
	l.log(INFO, nil, false, args...)
}

func (l *BasicLogger) Infof(format string, args ...interface{}) {
	l.log(INFO, &format, false, args...)
}

func (l *BasicLogger) Infoln(args ...interface{}) {
	l.log(INFO, nil, true, args...)
}

func (l *BasicLogger) Print(args ...interface{}) {
	l.log(PRINT, nil, false, args...)
}

func (l *BasicLogger) Printf(format string, args ...interface{}) {
	l.log(PRINT, &format, false, args...)
}

func (l *BasicLogger) Println(args ...interface{}) {
	l.log(PRINT, nil, true, args...)
}

func (l *BasicLogger) Warning(args ...interface{}) {
	l.log(WARNING, nil, false, args...)
}

func (l *BasicLogger) Warningf(format string, args ...interface{}) {
	l.log(WARNING, &format, false, args...)
}

func (l *BasicLogger) Warningln(args ...interface{}) {
	l.log(WARNING, nil, true, args...)
}

func (l *BasicLogger) Error(args ...interface{}) {
	l.log(ERROR, nil, false, args...)
}

func (l *BasicLogger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, &format, false, args...)
}

func (l *BasicLogger) Errorln(args ...interface{}) {
	l.log(ERROR, nil, true, args...)
}

func (l *BasicLogger) Fatal(args ...interface{}) {
	l.log(CRITICAL, nil, false, args...)
	os.Exit(1)
}

func (l *BasicLogger) Fatalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, false, args...)
	os.Exit(1)
}

func (l *BasicLogger) Fatalln(args ...interface{}) {
	l.log(CRITICAL, nil, true, args...)
	os.Exit(1)
}

func (l *BasicLogger) Panic(args ...interface{}) {
	l.log(CRITICAL, nil, false, args...)
	panic(fmt.Sprint(args...))
}

func (l *BasicLogger) Panicf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, false, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *BasicLogger) Panicln(args ...interface{}) {
	l.log(CRITICAL, nil, true, args...)
	panic(fmt.Sprintln(args...))
}

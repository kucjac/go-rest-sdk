package logger

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

type stdlogger struct{}

func (s *stdlogger) Print(args ...interface{})                 {}
func (s *stdlogger) Printf(format string, args ...interface{}) {}
func (s *stdlogger) Println(args ...interface{})               {}
func (s *stdlogger) Panic(args ...interface{})                 {}
func (s *stdlogger) Panicf(format string, args ...interface{}) {}
func (s *stdlogger) Panicln(args ...interface{})               {}
func (s *stdlogger) Fatal(args ...interface{})                 {}
func (s *stdlogger) Fatalf(format string, args ...interface{}) {}
func (s *stdlogger) Fatalln(args ...interface{})               {}

type leveledLogger struct{}

func (l *leveledLogger) Debugf(format string, args ...interface{})   {}
func (l *leveledLogger) Infof(format string, args ...interface{})    {}
func (l *leveledLogger) Warningf(format string, args ...interface{}) {}
func (l *leveledLogger) Errorf(format string, args ...interface{})   {}
func (l *leveledLogger) Fatalf(format string, args ...interface{})   {}
func (l *leveledLogger) Panicf(format string, args ...interface{})   {}
func (l *leveledLogger) Debug(args ...interface{})                   {}
func (l *leveledLogger) Info(args ...interface{})                    {}
func (l *leveledLogger) Warning(args ...interface{})                 {}
func (l *leveledLogger) Error(args ...interface{})                   {}
func (l *leveledLogger) Fatal(args ...interface{})                   {}
func (l *leveledLogger) Panic(args ...interface{})                   {}

type shortLeveledLogger struct{}

func (c *shortLeveledLogger) Debugf(format string, args ...interface{}) {}
func (c *shortLeveledLogger) Infof(format string, args ...interface{})  {}
func (c *shortLeveledLogger) Warnf(format string, args ...interface{})  {}
func (c *shortLeveledLogger) Errorf(format string, args ...interface{}) {}
func (c *shortLeveledLogger) Fatalf(format string, args ...interface{}) {}
func (c *shortLeveledLogger) Panicf(format string, args ...interface{}) {}
func (c *shortLeveledLogger) Debug(args ...interface{})                 {}
func (c *shortLeveledLogger) Info(args ...interface{})                  {}
func (c *shortLeveledLogger) Warn(args ...interface{})                  {}
func (c *shortLeveledLogger) Error(args ...interface{})                 {}
func (c *shortLeveledLogger) Fatal(args ...interface{})                 {}
func (c *shortLeveledLogger) Panic(args ...interface{})                 {}

type extendedLogger struct {
	leveledLogger
}

func (e *extendedLogger) Print(args ...interface{})                 {}
func (e *extendedLogger) Printf(format string, args ...interface{}) {}
func (e *extendedLogger) Println(args ...interface{})               {}
func (e *extendedLogger) Debugln(args ...interface{})               {}
func (e *extendedLogger) Infoln(args ...interface{})                {}
func (e *extendedLogger) Warningln(args ...interface{})             {}
func (e *extendedLogger) Errorln(args ...interface{})               {}
func (e *extendedLogger) Fatalln(args ...interface{})               {}
func (e *extendedLogger) Panicln(args ...interface{})               {}

type NonLogger struct{}

func TestNewLoggerWrapper(t *testing.T) {
	Convey("Subject: New Logger Wrapper.", t, func() {
		Convey("Having some loggers", func() {
			loggers := []interface{}{&stdlogger{}, &leveledLogger{}, &shortLeveledLogger{}, &extendedLogger{}}

			Convey(`If the logger implement possible interfaces,
			 wrapper handler should be returned`, func() {
				for i, logger := range loggers {
					wrapper, err := NewLoggerWrapper(logger)
					So(wrapper, ShouldHaveSameTypeAs, &LoggerWrapper{})
					So(err, ShouldBeNil)
					wrapper = MustGetLoggerWrapper(logger)
					So(wrapper, ShouldHaveSameTypeAs, &LoggerWrapper{})
					So(i+1, ShouldEqual, wrapper.currentLogger)
				}

				Convey("The loggers should enter its case ", func() {
					args := []interface{}{}
					format := "some format"
					for _, logger := range loggers {
						wrapper := MustGetLoggerWrapper(logger)
						wrapper.Print(args)
						wrapper.Printf(format, args)
						wrapper.Println(args)

						wrapper.Debug(args)
						wrapper.Debugf(format, args...)
						wrapper.Debugln(args)

						wrapper.Info(args)
						wrapper.Infof(format, args...)
						wrapper.Infoln(args)

						wrapper.Warning(args)
						wrapper.Warningf(format, args...)
						wrapper.Warningln(args)

						wrapper.Error(args)
						wrapper.Errorf(format, args)
						wrapper.Errorln(args)

						wrapper.Fatal(args)
						wrapper.Fatalf(format, args)
						wrapper.Fatalln(args)

						wrapper.Panic(args)
						wrapper.Panicf(format, args)
						wrapper.Panicln(args)
					}
				})
			})

			Convey(`If logger doesn't implement any known interface`, func() {
				unknownLogger := NonLogger{}
				wrapper, err := NewLoggerWrapper(unknownLogger)
				So(err, ShouldBeError)
				So(wrapper, ShouldBeNil)

				So(func() { MustGetLoggerWrapper(unknownLogger) }, ShouldPanic)
			})
		})
	})
}

func TestBuildLeveled(t *testing.T) {
	Convey("Having some logging parameters", t, func() {
		level := DEBUG
		format := "some format"
		arguments := []interface{}{"First", "Second"}

		Convey("Providing nil format should add level as first argument to args", func() {
			args := buildLeveled(level, nil, arguments...)
			So(args[0], ShouldEqual, fmt.Sprintf("%s: ", level))
		})

		Convey("buildLeveled with format should change the format string", func() {
			thisFormat := format
			args := buildLeveled(level, &thisFormat, arguments...)
			So(thisFormat, ShouldNotEqual, format)
			So(args, ShouldResemble, arguments)
		})
	})
}

func ExampleNewLoggerWrapper(t *testing.T) {
	// Having some logger (i.e. BasicLogger) that doesn't implement ExtendedLeveledLogger
	basic := NewBasicLogger(os.Stdout, "", 0)

	// In order to wrap it with LoggerWrapper use NewLoggerWrapper
	// or MustGetLoggerWrapper functions
	wrapper := MustGetLoggerWrapper(basic)

	// while having it wrapped by using LoggerWrapper we can use the methods of
	// ExtendedLeveledLogger
	wrapper.Println("Have fun")
	wrapper.Fatalln("This is the end...")

}

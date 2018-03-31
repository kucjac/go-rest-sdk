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

func TestNewGenericLogger(t *testing.T) {
	Convey("Subject: New Generic Logger.", t, func() {
		Convey("Having some loggers", func() {
			loggers := []interface{}{&stdlogger{}, &leveledLogger{}, &shortLeveledLogger{}, &extendedLogger{}}

			Convey(`If the logger implement possible interfaces,
			 generic handler should be returned`, func() {
				for i, logger := range loggers {
					generic, err := NewGenericLogger(logger)
					So(generic, ShouldHaveSameTypeAs, &GenericLogger{})
					So(err, ShouldBeNil)
					generic = MustGetGenericLogger(logger)
					So(generic, ShouldHaveSameTypeAs, &GenericLogger{})
					So(i+1, ShouldEqual, generic.currentLogger)
				}

				Convey("The loggers should enter its case ", func() {
					args := []interface{}{}
					format := "some format"
					for _, logger := range loggers {
						generic := MustGetGenericLogger(logger)
						generic.Print(args)
						generic.Printf(format, args)
						generic.Println(args)

						generic.Debug(args)
						generic.Debugf(format, args...)
						generic.Debugln(args)

						generic.Info(args)
						generic.Infof(format, args...)
						generic.Infoln(args)

						generic.Warning(args)
						generic.Warningf(format, args...)
						generic.Warningln(args)

						generic.Error(args)
						generic.Errorf(format, args)
						generic.Errorln(args)

						generic.Fatal(args)
						generic.Fatalf(format, args)
						generic.Fatalln(args)

						generic.Panic(args)
						generic.Panicf(format, args)
						generic.Panicln(args)
					}
				})
			})

			Convey(`If logger doesn't implement any known interface`, func() {
				unknownLogger := NonLogger{}
				generic, err := NewGenericLogger(unknownLogger)
				So(err, ShouldBeError)
				So(generic, ShouldBeNil)

				So(func() { MustGetGenericLogger(unknownLogger) }, ShouldPanic)
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

func ExampleNewGenericLogger(t *testing.T) {
	// Having some logger (i.e. BasicLogger)
	basic := NewBasicLogger(os.Stdout, "", 0)

	// It's worth to noting that BasicLogger doesn't implement ExtendedLeveledLogger

	// In order to wrap it with GenericLogger use NewGenericLogger
	// or MustGetGenericLogger functions
	generic := MustGetGenericLogger(basic)

	// while having it wrapped by using GenericLogger we can use the methods of
	// ExtendedLeveledLogger
	generic.Println("Have fun")
	generic.Fatalln("This is the end...")

}

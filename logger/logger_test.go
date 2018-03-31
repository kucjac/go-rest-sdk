package logger

import (
	"bytes"
	"fmt"
	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMessage(t *testing.T) {
	Convey("Subject: Logging Message.", t, func() {
		Convey("Message Method()", func() {
			Convey("Having a message that has nil fmt", func() {
				message := &Message{
					id:    1,
					level: INFO,
					fmt:   nil,
					args: []interface{}{
						"Some argument", "SecondArgument",
					},
				}
				Convey("The Message() should return a string as a contentation of args", func() {
					msgStr := message.Message()
					So(msgStr, ShouldEqual, fmt.Sprint(message.args...))
				})

			})
			Convey("Having a message with non-nil fmt", func() {
				format := "%s, %s"
				message := &Message{
					id:    2,
					level: DEBUG,
					fmt:   &format,
					args: []interface{}{
						"first", "second",
					},
				}
				Convey(`The Message() should return formatted 
					string of fmt and arguments`, func() {
					msgStr := message.Message()
					So(msgStr, ShouldEqual, fmt.Sprintf(*message.fmt, message.args...))
				})
				Convey("The String() method returns formatted Message", func() {
					str := message.String()
					So(str, ShouldEqual, fmt.Sprintf("%s|%04x: %s",
						message.level,
						message.id,
						message.getMessage()),
					)
				})
			})
		})
	})
}

func TestBasicLogger(t *testing.T) {
	Convey("Subject: BasicLogger tests", t, func() {
		Convey("NewBasicLogger() function creates a new BasicLogger", func() {
			var buf bytes.Buffer
			logger := NewBasicLogger(&buf, "", 0)
			So(logger, ShouldNotBeNil)
			So(logger, ShouldHaveSameTypeAs, &BasicLogger{})
			So(logger.stdLogger, ShouldNotBeNil)
			So(logger.level, ShouldEqual, DEBUG)

			Convey("SetLevel() method sets the logging level or logger", func() {
				logger.SetLevel(ERROR)
				So(logger.level, ShouldEqual, ERROR)
			})

			Convey("BasicLogger should implement ExtendedLeveledLogger interface{}", func() {
				So(logger, ShouldImplement, (*LeveledLogger)(nil))
				format := "%s-%s"
				args := []interface{}{"First", "Second"}
				msg := prepareMessage(logSequenceID, DEBUG, nil, args...)
				msgFmt := prepareMessage(logSequenceID, DEBUG, &format, args...)

				logger.SetLevel(DEBUG)

				buf.Reset()
				logger.Debug(args...)
				msg.id = logSequenceID
				So(buf.String(), ShouldEqual, fmtMsg(msg))

				buf.Reset()
				logger.Debugf(format, args...)
				msgFmt.id = logSequenceID
				So(buf.String(), ShouldEqual, fmtMsg(msgFmt))

				Convey("Setting Level above current would not allow logging lower level", func() {
					logger.SetLevel(INFO)

					buf.Reset()
					logger.Debug(args...)

					So(buf.String(), ShouldBeEmpty)
				})

				Convey("Testing Info methods", func() {
					logger.SetLevel(INFO)
					buf.Reset()
					logger.Info(args...)
					msg.level = INFO
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					logger.Infof(format, args...)
					msgFmt.level = INFO
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})

				Convey("Testing Warning methods", func() {
					logger.SetLevel(WARNING)
					buf.Reset()
					logger.Warning(args...)
					msg.level = WARNING
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					logger.Warningf(format, args...)
					msgFmt.level = WARNING
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})

				Convey("Testing Error methods", func() {
					logger.SetLevel(ERROR)
					buf.Reset()
					logger.Error(args...)
					msg.level = ERROR
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					logger.Errorf(format, args...)
					msgFmt.level = ERROR
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})

				Convey("Testing Panic methods", func() {
					logger.SetLevel(CRITICAL)
					buf.Reset()
					So(func() { logger.Panic(args...) }, ShouldPanic)
					msg.level = CRITICAL
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					msgFmt.level = CRITICAL
					So(func() { logger.Panicf(format, args...) }, ShouldPanic)
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})

				Convey("Testing Print methods, no matter what level of logging is set", func() {
					logger.SetLevel(CRITICAL)
					buf.Reset()

					logger.Print(args...)
					msg.level = PRINT
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					logger.Printf(format, args...)
					msgFmt.level = PRINT
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})

				Convey("Testing Fatal methods", func() {
					logger.SetLevel(INFO)
					buf.Reset()

					panicText := "os.Exit called"
					fakeExit := func(int) {
						panic(panicText)
					}
					patch := monkey.Patch(os.Exit, fakeExit)
					defer patch.Unpatch()

					So(func() { logger.Fatal(args...) }, ShouldPanicWith, panicText)
					msg.level = CRITICAL
					msg.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msg))

					buf.Reset()
					So(func() { logger.Fatalf(format, args...) }, ShouldPanicWith, panicText)
					msgFmt.level = CRITICAL
					msgFmt.id = logSequenceID
					So(buf.String(), ShouldEqual, fmtMsg(msgFmt))
				})
			})
		})
	})
}

func prepareMessage(id uint64, level Level, fmt *string, args ...interface{}) *Message {
	return &Message{id: id, level: level, fmt: fmt, args: args}
}

func fmtMsg(msg *Message) string {
	return fmt.Sprintf("%s\n", msg.String())
}

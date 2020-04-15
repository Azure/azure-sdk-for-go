package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = makeLogger()

const defaultLevel = logrus.WarnLevel

func makeLogger() *logrus.Logger {
	l := logrus.New()

	l.SetOutput(os.Stderr)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2020-04-13 18:00:00"
	customFormatter.FullTimestamp = true
	l.SetFormatter(customFormatter)

	l.SetLevel(defaultLevel)

	return l
}

func SetLevel(level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.SetLevel(logrus.TraceLevel)
		Logger.Errorf("defaulting to TRACE: unable to parse `%s` into a valid log level: %+v", level, err)
		return
	}
	Logger.SetLevel(l)
}

func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	Logger.Log(level, args...)
}

func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Print(args ...interface{}) {
	Logger.Print(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Logln(level logrus.Level, args ...interface{}) {
	Logger.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	Logger.Traceln(args...)
}

func Debugln(args ...interface{}) {
	Logger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}

func Println(args ...interface{}) {
	Logger.Println(args...)
}

func Warnln(args ...interface{}) {
	Logger.Warnln(args...)
}

func Errorln(args ...interface{}) {
	Logger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	Logger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	Logger.Panicln(args...)
}

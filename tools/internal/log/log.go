// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
)

// Logger ...
var Logger = makeLogger()

const defaultLevel = logrus.InfoLevel

func makeLogger() *logrus.Logger {
	l := logrus.New()

	l.SetOutput(os.Stderr)

	exampleFormatter := &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
	l.SetFormatter(exampleFormatter)

	l.SetLevel(defaultLevel)

	return l
}

// SetLevel sets the level of log
func SetLevel(level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.SetLevel(logrus.TraceLevel)
		Logger.Errorf("defaulting to TRACE: unable to parse `%s` into a valid log level: %+v", level, err)
		return
	}
	Logger.SetLevel(l)
}

// Tracef ...
func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Printf ...
func Printf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panicf ...
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// Log ...
func Log(level logrus.Level, args ...interface{}) {
	Logger.Log(level, args...)
}

// Trace ...
func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

// Debug ...
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Info ...
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Print ...
func Print(args ...interface{}) {
	Logger.Print(args...)
}

// Warn ...
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Error ...
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Fatal ...
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Panic ...
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Logln ...
func Logln(level logrus.Level, args ...interface{}) {
	Logger.Logln(level, args...)
}

// Traceln ...
func Traceln(args ...interface{}) {
	Logger.Traceln(args...)
}

// Debugln ...
func Debugln(args ...interface{}) {
	Logger.Debugln(args...)
}

// Infoln ...
func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}

// Println ...
func Println(args ...interface{}) {
	Logger.Println(args...)
}

// Warnln ...
func Warnln(args ...interface{}) {
	Logger.Warnln(args...)
}

// Errorln ...
func Errorln(args ...interface{}) {
	Logger.Errorln(args...)
}

// Fatalln ...
func Fatalln(args ...interface{}) {
	Logger.Fatalln(args...)
}

// Panicln ...
func Panicln(args ...interface{}) {
	Logger.Panicln(args...)
}

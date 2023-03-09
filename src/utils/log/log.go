// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package log initialize logurs
package log

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Source: https://github.com/sirupsen/logrus/blob/master/exported.go
var std = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.JSONFormatter),
	Hooks:     make(logrus.LevelHooks),
	// TODO(daniel): Add a log --loglevel=[DEBUG|INFO|...] to control this logging level
	Level:    logrus.InfoLevel,
	ExitFunc: os.Exit,
}

// Refer: https://go.dev/play/p/q0kyZvvbT0C
func logger() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("while logger invoke, failed to get context info")
	}

	filename := file + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return std.WithField("file", filename).WithField("function", fn)
}

// WithError creates an entry from the standard logger and adds an error to it,
// using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return logger().WithField(logrus.ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return logger().WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return logger().WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger().WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *logrus.Entry {
	return logger().WithTime(t)
}

// Trace logs a message at logrus. Level Trace on the standard logger.
func Trace(args ...interface{}) {
	logger().Trace(args...)
}

// Debug logs a message at logrus. Level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger().Debug(args...)
}

// Print logs a message at logrus. Level Info on the standard logger.
func Print(args ...interface{}) {
	logger().Print(args...)
}

// Info logs a message at logrus. Level Info on the standard logger.
func Info(args ...interface{}) {
	logger().Info(args...)
}

// Warn logs a message at logrus. Level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger().Warn(args...)
}

// Warning logs a message at logrus. Level Warn on the standard logger.
func Warning(args ...interface{}) {
	logger().Warning(args...)
}

// Error logs a message at logrus. Level Error on the standard logger.
func Error(args ...interface{}) {
	logger().Error(args...)
}

// Panic logs a message at logrus. Level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger().Panic(args...)
}

// Fatal logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger().Fatal(args...)
}

// Tracef logs a message at logrus. Level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	logger().Tracef(format, args...)
}

// Debugf logs a message at logrus. Level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger().Debugf(format, args...)
}

// Printf logs a message at logrus. Level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logger().Printf(format, args...)
}

// Infof logs a message at logrus. Level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger().Infof(format, args...)
}

// Warnf logs a message at logrus. Level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger().Warnf(format, args...)
}

// Warningf logs a message at logrus. Level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logger().Warningf(format, args...)
}

// Errorf logs a message at logrus. Level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger().Errorf(format, args...)
}

// Panicf logs a message at logrus. Level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logger().Panicf(format, args...)
}

// Fatalf logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logger().Fatalf(format, args...)
}

// Traceln logs a message at logrus. Level Trace on the standard logger.
func Traceln(args ...interface{}) {
	logger().Traceln(args...)
}

// Debugln logs a message at logrus. Level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger().Debugln(args...)
}

// Println logs a message at logrus. Level Info on the standard logger.
func Println(args ...interface{}) {
	logger().Println(args...)
}

// Infoln logs a message at logrus. Level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger().Infoln(args...)
}

// Warnln logs a message at logrus. Level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger().Warnln(args...)
}

// Warningln logs a message at logrus. Level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logger().Warningln(args...)
}

// Errorln logs a message at logrus. Level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger().Errorln(args...)
}

// Panicln logs a message at logrus. Level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger().Panicln(args...)
}

// Fatalln logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger().Fatalln(args...)
}

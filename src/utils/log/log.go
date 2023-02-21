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
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// source: https://github.com/sirupsen/logrus/blob/master/exported.go
func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	// TODO(daniel): Add a log --loglevel=[DEBUG|INFO|...] to control this logging level
}

// link: https://go.dev/play/p/q0kyZvvbT0C
func logger() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return logrus.WithField("file", filename).WithField("function", fn)
}

func StandardLogger() *logrus.Logger {
	return logger().Logger
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logger().Logger.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logger().Logger.SetFormatter(formatter)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	logger().Logger.SetReportCaller(include)
}

// SetLevel sets the standard logger logrus. Level.
func SetLevel(level logrus.Level) {
	logger().Logger.SetLevel(level)
}

// GetLevel returns the standard logger logrus. Level.
func GetLevel() logrus.Level {
	return logger().Logger.GetLevel()
}

// IsLevelEnabled checks if the log logrus. Level of the standard logger is greater than the logrus. Level param
func IsLevelEnabled(level logrus.Level) bool {
	return logger().Logger.IsLevelEnabled(level)
}

// AddHook adds a logrus. Hook to the standard logger logrus. Hooks.
func AddHook(hook logrus.Hook) {
	logger().Logger.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it,
// using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return logger().Logger.WithField(logrus.ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return logger().Logger.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return logger().Logger.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger().Logger.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *logrus.Entry {
	return logger().Logger.WithTime(t)
}

// Trace logs a message at logrus. Level Trace on the standard logger.
func Trace(args ...interface{}) {
	logger().Logger.Trace(args...)
}

// Debug logs a message at logrus. Level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger().Logger.Debug(args...)
}

// Print logs a message at logrus. Level Info on the standard logger.
func Print(args ...interface{}) {
	logger().Logger.Print(args...)
}

// Info logs a message at logrus. Level Info on the standard logger.
func Info(args ...interface{}) {
	logger().Logger.Info(args...)
}

// Warn logs a message at logrus. Level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger().Logger.Warn(args...)
}

// Warning logs a message at logrus. Level Warn on the standard logger.
func Warning(args ...interface{}) {
	logger().Logger.Warning(args...)
}

// Error logs a message at logrus. Level Error on the standard logger.
func Error(args ...interface{}) {
	logger().Logger.Error(args...)
}

// Panic logs a message at logrus. Level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger().Logger.Panic(args...)
}

// Fatal logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger().Logger.Fatal(args...)
}

// TraceFn logs a message from a func at logrus. Level Trace on the standard logger.
func TraceFn(fn logrus.LogFunction) {
	logger().Logger.TraceFn(fn)
}

// DebugFn logs a message from a func at logrus. Level Debug on the standard logger.
func DebugFn(fn logrus.LogFunction) {
	logger().Logger.DebugFn(fn)
}

// PrintFn logs a message from a func at logrus. Level Info on the standard logger.
func PrintFn(fn logrus.LogFunction) {
	logger().Logger.PrintFn(fn)
}

// InfoFn logs a message from a func at logrus. Level Info on the standard logger.
func InfoFn(fn logrus.LogFunction) {
	logger().Logger.InfoFn(fn)
}

// WarnFn logs a message from a func at logrus. Level Warn on the standard logger.
func WarnFn(fn logrus.LogFunction) {
	logger().Logger.WarnFn(fn)
}

// WarningFn logs a message from a func at logrus. Level Warn on the standard logger.
func WarningFn(fn logrus.LogFunction) {
	logger().Logger.WarningFn(fn)
}

// ErrorFn logs a message from a func at logrus. Level Error on the standard logger.
func ErrorFn(fn logrus.LogFunction) {
	logger().Logger.ErrorFn(fn)
}

// PanicFn logs a message from a func at logrus. Level Panic on the standard logger.
func PanicFn(fn logrus.LogFunction) {
	logger().Logger.PanicFn(fn)
}

// FatalFn logs a message from a func at logrus.
// Level Fatal on the standard logger then the process will exit with status set to 1.
func FatalFn(fn logrus.LogFunction) {
	logger().Logger.FatalFn(fn)
}

// Tracef logs a message at logrus. Level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	logger().Logger.Tracef(format, args...)
}

// Debugf logs a message at logrus. Level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger().Logger.Debugf(format, args...)
}

// Printf logs a message at logrus. Level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logger().Logger.Printf(format, args...)
}

// Infof logs a message at logrus. Level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger().Logger.Infof(format, args...)
}

// Warnf logs a message at logrus. Level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger().Logger.Warnf(format, args...)
}

// Warningf logs a message at logrus. Level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logger().Logger.Warningf(format, args...)
}

// Errorf logs a message at logrus. Level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger().Logger.Errorf(format, args...)
}

// Panicf logs a message at logrus. Level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logger().Logger.Panicf(format, args...)
}

// Fatalf logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logger().Logger.Fatalf(format, args...)
}

// Traceln logs a message at logrus. Level Trace on the standard logger.
func Traceln(args ...interface{}) {
	logger().Logger.Traceln(args...)
}

// Debugln logs a message at logrus. Level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger().Logger.Debugln(args...)
}

// Println logs a message at logrus. Level Info on the standard logger.
func Println(args ...interface{}) {
	//logger().Logger.Println(args...)
	logger().Println(args...)
}

// Infoln logs a message at logrus. Level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger().Logger.Infoln(args...)
}

// Warnln logs a message at logrus. Level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger().Logger.Warnln(args...)
}

// Warningln logs a message at logrus. Level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logger().Logger.Warningln(args...)
}

// Errorln logs a message at logrus. Level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger().Logger.Errorln(args...)
}

// Panicln logs a message at logrus. Level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger().Logger.Panicln(args...)
}

// Fatalln logs a message at logrus. Level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger().Logger.Fatalln(args...)
}

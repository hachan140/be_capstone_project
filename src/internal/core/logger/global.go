package logger

import (
	"be-capstone-project/src/internal/core/logger/internal"
	"context"
	"fmt"
)

var global internal.Logger

func init() {
	var err error
	if global, err = internal.NewDefaultLogger(&internal.Options{CallerSkip: 3}); err != nil {
		panic(fmt.Errorf("init global logger error [%v]", err))
	}
}

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger internal.Logger) {
	global = logger
}

// GetGlobal Get global logger instance
func GetGlobal() internal.Logger {
	return global
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	global.Info(args...)
}

// Infof uses fmt.Sprintf to log a template message.
func Infof(msgFormat string, args ...interface{}) {
	global.Infof(msgFormat, args...)
}

// Infow uses fmt.Sprintf to log a template message with extra context value.
func Infow(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Infow(keysAndValues, msgFormat, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	global.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a template message.
func Debugf(msgFormat string, args ...interface{}) {
	global.Debugf(msgFormat, args...)
}

// Debugw uses fmt.Sprintf to log a template message with extra context value.
func Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Debugw(keysAndValues, msgFormat, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	global.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a template message.
func Warnf(msgFormat string, args ...interface{}) {
	global.Warnf(msgFormat, args...)
}

// Warnw uses fmt.Sprintf to log a template message with extra context value.
func Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Warnw(keysAndValues, msgFormat, args...)
}

// Error uses fmt.Sprint to construct and log a message. Error level log with a stack trace
func Error(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Error(args...)
}

// Errorf uses fmt.Sprintf to log a template message. ErrorCtx level log with a stack trace.
func Errorf(msgFormat string, args ...interface{}) {
	global.Errorf(msgFormat, args...)
}

// Errorw uses fmt.Sprintf to log a template message with extra context value. ErrorCtx level log with a stack trace.
func Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Errorw(keysAndValues, msgFormat, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
func Fatalf(msgFormat string, args ...interface{}) {
	global.Fatalf(msgFormat, args...)
}

// Fatalw uses fmt.Sprintf to log a template message with extra context value, then calls os.Exit.
func Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Fatalw(keysAndValues, msgFormat, args...)
}

func FatalfIf(cond bool, fmt string, args ...interface{}) {
	global.FatalfIf(cond, fmt, args)
}

// FatalIfError if error and print error then calls os.Exit
func FatalIfError(err error) {
	global.FatalIfError(err)
}

func Write(p []byte) (n int, err error) {
	return global.Write(p)
}

// Copyright 2018 xgfone
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

package glog

import "io"

// LoggerWithoutError is equal to Logger, but not returning the error.
type LoggerWithoutError interface {
	Depth(stackDepth int) LoggerWithoutError
	Level(level Level) LoggerWithoutError
	Encoder(encoder Encoder) LoggerWithoutError
	Cxt(ctxs ...interface{}) LoggerWithoutError

	// Writer is the convenient function of GetEncoder().Writer().
	Writer() io.Writer
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder

	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type loggerWithoutError struct {
	Logger
}

func newLoggerWithoutError(logger Logger, depth bool) LoggerWithoutError {
	if depth {
		return loggerWithoutError{Logger: logger.Depth(logger.GetDepth() + 1)}
	}
	return loggerWithoutError{Logger: logger}
}

// ToLoggerWithoutError converts the Logger to LoggerWithoutError.
//
// If logger is missing, it will use the global logger by default.
func ToLoggerWithoutError(logger ...Logger) LoggerWithoutError {
	_logger := GetGlobalLogger()
	if len(logger) > 0 && logger[0] != nil {
		_logger = logger[0]
	}
	return newLoggerWithoutError(_logger, true)
}

// ToLogger converts the LoggerWithoutError to Logger.
//
// Notice: LoggerWithoutError must be the built-in implementation
// returned by ToLoggerWithoutError.
func ToLogger(logger LoggerWithoutError) Logger {
	return logger.(loggerWithoutError).Logger.Depth(logger.GetDepth() - 1)
}

func (l loggerWithoutError) Depth(stackDepth int) LoggerWithoutError {
	return newLoggerWithoutError(l.Logger.Depth(stackDepth), false)
}

func (l loggerWithoutError) Level(level Level) LoggerWithoutError {
	return newLoggerWithoutError(l.Logger.Level(level), false)
}

func (l loggerWithoutError) Encoder(encoder Encoder) LoggerWithoutError {
	return newLoggerWithoutError(l.Logger.Encoder(encoder), false)
}

func (l loggerWithoutError) Cxt(ctxs ...interface{}) LoggerWithoutError {
	return newLoggerWithoutError(l.Logger.Cxt(ctxs...), false)
}

func (l loggerWithoutError) Trace(msg string, args ...interface{}) {
	l.Logger.Trace(msg, args...)
}

func (l loggerWithoutError) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(msg, args...)
}

func (l loggerWithoutError) Info(msg string, args ...interface{}) {
	l.Logger.Info(msg, args...)
}

func (l loggerWithoutError) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(msg, args...)
}

func (l loggerWithoutError) Error(msg string, args ...interface{}) {
	l.Logger.Error(msg, args...)
}

func (l loggerWithoutError) Panic(msg string, args ...interface{}) {
	l.Logger.Panic(msg, args...)
}

func (l loggerWithoutError) Fatal(msg string, args ...interface{}) {
	l.Logger.Fatal(msg, args...)
}

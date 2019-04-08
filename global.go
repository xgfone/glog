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

package logger

import (
	"io"
	"os"
)

var last Logger
var root Logger

func init() {
	var defaultConf = EncoderConfig{IsLevel: true, IsTime: true}
	SetGlobalLogger(New(FmtTextEncoder(os.Stdout, defaultConf)))
}

// SetGlobalLogger sets the global logger to log.
//
// If log is nil, it will do nothing.
func SetGlobalLogger(log Logger) {
	if log != nil {
		last = log
		root = log.WithDepth(log.GetDepth() + 1)
	}
}

// GetGlobalLogger returns the global logger.
//
// The name of the default global logger is "root".
func GetGlobalLogger() Logger {
	return last
}

// WithName returns a new logger with the name.
func WithName(name string) Logger {
	return root.WithName(name)
}

// WithLevel returns a new logger with the level.
func WithLevel(level Level) Logger {
	return root.WithLevel(level)
}

// WithEncoder returns a new logger with the encoder.
func WithEncoder(encoder Encoder) Logger {
	return root.WithEncoder(encoder)
}

// WithCtx returns a new logger with the contexts.
func WithCtx(ctxs ...interface{}) Logger {
	return root.WithCxt(ctxs...)
}

// WithDepth returns a new logger with the caller depth.
func WithDepth(depth int) Logger {
	return root.WithDepth(depth)
}

// GetWriter returns the underlying writer of the global logger.
func GetWriter() Writer {
	return root.GetEncoder().Writer()
}

// GetName returns the name of the current global logger.
func GetName() string {
	return root.GetName()
}

// GetDepth returns the caller depth of the global logger.
func GetDepth() int {
	return root.GetDepth()
}

// GetLevel returns the level of the global logger.
func GetLevel() Level {
	return root.GetLevel()
}

// GetEncoder returns the encoder of the global logger.
func GetEncoder() Encoder {
	return root.GetEncoder()
}

// SetName sets the name of the current global logger.
func SetName(name string) {
	root.SetName(name)
}

// SetDepth sets the caller depth of the global logger.
func SetDepth(depth int) {
	root.SetDepth(depth)
}

// SetLevel sets the level of the global logger.
func SetLevel(level Level) {
	root.SetLevel(level)
}

// SetEncoder sets the encoder of the global logger.
func SetEncoder(encoder Encoder) {
	root.SetEncoder(encoder)
}

// Trace fires a TRACE log.
//
// The meaning of arguments is in accordance with the encoder.
func Trace(msg string, args ...interface{}) error {
	return root.Trace(msg, args...)
}

// Debug fires a DEBUG log.
//
// The meaning of arguments is in accordance with the encoder.
func Debug(msg string, args ...interface{}) error {
	return root.Debug(msg, args...)
}

// Info fires a INFO log.
//
// The meaning of arguments is in accordance with the encoder.
func Info(msg string, args ...interface{}) error {
	return root.Info(msg, args...)
}

// Warn fires a WARN log.
//
// The meaning of arguments is in accordance with the encoder.
func Warn(msg string, args ...interface{}) error {
	return root.Warn(msg, args...)
}

// Error fires a ERROR log.
//
// The meaning of arguments is in accordance with the encoder.
func Error(msg string, args ...interface{}) error {
	return root.Error(msg, args...)
}

// Panic fires a PANIC log then panic.
//
// The meaning of arguments is in accordance with the encoder.
func Panic(msg string, args ...interface{}) error {
	return root.Panic(msg, args...)
}

// Fatal fires a FATAL log then terminates the program.
//
// The meaning of arguments is in accordance with the encoder.
func Fatal(msg string, args ...interface{}) error {
	return root.Fatal(msg, args...)
}

// SimpleLogger returns a new Logger with the level and the writer will use
// os.Stdout if filepath is "", or use the file based on SizedRotatingFileWriter.
//
// Notice: the file size is 1GB and the number is 30 by default. But you can
// change it by passing the last two parameters as follow.
//
//     FileLogger(level, filepath, 1024*1024*1024) // 1GB each file
//     FileLogger(level, filepath, 1024*1024*1024, 10) // 1GB each file and 10 files
//
func SimpleLogger(level, filepath string, args ...int) (Logger, io.Closer, error) {
	logger := GetGlobalLogger().WithLevel(NameToLevel(level))
	if filepath == "" {
		return logger, nothingCloser{}, nil
	}

	size := 1024 * 1024 * 1024
	count := 30
	switch len(args) {
	case 1:
		size = args[0]
	case 2:
		size = args[0]
		count = args[1]
	}

	file, closer, err := SizedRotatingFileWriter(filepath, size, count)
	if err != nil {
		return nil, nil, err
	}
	logger.GetEncoder().ResetWriter(file)
	return logger, closer, nil
}

type nothingCloser struct{}

func (nc nothingCloser) Close() error { return nil }

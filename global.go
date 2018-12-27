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
//
// Notice: for the global logger, it must be the builtin implementation.
func SetGlobalLogger(log Logger) {
	if log != nil {
		last = log
		root = log.Depth(log.GetDepth() + 1)
	}
}

// GetGlobalLogger returns the global logger.
func GetGlobalLogger() Logger {
	return last
}

// WithLevel returns a new logger with the level.
//
// Since Level is the level type, we use WithLevel as the function name.
func WithLevel(level Level) Logger {
	return root.Level(level)
}

// WithEncoder returns a new logger with the encoder.
//
// Since Encoder is the encoder type, we use WithEncoder as the function name.
func WithEncoder(encoder Encoder) Logger {
	return root.Encoder(encoder)
}

// WithCtx returns a new logger with the contexts.
//
// In order to keep consistent with WithLevel and WithEncoder,
// we use WithCtx, not Ctx.
func WithCtx(ctxs ...interface{}) Logger {
	return root.Cxt(ctxs...)
}

// WithDepth returns a new logger with the caller depth.
func WithDepth(depth int) Logger {
	return root.Depth(depth)
}

// GetWriter returns the underlying writer of the global logger.
func GetWriter() Writer {
	return root.Writer()
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
//     FileLogger(level, filepath, 2*1024*1024*1024) // 2GB each file
//     FileLogger(level, filepath, 2*1024*1024*1024, 10) // 2GB each file and 10 files
//
func SimpleLogger(level, filepath string, args ...int) (Logger, io.Closer, error) {
	logger := GetGlobalLogger().Level(NameToLevel(level))
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

	file, err := SizedRotatingFileWriter(filepath, size, count)
	if err != nil {
		return nil, nil, err
	}
	logger.GetEncoder().ResetWriter(file)
	return logger, file, nil
}

type nothingCloser struct{}

func (nc nothingCloser) Close() error { return nil }

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

package miss

import "os"

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
	if last != nil {
		return last
	}
	return root
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

// Trace fires a TRACE log.
//
// The meaning of arguments is in accordance with the encoder.
func Trace(msg string, args ...interface{}) error {
	return root.Trace(msg, args)
}

// Debug fires a DEBUG log.
//
// The meaning of arguments is in accordance with the encoder.
func Debug(msg string, args ...interface{}) error {
	return root.Debug(msg, args)
}

// Info fires a INFO log.
//
// The meaning of arguments is in accordance with the encoder.
func Info(msg string, args ...interface{}) error {
	return root.Info(msg, args)
}

// Warn fires a WARN log.
//
// The meaning of arguments is in accordance with the encoder.
func Warn(msg string, args ...interface{}) error {
	return root.Warn(msg, args)
}

// Error fires a ERROR log.
//
// The meaning of arguments is in accordance with the encoder.
func Error(msg string, args ...interface{}) error {
	return root.Error(msg, args)
}

// Panic fires a PANIC log then panic.
//
// The meaning of arguments is in accordance with the encoder.
func Panic(msg string, args ...interface{}) error {
	return root.Panic(msg, args)
}

// Fatal fires a FATAL log then terminates the program.
//
// The meaning of arguments is in accordance with the encoder.
func Fatal(msg string, args ...interface{}) error {
	return root.Fatal(msg, args)
}

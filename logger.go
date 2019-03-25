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
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// ErrPanic will be used when firing a PANIC level log.
var ErrPanic = fmt.Errorf("the panic level log")

// DefaultLoggerDepth is the depth for the default implementing logger.
const DefaultLoggerDepth = 2

// Logger is an immutable logger interface.
type Logger interface {
	// Some methods to return a new Logger with the argument and the current state.
	Level(level Level) Logger
	Encoder(encoder Encoder) Logger
	Cxt(ctxs ...interface{}) Logger
	// stackDepth is the calling depth of the logger, which will be passed to
	// the encoder. The default depth is the global variable DefaultLoggerDepth
	// for the new Logger.
	//
	// It should be used typically when you wrap the logger. For example,
	//
	//   _logger := logger.New(logger.KvTextEncoder(os.Stdout))
	//   _logger = _logger.Depth(_logger.GetDepth() + 1)
	//
	//   func Debug(m string, args ...interface{}) { _logger.Debug(m, args...) }
	//   func Info(m string, args ...interface{}) { _logger.Debug(m, args...) }
	//   func Warn(m string, args ...interface{}) { _logger.Debug(m, args...) }
	//   ...
	//
	Depth(stackDepth int) Logger

	// Some methods to return the inner state.
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder
	Writer() io.Writer // [DEPRECATED] Please use GetEncoder().Writer().

	// Some Logs based on the level.
	Trace(msg string, args ...interface{}) error
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Panic(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{}) error
}

// Setter is an interface of the setter of Logger, which modify
// the inner state of Logger, not returning a new Logger.
//
// This interface indicates a mutable Logger.
//
// Notice: the builtin implementation has implemented this interface.
type Setter interface {
	SetDepth(depth int)
	SetLevel(level Level) // It should be thread-safe.
	SetEncoder(encoder Encoder)
}

type logger struct {
	enc Encoder
	lvl Level
	ctx []interface{}

	depth int
}

// New returns a new Logger.
func New(encoder Encoder) Logger {
	return &logger{
		lvl: LvlTrace,
		enc: encoder,
		ctx: make([]interface{}, 0),

		depth: DefaultLoggerDepth,
	}
}

func newLogger(l *logger) *logger {
	return &logger{
		enc: l.enc,
		ctx: l.ctx,
		lvl: l.GetLevel(),

		depth: l.depth,
	}
}

func (l *logger) Writer() io.Writer {
	return l.enc.Writer()
}

func (l *logger) GetDepth() int {
	return l.depth
}

func (l *logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&l.lvl)))
}

func (l *logger) GetEncoder() Encoder {
	return l.enc
}

func (l *logger) SetDepth(depth int) {
	l.depth = depth
}

func (l *logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&l.lvl), int32(level))
}

func (l *logger) SetEncoder(encoder Encoder) {
	l.enc = encoder
}

func (l *logger) Depth(depth int) Logger {
	log := newLogger(l)
	log.depth = depth
	return log
}

func (l *logger) Level(level Level) Logger {
	log := newLogger(l)
	log.lvl = level
	return log
}

func (l *logger) Encoder(encoder Encoder) Logger {
	log := newLogger(l)
	log.enc = encoder
	return log
}

func (l *logger) Cxt(ctxs ...interface{}) Logger {
	log := newLogger(l)
	log.ctx = append(l.ctx, ctxs...)
	return log
}

func (l *logger) log(level Level, msg string, args []interface{}) (err error) {
	if level < l.GetLevel() {
		return nil
	}
	err = l.enc.Encode(l.depth, level, msg, args, l.ctx)

	switch level {
	case LvlPanic:
		panic(ErrPanic)
	case LvlFatal:
		os.Exit(1)
	}

	return
}

func (l *logger) Trace(msg string, args ...interface{}) error {
	return l.log(LvlTrace, msg, args)
}

func (l *logger) Debug(msg string, args ...interface{}) error {
	return l.log(LvlDebug, msg, args)
}

func (l *logger) Info(msg string, args ...interface{}) error {
	return l.log(LvlInfo, msg, args)
}

func (l *logger) Warn(msg string, args ...interface{}) error {
	return l.log(LvlWarn, msg, args)
}

func (l *logger) Error(msg string, args ...interface{}) error {
	return l.log(LvlError, msg, args)
}

func (l *logger) Panic(msg string, args ...interface{}) error {
	return l.log(LvlPanic, msg, args)
}

func (l *logger) Fatal(msg string, args ...interface{}) error {
	return l.log(LvlFatal, msg, args)
}

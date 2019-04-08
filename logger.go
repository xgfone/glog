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

// LogGetter is an interface to return the inner information of Logger.
type LogGetter interface {
	GetName() string
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder
}

// LogSetter is an interface to modify the inner information of Logger.
type LogSetter interface {
	SetName(name string)
	SetDepth(depth int)
	SetLevel(level Level) // It should be thread-safe.
	SetEncoder(encoder Encoder)
}

// LogWither is an interface to return a new Logger based on the current logger
// with the new argument.
type LogWither interface {
	WithName(name string) Logger
	WithLevel(level Level) Logger
	WithEncoder(encoder Encoder) Logger
	WithCxt(ctxs ...interface{}) Logger
	WithDepth(stackDepth int) Logger
}

// LogOutputter is an interface to emit the log.
type LogOutputter interface {
	Trace(msg string, args ...interface{}) error
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Panic(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{}) error
}

// Logger is an compositive logger interface.
type Logger interface {
	LogGetter
	LogSetter
	LogWither
	LogOutputter
}

// WriterLogger is an Logger interfce with the method Writer().
type WriterLogger interface {
	Logger

	Writer() io.Writer
}

// ToWriterLogger converts Logger to WriterLogger.
//
// if logger is nil, use the default global logger.
func ToWriterLogger(logger ...Logger) WriterLogger {
	log := GetGlobalLogger()
	if len(logger) > 0 && logger[0] != nil {
		log = logger[0]
	}
	return wlogger{log}
}

type wlogger struct {
	Logger
}

func (wl wlogger) Writer() io.Writer {
	return wl.GetEncoder().Writer()
}

type logger struct {
	enc Encoder
	lvl Level
	ctx []interface{}

	name  string
	depth int
}

// New returns a new Logger.
//
// The name is "root" by default.
func New(encoder Encoder) Logger {
	return &logger{
		lvl: LvlTrace,
		enc: encoder,
		ctx: make([]interface{}, 0),

		name:  "root",
		depth: DefaultLoggerDepth,
	}
}

func newLogger(l *logger) *logger {
	return &logger{
		enc: l.enc,
		ctx: l.ctx,
		lvl: l.GetLevel(),

		name:  l.name,
		depth: l.depth,
	}
}

func (l *logger) GetName() string {
	return l.name
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

func (l *logger) SetName(name string) {
	l.name = name
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

func (l *logger) WithName(name string) Logger {
	log := newLogger(l)
	log.name = name
	return log
}

func (l *logger) WithDepth(depth int) Logger {
	log := newLogger(l)
	log.depth = depth
	return log
}

func (l *logger) WithLevel(level Level) Logger {
	log := newLogger(l)
	log.lvl = level
	return log
}

func (l *logger) WithEncoder(encoder Encoder) Logger {
	log := newLogger(l)
	log.enc = encoder
	return log
}

func (l *logger) WithCxt(ctxs ...interface{}) Logger {
	log := newLogger(l)
	log.ctx = append(l.ctx, ctxs...)
	return log
}

func (l *logger) log(lvl Level, msg string, args []interface{}) (err error) {
	if lvl < l.GetLevel() {
		return nil
	}

	err = l.enc.Encode(Record{
		Lvl:   lvl,
		Msg:   msg,
		Args:  args,
		Ctxs:  l.ctx,
		Name:  l.name,
		Depth: l.depth,
	})

	switch lvl {
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

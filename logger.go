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

import (
	"fmt"
	"io"
	"os"
)

// ErrPanic will be used when firing a PANIC level log.
var ErrPanic = fmt.Errorf("the panic level log")

// DefaultLoggerDepth is the depth for the default implementing logger.
const DefaultLoggerDepth = 2

// Logger is a logger interface.
type Logger interface {
	// Depth returns a new Logger with the stack depth.
	//
	// stackDepth is the calling depth of the logger, which will be passed to
	// the encoder. The default depth is the global variable DefaultLoggerDepth
	// for the new Logger.
	//
	// It should be used typically when you wrap the logger. For example,
	//
	//   logger := miss.New(miss.KvTextEncoder(os.Stdout))
	//   logger = logger.Depth(logger.GetDepth() + 1)
	//
	//   func Debug(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   func Info(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   func Warn(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   ...
	//
	Depth(stackDepth int) Logger

	// Level returns a new Logger with the new level.
	Level(level Level) Logger

	// Encoder returns a new logger with the new encoder.
	Encoder(encoder Encoder) Logger

	// Ctx returns a new logger with the new contexts.
	Cxt(ctxs ...interface{}) Logger

	// Writer is the convenient function of GetEncoder().Writer().
	Writer() io.Writer
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder

	Trace(msg string, args ...interface{}) error
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Panic(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{}) error
}

type logger struct {
	enc Encoder
	lvl Level
	ctx []interface{}

	depth int
}

// New returns a new Logger.
func New(encoder Encoder) Logger {
	return logger{
		lvl: TRACE,
		enc: encoder,
		ctx: make([]interface{}, 0),

		depth: DefaultLoggerDepth,
	}
}

func newLogger(l logger) logger {
	return logger{
		enc: l.enc,
		ctx: l.ctx,
		lvl: l.lvl,

		depth: l.depth,
	}
}

func (l logger) Writer() io.Writer {
	return l.enc.Writer()
}

func (l logger) GetDepth() int {
	return l.depth
}

func (l logger) GetLevel() Level {
	return l.lvl
}

func (l logger) GetEncoder() Encoder {
	return l.enc
}

func (l logger) Depth(depth int) Logger {
	log := newLogger(l)
	log.depth = depth
	return log
}

func (l logger) Level(level Level) Logger {
	log := newLogger(l)
	log.lvl = level
	return log
}

func (l logger) Encoder(encoder Encoder) Logger {
	log := newLogger(l)
	log.enc = encoder
	return log
}

func (l logger) Cxt(ctxs ...interface{}) Logger {
	log := newLogger(l)
	log.ctx = append(l.ctx, ctxs...)
	return log
}

func (l logger) log(level Level, msg string, args []interface{}) (err error) {
	if level < l.lvl {
		return nil
	}
	err = l.enc.Encode(l.depth, level, msg, args, l.ctx)

	switch level {
	case PANIC:
		panic(ErrPanic)
	case FATAL:
		os.Exit(1)
	}

	return
}

func (l logger) Trace(msg string, args ...interface{}) error {
	return l.log(TRACE, msg, args)
}

func (l logger) Debug(msg string, args ...interface{}) error {
	return l.log(DEBUG, msg, args)
}

func (l logger) Info(msg string, args ...interface{}) error {
	return l.log(INFO, msg, args)
}

func (l logger) Warn(msg string, args ...interface{}) error {
	return l.log(WARN, msg, args)
}

func (l logger) Error(msg string, args ...interface{}) error {
	return l.log(ERROR, msg, args)
}

func (l logger) Panic(msg string, args ...interface{}) error {
	return l.log(PANIC, msg, args)
}

func (l logger) Fatal(msg string, args ...interface{}) error {
	return l.log(FATAL, msg, args)
}

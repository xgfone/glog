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

import "fmt"

// ErrKeyValueNum will be used when the number of key-values is not even.
var ErrKeyValueNum = fmt.Errorf("the number of key-values must be even")

// Encoder is a log encoder.
//
// Notice: if the encoder implementation supports the level when writing data,
// it should firstly decide whether the writer is LevelWriter and use WriteLevel
// to write the log, not Write.
type Encoder interface {
	// Reset the underlying writer.
	ResetWriter(Writer)

	// Return the underlying writer.
	//
	// Notice: only the most underlying encoder requires it. For the inner
	// encoder, such as FilterEncoder and MultiEncoder, it may be nil.
	// So, at the moment, the log information should be passed to the next encoder.
	Writer() Writer

	// Encode the log record and write it into the underlying writer.
	Encode(Record) error
}

type funcEncoder struct {
	writer  Writer
	encoder func(Writer, Record) error
}

func (e *funcEncoder) Writer() Writer {
	return e.writer
}

func (e *funcEncoder) ResetWriter(w Writer) {
	e.writer = w
}

func (e *funcEncoder) Encode(record Record) error {
	record.Depth++
	return e.encoder(e.writer, record)
}

// EncoderFunc converts a function to an hashable Encoder.
func EncoderFunc(w Writer, f func(Writer, Record) error) Encoder {
	return &funcEncoder{writer: w, encoder: f}
}

// MultiEncoder uses many encoders to encode the log record.
//
// It will return a MultiError if there is a error returned by an encoder
// in the corresponding order. For example,
//
//     encoders = ["kvtext", "kvjson"]
//     enc1 := KvTextEncoder(os.Stdout)
//     enc2 := KvJsonEncoder(os.Stderr)
//     logger := New(MultiEncoder(enc1, enc2))
//     err := logger.Info("msg", "key", "value")
//     if err != nil {
//         errs := err.(MultiError)
//         for i, e := range errs {
//             if e != nil {
//                 fmt.Printf("%s: %s\n", encoders[i], e.Error())
//             }
//         }
//     }
func MultiEncoder(encoders ...Encoder) Encoder {
	if len(encoders) == 0 {
		panic(fmt.Errorf("multi-encoder has no encoders"))
	} else if encoders[0] == nil {
		panic(fmt.Errorf("the first encoder must not be nil"))
	}

	return EncoderFunc(encoders[0].Writer(), func(w Writer, r Record) error {
		r.Depth++
		var hasErr bool
		errs := make([]error, len(encoders))
		for i, encoder := range encoders {
			e := encoder.Encode(r)
			errs[i] = e
			if e != nil {
				hasErr = true
			}
		}

		if hasErr {
			return MultiError{errs}
		}
		return nil
	})
}

// FilterEncoder returns an encoder that only forwards logs
// to the wrapped encoder if the given function evaluates true.
//
// For example, filter those logs that the level is less than ERROR.
//
//    FilterEncoder(encoder, func(r Record) bool { return r.Lvl >= LvlError })
//
func FilterEncoder(encoder Encoder, f func(Record) bool) Encoder {
	return EncoderFunc(encoder.Writer(), func(w Writer, r Record) error {
		if f(r) {
			r.Depth++
			return encoder.Encode(r)
		}
		return nil
	})
}

// AllowLoggerFilterEncoder returns a new encoder that only emits the logs
// emitted by the loggers named in allow.
func AllowLoggerFilterEncoder(allow []string, encoder Encoder) Encoder {
	return FilterEncoder(encoder, func(r Record) bool {
		for _, name := range allow {
			if name == r.Name {
				return true
			}
		}
		return false
	})
}

// DenyLoggerFilterEncoder returns a new encoder that ignore the logs
// emitted by the loggers named in deny.
func DenyLoggerFilterEncoder(deny []string, encoder Encoder) Encoder {
	return FilterEncoder(encoder, func(r Record) bool {
		for _, name := range deny {
			if name == r.Name {
				return false
			}
		}
		return true
	})
}

// LevelFilterEncoder returns an encoder that only writes records which are
// greater than the given verbosity level to the wrapped Handler.
//
// For example, to only output ERROR/PANIC/FATAL logs:
//
//     LevelFilterEncoder(LvlError, KvTextEncoder(os.Stdout))
//
func LevelFilterEncoder(level Level, encoder Encoder) Encoder {
	return FilterEncoder(encoder, func(r Record) bool { return r.Lvl >= level })
}

// NothingEncoder returns an encoder that does nothing.
func NothingEncoder() Encoder {
	return EncoderFunc(DiscardWriter(), func(w Writer, r Record) error { return nil })
}

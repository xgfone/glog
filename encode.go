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
	"encoding/json"
	"fmt"
	"time"

	"github.com/xgfone/logger/utils"
)

// The separators of the KV and the KV pair.
const (
	TextKVSep     = "="
	TextKVPairSep = " "
)

// Some key names. You can modify them to redefine them.
const (
	LevelKey = "lvl"
	NameKey  = "log"
	TimeKey  = "t"
	MsgKey   = "msg"
)

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
//    FilterEncoder(func(r Record) bool { return r.Lvl >= LvlError })
//
func FilterEncoder(f func(Record) bool, encoder Encoder) Encoder {
	return EncoderFunc(encoder.Writer(), func(w Writer, r Record) error {
		if f(r) {
			r.Depth++
			return encoder.Encode(r)
		}
		return nil
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
	return FilterEncoder(func(r Record) bool { return r.Lvl >= level }, encoder)
}

// NothingEncoder returns an encoder that does nothing.
func NothingEncoder() Encoder {
	return EncoderFunc(DiscardWriter(), func(w Writer, r Record) error { return nil })
}

// EncoderConfig configures the encoder.
type EncoderConfig struct {
	// If true, the encoder disable appending a newline.
	NotNewLine bool

	// TimeLayout is used to format time.Time.
	//
	// The default is time.RFC3339Nano.
	TimeLayout string

	// If true, the time uses UTC.
	IsTimeUTC bool

	// If true, the encoder will encode the name of the logger.
	IsName bool

	// If ture, the encoder will encode the current time.
	IsTime bool

	// If ture, the encoder will encode the level.
	IsLevel bool

	// If ture, the encoder will encode the level as the short description
	// instead of the long description.
	//
	// Notice: if IsLevel is false, it will be ignored.
	IsShortLevel bool

	// For the Key-Value encoder, it represents the key name of the logger name.
	// The global constant, NameKey, will be used by default.
	NameKey string

	// For the Key-Value encoder, it represents the key name of the time.
	// The global constant, TimeKey, will be used by default.
	TimeKey string

	// For the Key-Value encoder, it represents the key name of the level.
	// The global constant, LvlKey, will be used by default.
	LevelKey string

	// For the Key-Value encoder, it represents the key name of the message.
	// The global constant, MsgKey, will be used by default.
	MsgKey string

	// The separator between key and value, such as "=".
	// The global constant, TextKVSep, will be used by default.
	TextKVSep string

	// The separator between the key-value pairs, such as " ".
	// The global constant, TextKVPairSep, will be used by default.
	TextKVPairSep string
}

func (ec EncoderConfig) init() EncoderConfig {
	if ec.TimeLayout == "" {
		ec.TimeLayout = time.RFC3339Nano
	}

	if ec.NameKey == "" {
		ec.NameKey = NameKey
	}
	if ec.TimeKey == "" {
		ec.TimeKey = TimeKey
	}
	if ec.LevelKey == "" {
		ec.LevelKey = LevelKey
	}
	if ec.MsgKey == "" {
		ec.MsgKey = MsgKey
	}

	if ec.TextKVSep == "" {
		ec.TextKVSep = TextKVSep
	}
	if ec.TextKVPairSep == "" {
		ec.TextKVPairSep = TextKVPairSep
	}

	return ec
}

func newKvEncoderConfig(conf ...EncoderConfig) EncoderConfig {
	var c EncoderConfig
	if len(conf) > 0 {
		c = conf[0]
	}
	return c.init()
}

// KvTextEncoder returns a text encoder based on the key-value pair,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func KvTextEncoder(out Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)

	return EncoderFunc(out, func(out Writer, r Record) error {
		r.Depth++
		arglen := len(r.Args)
		ctxlen := len(r.Ctxs)
		if arglen%2 != 0 || ctxlen%2 != 0 {
			return ErrKeyValueNum
		}

		var err error
		var v interface{}
		var sep bool
		w := utils.DefaultBufferPools.Get()
		defer utils.DefaultBufferPools.Put(w)

		if c.IsTime {
			w.WriteByte('t')
			w.WriteString(c.TextKVSep)
			w.Write(utils.EncodeNowTime(c.TimeLayout, c.IsTimeUTC))
			sep = true
		}

		if c.IsLevel {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			w.WriteString(c.LevelKey)
			w.WriteString(c.TextKVSep)
			r.Lvl.WriteTo(w, c.IsShortLevel)
			sep = true
		}

		if c.IsName {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			w.WriteString(c.NameKey)
			w.WriteString(c.TextKVSep)
			w.WriteString(r.Name)
			sep = true
		}

		for i := 0; i < ctxlen; i += 2 {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			if v, err = MayBeValuer(r, r.Ctxs[i]); err != nil {
				return err
			}
			if err = utils.WriteIntoBuffer(w, v, true); err != nil {
				return err
			}
			w.WriteString(c.TextKVSep)
			if v, err = MayBeValuer(r, r.Ctxs[i+1]); err != nil {
				return err
			}
			if err = utils.WriteIntoBuffer(w, v, true); err != nil {
				return err
			}
			sep = true
		}

		for i := 0; i < arglen; i += 2 {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			if v, err = MayBeValuer(r, r.Args[i]); err != nil {
				return err
			}
			if err = utils.WriteIntoBuffer(w, v, true); err != nil {
				return err
			}

			w.WriteString(c.TextKVSep)

			if v, err = MayBeValuer(r, r.Args[i+1]); err != nil {
				return err
			}
			if err = utils.WriteIntoBuffer(w, v, true); err != nil {
				return err
			}
			sep = true
		}

		if sep {
			w.WriteString(c.TextKVPairSep)
		}

		w.WriteString(c.MsgKey)
		w.WriteString(c.TextKVSep)
		w.WriteString(r.Msg)

		if !c.NotNewLine {
			w.WriteByte('\n')
		}

		_, err = MayWriteLevel(out, r.Lvl, w.Bytes())
		return err
	})
}

// FmtTextEncoder returns a text encoder based on the % formatter,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func FmtTextEncoder(out Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)

	return EncoderFunc(out, func(out Writer, r Record) error {
		r.Depth++
		var err error
		var sep bool
		w := utils.DefaultBufferPools.Get()
		defer utils.DefaultBufferPools.Put(w)

		if c.IsTime {
			w.Write(utils.EncodeNowTime(c.TimeLayout, c.IsTimeUTC))
			sep = true
		}

		if c.IsName {
			if sep {
				w.WriteByte(' ')
			}
			w.WriteByte('(')
			w.WriteString(r.Name)
			w.WriteByte(')')
			sep = true
		}

		ctxlen := len(r.Ctxs)
		if ctxlen > 0 {
			if sep {
				w.WriteByte(' ')
			}

			for _, v := range r.Ctxs {
				w.WriteByte('{')
				if v, err = MayBeValuer(r, v); err != nil {
					return err
				}
				if err = utils.WriteIntoBuffer(w, v, true); err != nil {
					return err
				}
				w.WriteByte('}')
			}

			sep = true
		}

		if c.IsLevel {
			if sep {
				w.WriteByte(' ')
			}
			w.WriteByte('[')
			r.Lvl.WriteTo(w, c.IsShortLevel)
			w.WriteByte(']')
			sep = true
		}

		if sep {
			w.WriteString(": ")
		}

		for i := range r.Args {
			if r.Args[i], err = MayBeValuer(r, r.Args[i]); err != nil {
				return err
			}
		}

		if _, err = fmt.Fprintf(w, r.Msg, r.Args...); err != nil {
			return err
		}

		if !c.NotNewLine {
			w.WriteByte('\n')
		}

		_, err = MayWriteLevel(out, r.Lvl, w.Bytes())
		return err
	})
}

// KvJSONEncoder encodes the log as the JSON and outputs it to w.
//
// KvStdJSONEncoder and KvSimpleJSONEncoder will use this encoder.
//
// Notice: KvJSONEncoder doesn't append the newline.
func KvJSONEncoder(encodeJSON func(w Writer, newline bool, v interface{}) error,
	w Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)

	return EncoderFunc(w, func(out Writer, r Record) (err error) {
		r.Depth++
		_len := 3
		argslen := len(r.Args)
		ctxslen := len(r.Ctxs)
		if argslen%2 != 0 || ctxslen%2 != 0 {
			return ErrKeyValueNum
		}
		_len += argslen/2 + ctxslen/2
		maps := make(map[string]interface{}, _len)
		maps[c.MsgKey] = r.Msg

		if c.IsName {
			maps[c.NameKey] = r.Name
		}
		if c.IsLevel {
			if c.IsShortLevel {
				maps[c.LevelKey] = r.Lvl.ShortString()
			} else {
				maps[c.LevelKey] = r.Lvl.String()
			}
		}
		if c.IsTime {
			now := time.Now()
			if c.IsTimeUTC {
				now = now.UTC()
			}
			maps[c.TimeKey] = now
		}

		var v1, v2 interface{}
		for i := 0; i < ctxslen; i += 2 {
			if v1, err = MayBeValuer(r, r.Ctxs[i]); err != nil {
				return err
			}
			if v2, err = MayBeValuer(r, r.Ctxs[i+1]); err != nil {
				return err
			}
			maps[utils.ToString(v1)] = v2
		}
		for i := 0; i < argslen; i += 2 {
			if v1, err = MayBeValuer(r, r.Args[i]); err != nil {
				return err
			}
			if v2, err = MayBeValuer(r, r.Args[i+1]); err != nil {
				return err
			}
			maps[utils.ToString(v1)] = v2
		}

		return encodeJSON(w, !c.NotNewLine, maps)
	})
}

// KvStdJSONEncoder returns a new JSON encoder using the standard library,
// json, to encode the log record.
func KvStdJSONEncoder(w Writer, conf ...EncoderConfig) Encoder {
	return KvJSONEncoder(func(out Writer, newline bool, v interface{}) error {
		bs, err := json.Marshal(v)
		if err == nil {
			if newline {
				bs = append(bs, '\n')
			}
			_, err = out.Write(bs)
		}
		return err
	}, w, conf...)
}

// KvSimpleJSONEncoder returns a new JSON encoder using the funcion MarshalJSON
// to encode the log record.
//
// Except for the type of Array and Slice, it does not use the reflection.
// So it's faster than the standard library json.
func KvSimpleJSONEncoder(w Writer, conf ...EncoderConfig) Encoder {
	return KvJSONEncoder(func(out Writer, newline bool, v interface{}) error {
		buf := utils.DefaultBufferPools.Get()
		_, err := utils.MarshalJSON(buf, v)
		if err == nil {
			if newline {
				buf.WriteByte('\n')
			}
			_, err = w.Write(buf.Bytes())
		}
		utils.DefaultBufferPools.Put(buf)
		return err
	}, w, conf...)
}

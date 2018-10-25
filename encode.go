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
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// The separators of the KV and the KV pair.
const (
	TextKVSep     = "="
	TextKVPairSep = " "
)

// Some key names. You can modify them to redefine them.
const (
	LevelKey = "lvl"
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
	Encode(level Level, msg string, args []interface{}, ctx []interface{}) error
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
	return EncoderFunc(func(l Level, m string, a, c []interface{}) error {
		var hasErr bool
		errs := make([]error, len(encoders))
		for i, encoder := range encoders {
			e := encoder.Encode(l, m, a, c)
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

type encoderFunc func(Level, string, []interface{}, []interface{}) error

func (e encoderFunc) Encode(l Level, m string, args, ctx []interface{}) error {
	return e(l, m, args, ctx)
}

// EncoderFunc converts a function to an Encoder.
func EncoderFunc(f func(Level, string, []interface{}, []interface{}) error) Encoder {
	return encoderFunc(f)
}

// FilterEncoder returns an encoder that only forwards logs
// to the wrapped encoder if the given function evaluates true.
//
// For example, filter those logs that the level is less than ERROR.
//
//    FilterEncoder(func(lvl Level, msg string, args []interface{},
//                       ctxs []interface{}) bool {
//        return level >= ERROR
//    })
//
func FilterEncoder(f func(Level, string, []interface{}, []interface{}) bool,
	encoder Encoder) Encoder {
	return EncoderFunc(func(l Level, m string, args []interface{},
		ctxs []interface{}) error {
		if f(l, m, args, ctxs) {
			return encoder.Encode(l, m, args, ctxs)
		}
		return nil
	})
}

// LevelFilterEncoder returns an encoder that only writes records which are
// greater than the given verbosity level to the wrapped Handler.
//
// For example, to only output Error/PANIC/FATAL logs:
//
//     miss.LevelFilterEncoder(miss.ERROR, miss.KvTextEncoder(os.Stdout))
//
func LevelFilterEncoder(level Level, encoder Encoder) Encoder {
	return FilterEncoder(func(l Level, m string, args, ctxs []interface{}) bool {
		return l >= level
	}, encoder)
}

// NothingEncoder returns an encoder that does nothing.
func NothingEncoder() Encoder {
	return EncoderFunc(func(l Level, m string, args, ctx []interface{}) error {
		return nil
	})
}

// EncoderConfig configures the encoder.
type EncoderConfig struct {
	Slice []interface{}
	Map   map[string]interface{}

	// If true, the encoder disable appending a newline.
	NotNewLine bool

	// TimeLayout is used to format time.Time.
	//
	// The default is time.RFC3339Nano.
	TimeLayout string

	// If true, the time uses UTC.
	IsTimeUTC bool

	// If ture, the encoder will encode the current time.
	IsTime bool

	// If ture, the encoder will encode the level.
	IsLevel bool

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

	if ec.Slice == nil {
		ec.Slice = make([]interface{}, 0)
	}

	if ec.Map == nil {
		ec.Map = make(map[string]interface{})
	}

	return ec
}

func newKvEncoderConfig(conf ...EncoderConfig) EncoderConfig {
	var c EncoderConfig
	if len(conf) > 0 {
		c = conf[0]
	}

	if c.TimeKey == "" {
		c.TimeKey = TimeKey
	}
	if c.LevelKey == "" {
		c.LevelKey = LevelKey
	}
	if c.MsgKey == "" {
		c.MsgKey = MsgKey
	}

	if c.TextKVSep == "" {
		c.TextKVSep = TextKVSep
	}
	if c.TextKVPairSep == "" {
		c.TextKVPairSep = TextKVPairSep
	}

	if len(c.Slice) > 0 {
		c.Slice = append([]interface{}{}, c.Slice...)
	}

	if len(c.Map) > 0 {
		maps := make(map[string]interface{}, len(c.Map))
		for k, v := range c.Map {
			maps[k] = v
		}
		c.Map = maps
	}

	return c.init()
}

// KvTextEncoder returns a text encoder based on the key-value pair,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func KvTextEncoder(out io.Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)
	return EncoderFunc(func(l Level, m string, args, ctxs []interface{}) error {
		arglen := len(args)
		ctxlen := len(ctxs)
		if arglen%2 != 0 || ctxlen%2 != 0 {
			return ErrKeyValueNum
		}

		var err error
		var sep bool
		w := DefaultBufferPools.Get()
		defer DefaultBufferPools.Put(w)

		if c.IsTime {
			w.WriteByte('t')
			w.WriteString(c.TextKVSep)
			w.Write(getNowTime(c.TimeLayout, c.IsTimeUTC))
			sep = true
		}

		if c.IsLevel {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			w.WriteString(c.LevelKey)
			w.WriteString(c.TextKVSep)
			w.Write(l.Bytes())
			sep = true
		}

		for i := 0; i < ctxlen; i += 2 {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			if err = WriteIntoBufferErr(w, ctxs[i]); err != nil {
				return err
			}
			w.WriteString(c.TextKVSep)
			if err = WriteIntoBufferErr(w, ctxs[i+1]); err != nil {
				return err
			}
			sep = true
		}

		for i := 0; i < arglen; i += 2 {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			if err = WriteIntoBufferErr(w, args[i]); err != nil {
				return err
			}
			w.WriteString(c.TextKVSep)
			if err = WriteIntoBufferErr(w, args[i+1]); err != nil {
				return err
			}
			sep = true
		}

		if sep {
			w.WriteString(c.TextKVPairSep)
		}

		w.WriteString(c.MsgKey)
		w.WriteString(c.TextKVSep)
		w.WriteString(m)

		if !c.NotNewLine {
			w.WriteByte('\n')
		}

		_, err = MayWriteLevel(out, l, w.Bytes())
		return err
	})
}

// FmtTextEncoder returns a text encoder based on the % formatter,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func FmtTextEncoder(out io.Writer, conf ...EncoderConfig) Encoder {
	var c EncoderConfig
	if len(conf) > 0 {
		c = conf[0]

		if len(c.Slice) > 0 {
			c.Slice = append([]interface{}{}, c.Slice...)
		}

		if len(c.Map) > 0 {
			maps := make(map[string]interface{}, len(c.Map))
			for k, v := range c.Map {
				maps[k] = v
			}
			c.Map = maps
		}
	}
	c = c.init()

	return EncoderFunc(func(l Level, m string, args, ctxs []interface{}) error {
		var err error
		var sep bool
		w := DefaultBufferPools.Get()
		defer DefaultBufferPools.Put(w)

		if c.IsTime {
			w.Write(getNowTime(c.TimeLayout, c.IsTimeUTC))
			sep = true
		}

		if c.IsLevel {
			if sep {
				w.WriteByte(' ')
			}
			w.Write(l.Bytes())
			sep = true
		}

		ctxlen := len(ctxs)
		if ctxlen > 0 {
			if sep {
				w.WriteByte(' ')
			}

			for _, v := range ctxs {
				w.WriteByte('[')
				if err = WriteIntoBufferErr(w, v); err != nil {
					return err
				}
				w.WriteByte(']')
			}

			sep = true
		}

		if sep {
			w.WriteString(" :=>: ")
		}

		w.WriteString(fmt.Sprintf(m, args...))

		if !c.NotNewLine {
			w.WriteByte('\n')
		}

		_, err = MayWriteLevel(out, l, w.Bytes())
		return err
	})

}

// KvStdJSONEncoder returns a new encoder using the standard library, json,
// to encode the log record.
func KvStdJSONEncoder(w io.Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)

	return EncoderFunc(func(l Level, m string, args, ctxs []interface{}) error {
		_len := 3
		argslen := len(args)
		ctxslen := len(ctxs)
		if argslen%2 != 0 || ctxslen%2 != 0 {
			return ErrKeyValueNum
		}
		_len = _len + argslen/2 + ctxslen/2
		maps := make(map[string]interface{}, _len)
		maps[c.MsgKey] = m

		if c.IsLevel {
			maps[c.LevelKey] = l.String()
		}
		if c.IsTime {
			maps[c.TimeKey] = string(getNowTime(c.TimeLayout, c.IsTimeUTC))
		}

		for i := 0; i < argslen; i += 2 {
			maps[ToString(args[i])] = args[i+1]
		}
		for i := 0; i < ctxslen; i += 2 {
			maps[ToString(ctxs[i])] = ctxs[i+1]
		}

		bs, err := json.Marshal(maps)
		if err == nil {
			_, err = w.Write(bs)
		}
		return err
	})
}

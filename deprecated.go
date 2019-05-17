// Copyright 2019 xgfone
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

	"github.com/xgfone/go-tools/json2"
)

// Some key names. You can modify them to redefine them. DEPRECATED!!!
const (
	LevelKey = "lvl"
	NameKey  = "log"
	TimeKey  = "t"
	MsgKey   = "msg"

	TextKVSep     = "="
	TextKVPairSep = " "
)

// EncoderConfig configures the encoder.
//
// DEPRECATED!!!
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
		w := DefaultBufferPool.Get()
		defer DefaultBufferPool.Put(w)

		if c.IsTime {
			w.WriteByte('t')
			w.WriteString(c.TextKVSep)
			w.Write(json2.EncodeNowTime(c.TimeLayout, c.IsTimeUTC))
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
			if err = json2.Write(w, v, true); err != nil {
				return err
			}
			w.WriteString(c.TextKVSep)
			if v, err = MayBeValuer(r, r.Ctxs[i+1]); err != nil {
				return err
			}
			if err = json2.Write(w, v, true); err != nil {
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
			if err = json2.Write(w, v, true); err != nil {
				return err
			}

			w.WriteString(c.TextKVSep)

			if v, err = MayBeValuer(r, r.Args[i+1]); err != nil {
				return err
			}
			if err = json2.Write(w, v, true); err != nil {
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
//
// DEPRECATED!!! Please use NewFmtEncoder().
func FmtTextEncoder(out Writer, conf ...EncoderConfig) Encoder {
	c := newKvEncoderConfig(conf...)

	return EncoderFunc(out, func(out Writer, r Record) error {
		r.Depth++
		var err error
		var sep bool
		w := DefaultBufferPool.Get()
		defer DefaultBufferPool.Put(w)

		if c.IsTime {
			w.Write(json2.EncodeNowTime(c.TimeLayout, c.IsTimeUTC))
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
				if err = json2.Write(w, v, true); err != nil {
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
//
// Deprecated!!! Please use NewJSONEncoder().
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
			maps[json2.ToString(v1)] = v2
		}
		for i := 0; i < argslen; i += 2 {
			if v1, err = MayBeValuer(r, r.Args[i]); err != nil {
				return err
			}
			if v2, err = MayBeValuer(r, r.Args[i+1]); err != nil {
				return err
			}
			maps[json2.ToString(v1)] = v2
		}

		return encodeJSON(w, !c.NotNewLine, maps)
	})
}

// KvStdJSONEncoder returns a new JSON encoder using the standard library,
// json, to encode the log record.
//
// Deprecated!!! Please use NewStdJSONEncoder().
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
//
// Deprecated!!! Please use NewSimpleJSONEncoder().
func KvSimpleJSONEncoder(w Writer, conf ...EncoderConfig) Encoder {
	return KvJSONEncoder(func(out Writer, newline bool, v interface{}) error {
		buf := DefaultBufferPool.Get()
		_, err := json2.MarshalJSON(buf, v)
		if err == nil {
			if newline {
				buf.WriteByte('\n')
			}
			_, err = w.Write(buf.Bytes())
		}
		DefaultBufferPool.Put(buf)
		return err
	}, w, conf...)
}

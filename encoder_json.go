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

	"github.com/xgfone/go-tools/json2"
)

// JSONEncoderConfig is used to configure the json encoder.
type JSONEncoderConfig struct {
	// The key name of the time, which the encoder will extract its value
	// from the global Valuers and output them as the key-value.
	//
	// It is "time" by default. You can set it to "utctime" to output UTC time.
	TimeKey string

	// The key name of the level, which the encoder will extract its value
	// from the global Valuers and output them as the key-value.
	//
	// It is "level" by default. You can set it to "short_level" to output
	// the short name of level.
	LevelKey string

	// The name of the message, which is "msg" by default.
	MsgKey string

	// If true, the encoder won't append a newline.
	NoNewLine bool

	// Valuers can be used to override the valuer in the global Valuers.
	Valuers map[string]Valuer

	// The separators between key and value or key-value pairs.
	//
	// Notice: it's only used by the NewTextJSONEncoder encoder.
	TextKVSep     string // The default is "=".
	TextKVPairSep string // The default is " ".
}

func (c *JSONEncoderConfig) init() {
	if c.TimeKey == "" {
		c.TimeKey = "time"
	}
	if c.LevelKey == "" {
		c.LevelKey = "level"
	}
	if c.MsgKey == "" {
		c.MsgKey = "msg"
	}

	if c.TextKVSep == "" {
		c.TextKVSep = "="
	}
	if c.TextKVPairSep == "" {
		c.TextKVPairSep = " "
	}

	if c.Valuers == nil {
		c.Valuers = make(map[string]Valuer, len(Valuers)*2)
	}
	for k, v := range Valuers {
		if _, ok := c.Valuers[k]; !ok {
			c.Valuers[k] = v
		}
	}
}

// NewJSONEncoder encodes the log as the JSON and outputs it to w.
//
// KvStdJSONEncoder and KvSimpleJSONEncoder will use this encoder.
//
// Notice: KvJSONEncoder doesn't append the newline.
func NewJSONEncoder(encodeJSON func(w Writer, newline bool, v interface{}) error,
	w Writer, conf ...JSONEncoderConfig) Encoder {

	var c JSONEncoderConfig
	if len(conf) > 0 {
		c = conf[0]
	}
	c.init()

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

		if f, ok := c.Valuers[c.LevelKey]; ok {
			lvl, _ := f(r)
			maps[c.LevelKey] = lvl
		}

		if f, ok := c.Valuers[c.TimeKey]; ok {
			now, _ := f(r)
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

		return encodeJSON(w, !c.NoNewLine, maps)
	})
}

// NewStdJSONEncoder returns a new JSON encoder using the standard library,
// json, to encode the log record.
func NewStdJSONEncoder(w Writer, conf ...JSONEncoderConfig) Encoder {
	return NewJSONEncoder(func(out Writer, newline bool, v interface{}) error {
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

// NewSimpleJSONEncoder returns a new JSON encoder using the funcion
// json2.MarshalJSON to encode the log record.
//
// Except for the type of Array and Slice, it does not use the reflection.
// So it's faster than the standard library json.
func NewSimpleJSONEncoder(w Writer, conf ...JSONEncoderConfig) Encoder {
	return NewJSONEncoder(func(out Writer, newline bool, v interface{}) error {
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

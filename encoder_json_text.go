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

import "github.com/xgfone/go-tools/json2"

// NewTextJSONEncoder returns a text encoder based on the key-value pair,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func NewTextJSONEncoder(out Writer, conf ...JSONEncoderConfig) Encoder {
	var c JSONEncoderConfig
	if len(conf) > 0 {
		c = conf[0]
	}
	c.init()

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

		if f, ok := c.Valuers[c.TimeKey]; ok {
			now, _ := f(r)
			w.WriteString(c.TimeKey)
			w.WriteString(c.TextKVSep)
			json2.Write(w, now, true)
			sep = true
		}

		if f, ok := c.Valuers[c.LevelKey]; ok {
			if sep {
				w.WriteString(c.TextKVPairSep)
			}

			lvl, _ := f(r)
			w.WriteString(c.LevelKey)
			w.WriteString(c.TextKVSep)
			json2.Write(w, lvl, true)
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

		if !c.NoNewLine {
			w.WriteByte('\n')
		}

		_, err = MayWriteLevel(out, r.Lvl, w.Bytes())
		return err
	})
}

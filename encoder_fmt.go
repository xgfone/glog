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
	"fmt"
	"strings"

	"github.com/xgfone/go-tools/json2"
	"github.com/xgfone/go-tools/strings2"
)

// FmtEncoderConfig is used to configure the fmt encoder.
type FmtEncoderConfig struct {
	// The template of the log message.
	//
	// The default is "{time} {ctx} {caller} [{level}]: {msg}",
	// which outputs the log like this:
	//
	//   2019-05-17T10:24:02.5798087+08:00 ctx1|ctx2 github.com/xgfone/logger/encoder_example_test.go:41 [INFO]: msg
	//
	Tmpl string

	// Valuers is used to override the valuers.
	//
	// It will use the global Valuers by default, and add the two valuers,
	// that's, "ctx" and "msg".
	Valuers map[string]Valuer

	// If true, the encoder won't append a newline.
	NoNewLine bool

	// The left delimiter of the context placeholder, which is "{" by default.
	Left string

	// The right delimiter of the context placeholder, which is "}" by default.
	Right string
}

// NewFmtEncoder returns a text encoder based on the % formatter by the template,
// which will output the result into out.
//
// This encoder will add not only the global Valuers but also the customized
// valuer "ctx" and "msg" into conf.Valuers if they doesn't exist. Thereinto,
// "ctx" formats the contexts and "msg" formats the message by using fmt.Sprintf
// with the "%" formatter.
//
// Notice: This encoder supports LevelWriter.
func NewFmtEncoder(out Writer, conf ...FmtEncoderConfig) Encoder {
	var c FmtEncoderConfig
	if len(conf) > 0 {
		c = conf[0]
	}

	var defaultTmpl bool
	if c.Tmpl = strings.TrimSpace(c.Tmpl); c.Tmpl == "" {
		c.Tmpl = "{time} {ctx} {caller} [{level}]: {msg}"
		defaultTmpl = true
	}

	if !c.NoNewLine {
		if c.Tmpl[len(c.Tmpl)-1] != '\n' {
			c.Tmpl += "\n"
		}
	}

	if c.Left == "" {
		c.Left = "{"
	}
	if c.Right == "" {
		c.Right = "}"
	}

	if defaultTmpl && (c.Left != "{" || c.Right != "}") {
		panic("must not use the default template when customizing left or right delimiters")
	}

	if c.Valuers == nil {
		c.Valuers = make(map[string]Valuer, len(Valuers)+4)
	}

	for k, v := range Valuers {
		if _, ok := c.Valuers[k]; !ok {
			c.Valuers[k] = v
		}
	}

	if _, ok := c.Valuers["msg"]; !ok {
		c.Valuers["msg"] = func(r Record) (v interface{}, err error) {
			r.Depth++
			for i := range r.Args {
				if r.Args[i], err = MayBeValuer(r, r.Args[i]); err != nil {
					return
				}
			}
			return fmt.Sprintf(r.Msg, r.Args...), nil
		}
	}

	if _, ok := c.Valuers["ctx"]; !ok {
		c.Valuers["ctx"] = func(r Record) (v interface{}, err error) {
			if len(r.Ctxs) == 0 {
				return "", nil
			}

			buf := DefaultBufferPool.Get()
			defer DefaultBufferPool.Put(buf)

			r.Depth++
			for i, ctx := range r.Ctxs {
				if i > 0 {
					buf.WriteByte('|')
				}
				if ctx, err = MayBeValuer(r, ctx); err != nil {
					return
				}
				if err = json2.Write(buf, ctx, true); err != nil {
					return
				}
			}
			return buf.String(), nil
		}
	}

	formatTemplate := strings2.NewFormat(c.Left, c.Right).FormatOutput
	return EncoderFunc(out, func(w Writer, r Record) (err error) {
		r.Depth += 3

		buf := DefaultBufferPool.Get()
		formatTemplate(buf, c.Tmpl, func(key string) (interface{}, bool) {
			if f, ok := c.Valuers[key]; ok {
				v, err := f(r)
				if err != nil {
					return fmt.Sprintf("value error: %s", err.Error()), true
				}
				return v, true
			}
			return nil, false
		})
		_, err = MayWriteLevel(w, r.Lvl, buf.Bytes())
		DefaultBufferPool.Put(buf)
		return
	})
}

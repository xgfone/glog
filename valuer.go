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
	"time"

	"github.com/go-stack/stack"
)

// Predefine some valuers.
var (
	CallerStackValuer = CallerStack()
)

// Valuers is used to store the context values of the format template.
var Valuers = map[string]Valuer{
	"name":          func(r Record) (interface{}, error) { return r.Name, nil },
	"level":         func(r Record) (interface{}, error) { return r.Lvl.String(), nil },
	"short_level":   func(r Record) (interface{}, error) { return r.Lvl.ShortString(), nil },
	"time":          func(r Record) (interface{}, error) { return time.Now().Format(time.RFC3339Nano), nil },
	"utctime":       func(r Record) (interface{}, error) { return time.Now().UTC().Format(time.RFC3339Nano), nil },
	"line":          func(r Record) (interface{}, error) { r.Depth++; return r.Line(), nil },
	"lineno":        func(r Record) (interface{}, error) { r.Depth++; return r.LineAsInt(), nil },
	"funcname":      func(r Record) (interface{}, error) { r.Depth++; return r.FuncName(), nil },
	"filename":      func(r Record) (interface{}, error) { r.Depth++; return r.FileName(), nil },
	"long_filename": func(r Record) (interface{}, error) { r.Depth++; return r.LongFileName(), nil },
	"package":       func(r Record) (interface{}, error) { r.Depth++; return r.Package(), nil },
	"caller":        func(r Record) (interface{}, error) { r.Depth++; return r.Caller(), nil },
	"long_caller":   func(r Record) (interface{}, error) { r.Depth++; return r.LongCaller(), nil },
}

// A Valuer generates a log value, which represents a dynamic value
// that is re-evaluated with each log event before firing it.
type Valuer func(Record) (interface{}, error)

// MayBeValuer calls it and returns the result if v is Valuer.
// Or returns v without change.
func MayBeValuer(record Record, v interface{}) (interface{}, error) {
	switch f := v.(type) {
	case Valuer:
		record.Depth++
		return f(record)
	case func(Record) (interface{}, error):
		record.Depth++
		return f(record)
	}
	return v, nil
}

// Caller returns a Valuer that returns the caller "file:line".
//
// If fullPath is true, the file is the full path but removing the GOPATH prefix.
func Caller(fullPath ...bool) Valuer {
	format := "%v"
	if len(fullPath) > 0 && fullPath[0] {
		format = "%+v"
	}

	return func(r Record) (interface{}, error) {
		return fmt.Sprintf(format, stack.Caller(r.Depth+1)), nil
	}
}

// CallerStack returns a Valuer returning the caller stack without runtime.
//
// If fullPath is true, the file is the full path but removing the GOPATH prefix.
func CallerStack(fullPath ...bool) Valuer {
	format := "%v"
	if len(fullPath) > 0 && fullPath[0] {
		format = "%+v"
	}

	return func(r Record) (interface{}, error) {
		s := stack.Trace().TrimBelow(stack.Caller(r.Depth + 1)).TrimRuntime()
		if len(s) > 0 {
			return fmt.Sprintf(format, s), nil
		}
		return "", nil
	}
}

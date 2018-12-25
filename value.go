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
	"runtime"
	"strconv"
	"strings"
)

// A Valuer generates a log value, which represents a dynamic value
// that is re-evaluated with each log event before firing it.
type Valuer func(depth int, level Level) (interface{}, error)

// MayBeValuer calls it and returns the result if v is Valuer.
// Or returns v without change.
func MayBeValuer(depth int, lvl Level, v interface{}) (interface{}, error) {
	if f, ok := v.(Valuer); ok {
		return f(depth+1, lvl)
	}
	return v, nil
}

// Caller returns a Valuer that returns a file and line from a specified depth
// in the callstack.
//
// If fullPath is true, it will return the full path of the file.
func Caller(fullPath ...bool) Valuer {
	return func(depth int, level Level) (interface{}, error) {
		_, file, line, _ := runtime.Caller(depth + 1)
		if len(fullPath) == 0 || !fullPath[0] {
			idx := strings.LastIndexByte(file, '/')
			// using idx+1 below handles both of following cases:
			// idx == -1 because no "/" was found, or
			// idx >= 0 and we want to start at the character after the found "/".
			file = file[idx+1:]
		}
		return file + ":" + strconv.Itoa(line), nil
	}
}

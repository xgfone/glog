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
	"io"
	"reflect"
)

// Predefine some json mark
var (
	NullBytes  = []byte("null")
	TrueBytes  = []byte("true")
	FalseBytes = []byte("false")

	CommaBytes        = []byte{','}
	ColonBytes        = []byte{':'}
	LeftBracketBytes  = []byte{'['}
	RightBracketBytes = []byte{']'}
	LeftBraceBytes    = []byte{'{'}
	RightBraceBytes   = []byte{'}'}
)

// MarshalJSON marshals a value v as JSON into w.
//
// Support the types:
//   nil
//   bool
//   string | error
//   float32
//   float64
//   int
//   int8
//   int16
//   int32
//   int64
//   uint
//   uint8
//   uint16
//   uint32
//   uint64
//   time.Time
//   map[string]interface{} for json object
//   json.Marshaler
//   fmt.Stringer
//   Array or Slice of the type above
func MarshalJSON(w io.Writer, v interface{}) (n int, err error) {
	switch _v := v.(type) {
	case nil:
		return w.Write(NullBytes)
	case bool:
		if _v {
			return w.Write(TrueBytes)
		}
		return w.Write(FalseBytes)
	case string:
		return WriteString(w, _v, true)
	case error:
		return WriteString(w, _v.Error(), true)
	case fmt.Stringer:
		return WriteString(w, _v.String(), true)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
		float32, float64:
	case map[string]interface{}:
		// Write {
		if n, err = w.Write(LeftBraceBytes); err != nil {
			return n, err
		}
		total := n

		count := 0
		for key, value := range _v {
			if count > 0 {
				// Write comma
				if n, err = w.Write(CommaBytes); err != nil {
					return total, err
				}
				total += n
			}

			// Write key
			if n, err = WriteString(w, key, true); err != nil {
				return total, err
			}
			total += n

			// Write :
			if n, err = w.Write(ColonBytes); err != nil {
				return total, err
			}
			total += n

			// Write value
			if n, err = MarshalJSON(w, value); err != nil {
				return total, err
			}
			total += n

			count++
		}

		// Write }
		if n, err = w.Write(RightBraceBytes); err != nil {
			return n, err
		}
		return total + 1, nil
	case json.Marshaler:
		bs, err := _v.MarshalJSON()
		if err != nil {
			return 0, err
		}
		return w.Write(bs)
	case []string: // Optimzie []string
		if n, err = w.Write(LeftBracketBytes); err != nil {
			return n, err
		}

		total := n
		for i, _len := 0, len(_v); i < _len; i++ {
			if i > 0 {
				if n, err = w.Write(CommaBytes); err != nil {
					return total, err
				}
				total += n
			}

			if n, err = WriteString(w, _v[i], true); err != nil {
				return total, err
			}
			total += n
		}

		if n, err = w.Write(RightBracketBytes); err != nil {
			return total, err
		}
	case []interface{}: // Optimzie []interface{}
		if n, err = w.Write(LeftBracketBytes); err != nil {
			return n, err
		}

		total := n
		for i, _len := 0, len(_v); i < _len; i++ {
			if i > 0 {
				if n, err = w.Write(CommaBytes); err != nil {
					return total, err
				}
				total += n
			}

			if n, err = MarshalJSON(w, _v[i]); err != nil {
				return total, err
			}
			total += n
		}

		if n, err = w.Write(RightBracketBytes); err != nil {
			return total, err
		}
	default:
		// Check whether it's an array or slice.
		value := reflect.ValueOf(v)
		kind := value.Kind()
		if kind != reflect.Array && kind != reflect.Slice {
			return 0, fmt.Errorf("unknown type '%s'", value.Type().String())
		}

		if n, err = w.Write(LeftBracketBytes); err != nil {
			return n, err
		}

		total := n
		_len := value.Len()
		for i := 0; i < _len; i++ {
			if i > 0 {
				if n, err = w.Write(CommaBytes); err != nil {
					return total, err
				}
				total += n
			}

			if n, err = MarshalJSON(w, value.Index(i).Interface()); err != nil {
				return total, err
			}
			total += n
		}

		if n, err = w.Write(RightBracketBytes); err != nil {
			return total, err
		}
		return total + 1, nil
	}

	return w.Write(ToBytes(v))
}

// MarshalKvJSON marshals some key-value pairs as JSON into w.
//
// Notice: the key must be string, and the value may be one of the following:
//   nil
//   bool
//   string | error
//   float32
//   float64
//   int
//   int8
//   int16
//   int32
//   int64
//   uint
//   uint8
//   uint16
//   uint32
//   uint64
//   time.Time
//   map[string]interface{} for json object
//   json.Marshaler
//   fmt.Stringer
//   Array or Slice of the type above
func MarshalKvJSON(w io.Writer, args ...interface{}) (n int, err error) {
	_len := len(args)
	if _len == 0 {
		return
	} else if _len%2 == 1 {
		return 0, fmt.Errorf("args must be even")
	}

	var m int

	// Write {
	if m, err = w.Write(LeftBraceBytes); err != nil {
		return
	}
	n += m

	for i := 0; i < _len; i += 2 {
		// Write comma
		if i > 0 {
			if m, err = w.Write(CommaBytes); err != nil {
				return
			}
			n += m
		}

		// Write Key
		key, ok := args[i].(string)
		if !ok {
			return 0, fmt.Errorf("the %dth key is not string", i/2)
		}
		if m, err = WriteString(w, key, true); err != nil {
			return
		}
		n += m

		// Write :
		if m, err = WriteString(w, ":"); err != nil {
			return
		}
		n += m

		// Write Value
		switch v := args[i+1].(type) {
		case nil:
			if m, err = w.Write(NullBytes); err != nil {
				return
			}
			n += m
		case bool:
			if v {
				m, err = w.Write(TrueBytes)
			} else {
				m, err = w.Write(FalseBytes)
			}
			if err != nil {
				return
			}
			n += m
		case string:
			if m, err = WriteString(w, v, true); err != nil {
				return
			}
			n += m
		case error:
			if m, err = WriteString(w, v.Error(), true); err != nil {
				return
			}
			n += m
		case fmt.Stringer:
			if m, err = WriteString(w, v.String(), true); err != nil {
				return
			}
			n += m
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			if m, err = w.Write(ToBytes(v)); err != nil {
				return
			}
			n += m
		case json.Marshaler:
			var bs []byte
			if bs, err = v.MarshalJSON(); err != nil {
				return
			} else if m, err = w.Write(bs); err != nil {
				return
			}
			n += m
		default: // For array, slice or map
			if m, err = MarshalJSON(w, v); err != nil {
				return
			}
			n += m
		}
	}

	// Write }
	if m, err = w.Write(RightBraceBytes); err != nil {
		return
	}
	n += m
	return
}

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
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	nilBytes = []byte("nil")

	doubleQuotationByte = []byte{'"'}
	singleQuotationByte = []byte("'")
)

// Predefine some errors.
var (
	ErrType = fmt.Errorf("not support the type")
)

// Byter returns a []byte.
type Byter interface {
	Bytes() []byte
}

// MarshalText is an interface to marshal a value to text.
type MarshalText interface {
	MarshalText() ([]byte, error)
}

// StringWriter is a WriteString interface.
type StringWriter interface {
	WriteString(string) (int, error)
}

// WriteString writes s into w.
func WriteString(w io.Writer, s string, quote ...bool) (n int, err error) {
	var quotation bool
	if len(quote) > 0 && quote[0] {
		quotation = true
		if n, err = w.Write(doubleQuotationByte); err != nil {
			return
		}
	}

	if ws, ok := w.(StringWriter); ok {
		if n, err = ws.WriteString(s); err != nil {
			return
		}
	} else {
		if n, err = w.Write([]byte(s)); err != nil {
			return
		}
	}

	if quotation {
		if n, err = w.Write(doubleQuotationByte); err != nil {
			return
		}
	}

	return len(s), nil
}

// MultiError represents more than one error.
type MultiError struct {
	errs []error
}

func (m MultiError) Error() string {
	return ""
}

// Errors returns a list of errors.
func (m MultiError) Errors() []error {
	return m.errs
}

// ToBytesErr encodes a value to []byte.
//
// For the time.Time, it uses time.RFC3339Nano to format it.
//
// Support the types:
//   nil
//   bool
//   []byte
//   string
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
//   interface error
//   interface fmt.Stringer
//   interface Valuer
//   interface Byter
//   interface MarshalText
//
// For other types, return the error ErrType.
func ToBytesErr(i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case nil:
		return nilBytes, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case bool:
		if v {
			return TrueBytes, nil
		}
		return FalseBytes, nil
	case float32:
		return strconv.AppendFloat(make([]byte, 0, 24), float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(make([]byte, 0, 24), v, 'f', -1, 64), nil
	case int:
		return strconv.AppendInt(make([]byte, 0, 20), int64(v), 10), nil
	case int8:
		return strconv.AppendInt(make([]byte, 0, 20), int64(v), 10), nil
	case int16:
		return strconv.AppendInt(make([]byte, 0, 20), int64(v), 10), nil
	case int32:
		return strconv.AppendInt(make([]byte, 0, 20), int64(v), 10), nil
	case int64:
		return strconv.AppendInt(make([]byte, 0, 20), v, 10), nil
	case uint:
		return strconv.AppendUint(make([]byte, 0, 20), uint64(v), 10), nil
	case uint8:
		return strconv.AppendUint(make([]byte, 0, 20), uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(make([]byte, 0, 20), uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(make([]byte, 0, 20), uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(make([]byte, 0, 20), v, 10), nil
	case time.Time:
		return encodeTime(v, time.RFC3339Nano), nil
	case Valuer:
		i, err := v()
		if err != nil {
			return nil, nil
		}
		return ToBytesErr(i)
	case Byter:
		return v.Bytes(), nil
	case MarshalText:
		return v.MarshalText()
	case error:
		return []byte(v.Error()), nil
	case fmt.Stringer:
		return []byte(v.String()), nil
	default:
		return nil, ErrType
	}
}

// ToBytes is the same as ToBytesErr, but ignoring the error.
func ToBytes(i interface{}) []byte {
	bs, _ := ToBytesErr(i)
	return bs
}

// ToStringErr is the same as ToBytesErr, but returns string.
func ToStringErr(i interface{}) (string, error) {
	switch v := i.(type) {
	case nil:
		return "nil", nil
	case string:
		return v, nil
	case error:
		return v.Error(), nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		bs, err := ToBytesErr(i)
		return string(bs), err
	}
}

// ToString is the same as ToBytesErr, but returns string and ignores the error.
func ToString(i interface{}) string {
	s, _ := ToStringErr(i)
	return s
}

// WriteIntoBuffer is the same as ToBytesErr, but writes the result into w.
func WriteIntoBuffer(w *bytes.Buffer, i interface{}) error {
	switch v := i.(type) {
	case nil:
		w.WriteString("nil")
	case []byte:
		w.Write(v)
	case string:
		w.WriteString(v)
	case error:
		w.WriteString(v.Error())
	case fmt.Stringer:
		w.WriteString(v.String())
	default:
		bs, err := ToBytesErr(i)
		if err != nil {
			return err
		}
		w.Write(bs)
	}
	return nil
}

func encodeNowTime(layout string, utc ...bool) []byte {
	return encodeTime(time.Now(), layout, utc...)
}

func encodeTime(t time.Time, layout string, utc ...bool) []byte {
	if len(utc) > 0 && utc[0] {
		t = t.UTC()
	}
	return t.AppendFormat(make([]byte, 0, 36), layout)
}

// Range returns a integer range between start and stop, which progressively
// increase or descrease by step.
//
// If step is positive, r[i] = start + step*i when i>0 and r[i]<stop.
//
// If step is negative, r[i] = start + step*i but when i>0 and r[i]>stop.
//
// If step is 0, it will panic.
func Range(start, stop, step int) (r []int) {
	if step > 0 {
		for start < stop {
			r = append(r, start)
			start += step
		}
		return
	} else if step < 0 {
		for start > stop {
			r = append(r, start)
			start += step
		}
		return
	}

	panic(fmt.Errorf("The step must not be 0"))
}

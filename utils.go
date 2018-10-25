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
	"strconv"
	"time"
)

// MarshalText is an interface to marshal a value to text.
type MarshalText interface {
	MarshalText() ([]byte, error)
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

// ToString encodes a value to string,
//
// For the time.Time, it uses time.RFC3339Nano to format it.
func ToString(i interface{}) string {
	switch v := i.(type) {
	case nil:
		return "nil"
	case []byte:
		return string(v)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 3, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 3, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case time.Time:
		var _bs [64]byte
		return string(v.AppendFormat(_bs[:0], time.RFC3339Nano))
	case MarshalText:
		b, _ := v.MarshalText()
		return string(b)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%+v", v)
	}
}

// WriteIntoBytes encodes a value to []byte,
// then writes it into bs and returns bs.
//
// For the time.Time, it uses time.RFC3339Nano to format it.
func WriteIntoBytes(bs []byte, i interface{}) []byte {
	switch v := i.(type) {
	case nil:
		bs = append(bs, "nil"...)
	case []byte:
		bs = append(bs, v...)
	case string:
		bs = append(bs, v...)
	case bool:
		bs = append(bs, strconv.FormatBool(v)...)
	case float32:
		bs = append(bs, strconv.FormatFloat(float64(v), 'f', 3, 64)...)
	case float64:
		bs = append(bs, strconv.FormatFloat(v, 'f', 3, 64)...)
	case int64:
		bs = append(bs, strconv.FormatInt(v, 10)...)
	case int:
		bs = append(bs, strconv.FormatInt(int64(v), 10)...)
	case int8:
		bs = append(bs, strconv.FormatInt(int64(v), 10)...)
	case int16:
		bs = append(bs, strconv.FormatInt(int64(v), 10)...)
	case int32:
		bs = append(bs, strconv.FormatInt(int64(v), 10)...)
	case uint64:
		bs = append(bs, strconv.FormatUint(v, 10)...)
	case uint:
		bs = append(bs, strconv.FormatUint(uint64(v), 10)...)
	case uint8:
		bs = append(bs, strconv.FormatUint(uint64(v), 10)...)
	case uint16:
		bs = append(bs, strconv.FormatUint(uint64(v), 10)...)
	case uint32:
		bs = append(bs, strconv.FormatUint(uint64(v), 10)...)
	case time.Time:
		var _bs [64]byte
		bs = append(bs, v.AppendFormat(_bs[:0], time.RFC3339Nano)...)
	case MarshalText:
		b, _ := v.MarshalText()
		bs = append(bs, b...)
	case error:
		bs = append(bs, v.Error()...)
	case fmt.Stringer:
		bs = append(bs, v.String()...)
	default:
		bs = append(bs, fmt.Sprintf("%+v", v)...)
	}
	return bs
}

// WriteIntoBuffer encodes a value to byte.Buffer,
//
// For the time.Time, it uses time.RFC3339Nano to format it.
func WriteIntoBuffer(w *bytes.Buffer, i interface{}) {
	switch v := i.(type) {
	case nil:
		w.WriteString("nil")
	case []byte:
		w.Write(v)
	case string:
		w.WriteString(v)
	case bool:
		w.WriteString(strconv.FormatBool(v))
	case float32:
		w.WriteString(strconv.FormatFloat(float64(v), 'f', 3, 64))
	case float64:
		w.WriteString(strconv.FormatFloat(v, 'f', 3, 64))
	case int64:
		w.WriteString(strconv.FormatInt(v, 10))
	case int:
		w.WriteString(strconv.FormatInt(int64(v), 10))
	case int8:
		w.WriteString(strconv.FormatInt(int64(v), 10))
	case int16:
		w.WriteString(strconv.FormatInt(int64(v), 10))
	case int32:
		w.WriteString(strconv.FormatInt(int64(v), 10))
	case uint64:
		w.WriteString(strconv.FormatUint(v, 10))
	case uint:
		w.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint8:
		w.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint16:
		w.WriteString(strconv.FormatUint(uint64(v), 10))
	case uint32:
		w.WriteString(strconv.FormatUint(uint64(v), 10))
	case time.Time:
		var _bs [64]byte
		w.Write(v.AppendFormat(_bs[:0], time.RFC3339Nano))
	case MarshalText:
		b, _ := v.MarshalText()
		w.Write(b)
	case error:
		w.WriteString(v.Error())
	case fmt.Stringer:
		w.WriteString(v.String())
	default:
		w.WriteString(fmt.Sprintf("%+v", v))
	}
}

// WriteIntoBufferErr is the version returning error of WriterIntoBuffer.
func WriteIntoBufferErr(w *bytes.Buffer, i interface{}) error {
	switch v := i.(type) {
	case MarshalText:
		b, err := v.MarshalText()
		if err != nil {
			return err
		}
		w.Write(b)
	case Valuer:
		i, err := v()
		if err != nil {
			return err
		}
		WriteIntoBuffer(w, i)
	default:
		WriteIntoBuffer(w, i)
	}
	return nil
}

// WriteIntoBytesErr is the version returning error of WriterIntoBytes.
func WriteIntoBytesErr(bs []byte, i interface{}) ([]byte, error) {
	switch v := i.(type) {
	case MarshalText:
		b, err := v.MarshalText()
		if err != nil {
			return bs, err
		}
		bs = append(bs, b...)
	case Valuer:
		i, err := v()
		if err != nil {
			return bs, err
		}
		WriteIntoBytes(bs, i)
	default:
		WriteIntoBytes(bs, i)
	}
	return bs, nil
}

func getNowTime(layout string, utc ...bool) []byte {
	var _bs [64]byte
	now := time.Now()
	if len(utc) > 0 && utc[0] {
		now = now.UTC()
	}
	return now.AppendFormat(_bs[:0], layout)
}

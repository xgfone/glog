package miss

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

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
	if f, ok := i.(Valuer); ok {
		v, err := f()
		if err != nil {
			return err
		}
		WriteIntoBuffer(w, v)
	} else {
		WriteIntoBuffer(w, i)
	}
	return nil
}

// WriteIntoBytesErr is the version returning error of WriterIntoBytes.
func WriteIntoBytesErr(bs []byte, i interface{}) ([]byte, error) {
	if f, ok := i.(Valuer); ok {
		v, err := f()
		if err != nil {
			return bs, err
		}
		bs = WriteIntoBytes(bs, v)
	} else {
		bs = WriteIntoBytes(bs, i)
	}
	return bs, nil
}

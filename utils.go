package miss

import (
	"fmt"
	"strconv"
	"time"
)

// WriteIntoBytes encodes a value to []byte,
// then writes it into bs and returns bs.
//
// For the time.Time, it uses time.RFC3339Nano to format it.
func WriteIntoBytes(bs []byte, i interface{}) ([]byte, error) {
	var err error
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
	return bs, err
}

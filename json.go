package miss

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
//   map[string]interface{} for json object
//   json.Marshaler
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
			if n, err = WriteString(w, key); err != nil {
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

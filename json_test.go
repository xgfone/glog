package miss

import (
	"bytes"
	"fmt"
	"time"
)

func ExampleMarshalJSON() {
	buf := bytes.NewBuffer(nil)

	MarshalJSON(buf, 123)
	buf.WriteByte('\n')
	MarshalJSON(buf, 1.23)
	buf.WriteByte('\n')
	MarshalJSON(buf, "123")
	buf.WriteByte('\n')
	MarshalJSON(buf, time.Now())
	buf.WriteByte('\n')
	MarshalJSON(buf, []int{1, 2, 3})
	buf.WriteByte('\n')
	MarshalJSON(buf, []string{"a", "b", "c"})
	buf.WriteByte('\n')
	MarshalJSON(buf, []float64{1.2, 1.4, 1.6})
	buf.WriteByte('\n')
	MarshalJSON(buf, map[string]interface{}{"number": 123, "name": "abc"})
	buf.WriteByte('\n')

	fmt.Printf("%s\n", buf.String())

	// Output:
	//
}

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

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
	"os"
	"testing"
)

func TestCaller(t *testing.T) {
	if v, err := Caller(true)(Record{}); err != nil {
		t.Error(err)
	} else if v.(string) != "github.com/xgfone/logger/value_test.go:24" {
		t.Error(v)
	}
}

func ExampleCaller() {
	logger := New(KvTextEncoder(os.Stdout)).WithCxt("caller1", Caller())
	logger.Info("msg", "caller2", Caller())

	// Output:
	// caller1=value_test.go:33 caller2=value_test.go:33 msg=msg
}

func ExampleCallerStack() {
	if v, err := CallerStack(true)(Record{}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v.(string)[1:42])
		// v is like
		// [github.com/xgfone/logger/value_test.go:40 testing/example.go:121 testing/example.go:45 testing/testing.go:1073 _testmain.go:80]
	}

	// Output:
	// github.com/xgfone/logger/value_test.go:40
}

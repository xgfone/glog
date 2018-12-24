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

package callerstack

import (
	"fmt"

	"github.com/go-stack/stack"
	"github.com/xgfone/glog"
)

// Caller is the same as glog.Caller(true), but removing the GOPATH prefix.
func Caller(format ...string) glog.Valuer {
	return func(depth int, level glog.Level) (interface{}, error) {
		return fmt.Sprintf("%+v", stack.Caller(depth+1)), nil
	}
}

// CallerStack returns a Valuer returning the caller stack without runtime.
//
// The default is using "%+s:%d:%n" as the format. See github.com/go-stack/stack
func CallerStack(format ...string) glog.Valuer {
	return func(depth int, level glog.Level) (interface{}, error) {
		s := stack.Trace().TrimBelow(stack.Caller(depth + 1)).TrimRuntime()
		if len(s) > 0 {
			return fmt.Sprintf("%+v", s), nil
		}
		return "", nil
	}
}

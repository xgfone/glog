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
	"bytes"
	"testing"
)

func TestNoErrorLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	errlog := New(FmtTextEncoder(buf)).Cxt(Caller())
	logger := ToNoErrorLogger(errlog)

	errlog.Info("hello, %s", "abc")
	logger.Info("hello, %s", "xyz")
	if buf.String() != "{noerrorlog_test.go:27}: hello, abc\n{noerrorlog_test.go:28}: hello, xyz\n" {
		t.Error(buf.String())
	}

	if ToLogger(logger).GetDepth() != errlog.GetDepth() {
		t.Fail()
	}
}

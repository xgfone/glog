// Copyright 2019 xgfone
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
	"log"
	"os"
	"testing"
)

func TestStdLog(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New(FmtTextEncoder(buf))
	stdlog := log.New(logger.GetEncoder().Writer(), "[stdlog] ", log.Lshortfile)
	stdlog.Printf("hello, %s\n", "world")

	if buf.String() != "[stdlog] logger_test.go:28: hello, world\n" {
		t.Error(buf.String())
	}
}

func ExampleLevelFilterWriter() {
	logger1 := New(FmtTextEncoder(os.Stdout))
	logger1.Info("will output")

	writer := LevelFilterWriter(LvlError, os.Stdout)
	logger2 := New(FmtTextEncoder(writer))
	logger2.Info("won't output")

	// Output:
	// will output
}

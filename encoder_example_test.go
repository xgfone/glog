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

import "os"

func ExampleNewFmtEncoder() {
	// Only for test to replace the `time` context.
	valuers := map[string]Valuer{"time": func(r Record) (interface{}, error) { return "2019-05-16 17:29:12", nil }}

	encoder1 := NewFmtEncoder(os.Stdout, FmtEncoderConfig{Valuers: valuers})
	log1 := New(encoder1).WithCxt("fmt", func(r Record) (interface{}, error) { return "example", nil })
	log1.Info("test %s %s", "fmt", func(r Record) (interface{}, error) { return "encoder", nil })

	conf := FmtEncoderConfig{Tmpl: "{time} {filename}:{lineno:-6d} {level}: {msg}", Valuers: valuers}
	encoder2 := NewFmtEncoder(os.Stdout, conf)
	log2 := New(encoder2).WithCxt("fmt", func(r Record) (interface{}, error) { return "example", nil })
	log2.Info("test %s %s", "fmt", func(r Record) (interface{}, error) { return "encoder", nil })

	// Output:
	// 2019-05-16 17:29:12 fmt|example encoder_example_test.go:25 [INFO]: test fmt encoder
	// 2019-05-16 17:29:12 encoder_example_test.go:30     INFO: test fmt encoder
}

func ExampleNewStdJSONEncoder() {
	// Only for test to replace the `time` context.
	valuers := map[string]Valuer{"time": func(r Record) (interface{}, error) { return "2019-05-16 17:29:12", nil }}

	encoder := NewStdJSONEncoder(os.Stdout, JSONEncoderConfig{Valuers: valuers})
	log := New(encoder).WithCxt("caller", Caller())
	log.Info("test encoder", "encoder", "json", "type", func(r Record) (interface{}, error) { return "std", nil })

	// Output:
	// {"caller":"encoder_example_test.go:43","encoder":"json","level":"INFO","msg":"test encoder","time":"2019-05-16 17:29:12","type":"std"}
}

func ExampleNewTextJSONEncoder() {
	// Only for test to replace the `time` context.
	valuers := map[string]Valuer{"time": func(r Record) (interface{}, error) { return "2019-05-16T17:29:12Z", nil }}

	encoder := NewTextJSONEncoder(os.Stdout, JSONEncoderConfig{Valuers: valuers})
	log := New(encoder).WithCxt("caller", Caller())
	log.Info("test encoder", "encoder", "json", "type", func(r Record) (interface{}, error) { return "std", nil })

	// Output:
	// time=2019-05-16T17:29:12Z level=INFO caller=encoder_example_test.go:55 encoder=json type=std msg=test encoder
}

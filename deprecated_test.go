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

func ExampleKvTextEncoder() {
	conf := EncoderConfig{IsLevel: true}
	encoder := KvTextEncoder(os.Stdout, conf)
	log := New(encoder).WithCxt("name", "example", "id", 123, "caller", Caller())
	log.Info("test", "key1", "value1", "key2", "value2")

	// Output:
	// lvl=INFO name=example id=123 caller=deprecated_test.go:23 key1=value1 key2=value2 msg=test
}

func ExampleKvStdJSONEncoder() {
	conf := EncoderConfig{IsLevel: true}
	encoder := KvStdJSONEncoder(os.Stdout, conf)
	log := New(encoder).WithCxt("name", "example", "id", 123, "caller", Caller())
	log.Info("test", "key1", "value1", "key2", "value2")

	// Output:
	// {"caller":"deprecated_test.go:33","id":123,"key1":"value1","key2":"value2","lvl":"INFO","msg":"test","name":"example"}
}

func ExampleFmtTextEncoder() {
	conf := EncoderConfig{IsLevel: true}
	encoder := FmtTextEncoder(os.Stdout, conf)
	log := New(encoder).WithCxt("kv", "text", "example", Caller())
	log.Info("test %s %s", "value1", "value2")

	// Output:
	// {kv}{text}{example}{deprecated_test.go:43} [INFO]: test value1 value2
}

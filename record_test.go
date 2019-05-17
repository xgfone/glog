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

import "fmt"

func ExampleRecord_LongCaller() {
	caller := new(Record).LongCaller()
	fmt.Println(caller)

	// Output:
	// github.com/xgfone/logger/record_test.go:20
}

func ExampleRecord_CallerStack() {
	caller := new(Record).CallerStack()
	fmt.Println(caller[1:43])
	// caller is like
	// [github.com/xgfone/logger/record_test.go:28 testing/example.go:121 testing/example.go:45 testing/testing.go:1073 _testmain.go:84]

	// Output:
	// github.com/xgfone/logger/record_test.go:28
}

func ExampleRecord_Line() {
	line := new(Record).Line()
	fmt.Println(line)

	// Output:
	// 38
}

func ExampleRecord_FileName() {
	filename := new(Record).FileName()
	fmt.Println(filename)

	// Output:
	// record_test.go
}

func ExampleRecord_LongFileName() {
	filename := new(Record).LongFileName()
	fmt.Println(filename)

	// Output:
	// github.com/xgfone/logger/record_test.go
}

func ExampleRecord_QualifiedFuncName() {
	funcname := new(Record).QualifiedFuncName()
	fmt.Println(funcname)

	// Output:
	// github.com/xgfone/logger.ExampleRecord_QualifiedFuncName
}

func ExampleRecord_FuncName() {
	funcname := new(Record).FuncName()
	fmt.Println(funcname)

	// Output:
	// ExampleRecord_FuncName
}

func ExampleRecord_Package() {
	funcname := new(Record).Package()
	fmt.Println(funcname)

	// Output:
	// github.com/xgfone/logger
}

func ExampleRecord_Caller() {
	caller := new(Record).Caller()
	fmt.Println(caller)

	// Output:
	// record_test.go:86
}

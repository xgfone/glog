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
	"strconv"
	"strings"

	"github.com/go-stack/stack"
)

// Record stands for a log record.
type Record struct {
	// Name is the name of the logger to emit the log record.
	Name string

	// Depth is the depth of the caller.
	Depth int

	// Lvl is the level of the emitted log.
	Lvl Level

	// Msg and Args are the arguments of the emitted log.
	Msg  string
	Args []interface{}

	// Ctxs is the contexts of the Logger instance.
	Ctxs []interface{}

	caller stack.Call
	okCall bool
}

func (r *Record) getCaller() {
	if !r.okCall {
		r.caller = stack.Caller(r.Depth + 2)
		r.okCall = true
	}
}

// Caller is the same as LongCaller(), but using the short filename.
func (r *Record) Caller() string {
	r.getCaller()
	return fmt.Sprintf("%v", r.caller)
}

// LongCaller returns the caller "file:line", which is equal to
//     Caller(true)(r)
// or
//     r.LongFileName() + ":" + r.Line()
func (r *Record) LongCaller() string {
	r.getCaller()
	return fmt.Sprintf("%+v", r.caller)
}

// CallerStack returns the caller stack, which is equal to
//     CallerStack(true)(r)
func (r *Record) CallerStack() string {
	r.getCaller()
	s := stack.Trace().TrimBelow(r.caller).TrimRuntime()
	if len(s) > 0 {
		return fmt.Sprintf("%+v", s)
	}
	return ""
}

// LongFileName returns the long filename where the caller is in.
func (r *Record) LongFileName() string {
	r.getCaller()
	caller := r.LongCaller()
	if index := strings.IndexByte(caller, ':'); index > -1 {
		return caller[:index]
	}
	return ""
}

// FileName returns the short filename where the caller is in.
func (r *Record) FileName() string {
	r.getCaller()
	filename := r.LongFileName()
	if index := strings.LastIndexByte(filename, '/'); index > -1 {
		return filename[index+1:]
	}
	return filename
}

// Line returns the line number where the caller is on.
func (r *Record) Line() string {
	r.getCaller()
	caller := r.Caller()
	if index := strings.IndexByte(caller, ':'); index > -1 {
		return caller[index+1:]
	}
	return ""
}

// LineAsInt is the same as Line(), but returns the integer.
//
// Return 0 if the line is missing.
func (r *Record) LineAsInt() int {
	r.getCaller()
	if line := r.Line(); line != "" {
		v, _ := strconv.ParseInt(line, 10, 32)
		return int(v)
	}
	return 0
}

// QualifiedFuncName returns the qualified function name where the caller is in,
// which is equal to
//     r.Package() + "." + r.FuncName()
func (r *Record) QualifiedFuncName() string {
	r.getCaller()
	return fmt.Sprintf("%+n", r.caller)
}

// FuncName returns the function name where the caller is in.
func (r *Record) FuncName() string {
	r.getCaller()
	funcName := r.QualifiedFuncName()
	if index := strings.LastIndexByte(funcName, '.'); index > -1 {
		return funcName[index+1:]
	}
	return funcName
}

// Package returns the package where the callee is called.
func (r *Record) Package() string {
	r.getCaller()
	funcName := r.QualifiedFuncName()
	if index := strings.LastIndexByte(funcName, '.'); index > -1 {
		return funcName[:index]
	}
	return funcName
}

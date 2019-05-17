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

import "testing"

func BenchmarkLoggerNothingEncoderNoArgs(b *testing.B) {
	logger := New(NothingEncoder())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerNothingEncoderOneArg(b *testing.B) {
	logger := New(NothingEncoder())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "arg")
		}
	})
}

func BenchmarkLoggerNothingEncoderTwoArgs(b *testing.B) {
	logger := New(NothingEncoder())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "arg1", "arg2")
		}
	})
}

func BenchmarkLoggerNewTextJSONEncoderNoArgs(b *testing.B) {
	logger := New(NewTextJSONEncoder(DiscardWriter())).WithCxt("name", "bench")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerNewTextJSONEncoderArgs(b *testing.B) {
	logger := New(NewTextJSONEncoder(DiscardWriter())).WithCxt("name", "bench")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "key1", "value1", "key2", "value2")
		}
	})
}

func BenchmarkLoggerNewFmtEncoderNoArgs(b *testing.B) {
	conf := FmtEncoderConfig{Tmpl: "{time} {ctx} [{level}]: {msg}"}
	logger := New(NewFmtEncoder(DiscardWriter(), conf)).WithCxt("name", "bench")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerNewFmtEncoderArgs(b *testing.B) {
	conf := FmtEncoderConfig{Tmpl: "{time} {ctx} [{level}]: {msg}"}
	logger := New(NewFmtEncoder(DiscardWriter(), conf)).WithCxt("name", "bench")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test %s %s %s", "value1", "value2", "value3")
		}
	})
}

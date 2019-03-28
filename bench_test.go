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
	"testing"

	"github.com/xgfone/logger/utils"
)

func BenchmarkWriteIntoBuffer(b *testing.B) {
	bufPools := utils.NewBufferPool()
	bufPools.Put(bufPools.Get())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bufPools.Get()
			WriteIntoBuffer(buf, 123)
			bufPools.Put(buf)
		}
	})
}

func BenchmarkLoggerNothingEncoderNoArgs(b *testing.B) {
	logger := New(NothingEncoder())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerNothingEncoderOneArg(b *testing.B) {
	logger := New(NothingEncoder())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "arg")
		}
	})
}

func BenchmarkLoggerNothingEncoderTwoArgs(b *testing.B) {
	logger := New(NothingEncoder())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "arg1", "arg2")
		}
	})
}

func BenchmarkLoggerKvTextEncoderNoArgs(b *testing.B) {
	conf := EncoderConfig{IsLevel: true, IsTime: true}
	logger := New(KvTextEncoder(DiscardWriter(), conf)).Cxt("name", "bench")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerKvTextEncoderArgs(b *testing.B) {
	conf := EncoderConfig{IsLevel: true, IsTime: true}
	logger := New(KvTextEncoder(DiscardWriter(), conf)).Cxt("name", "bench")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test", "key1", "value1", "key2", "value2")
		}
	})
}

func BenchmarkLoggerFmtTextEncoderNoArgs(b *testing.B) {
	conf := EncoderConfig{IsLevel: true, IsTime: true}
	logger := New(FmtTextEncoder(DiscardWriter(), conf)).Cxt("name", "bench")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test")
		}
	})
}

func BenchmarkLoggerFmtTextEncoderArgs(b *testing.B) {
	conf := EncoderConfig{IsLevel: true, IsTime: true}
	logger := New(FmtTextEncoder(DiscardWriter(), conf)).Cxt("name", "bench")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test %s %s %s", "value1", "value2", "value3")
		}
	})
}

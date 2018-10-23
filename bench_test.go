package miss

import (
	"testing"
)

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

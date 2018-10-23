package miss

import "testing"

func BenchmarkBytesPool(b *testing.B) {
	bp := NewBytesPool()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bp.Put(bp.Get())
		}
	})
}

func BenchmarkBufferPool(b *testing.B) {
	bp := NewBufferPool()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bp.Put(bp.Get())
		}
	})
}

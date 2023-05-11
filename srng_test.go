package wyhash

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestOne(t *testing.T) {
	rng := SRNG(time.Now().UnixNano())
	t.Log(_wyp0, _wyp1)
	t.Log(rng.Uint64())
	t.Log(rng.Uint64())
	t.Log(rng.Uint64())

	b := make([]byte, 64)
	rng.ReadN(b, 32, 126)
	t.Log(string(b))
	rng.Read(b)
	t.Log(b)
}

func BenchmarkRead(b *testing.B) {
	rng := SRNG(time.Now().UnixNano())
	bf := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.Read(bf)
	}
}

func BenchmarkReadN(b *testing.B) {
	rng := SRNG(time.Now().UnixNano())
	bf := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.ReadN(bf, 32, 126)
	}
}

func BenchmarkReadGo(b *testing.B) {
	bf := make([]byte, 64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Read(bf)
	}
}

func BenchmarkReadConcurrent(b *testing.B) {
	rng := SRNG(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			bf := make([]byte, 32)
			rng.Read(bf)
		}()
	}
	wg.Wait()
}

func BenchmarkReadConcurrentGo(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			bf := make([]byte, 32)
			rand.Read(bf)
		}()
	}
	wg.Wait()
}

func TestSRng(t *testing.T) {
	var rng SRNG
	for i, want := range rngvecs {
		got := rng.Uint64()
		if got != want {
			t.Errorf("rng.Next()[%d]=%x, want %x", i, got, want)
		}
	}
}

func BenchmarkSRNG(b *testing.B) {
	var blackholeInt int
	var blackholeUint64 uint64
	var blackholeFloat64 float64

	b.Run("Int", func(b *testing.B) {
		rng := SRNG(2345)
		for i := 0; i < b.N; i++ {
			blackholeInt += rng.Int()
		}
	})

	b.Run("Intn", func(b *testing.B) {
		rng := SRNG(2345)
		for i := 0; i < b.N; i++ {
			blackholeInt += rng.Intn(1000)
		}
	})

	b.Run("Uint64", func(b *testing.B) {
		rng := SRNG(2345)
		for i := 0; i < b.N; i++ {
			blackholeUint64 += rng.Uint64()
		}
	})

	b.Run("Uint64n", func(b *testing.B) {
		b.Run("Large", func(b *testing.B) {
			rng := SRNG(2345)
			for i := 0; i < b.N; i++ {
				blackholeUint64 += rng.Uint64n(1<<63 + 1)
			}
		})

		b.Run("Med", func(b *testing.B) {
			rng := SRNG(2345)
			for i := 0; i < b.N; i++ {
				blackholeUint64 += rng.Uint64n(1<<31 + 1)
			}
		})

		b.Run("Small", func(b *testing.B) {
			rng := SRNG(2345)
			for i := 0; i < b.N; i++ {
				blackholeUint64 += rng.Uint64n(1000)
			}
		})
	})

	b.Run("Float64", func(b *testing.B) {
		rng := SRNG(2345)
		for i := 0; i < b.N; i++ {
			blackholeFloat64 += rng.Float64()
		}
	})

	runtime.KeepAlive(blackholeUint64)
	runtime.KeepAlive(blackholeFloat64)
}

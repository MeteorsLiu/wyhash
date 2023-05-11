package wyhash

import "math/bits"

// RNG is a random number generator.
// The zero value is valid.
type RNG uint64

// Int returns a random positive int.
// Not safe for concurrent callers.
func (r *RNG) Int() int {
	return int(uint(r.Uint64()) >> 1)
}

// Intn returns an int uniformly in [0, n).
// Not safe for concurrent callers.
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.Uint64n(uint64(n)))
}

// Uint64 returns a random uint64.
// Not safe for concurrent callers.
func (r *RNG) Uint64() uint64 {
	*r += RNG(_wyp0)
	return _wymum(u64(*r)^_wyp1, u64(*r))
}

// Uint64n returns a uint64 uniformly in [0, n).
// Not safe for concurrent callers.
func (r *RNG) Uint64n(n uint64) uint64 {
	if n == 0 {
		return 0
	}

	x := r.Uint64()
	h, l := bits.Mul64(x, n)

	if l < n {
		t := -n
		if t >= n {
			t -= n
			if t >= n {
				t = t % n
			}
		}

	again:
		if l < t {
			x = r.Uint64()
			h, l = bits.Mul64(x, n)
			goto again
		}
	}

	return h
}

// Float64 returns a float64 uniformly in [0, 1).
// Not safe for concurrent callers.
func (r *RNG) Float64() (v float64) {
again:
	v = float64(r.Uint64()>>11) / (1 << 53)
	if v == 1 {
		goto again
	}
	return
}

func (r *RNG) ReadN(b []byte, min, max int) (i int) {
	width := byte(max - min)
	minN := byte(min)
	isPowerofTwo := width&(width-1) == 0
	var pr uint64
	var each byte
	shift := 0
	for i = 0; i < len(b); i++ {
		if shift == 0 {
			pr = r.Uint64()
			shift = 7
		}
		each = byte(pr)
		if isPowerofTwo {
			each &= (width - 1)
		} else {
			each %= width
		}
		b[i] = each + minN
		pr >>= 8
		shift--
	}
	return
}

func (r *RNG) Read(b []byte) (i int) {
	var pr uint64
	shift := 0
	for i = 0; i < len(b); i++ {
		if shift == 0 {
			pr = r.Uint64()
			shift = 7
		}
		b[i] = byte(pr)
		pr >>= 8
		shift--
	}
	return
}

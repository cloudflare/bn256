package bn256

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"runtime"
	"sync"
	"testing"
)

// randomGF returns a random integer between 0 and p-1.
func randomGF(r io.Reader) *big.Int {
	k, err := rand.Int(r, p)
	if err != nil {
		panic(err)
	}
	return k
}

// toBigInt converts a field element into its reduced (mod p)
// integer representation.
func toBigInt(a *gfP) *big.Int {
	v := &gfP{}
	montDecode(v, a)
	c := new(big.Int)
	for i := len(v) - 1; i >= 0; i-- {
		c.Lsh(c, 64)
		c.Add(c, new(big.Int).SetUint64(v[i]))
	}
	return c
}

// togfP converts an integer into a field element (in
// Montgomery representation). This function assumes the
// input is between 0 and p-1; otherwise it panics.
func togfP(k *big.Int) *gfP {
	if k.Cmp(p) >= 0 {
		panic("not in the range 0 to p-1")
	}
	v := k.Bytes()
	v32 := [32]byte{}
	for i := len(v) - 1; i >= 0; i-- {
		v32[len(v)-1-i] = v[i]
	}
	u := &gfP{
		binary.LittleEndian.Uint64(v32[0*8 : 1*8]),
		binary.LittleEndian.Uint64(v32[1*8 : 2*8]),
		binary.LittleEndian.Uint64(v32[2*8 : 3*8]),
		binary.LittleEndian.Uint64(v32[3*8 : 4*8]),
	}
	montEncode(u, u)
	return u
}

func TestGFp(t *testing.T) {
	const testTimes = 1 << 8

	t.Run("add", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			bigB := randomGF(rand.Reader)
			want := bigC.Add(bigA, bigB).Mod(bigC, p)

			a := togfP(bigA)
			b := togfP(bigB)
			gfpAdd(c, a, b)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})

	t.Run("sub", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			bigB := randomGF(rand.Reader)
			want := bigC.Sub(bigA, bigB).Mod(bigC, p)

			a := togfP(bigA)
			b := togfP(bigB)
			gfpSub(c, a, b)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})

	t.Run("mul", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			bigB := randomGF(rand.Reader)
			want := bigC.Mul(bigA, bigB).Mod(bigC, p)

			a := togfP(bigA)
			b := togfP(bigB)
			gfpMul(c, a, b)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})

	t.Run("mul_fp_corruption", func(t *testing.T) {
		// By enabling the mutex profiling, the go runtime will traverse the
		// stack to measure mutex operations when a mutex is unlocked when a
		// goroutine is blocking on it.
		runtime.SetMutexProfileFraction(1)

		var wg sync.WaitGroup
		wg.Add(testTimes)
		for i := 0; i < testTimes; i++ {
			// If multiple goroutines interact with a global sync pool, the
			// goroutines may block on acquiring a lock. When that happens, the
			// mutex profiler will traverse the stack.
			go func() {
				defer wg.Done()

				a := togfP(randomGF(rand.Reader))
				b := togfP(randomGF(rand.Reader))
				c := &gfP{}
				gfpMul(c, a, b)

				// Print ends up interacting with a global sync pool.
				fmt.Print("")
			}()
		}
		wg.Wait()
	})

	t.Run("neg", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			want := bigC.Neg(bigA).Mod(bigC, p)

			a := togfP(bigA)
			gfpNeg(c, a)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})

	t.Run("inv", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			want := bigC.ModInverse(bigA, p)

			a := togfP(bigA)
			c.Invert(a)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})

	t.Run("sqrt", func(t *testing.T) {
		c := &gfP{}
		bigC := new(big.Int)
		for i := 0; i < testTimes; i++ {
			bigA := randomGF(rand.Reader)
			bigA.Mul(bigA, bigA).Mod(bigA, p)
			want := bigC.ModSqrt(bigA, p)

			a := togfP(bigA)
			c.Sqrt(a)
			got := toBigInt(c)

			if got.Cmp(want) != 0 {
				t.Errorf("got: %v want:%v", got, want)
			}
		}
	})
}

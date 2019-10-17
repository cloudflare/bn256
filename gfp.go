package bn256

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	"golang.org/x/crypto/hkdf"
)

type gfP [4]uint64

func newGFp(x int64) (out *gfP) {
	if x >= 0 {
		out = &gfP{uint64(x)}
	} else {
		out = &gfP{uint64(-x)}
		gfpNeg(out, out)
	}

	montEncode(out, out)
	return out
}

// hashToBase implements hashing a message to an element of the field.
// It follows the recommendations from https://tools.ietf.org/pdf/draft-irtf-cfrg-hash-to-curve-04.pdf
// L = ceil((256+128)/8)=48, ctr = 0, i = 1
func hashToBase(msg, dst []byte) *gfP {
	var t [48]byte
	info := []byte{'H', '2', 'C', byte(0), byte(1)}
	r := hkdf.New(sha256.New, msg, dst, info)
	if _, err := r.Read(t[:]); err != nil {
		panic(err)
	}
	var x big.Int
	v := x.SetBytes(t[:]).Mod(&x, p).Bytes()
	u := &gfP{
		binary.LittleEndian.Uint64(v[0*8 : 1*8]),
		binary.LittleEndian.Uint64(v[1*8 : 2*8]),
		binary.LittleEndian.Uint64(v[2*8 : 3*8]),
		binary.LittleEndian.Uint64(v[3*8 : 4*8]),
	}
	montEncode(u, u)
	return u
}

func (e *gfP) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x", e[3], e[2], e[1], e[0])
}

func (e *gfP) Set(f *gfP) {
	e[0] = f[0]
	e[1] = f[1]
	e[2] = f[2]
	e[3] = f[3]
}

func (e *gfP) Exp(f *gfP, bits [4]uint64) {
	sum, power := &gfP{}, &gfP{}
	sum.Set(rN1)
	power.Set(f)

	for word := 0; word < 4; word++ {
		for bit := uint(0); bit < 64; bit++ {
			if (bits[word]>>bit)&1 == 1 {
				gfpMul(sum, sum, power)
			}
			gfpMul(power, power, power)
		}
	}

	gfpMul(sum, sum, r3)
	e.Set(sum)
}

func (e *gfP) Invert(f *gfP) {
	e.Exp(f, pMinus2)
}

func (e *gfP) Sqrt(f *gfP) {
	// Since p = 4k+3, then e = f^(k+1) is a root of f.
	e.Exp(f, pPlus1over4)
}

func (e *gfP) Marshal(out []byte) {
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			out[8*w+b] = byte(e[3-w] >> (56 - 8*b))
		}
	}
}

func (e *gfP) Unmarshal(in []byte) {
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			e[3-w] += uint64(in[8*w+b]) << (56 - 8*b)
		}
	}
}

func montEncode(c, a *gfP) { gfpMul(c, a, r2) }
func montDecode(c, a *gfP) { gfpMul(c, a, &gfP{1}) }

func sign0(e *gfP) int {
	var x [4]uint64
	montDecode((*gfP)(&x), e)
	for w := 3; w >= 0; w-- {
		if x[w] > pMinus1Over2[w] {
			return 1
		} else if x[w] < pMinus1Over2[w] {
			return -1
		}
	}
	return 1
}

func legendre(e *gfP) int {
	f := &gfP{}
	// Since p = 4k+3, then e^(2k+1) is the Legendre symbol of e.
	f.Exp(e, pMinus1Over2)

	montDecode(f, f)

	if *f != [4]uint64{} {
		return 2*int(f[0]&1) - 1
	}

	return 0
}

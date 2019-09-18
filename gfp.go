package bn256

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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

func os2ip(b [32]byte) [4]uint64 {
	t := [4]uint64{}
	for w := 0; w < 4; w++ {
		t[w] = binary.LittleEndian.Uint64(b[8*w : 8*w+8])
	}
	return t
}

// hashToBase implements hashing a message to an element of the field.
// It follows the recommendations from https://tools.ietf.org/pdf/draft-irtf-cfrg-hash-to-curve-04.pdf
// Note that:
//      - we don't use HKDF-Extract and HKDF-Expand procedures as the field is prime,
//      - the bias introduced by interpreting a 256-bit hash as an integer modulo p is negligeble.
func hashToBase(msg []byte) *gfP {
	dstMsg := make([]byte, 16+len(msg))
	copy(dstMsg[:16], []byte("H2C-BN256-SHA256"))
	copy(dstMsg[16:], msg)
	h := sha256.Sum256(dstMsg)
	e := &gfP{}
	*e = gfP(os2ip(h))
	montEncode(e, e)
	return e
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

var one = *newGFp(1)

func legendre(e *gfP) int {
	if *e == [4]uint64{} {
		return 0
	}
	f := &gfP{}
	f.Exp(e, pMinus1Over2)

	if *f == one {
		return 1
	}

	return -1
}

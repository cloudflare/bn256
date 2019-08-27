package bn256

import (
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

func fromBytes(b []byte) *gfP {
    if len(b) != 32 {
        return nil
    }
    e := &gfP{}
    for w:=0;w<4;w++ {
        e[w] = binary.LittleEndian.Uint64(b[8*w:8*w+8])
    }
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

func (e *gfP) Pow(f *gfP, bits [4]uint64) {
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
    // pI is (p-2)
	pI := [4]uint64{0x185cac6c5e089665, 0xee5b88d120b5b59e, 0xaa6fecb86184dc21, 0x8fb501e34aa387f9}
    e.Pow(f, pI)
}

func (e *gfP) Sqrt(f *gfP) {
    // since s = (p-1)/2 is odd, then q=(s+1)/2 is a positive integer, 
    // and we can define e = f^q that yields e^2=f^s·f=f.
    q := [4]uint64{0x86172b1b1782259a, 0x7b96e234482d6d67, 0x6a9bfb2e18613708, 0x23ed4078d2a8e1fe}
    e.Pow(f, q)
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
    // pL is (p-1)/2
    pL := [4]uint64{0x0c2e56362f044b33, 0xf72dc468905adacf, 0xd537f65c30c26e10, 0x47da80f1a551c3fc}
    f := &gfP{}
    f.Pow(e, pL)

    if *f == one {
        return 1
    }

    return -1
}

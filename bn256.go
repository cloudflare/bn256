// Package bn256 implements a particular bilinear group at the 128-bit security level.
//
// Bilinear groups are the basis of many of the new cryptographic protocols that have been proposed over the past
// decade. They consist of a triplet of groups (G₁, G₂ and GT) such that there exists a function e(g₁ˣ,g₂ʸ)=gTˣʸ (where
// gₓ is a generator of the respective group). That function is called a pairing function.
//
// This package specifically implements the Optimal Ate pairing over a 256-bit Barreto-Naehrig curve as described in
// http://cryptojedi.org/papers/dclxvi-20100714.pdf. Its output is compatible with the implementation described in that
// paper.
package bn256

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

func randomK(r io.Reader) (k *big.Int, err error) {
	for {
		k, err = rand.Int(r, Order)
		if k.Sign() > 0 || err != nil {
			return
		}
	}

	return
}

// G1 is an abstract cyclic group. The zero value is suitable for use as the output of an operation, but cannot be used
// as an input.
type G1 struct {
	p *curvePoint
}

// RandomG1 returns x and g₁ˣ where x is a random, non-zero number read from r.
func RandomG1(r io.Reader) (*big.Int, *G1, error) {
	k, err := randomK(r)
	if err != nil {
		return nil, nil, err
	}

	return k, new(G1).ScalarBaseMult(k), nil
}

func (g *G1) String() string {
	return "bn256.G1" + g.p.String()
}

// ScalarBaseMult sets e to g*k where g is the generator of the group and then returns e.
func (e *G1) ScalarBaseMult(k *big.Int) *G1 {
	if e.p == nil {
		e.p = &curvePoint{}
	}
	e.p.Mul(curveGen, k)
	return e
}

// ScalarMult sets e to a*k and then returns e.
func (e *G1) ScalarMult(a *G1, k *big.Int) *G1 {
	if e.p == nil {
		e.p = &curvePoint{}
	}
	e.p.Mul(a.p, k)
	return e
}

// Add sets e to a+b and then returns e.
func (e *G1) Add(a, b *G1) *G1 {
	if e.p == nil {
		e.p = &curvePoint{}
	}
	e.p.Add(a.p, b.p)
	return e
}

// Neg sets e to -a and then returns e.
func (e *G1) Neg(a *G1) *G1 {
	if e.p == nil {
		e.p = &curvePoint{}
	}
	e.p.Neg(a.p)
	return e
}

// Set sets e to a and then returns e.
func (e *G1) Set(a *G1) *G1 {
	if e.p == nil {
		e.p = &curvePoint{}
	}
	e.p.Set(a.p)
	return e
}

// Marshal converts n to a byte slice.
func (n *G1) Marshal() []byte {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	n.p.MakeAffine()
	ret := make([]byte, numBytes*2)
	if n.p.IsInfinity() {
		return ret
	}

	x, y := &gfP{}, &gfP{}
	montDecode(x, &n.p.x)
	montDecode(y, &n.p.y)

	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			ret[8*w+b] = byte(x[3-w] >> (56 - 8*b))
		}
	}
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			ret[8*w+b+32] = byte(y[3-w] >> (56 - 8*b))
		}
	}

	return ret
}

// Unmarshal sets e to the result of converting the output of Marshal back into a group element and then returns e.
func (e *G1) Unmarshal(m []byte) ([]byte, error) {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	if len(m) < 2*numBytes {
		return nil, errors.New("bn256: not enough data")
	}

	if e.p == nil {
		e.p = &curvePoint{}
	} else {
		e.p.x, e.p.y = gfP{0}, gfP{0}
	}

	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			e.p.x[3-w] += uint64(m[8*w+b]) << (56 - 8*b)
		}
	}
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			e.p.y[3-w] += uint64(m[8*w+b+32]) << (56 - 8*b)
		}
	}
	montEncode(&e.p.x, &e.p.x)
	montEncode(&e.p.y, &e.p.y)

	zero := gfP{0}
	if e.p.x == zero && e.p.y == zero {
		// This is the point at infinity.
		e.p.y = *newGFp(1)
		e.p.z = gfP{0}
		e.p.t = gfP{0}
	} else {
		e.p.z = *newGFp(1)
		e.p.t = *newGFp(1)

		if !e.p.IsOnCurve() {
			return nil, errors.New("bn256: malformed point")
		}
	}

	return m[2*numBytes:], nil
}

// G2 is an abstract cyclic group. The zero value is suitable for use as the output of an operation, but cannot be used
// as an input.
type G2 struct {
	p *twistPoint
}

// RandomG2 returns x and g₂ˣ where x is a random, non-zero number read from r.
func RandomG2(r io.Reader) (*big.Int, *G2, error) {
	k, err := randomK(r)
	if err != nil {
		return nil, nil, err
	}

	return k, new(G2).ScalarBaseMult(k), nil
}

func (g *G2) String() string {
	return "bn256.G2" + g.p.String()
}

// ScalarBaseMult sets e to g*k where g is the generator of the group and then returns out.
func (e *G2) ScalarBaseMult(k *big.Int) *G2 {
	if e.p == nil {
		e.p = &twistPoint{}
	}
	e.p.Mul(twistGen, k)
	return e
}

// ScalarMult sets e to a*k and then returns e.
func (e *G2) ScalarMult(a *G2, k *big.Int) *G2 {
	if e.p == nil {
		e.p = &twistPoint{}
	}
	e.p.Mul(a.p, k)
	return e
}

// Add sets e to a+b and then returns e.
func (e *G2) Add(a, b *G2) *G2 {
	if e.p == nil {
		e.p = &twistPoint{}
	}
	e.p.Add(a.p, b.p)
	return e
}

// Neg sets e to -a and then returns e.
func (e *G2) Neg(a *G2) *G2 {
	if e.p == nil {
		e.p = &twistPoint{}
	}
	e.p.Neg(a.p)
	return e
}

// Set sets e to a and then returns e.
func (e *G2) Set(a *G2) *G2 {
	if e.p == nil {
		e.p = &twistPoint{}
	}
	e.p.Set(a.p)
	return e
}

// // Marshal converts n into a byte slice.
// func (n *G2) Marshal() []byte {
// 	// Each value is a 256-bit number.
// 	const numBytes = 256 / 8
//
// 	n.p.MakeAffine(nil)
// 	ret := make([]byte, numBytes*4)
// 	if n.p.IsInfinity() {
// 		return ret
// 	}
//
// 	xxBytes := new(big.Int).Mod(n.p.x.x, p).Bytes()
// 	xyBytes := new(big.Int).Mod(n.p.x.y, p).Bytes()
// 	yxBytes := new(big.Int).Mod(n.p.y.x, p).Bytes()
// 	yyBytes := new(big.Int).Mod(n.p.y.y, p).Bytes()
//
// 	copy(ret[1*numBytes-len(xxBytes):], xxBytes)
// 	copy(ret[2*numBytes-len(xyBytes):], xyBytes)
// 	copy(ret[3*numBytes-len(yxBytes):], yxBytes)
// 	copy(ret[4*numBytes-len(yyBytes):], yyBytes)
//
// 	return ret
// }
//
// // Unmarshal sets e to the result of converting the output of Marshal back into a group element and then returns e.
// func (e *G2) Unmarshal(m []byte) ([]byte, error) {
// 	// Each value is a 256-bit number.
// 	const numBytes = 256 / 8
//
// 	if len(m) < 4*numBytes {
// 		return nil, errors.New("bn256: not enough data")
// 	}
//
// 	if e.p == nil {
// 		e.p = newTwistPoint(nil)
// 	}
//
// 	e.p.x.x.SetBytes(m[0*numBytes : 1*numBytes])
// 	e.p.x.y.SetBytes(m[1*numBytes : 2*numBytes])
// 	e.p.y.x.SetBytes(m[2*numBytes : 3*numBytes])
// 	e.p.y.y.SetBytes(m[3*numBytes : 4*numBytes])
//
// 	if e.p.x.IsZero() && e.p.y.IsZero() {
// 		// This is the point at infinity.
// 		e.p.y.SetOne()
// 		e.p.z.SetZero()
// 		e.p.t.SetZero()
// 	} else {
// 		e.p.z.SetOne()
// 		e.p.t.SetOne()
//
// 		if !e.p.IsOnCurve() {
// 			return nil, errors.New("bn256: malformed point")
// 		}
// 	}
//
// 	return m[4*numBytes:], nil
// }

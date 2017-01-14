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
// BUG(agl): this function is not complete: a==b fails.
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

// // Marshal converts n to a byte slice.
// func (n *G1) Marshal() []byte {
// 	// Each value is a 256-bit number.
// 	const numBytes = 256 / 8
//
// 	n.p.MakeAffine()
// 	ret := make([]byte, numBytes*2)
// 	if n.p.IsInfinity() {
// 		return ret
// 	}
//
// 	xBytes := new(big.Int).Mod(n.p.x, p).Bytes()
// 	yBytes := new(big.Int).Mod(n.p.y, p).Bytes()
//
// 	copy(ret[1*numBytes-len(xBytes):], xBytes)
// 	copy(ret[2*numBytes-len(yBytes):], yBytes)
//
// 	return ret
// }
//
// // Unmarshal sets e to the result of converting the output of Marshal back into a group element and then returns e.
// func (e *G1) Unmarshal(m []byte) ([]byte, error) {
// 	// Each value is a 256-bit number.
// 	const numBytes = 256 / 8
//
// 	if len(m) < 2*numBytes {
// 		return nil, errors.New("bn256: not enough data")
// 	}
//
// 	if e.p == nil {
// 		e.p = &curvePoint{}
// 	}
//
// 	e.p.x.SetBytes(m[0*numBytes : 1*numBytes])
// 	e.p.y.SetBytes(m[1*numBytes : 2*numBytes])
//
// 	if e.p.x.Sign() == 0 && e.p.y.Sign() == 0 {
// 		// This is the point at infinity.
// 		e.p.y.SetInt64(1)
// 		e.p.z.SetInt64(0)
// 		e.p.t.SetInt64(0)
// 	} else {
// 		e.p.z.SetInt64(1)
// 		e.p.t.SetInt64(1)
//
// 		if !e.p.IsOnCurve() {
// 			return nil, errors.New("bn256: malformed point")
// 		}
// 	}
//
// 	return m[2*numBytes:], nil
// }

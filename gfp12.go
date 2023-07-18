package bn256

// For details of the algorithms used, see "Multiplication and Squaring on
// Pairing-Friendly Fields, Devegili et al.
// http://eprint.iacr.org/2006/471.pdf.

import (
	"math/big"
)

// gfP12 implements the field of size p¹² as a quadratic extension of gfP6
// where ω²=τ.
type gfP12 struct {
	x, y gfP6 // value is xω + y
}

var gfP12Gen *gfP12 = &gfP12{
	x: gfP6{
		x: gfP2{
			x: gfP{0x62d608d6bb67a4fb, 0x9a66ec93f0c2032f, 0x5391628e924e1a34, 0x2162dbf7de801d0e},
			y: gfP{0x3e0c1a72bf08eb4f, 0x4972ec05990a5ecc, 0xf7b9a407ead8007e, 0x3ca04c613572ce49},
		},
		y: gfP2{
			x: gfP{0xace536a5607c910e, 0xda93774a941ddd40, 0x5de0e9853b7593ad, 0xe05bb926f513153},
			y: gfP{0x3f4c99f8abaf1a22, 0x66d5f6121f86dc33, 0x8e0a82f68a50abba, 0x819927d1eebd0695},
		},
		z: gfP2{
			x: gfP{0x7cdef49c5477faa, 0x40eb71ffedaa199d, 0xbc896661f17c9b8f, 0x3144462983c38c02},
			y: gfP{0xcd09ee8dd8418013, 0xf8d050d05faa9b11, 0x589e90a555507ee1, 0x58e4ab25f9c49c15},
		},
	},
	y: gfP6{
		x: gfP2{
			x: gfP{0x7e76809b142d020b, 0xd9949d1b2822e995, 0x3de93d974f84b076, 0x144523477028928d},
			y: gfP{0x79952799f9ef4b0, 0x4102c47aa3df01c6, 0xfa82a633c53da2e1, 0x54c3f0392f9f7e0e},
		},
		y: gfP2{
			x: gfP{0xd3432a335533272b, 0xa008fbbdc7d74f4a, 0x68e3c81eb7295ed9, 0x17fe34c21fdecef2},
			y: gfP{0xfb0bc4c0ef6df55f, 0x8bdc585b70bc2120, 0x17d498d2cb720def, 0x2a368248319b899c},
		},
		z: gfP2{
			x: gfP{0xf8487d81cb354c6c, 0x7421be69f1522caa, 0x6940c778b9fb2d54, 0x7da4b04e102bb621},
			y: gfP{0x97b91989993e7be4, 0x8526545356eab684, 0xb050073022eb1892, 0x658b432ad09939c0},
		},
	},
}

func (e *gfP12) String() string {
	return "(" + e.x.String() + "," + e.y.String() + ")"
}

func (e *gfP12) Set(a *gfP12) *gfP12 {
	e.x.Set(&a.x)
	e.y.Set(&a.y)
	return e
}

func (e *gfP12) SetZero() *gfP12 {
	e.x.SetZero()
	e.y.SetZero()
	return e
}

func (e *gfP12) SetOne() *gfP12 {
	e.x.SetZero()
	e.y.SetOne()
	return e
}

func (e *gfP12) IsZero() bool {
	return e.x.IsZero() && e.y.IsZero()
}

func (e *gfP12) IsOne() bool {
	return e.x.IsZero() && e.y.IsOne()
}

func (e *gfP12) Conjugate(a *gfP12) *gfP12 {
	e.x.Neg(&a.x)
	e.y.Set(&a.y)
	return e
}

func (e *gfP12) Neg(a *gfP12) *gfP12 {
	e.x.Neg(&a.x)
	e.y.Neg(&a.y)
	return e
}

// Frobenius computes (xω+y)^p = x^p ω·ξ^((p-1)/6) + y^p
func (e *gfP12) Frobenius(a *gfP12) *gfP12 {
	e.x.Frobenius(&a.x)
	e.y.Frobenius(&a.y)
	e.x.MulScalar(&e.x, xiToPMinus1Over6)
	return e
}

// FrobeniusP2 computes (xω+y)^p² = x^p² ω·ξ^((p²-1)/6) + y^p²
func (e *gfP12) FrobeniusP2(a *gfP12) *gfP12 {
	e.x.FrobeniusP2(&a.x)
	e.x.MulGFP(&e.x, xiToPSquaredMinus1Over6)
	e.y.FrobeniusP2(&a.y)
	return e
}

func (e *gfP12) FrobeniusP4(a *gfP12) *gfP12 {
	e.x.FrobeniusP4(&a.x)
	e.x.MulGFP(&e.x, xiToPSquaredMinus1Over3)
	e.y.FrobeniusP4(&a.y)
	return e
}

func (e *gfP12) Add(a, b *gfP12) *gfP12 {
	e.x.Add(&a.x, &b.x)
	e.y.Add(&a.y, &b.y)
	return e
}

func (e *gfP12) Sub(a, b *gfP12) *gfP12 {
	e.x.Sub(&a.x, &b.x)
	e.y.Sub(&a.y, &b.y)
	return e
}

func (e *gfP12) Mul(a, b *gfP12) *gfP12 {
	tx := (&gfP6{}).Mul(&a.x, &b.y)
	t := (&gfP6{}).Mul(&b.x, &a.y)
	tx.Add(tx, t)

	ty := (&gfP6{}).Mul(&a.y, &b.y)
	t.Mul(&a.x, &b.x).MulTau(t)

	e.x.Set(tx)
	e.y.Add(ty, t)
	return e
}

func (e *gfP12) MulScalar(a *gfP12, b *gfP6) *gfP12 {
	e.x.Mul(&a.x, b)
	e.y.Mul(&a.y, b)
	return e
}

func (c *gfP12) Exp(a *gfP12, power *big.Int) *gfP12 {
	sum := (&gfP12{}).SetOne()
	t := &gfP12{}

	for i := power.BitLen() - 1; i >= 0; i-- {
		t.Square(sum)
		if power.Bit(i) != 0 {
			sum.Mul(t, a)
		} else {
			sum.Set(t)
		}
	}

	c.Set(sum)
	return c
}

// "New software speed records for cryptographic pairings"
// Section 3.3, Final exponentiation - 
// Algorithm 2 Exponentiation by v = 1868033.
// https://cryptojedi.org/papers/dclxvi-20100714.pdf
func (e *gfP12) powToVCyclo6(a *gfP12) *gfP12 {
	// The sequence of 21 special squarings and 4 multiplications
	t0, t1, t2 := &gfP12{}, &gfP12{}, &gfP12{}
	
	t0.SquareCyclo6(a)
	t0.SquareCyclo6(t0)
	t0.SquareCyclo6(t0) // t0 = a ^ 8
	t1.SquareCyclo6(t0)
	t1.SquareCyclo6(t1)
	t1.SquareCyclo6(t1) // t1 = a ^ 64
	t2.Conjugate(t0)    // t2 = a ^ -8
	t2.Mul(t2, a)       // t2 = a ^ -7
	t2.Mul(t2, t1)      // t2 = a ^ 57
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2) // t2 = a ^ (2^7 * 57) = a ^ 7296
	t2.Mul(t2, a)       // t2 = a ^ 7297
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2)
	t2.SquareCyclo6(t2) // t2 = a ^ (7297 * 256) = a ^ 1868032
	e.Mul(t2, a)
	return e
}

// PowToUCyclo6 is used in final exponentiation after easy part(a ^ ((p^2 + 1)(p^6-1))).
// Due to  u = v^3, so a^u can be implemented as three [powToVCyclo6].
func (e *gfP12) PowToUCyclo6(a *gfP12) *gfP12 {
	e.powToVCyclo6(a)
	e.powToVCyclo6(e)
	e.powToVCyclo6(e)
	return e
}

func (e *gfP12) Square(a *gfP12) *gfP12 {
	// Complex squaring algorithm
	v0 := (&gfP6{}).Mul(&a.x, &a.y)

	t := (&gfP6{}).MulTau(&a.x)
	t.Add(&a.y, t)
	ty := (&gfP6{}).Add(&a.x, &a.y)
	ty.Mul(ty, t).Sub(ty, v0)
	t.MulTau(v0)
	ty.Sub(ty, t)

	e.x.Add(v0, v0)
	e.y.Set(ty)
	return e
}

// SquareCyclo6 is used in final exponentiation after easy part(a ^ ((p^2 + 1)(p^6-1))).
// Note that after the easy part of the final exponentiation, 
// the resulting element lies in cyclotomic subgroup. 
// "New software speed records for cryptographic pairings"
// Section 3.3, Final exponentiation
// https://cryptojedi.org/papers/dclxvi-20100714.pdf
// The fomula reference:
// Granger/Scott (PKC2010). 
// Section 3.2
// https://eprint.iacr.org/2009/565.pdf
func (e *gfP12) SquareCyclo6(a *gfP12) *gfP12 {
	// f = xω + y = (h2τ² + h1τ + h0)ω + (g2τ² + g1τ + g0) = h2ω^5 + g2ω^4 + h1 ω^3 + g1ω² + h0ω + g0
	// we can also represets f as a cubic over a quadartic extension:
	// Fp4[s]=Fp2[s]/(s^2-ξ), Fp12[t]=Fp4[t]/(t^3-s), s^2=ξ, t^3=s then
	// f = ct² + bt + a = (c0 + c1s)t² + (b0 + b1s)t + (a0 + a1s) = c1t^5 + b1t^4 + a1t^3 + c0t^2 + b0t + a0
	// both extensions are based on Fp2, so we got t^6 = ω^6 = ξ, and 
	// a0 = g0, a1 = h1, b0 = h0, b1 = g2, c0 = g1, c1 = h2
	// g0 = a.y.z, h1 = a.x.y, h0 = a.x.z, g2 = a.y.x, g1 = a.y.y, h2 = a.x.x
	tmp := &gfP12{}

	f02 := &tmp.y.x
	f01 := &tmp.y.y
	f00 := &tmp.y.z
	f12 := &tmp.x.x
	f11 := &tmp.x.y
	f10 := &tmp.x.z

	t00, t01, t02, t10, t11, t12 := &gfP2{}, &gfP2{}, &gfP2{}, &gfP2{}, &gfP2{}, &gfP2{}

	gfP4Square(t11, t00, &a.x.y, &a.y.z) // (t00 + t11s) = (a0 + a1s)^2 = a²
	gfP4Square(t12, t01, &a.y.x, &a.x.z) // (t01 + t12s) = (b0 + b1s)^2 = b²
	gfP4Square(t02, t10, &a.x.x, &a.y.y) // (t10 + t02s) = (c0 + c1s)^2 = c²

	// t02 + t10s = (t10 + t02s)s
	f00.MulXi(t02)
	t02.Set(t10)
	t10.Set(f00)

	// triples
	// (t00 + t11s) = 3a²
	// (t01 + t12s) = 3b²
	// (t02 + t10s) = 3c²s
	f00.Add(t00, t00)
	t00.Add(f00, t00)
	f00.Add(t01, t01)
	t01.Add(f00, t01)
	f00.Add(t02, t02)
	t02.Add(f00, t02)
	f00.Add(t10, t10)
	t10.Add(f00, t10)
	f00.Add(t11, t11)
	t11.Add(f00, t11)
	f00.Add(t12, t12)
	t12.Add(f00, t12)

	// (f00 + f11s) = -2Conjugate(a0 + a1s) = -2Conjugate(a)
	// (f01 + f12s) = -2Conjugate(c0 + c1s) = -2Conjugate(c)
	// (f02 + f10s) = 2Conjugate(b0 + b1s) = 2Conjugate(b)
	f00.Add(&a.y.z, &a.y.z)
	f00.Neg(f00)
	f01.Add(&a.y.y, &a.y.y)
	f01.Neg(f01)
	f02.Add(&a.y.x, &a.y.x)
	f02.Neg(f02)
	f10.Add(&a.x.z, &a.x.z)
	f11.Add(&a.x.y, &a.x.y)
	f12.Add(&a.x.x, &a.x.x)

	// A = (f00 + f11s) = 3a² - 2Conjugate(a)
	// C = (f01 + f12s) = 3b² - 2Conjugate(c)
	// B = (f02 + f10s) = 3c² + 2Conjugate(b)
	f00.Add(f00, t00)
	f01.Add(f01, t01)
	f02.Add(f02, t02)
	f10.Add(f10, t10)
	f11.Add(f11, t11)
	f12.Add(f12, t12)

	return e.Set(tmp)
}

// Implicit gfP4 squaring for Granger/Scott special squaring in final expo
// gfP4Square takes two gfP2 x, y representing the gfP4 element xu+y, where
// u²=ξ.
func gfP4Square(retX, retY, x, y *gfP2) {
	t1, t2 := &gfP2{}, &gfP2{}

	t1.Square(x)
	t2.Square(y)

	retX.Add(x, y)
	retX.Square(retX)
	retX.Sub(retX, t1)
	retX.Sub(retX, t2) // retX = 2xy

	retY.MulXi(t1)
	retY.Add(retY, t2) // retY = x^2*xi + y^2
}

func (e *gfP12) Invert(a *gfP12) *gfP12 {
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t1, t2 := &gfP6{}, &gfP6{}

	t1.Square(&a.x)
	t2.Square(&a.y)
	t1.MulTau(t1).Sub(t2, t1)
	t2.Invert(t1)

	e.x.Neg(&a.x)
	e.y.Set(&a.y)
	e.MulScalar(e, t2)
	return e
}

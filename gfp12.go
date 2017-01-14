// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	x, y *gfP6 // value is xω + y
}

var gfP12Gen *gfP12 = &gfP12{
	x: &gfP6{
		x: &gfP2{
			x: bigFromBase10("21196725492326442667627602847448291017588785091108454726553939930507871042194"),
			y: bigFromBase10("43377286344385619261366633786771510260308325523330699049673437042931316243080"),
		},
		y: &gfP2{
			x: bigFromBase10("6519643025219345838017222895845262421933106748812739147102715685220942156402"),
			y: bigFromBase10("49167428805016148020211692433963935169227827998815966523759923767649702416286"),
		},
		z: &gfP2{
			x: bigFromBase10("20859153989312102139134197789193625833054352751342763983124765227328906753159"),
			y: bigFromBase10("54487609867103086000679472636440811221594782878893833473117980998935745956453"),
		},
	},
	y: &gfP6{
		x: &gfP2{
			x: bigFromBase10("39332187218762343173097683922586548248512461497033773500717078587710862589062"),
			y: bigFromBase10("53094021597187285527531149248961924798783924165048746910730430368152798446315"),
		},
		y: &gfP2{
			x: bigFromBase10("30733062774817315099333283633560206070773769397463591953591634872711994123707"),
			y: bigFromBase10("13560812839206871407210482171294630929511117297628119163695762035076749250365"),
		},
		z: &gfP2{
			x: bigFromBase10("57080396271372976186541856910255122494067079575964633601781325258931774013251"),
			y: bigFromBase10("60034081832369659851990276215766505463071420460830584339216728009871686767851"),
		},
	},
}

var orderGFp12 = bigFromBase10("5688586322044912818495476059952511913220161997008513033865263024457511546372336486465384053741177870122238790378988248492868967549063121515457048454734194860812915217459153099877897619772319898962451961994821938234951064262135307282633084386574670818976195275817753868262711253608580076303692877822251735983080479554213708955729530459705478609314292725579946637105704473921113038440759553251717079334665650011653184035065145535675670079086534589833773421356378413041225756374314356529606095183416556985272910233411621847595581523251820488259252312355448996015393109412004252136619188354139665479991056444730710187423766644386577604179020346526462921609605283518468180269339871327133051648853282457236848149703432532634710662240185522331738439942350111761168708496315643894145514267625862975955882102721485150562730926204774508417862029733749241093779742204934881384086327325872165686307524282561505309814589417379018796960")

func newGFp12(pool *bnPool) *gfP12 {
	return &gfP12{newGFp6(pool), newGFp6(pool)}
}

func (e *gfP12) String() string {
	e.Minimal()
	return "(" + e.x.String() + "," + e.y.String() + ")"
}

func (e *gfP12) Put(pool *bnPool) {
	e.x.Put(pool)
	e.y.Put(pool)
}

func (e *gfP12) Set(a *gfP12) *gfP12 {
	e.x.Set(a.x)
	e.y.Set(a.y)
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

func (e *gfP12) Minimal() {
	e.x.Minimal()
	e.y.Minimal()
}

func (e *gfP12) IsZero() bool {
	e.Minimal()
	return e.x.IsZero() && e.y.IsZero()
}

func (e *gfP12) IsOne() bool {
	e.Minimal()
	return e.x.IsZero() && e.y.IsOne()
}

func (e *gfP12) Conjugate(a *gfP12) *gfP12 {
	e.x.Neg(a.x)
	e.y.Set(a.y)
	return a
}

func (e *gfP12) Neg(a *gfP12) *gfP12 {
	e.x.Neg(a.x)
	e.y.Neg(a.y)
	return e
}

// Frobenius computes (xω+y)^p = x^p ω·ξ^((p-1)/6) + y^p
func (e *gfP12) Frobenius(a *gfP12, pool *bnPool) *gfP12 {
	e.x.Frobenius(a.x, pool)
	e.y.Frobenius(a.y, pool)
	e.x.MulScalar(e.x, xiToPMinus1Over6, pool)
	return e
}

// FrobeniusP2 computes (xω+y)^p² = x^p² ω·ξ^((p²-1)/6) + y^p²
func (e *gfP12) FrobeniusP2(a *gfP12, pool *bnPool) *gfP12 {
	e.x.FrobeniusP2(a.x)
	e.x.MulGFP(e.x, xiToPSquaredMinus1Over6)
	e.y.FrobeniusP2(a.y)
	return e
}

func (e *gfP12) FrobeniusP4(a *gfP12) *gfP12 {
	e.x.FrobeniusP4(a.x)
	e.x.MulGFP(e.x, xiToPSquaredMinus1Over3)

	e.y.FrobeniusP4(a.y)
	return e
}

func (e *gfP12) Add(a, b *gfP12) *gfP12 {
	e.x.Add(a.x, b.x)
	e.y.Add(a.y, b.y)
	return e
}

func (e *gfP12) Sub(a, b *gfP12) *gfP12 {
	e.x.Sub(a.x, b.x)
	e.y.Sub(a.y, b.y)
	return e
}

func (e *gfP12) Mul(a, b *gfP12, pool *bnPool) *gfP12 {
	tx := newGFp6(pool)
	tx.Mul(a.x, b.y, pool)
	t := newGFp6(pool)
	t.Mul(b.x, a.y, pool)
	tx.Add(tx, t)

	ty := newGFp6(pool)
	ty.Mul(a.y, b.y, pool)
	t.Mul(a.x, b.x, pool)
	t.MulTau(t, pool)
	e.y.Add(ty, t)
	e.x.Set(tx)

	tx.Put(pool)
	ty.Put(pool)
	t.Put(pool)
	return e
}

func (e *gfP12) MulScalar(a *gfP12, b *gfP6, pool *bnPool) *gfP12 {
	e.x.Mul(e.x, b, pool)
	e.y.Mul(e.y, b, pool)
	return e
}

func (c *gfP12) Exp(a *gfP12, power *big.Int, pool *bnPool) *gfP12 {
	sum := newGFp12(pool)
	sum.SetOne()
	t := newGFp12(pool)

	for i := power.BitLen() - 1; i >= 0; i-- {
		t.Square(sum, pool)
		if power.Bit(i) != 0 {
			sum.Mul(t, a, pool)
		} else {
			sum.Set(t)
		}
	}

	c.Set(sum)

	sum.Put(pool)
	t.Put(pool)

	return c
}

func (e *gfP12) Square(a *gfP12, pool *bnPool) *gfP12 {
	// Complex squaring algorithm
	v0 := newGFp6(pool)
	v0.Mul(a.x, a.y, pool)

	t := newGFp6(pool)
	t.MulTau(a.x, pool)
	t.Add(a.y, t)
	ty := newGFp6(pool)
	ty.Add(a.x, a.y)
	ty.Mul(ty, t, pool)
	ty.Sub(ty, v0)
	t.MulTau(v0, pool)
	ty.Sub(ty, t)

	e.y.Set(ty)
	e.x.Double(v0)

	v0.Put(pool)
	t.Put(pool)
	ty.Put(pool)

	return e
}

func (e *gfP12) Invert(a *gfP12, pool *bnPool) *gfP12 {
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t1 := newGFp6(pool)
	t2 := newGFp6(pool)

	t1.Square(a.x, pool)
	t2.Square(a.y, pool)
	t1.MulTau(t1, pool)
	t1.Sub(t2, t1)
	t2.Invert(t1, pool)

	e.x.Neg(a.x)
	e.y.Set(a.y)
	e.MulScalar(e, t2, pool)

	t1.Put(pool)
	t2.Put(pool)

	return e
}

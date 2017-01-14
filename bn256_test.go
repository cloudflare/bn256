// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bn256

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

var nonResidueGFp2 *gfP2 = &gfP2{
	x: bigFromBase10("53218179051328405336552655779633645459428885067951312777116768657680400902477"),
	y: bigFromBase10("22064341973762970966695021897037353416356658909813054882328212721160729323211"),
}

func (e *gfP2) Jacobi(a *gfP2, pool *bnPool) *gfP2 {
	exp := pool.Get().Rsh(orderGFp2, 1)
	e.Exp(a, exp, pool)
	pool.Put(exp)

	return e
}

func (e *gfP2) Sqrt(a *gfP2, pool *bnPool) *gfP2 {
	t := pool.Get().Rsh(orderGFp2, 1)
	exp, temp := pool.Get().SetInt64(0), pool.Get()

	for i := uint(1); t.Bit(0) == 0; i++ {
		t.Rsh(t, 1)

		x := newGFp2(pool).Invert(nonResidueGFp2, pool)
		x.Exp(x, exp, pool)
		x.Mul(x, a, pool)
		x.Exp(x, t, nil)

		if x.String() != "(0,1)" {
			temp.Lsh(one, i)
			exp.Add(exp, temp)
		}

		x.Put(pool)
	}

	temp.Add(t, one)
	temp.Rsh(temp, 1)

	h := newGFp2(pool).Invert(nonResidueGFp2, pool)
	h.Exp(h, exp, nil)
	h.Mul(h, a, nil)
	h.Exp(h, temp, nil)

	temp.Rsh(exp, 1)

	e.Exp(nonResidueGFp2, temp, nil)
	e.Mul(e, h, nil)

	pool.Put(t)
	pool.Put(exp)
	pool.Put(temp)

	return e
}

func TestGenTwistGen(t *testing.T) {
	// t.Skip()

	tr := make([]*big.Int, 13)
	tr[0] = big.NewInt(2)
	tr[1] = new(big.Int).Add(p, one)
	tr[1].Sub(tr[1], Order)

	for i := 2; i < 13; i++ {
		tr[i] = new(big.Int).Mul(tr[1], tr[i-1])
		temp := new(big.Int).Mul(p, tr[i-2])
		tr[i].Sub(tr[i], temp)
	}

	exp := big.NewInt(1)
	for i := 0; i < 12; i++ {
		exp.Mul(exp, p)
	}
	exp.Add(exp, one).Sub(exp, tr[12])
	exp.Div(exp, Order)
	exp.Div(exp, Order)
	// t.Log(new(big.Int).Mod(exp, Order))
	// return

	for i := 0; i < 100; i++ {
		xx, _ := rand.Int(rand.Reader, p)
		xy, _ := rand.Int(rand.Reader, p)
		x := &gfP2{xx, xy}

		temp := newGFp2(nil).Mul(x, x, nil)
		temp.Mul(temp, x, nil).Add(temp, twistB)
		temp.Minimal()

		jacobi := newGFp2(nil).Jacobi(temp, nil)
		jacobi.Minimal()
		if jacobi.y.BitLen() > 1 {
			// t.Log("qnr")
			continue
		}

		pt := &twistPoint{
			x: x,
			y: newGFp2(nil).Sqrt(temp, nil),
			z: newGFp2(nil).SetOne(),
			t: newGFp2(nil).SetOne(),
		}
		out := newTwistPoint(nil)
		out.Mul(pt, exp, nil)
		if out.IsInfinity() {
			// t.Log("inf")
			continue
		}

		t.Log(pt)
		break
	}

}

func isZero(n *big.Int) bool {
	return new(big.Int).Mod(n, p).Int64() == 0
}

func isOne(n *big.Int) bool {
	return new(big.Int).Mod(n, p).Int64() == 1
}

func TestGFp2Invert(t *testing.T) {
	pool := new(bnPool)

	a := newGFp2(pool)
	a.x.SetString("23423492374", 10)
	a.y.SetString("12934872398472394827398470", 10)

	inv := newGFp2(pool)
	inv.Invert(a, pool)

	b := newGFp2(pool).Mul(inv, a, pool)
	if b.x.Int64() != 0 || b.y.Int64() != 1 {
		t.Fatalf("bad result for a^-1*a: %s %s", b.x, b.y)
	}

	a.Put(pool)
	b.Put(pool)
	inv.Put(pool)

	if c := pool.Count(); c > 0 {
		t.Errorf("Pool count non-zero: %d\n", c)
	}
}

func TestGFp6Invert(t *testing.T) {
	pool := new(bnPool)

	a := newGFp6(pool)
	a.x.x.SetString("239487238491", 10)
	a.x.y.SetString("2356249827341", 10)
	a.y.x.SetString("082659782", 10)
	a.y.y.SetString("182703523765", 10)
	a.z.x.SetString("978236549263", 10)
	a.z.y.SetString("64893242", 10)

	inv := newGFp6(pool)
	inv.Invert(a, pool)

	b := newGFp6(pool).Mul(inv, a, pool)
	if !isZero(b.x.x) ||
		!isZero(b.x.y) ||
		!isZero(b.y.x) ||
		!isZero(b.y.y) ||
		!isZero(b.z.x) ||
		!isOne(b.z.y) {
		t.Fatalf("bad result for a^-1*a: %s", b)
	}

	a.Put(pool)
	b.Put(pool)
	inv.Put(pool)

	if c := pool.Count(); c > 0 {
		t.Errorf("Pool count non-zero: %d\n", c)
	}
}

func TestGFp12Invert(t *testing.T) {
	pool := new(bnPool)

	a := newGFp12(pool)
	a.x.x.x.SetString("239846234862342323958623", 10)
	a.x.x.y.SetString("2359862352529835623", 10)
	a.x.y.x.SetString("928836523", 10)
	a.x.y.y.SetString("9856234", 10)
	a.x.z.x.SetString("235635286", 10)
	a.x.z.y.SetString("5628392833", 10)
	a.y.x.x.SetString("252936598265329856238956532167968", 10)
	a.y.x.y.SetString("23596239865236954178968", 10)
	a.y.y.x.SetString("95421692834", 10)
	a.y.y.y.SetString("236548", 10)
	a.y.z.x.SetString("924523", 10)
	a.y.z.y.SetString("12954623", 10)

	inv := newGFp12(pool)
	inv.Invert(a, pool)

	b := newGFp12(pool).Mul(inv, a, pool)
	if !isZero(b.x.x.x) ||
		!isZero(b.x.x.y) ||
		!isZero(b.x.y.x) ||
		!isZero(b.x.y.y) ||
		!isZero(b.x.z.x) ||
		!isZero(b.x.z.y) ||
		!isZero(b.y.x.x) ||
		!isZero(b.y.x.y) ||
		!isZero(b.y.y.x) ||
		!isZero(b.y.y.y) ||
		!isZero(b.y.z.x) ||
		!isOne(b.y.z.y) {
		t.Fatalf("bad result for a^-1*a: %s", b)
	}

	a.Put(pool)
	b.Put(pool)
	inv.Put(pool)

	if c := pool.Count(); c > 0 {
		t.Errorf("Pool count non-zero: %d\n", c)
	}
}

func TestCurveImpl(t *testing.T) {
	pool := new(bnPool)

	g := &curvePoint{
		pool.Get().SetInt64(1),
		pool.Get().SetInt64(-2),
		pool.Get().SetInt64(1),
		pool.Get().SetInt64(0),
	}

	x := pool.Get().SetInt64(32498273234)
	X := newCurvePoint(pool)
	X.Mul(g, x, pool)

	y := pool.Get().SetInt64(98732423523)
	Y := newCurvePoint(pool)
	Y.Mul(g, y, pool)

	s1 := newCurvePoint(pool)
	s1.Mul(X, y, pool)
	s1.MakeAffine(pool)

	s2 := newCurvePoint(pool)
	s2.Mul(Y, x, pool)
	s2.MakeAffine(pool)

	if s1.x.Cmp(s2.x) != 0 ||
		s2.x.Cmp(s1.x) != 0 {
		t.Errorf("DH points don't match: (%s, %s) (%s, %s)", s1.x, s1.y, s2.x, s2.y)
	}

	pool.Put(x)
	X.Put(pool)
	pool.Put(y)
	Y.Put(pool)
	s1.Put(pool)
	s2.Put(pool)
	g.Put(pool)

	if c := pool.Count(); c > 0 {
		t.Errorf("Pool count non-zero: %d\n", c)
	}
}

func TestOrderG1(t *testing.T) {
	g := new(G1).ScalarBaseMult(Order)
	if !g.p.IsInfinity() {
		t.Error("G1 has incorrect order")
	}

	one := new(G1).ScalarBaseMult(new(big.Int).SetInt64(1))
	g.Add(g, one)
	g.p.MakeAffine(nil)
	if g.p.x.Cmp(one.p.x) != 0 || g.p.y.Cmp(one.p.y) != 0 {
		t.Errorf("1+0 != 1 in G1")
	}
}

func TestOrderG2(t *testing.T) {
	g := new(G2).ScalarBaseMult(Order)
	if !g.p.IsInfinity() {
		t.Error("G2 has incorrect order")
	}

	one := new(G2).ScalarBaseMult(new(big.Int).SetInt64(1))
	g.Add(g, one)
	g.p.MakeAffine(nil)
	if g.p.x.x.Cmp(one.p.x.x) != 0 ||
		g.p.x.y.Cmp(one.p.x.y) != 0 ||
		g.p.y.x.Cmp(one.p.y.x) != 0 ||
		g.p.y.y.Cmp(one.p.y.y) != 0 {
		t.Errorf("1+0 != 1 in G2")
	}
}

func TestOrderGT(t *testing.T) {
	g := new(GT).ScalarBaseMult(Order)
	if !g.p.IsOne() {
		t.Error("GT has incorrect order")
	}
}

func TestBilinearity(t *testing.T) {
	for i := 0; i < 2; i++ {
		a, p1, _ := RandomG1(rand.Reader)
		b, p2, _ := RandomG2(rand.Reader)
		e1 := Pair(p1, p2)

		e2 := Pair(&G1{curveGen}, &G2{twistGen})
		e2.ScalarMult(e2, a)
		e2.ScalarMult(e2, b)

		minusE2 := new(GT).Neg(e2)
		e1.Add(e1, minusE2)

		if !e1.p.IsOne() {
			t.Fatalf("bad pairing result: %s", e1)
		}
	}
}

func TestG1Marshal(t *testing.T) {
	g, g2 := new(G1).ScalarBaseMult(new(big.Int).SetInt64(1)), new(G1)
	form := g.Marshal()
	rest, err := new(G1).Unmarshal(form)
	if err != nil || len(rest) > 0 {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	g.ScalarBaseMult(Order)
	form = g.Marshal()
	rest, err = g2.Unmarshal(form)
	if err != nil || len(rest) > 0 {
		t.Fatalf("failed to unmarshal ∞: %v", err)
	} else if !g2.p.IsInfinity() {
		t.Fatalf("∞ unmarshaled incorrectly")
	}
}

func TestG2Marshal(t *testing.T) {
	g, g2 := new(G2).ScalarBaseMult(new(big.Int).SetInt64(1)), new(G2)
	form := g.Marshal()
	rest, err := new(G2).Unmarshal(form)
	if err != nil || len(rest) > 0 {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	g.ScalarBaseMult(Order)
	form = g.Marshal()
	rest, err = g2.Unmarshal(form)
	if err != nil || len(rest) > 0 {
		t.Fatalf("failed to unmarshal ∞: %v", err)
	} else if !g2.p.IsInfinity() {
		t.Fatalf("∞ unmarshaled incorrectly")
	}
}

func TestG1Identity(t *testing.T) {
	g := new(G1).ScalarBaseMult(new(big.Int).SetInt64(0))
	if !g.p.IsInfinity() {
		t.Error("failure")
	}
}

func TestG2Identity(t *testing.T) {
	g := new(G2).ScalarBaseMult(new(big.Int).SetInt64(0))
	if !g.p.IsInfinity() {
		t.Error("failure")
	}
}

func TestTripartiteDiffieHellman(t *testing.T) {
	a, _ := rand.Int(rand.Reader, Order)
	b, _ := rand.Int(rand.Reader, Order)
	c, _ := rand.Int(rand.Reader, Order)

	pa, pb, pc := new(G1), new(G1), new(G1)
	qa, qb, qc := new(G2), new(G2), new(G2)

	pa.Unmarshal(new(G1).ScalarBaseMult(a).Marshal())
	qa.Unmarshal(new(G2).ScalarBaseMult(a).Marshal())
	pb.Unmarshal(new(G1).ScalarBaseMult(b).Marshal())
	qb.Unmarshal(new(G2).ScalarBaseMult(b).Marshal())
	pc.Unmarshal(new(G1).ScalarBaseMult(c).Marshal())
	qc.Unmarshal(new(G2).ScalarBaseMult(c).Marshal())

	k1 := Pair(pb, qc)
	k1.ScalarMult(k1, a)
	k1Bytes := k1.Marshal()

	k2 := Pair(pc, qa)
	k2.ScalarMult(k2, b)
	k2Bytes := k2.Marshal()

	k3 := Pair(pa, qb)
	k3.ScalarMult(k3, c)
	k3Bytes := k3.Marshal()

	if !bytes.Equal(k1Bytes, k2Bytes) || !bytes.Equal(k2Bytes, k3Bytes) {
		t.Errorf("keys didn't agree")
	}
}

func BenchmarkGT(b *testing.B) {
	n, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		new(GT).ScalarBaseMult(n)
	}
}

func BenchmarkPairing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pair(&G1{curveGen}, &G2{twistGen})
	}
}

func BenchmarkG1(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G1).ScalarBaseMult(x)
	}
}

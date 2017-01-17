package bn256

import (
	"testing"

	"crypto/rand"
	"fmt"
	"math/big"
)

func TestGFPMul(t *testing.T) {
	for i := 0; i < 100000; i++ {
		A, _ := rand.Int(rand.Reader, p)
		B, _ := rand.Int(rand.Reader, p)
		A.Lsh(A, 256).Mod(A, p)
		AB := A.Bits()

		C := new(big.Int).Mul(A, B)
		C.Mod(C, p)
		CB := C.Bits()

		B.Lsh(B, 256).Mod(B, p)
		BB := B.Bits()

		a := &gfP{uint64(AB[0]), uint64(AB[1]), uint64(AB[2]), uint64(AB[3])}
		b := &gfP{uint64(BB[0]), uint64(BB[1]), uint64(BB[2]), uint64(BB[3])}
		c := &gfP{}

		gfpMul(c, a, b)

		if fmt.Sprintf("%x", CB) != fmt.Sprint(c) {
			t.Logf("%x", CB)
			t.Log(c)
			t.Fatal()
		}
	}
}

func TestGFPAdd(t *testing.T) {
	for i := 0; i < 100000; i++ {
		A, _ := rand.Int(rand.Reader, p)
		B, _ := rand.Int(rand.Reader, p)
		AB, BB := A.Bits(), B.Bits()

		a := &gfP{uint64(AB[0]), uint64(AB[1]), uint64(AB[2]), uint64(AB[3])}
		b := &gfP{uint64(BB[0]), uint64(BB[1]), uint64(BB[2]), uint64(BB[3])}
		c := &gfP{}

		gfpAdd(c, a, b)
		C := new(big.Int).Add(A, B)
		C.Mod(C, p)
		CB := C.Bits()

		if fmt.Sprintf("%x", CB) != fmt.Sprint(c) {
			t.Logf("%x", CB)
			t.Log(c)
			t.Fatal()
		}
	}
}

// func TestCurveImpl(t *testing.T) {
// 	g := &curvePoint{}
// 	g.Set(curveGen)
//
// 	x := pool.Get().SetInt64(32498273234)
// 	X := &curvePoint{}
// 	X.Mul(g, x, pool)
//
// 	y := pool.Get().SetInt64(98732423523)
// 	Y := &curvePoint{}
// 	Y.Mul(g, y, pool)
//
// 	s1 := &curvePoint{}
// 	s1.Mul(X, y, pool)
// 	s1.MakeAffine(pool)
//
// 	s2 := &curvePoint{}
// 	s2.Mul(Y, x, pool)
// 	s2.MakeAffine(pool)
//
// 	if s1.x.Cmp(s2.x) != 0 ||
// 		s2.x.Cmp(s1.x) != 0 {
// 		t.Errorf("DH points don't match: (%s, %s) (%s, %s)", s1.x, s1.y, s2.x, s2.y)
// 	}
//
// 	pool.Put(x)
// 	X.Put(pool)
// 	pool.Put(y)
// 	Y.Put(pool)
// 	s1.Put(pool)
// 	s2.Put(pool)
// 	g.Put(pool)
//
// 	if c := pool.Count(); c > 0 {
// 		t.Errorf("Pool count non-zero: %d\n", c)
// 	}
// }

func TestOrderG1(t *testing.T) {
	g := new(G1).ScalarBaseMult(Order)
	if !g.p.IsInfinity() {
		t.Error("G1 has incorrect order")
	}

	one := new(G1).ScalarBaseMult(new(big.Int).SetInt64(1))
	g.Add(g, one)
	g.p.MakeAffine()
	if g.p.x != one.p.x || g.p.y != one.p.y {
		t.Errorf("1+0 != 1 in G1")
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

func TestG1Identity(t *testing.T) {
	g := new(G1).ScalarBaseMult(new(big.Int).SetInt64(0))
	if !g.p.IsInfinity() {
		t.Error("failure")
	}
}

func BenchmarkG1(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G1).ScalarBaseMult(x)
	}
}

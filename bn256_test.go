package bn256

import (
	"testing"

	"crypto/rand"
	"fmt"
	"math/big"
)

func TestGFPAdd(t *testing.T) {
	for i := 0; i < 10000; i++ {
		A, _ := rand.Int(rand.Reader, p)
		B, _ := rand.Int(rand.Reader, p)
		AB, BB := A.Bits(), B.Bits()

		a := &gfP{uint64(AB[0]), uint64(AB[1]), uint64(AB[2]), uint64(AB[3])}
		b := &gfP{uint64(BB[0]), uint64(BB[1]), uint64(BB[2]), uint64(BB[3])}
		c := &gfP{}

		gfpAdd(c, a, b)
		C := new(big.Int).Add(A, B)
		C.Mod(C, p)

		if fmt.Sprintf("%64.64x", C) != fmt.Sprint(c) {
			t.Logf("%64.64x", C)
			t.Log(c)
			t.Fatal()
		}
	}
}

func TestGFPSub(t *testing.T) {
	for i := 0; i < 10000; i++ {
		A, _ := rand.Int(rand.Reader, p)
		B, _ := rand.Int(rand.Reader, p)
		AB, BB := A.Bits(), B.Bits()

		a := &gfP{uint64(AB[0]), uint64(AB[1]), uint64(AB[2]), uint64(AB[3])}
		b := &gfP{uint64(BB[0]), uint64(BB[1]), uint64(BB[2]), uint64(BB[3])}
		c := &gfP{}

		gfpSub(c, a, b)
		C := new(big.Int).Sub(A, B)
		C.Mod(C, p)

		if fmt.Sprintf("%64.64x", C) != fmt.Sprint(c) {
			t.Logf("%64.64x", C)
			t.Log(c)
			t.Fatal()
		}
	}
}

func TestGFPMul(t *testing.T) {
	for i := 0; i < 10000; i++ {
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

		if fmt.Sprintf("%64.64x", C) != fmt.Sprint(c) {
			t.Logf("%64.64x", CB)
			t.Log(c)
			t.Fatal()
		}
	}
}

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

func TestOrderG2(t *testing.T) {
	g := new(G2).ScalarBaseMult(Order)
	if !g.p.IsInfinity() {
		t.Error("G2 has incorrect order")
	}

	one := new(G2).ScalarBaseMult(new(big.Int).SetInt64(1))
	g.Add(g, one)
	g.p.MakeAffine()

	if *g.p != *twistGen {
		t.Errorf("1+0 != 1 in G2")
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

func TestG2Identity(t *testing.T) {
	g := new(G2).ScalarBaseMult(new(big.Int).SetInt64(0))
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

func BenchmarkG2(b *testing.B) {
	x, _ := rand.Int(rand.Reader, Order)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new(G2).ScalarBaseMult(x)
	}
}

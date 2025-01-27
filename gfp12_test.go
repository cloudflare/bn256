package bn256

import (
	"math/big"
	"testing"
)

func TestGfP12SquareCyclo6(t *testing.T) {
	// in MUST be an element of the 6-th cyclotomic group.
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.SquareCyclo6(in)
	expected.Square(in)

	if *got != *expected {
		t.Errorf("not same got=%v, expected=%v", got, expected)
	}
}

func TestGfp12PowToVCyclo6(t *testing.T) {
	// in MUST be an element of the 6-th cyclotomic group.
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.powToVCyclo6(in)
	expected.Exp(in, big.NewInt(1868033))

	if *got != *expected {
		t.Errorf("not same got=%v, expected=%v", got, expected)
	}
}

func TestGfp12PowToUCyclo6(t *testing.T) {
	// in MUST be an element of the 6-th cyclotomic group.
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.PowToUCyclo6(in)
	expected.Exp(in, u)

	if *got != *expected {
		t.Errorf("not same got=%v, expected=%v", got, expected)
	}
}

func BenchmarkGfp12Square(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.Square(gfP12Gen)
	}
}

func BenchmarkGfp12SquareCyclo6(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.SquareCyclo6(gfP12Gen)
	}
}

func BenchmarkGfp12ExpU(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.Exp(gfP12Gen, u)
	}
}

func BenchmarkGfp12PowToUCyclo6(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.PowToUCyclo6(gfP12Gen)
	}
}

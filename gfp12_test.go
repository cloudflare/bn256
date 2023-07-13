package bn256

import (
	"math/big"
	"testing"
)

func TestGfP12SpecialSquare(t *testing.T) {
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.SpecialSquare(in)
	expected.Square(in)

	if *got != *expected {
		t.Errorf("not same got=%v, expected=%v", got, expected)
	}
}

func TestGfp12SpecialPowV(t *testing.T) {
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.SpecialPowV(in)
	expected.Exp(in, big.NewInt(1868033))

	if *got != *expected {
		t.Errorf("not same got=%v, expected=%v", got, expected)
	}
}

func TestGfp12SpecialPowU(t *testing.T) {
	in := gfP12Gen

	got := &gfP12{}
	expected := &gfP12{}

	got.SpecialPowU(in)
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

func BenchmarkGfp12SpecialSquare(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.SpecialSquare(gfP12Gen)
	}
}

func BenchmarkGfp12ExpU(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.Exp(gfP12Gen, u)
	}
}

func BenchmarkGfp12SpecialPowU(b *testing.B) {
	got := &gfP12{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		got.SpecialPowU(gfP12Gen)
	}
}

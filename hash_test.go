package bn256

import (
	"testing"

	"math/rand"
    "strconv"
)

func TestHashCollision(t *testing.T) {
    g := Hash([]byte(strconv.Itoa(rand.Int())))
    h := Hash([]byte(strconv.Itoa(rand.Int())))
    if  *(g.p) == *(h.p) {
		t.Fatal("found a collision of hashes ")
    }
}

func TestHashTAICollision(t *testing.T) {
    g := HashTAI([]byte(strconv.Itoa(rand.Int())))
    h := HashTAI([]byte(strconv.Itoa(rand.Int())))
    if  *(g.p) == *(h.p) {
		t.Fatal("found a collision of hashes ")
    }
}

func BenchmarkHash(b *testing.B) {
    data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
        data[i] = []byte(strconv.Itoa(i))
    }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Hash(data[i])
	}
}

func BenchmarkHashTAI(b *testing.B) {
    data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
        data[i] = []byte(strconv.Itoa(i))
    }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		HashTAI(data[i])
	}
}

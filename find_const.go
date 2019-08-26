package bn256

import (
    "fmt"
    "math/big"
)

func code(e [4]uint64) string {
	return fmt.Sprintf("0x%16.16x, 0x%16.16x, 0x%16.16x, 0x%16.16x", e[0], e[1], e[2], e[3])
}

func bi2a (y *big.Int) [4]uint64{
    eb := new(big.Int).Set(y)
    x := new(big.Int)
    e := [4]uint64{}
    for w := 0; w<4; w++ {
        if eb.IsUint64() {
            e[w] = eb.Uint64()
            break
        }
        e[w] = x.Mod(eb, x.Lsh(new(big.Int).SetUint64(1), 64)).Uint64()
        eb.Rsh(eb, 64)
    }


    return e
}

func Find() {
    pL := new(big.Int).Set(p)
    one := new(big.Int).SetUint64(1)
    fmt.Println("(p-1)/2")
    pL.Sub(pL, one)
    fmt.Println(pL)
    pL.Rsh(pL, 1)
    fmt.Println(pL)
    fmt.Println(code(bi2a(pL)))
    fmt.Println("(s+1)/2")
    pL.Add(pL, one)
    pL.Rsh(pL, 1)
    fmt.Println(pL)
    fmt.Println(code(bi2a(pL)))
    fmt.Println()

    fmt.Println("sqrt(-3)")
    e := newGFp(-3)
    e.Sqrt(e)
    fmt.Println(code([4]uint64(*e)))
}

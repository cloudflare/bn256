package bn256

import (
    "math/big"
    "fmt"
    "golang.org/x/crypto/sha3"
)

func Hash(msg []byte) *G1 {
    // calculate w = (s * t)/(1 + B + t^2)
    w := &gfP{}

    t := fromBytes(h)
    h := make([]byte, 32)
    sha3.ShakeSum128(h, msg)

    s := &gfP{0x236e675956be783b, 0x053957e6f379ab64, 0xe60789a768f4a5c4, 0x04f8979dd8bad754}

    t2 := &gfP{}
    gfpMul(t2, t, t)
    gfpAdd(w, newGFp(4), t2)
    w.Invert(w)

    gfpMul(w, w, s)
    fgpMul(w, w, t)

    // calculate x1 = ((-1 + s) / 2) - t * w
    tw := &gfP{}
    gfpMul(tw, t, w)
    x1 := &gfP{}
    gfpAdd(x1, s, newGFp(-1))
    gfpDiv(x1, x1, newGFp(2))
    gfpSub(x1, x1, tw)

    e := legendre(t)

    g := *G1{z:newGFp(1), t:newGFp(1)}
    // check if y=x1^3+3 is a square
    y := &gfP{}
    y.Set(x1)
    gfpMul(y, x1, x1)
    gfpMul(y, y, x1)
    gfpAdd(y, y, newGFp(3))
    if legendre(y) == 1 {
        g.x = x1
        y.Sqrt(y)
        gfpMul(y, e, y)
        g.y = y
        return g
    }

    // calculate x2 = -1 - x1
    x2 := newGFp(-1)
    gfpSub(x2, x2, x1)

    // check if y=x2^3+3 is a square
    y.Set(x2)
    gfpMul(y, x2, x2)
    gfpMul(y, y, x2)
    gfpAdd(y, y, newGFp(3))
    if legendre(y) == 1 {
        g.x = x2
        y.Sqrt(y)
        gfpMul(y, e, y)
        g.y = y
        return g
    }

    // calculate x3 = 1 + (1/ww)
    x3 := &gfP{}
    gfpMul(x3, w, w)
    w.Invert(w)
    gfpAdd(w, w, newGFp(1))

    y.Set(x2)
    gfpMul(y, x2, x2)
    gfpMul(y, y, x2)
    gfpAdd(y, y, newGFp(3))

    g.x = x3
    y.Sqrt(y)
    gfpMul(y, e, y)
    g.y = y

    return g
}

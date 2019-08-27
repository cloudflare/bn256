package bn256

import (
    "golang.org/x/crypto/sha3"
)

func Hash(msg []byte) *G1 {
    // calculate w = (s * t)/(1 + B + t^2)
    w := &gfP{}

    h := make([]byte, 32)
    sha3.ShakeSum128(h, msg)
    t := fromBytes(h)
    montEncode(t, t)

    s := &gfP{0x236e675956be783b, 0x053957e6f379ab64, 0xe60789a768f4a5c4, 0x04f8979dd8bad754}

    t2 := &gfP{}
    gfpMul(t2, t, t)
    gfpAdd(w, curveB, t2)
    gfpAdd(w, w, &one)
    w.Invert(w)

    gfpMul(w, w, s)
    gfpMul(w, w, t)

    // calculate x1 = ((-1 + s) / 2) - t * w
    tw := &gfP{}
    gfpMul(tw, t, w)
    x1 := &gfP{}
    gfpAdd(x1, s, newGFp(-1))
    half := newGFp(2)
    half.Invert(half)
    gfpMul(x1, x1, half)
    gfpSub(x1, x1, tw)

    e := legendre(t)

    cp := &curvePoint{z:one, t:one}
    // check if y=x1^3+3 is a square
    y := &gfP{}
    y.Set(x1)
    gfpMul(y, x1, x1)
    gfpMul(y, y, x1)
    gfpAdd(y, y, curveB)
    if legendre(y) == 1 {
        cp.x = *x1
        y.Sqrt(y)
        if e == -1 {
            gfpNeg(y, y)
        }
        cp.y = *y
        return &G1{cp}
    }

    // calculate x2 = -1 - x1
    x2 := newGFp(-1)
    gfpSub(x2, x2, x1)

    // check if y=x2^3+3 is a square
    y.Set(x2)
    gfpMul(y, x2, x2)
    gfpMul(y, y, x2)
    gfpAdd(y, y, curveB)
    if legendre(y) == 1 {
        cp.x = *x2
        y.Sqrt(y)
        if e == -1 {
            gfpNeg(y, y)
        }
        cp.y = *y
        return &G1{cp}
    }

    // calculate x3 = 1 + (1/ww)
    x3 := &gfP{}
    gfpMul(x3, w, w)
    w.Invert(w)
    gfpAdd(w, w, &one)

    y.Set(x3)
    gfpMul(y, x3, x3)
    gfpMul(y, y, x3)
    gfpAdd(y, y, curveB)

    cp.x = *x3
    y.Sqrt(y)
    if e == -1 {
        gfpNeg(y, y)
    }
    cp.y = *y

    return &G1{cp}
}

func TryAndIncrement(msg []byte) *G1 {
    h := make([]byte, 32)
    sha3.ShakeSum128(h, msg)
    x := fromBytes(h)
    montEncode(x, x)
    for {
        t := &gfP{}
        gfpMul(t, x, x)
        gfpMul(t, t, x)
        gfpAdd(t, t, curveB)
        if legendre(t) == 1 {
            y := &gfP{}
            y.Sqrt(t)
            return &G1{&curvePoint{*x, *y, one, one}}
        }
        gfpAdd(x, x, &one)
    }
}

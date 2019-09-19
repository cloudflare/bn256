package bn256

// HashG1 implements a hashing function into the G1 group. It uses Fouque-Tibouchi encoding described here
// https://tools.ietf.org/pdf/draft-irtf-cfrg-hash-to-curve-04.pdf
func HashG1(msg []byte) *G1 {
	return mapToCurve(hashToBase(msg))
}

func mapToCurve(t *gfP) *G1 {
	one := *newGFp(1)

	// calculate w = (s * t)/(1 + B + t^2)
	w := &gfP{}

	t2 := &gfP{}
	gfpMul(t2, t, t)
	gfpAdd(w, curveB, t2)
	gfpAdd(w, w, &one)
	w.Invert(w)

	gfpMul(w, w, s)
	gfpMul(w, w, t)

	e := legendre(t)
	cp := &curvePoint{z: one, t: one}

	// calculate x1 = ((-1 + s) / 2) - t * w
	tw := &gfP{}
	gfpMul(tw, t, w)
	x1 := &gfP{}
	gfpSub(x1, sMinus1Over2, tw)

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

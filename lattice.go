package bn256

import (
	"math/big"
)

var half = new(big.Int).Rsh(Order, 1)

var curveLattice = &lattice{
	vectors: [][]*big.Int{
		[]*big.Int{bigFromBase10("254952053719217182009119236802174855688"), bigFromBase10("254952053719217182022156415784332439563")},
		[]*big.Int{bigFromBase10("254952053719217182035193594766490023438"), bigFromBase10("-254952053719217181996082057820017271813")},
	},
	inverse: []*big.Int{
		bigFromBase10("254952053719217181996082057820017271813"),
		bigFromBase10("254952053719217182022156415784332439563"),
	},
	det: bigFromBase10("130001099391293207465592877484719811485140812107807572779762125938088333599938"),
}

var targetLattice = &lattice{
	vectors: [][]*big.Int{
		[]*big.Int{bigFromBase10("13037178982157583874"), bigFromBase10("13037178982157583875"), bigFromBase10("13037178982157583875"), bigFromBase10("13037178982157583875")},
		[]*big.Int{bigFromBase10("6518589491078791938"), bigFromBase10("6518589491078791937"), bigFromBase10("6518589491078791937"), bigFromBase10("-13037178982157583874")},
		[]*big.Int{bigFromBase10("13037178982157583875"), bigFromBase10("-6518589491078791937"), bigFromBase10("-6518589491078791938"), bigFromBase10("-6518589491078791937")},
		[]*big.Int{bigFromBase10("6518589491078791936"), bigFromBase10("26074357964315167750"), bigFromBase10("-13037178982157583873"), bigFromBase10("6518589491078791936")},
	},
	inverse: []*big.Int{
		bigFromBase10("1661927778103044753630928116985183243902389359598104203531"),
		bigFromBase10("84984017906405727351583121079908799750"),
		bigFromBase10("3323855556206089507261856233970366487798260129705129615125"),
		bigFromBase10("-84984017906405727338545942097751215875"),
	},
	det: new(big.Int).Set(Order),
}

type lattice struct {
	vectors [][]*big.Int
	inverse []*big.Int
	det     *big.Int
}

// decompose takes a scalar mod Order as input and finds a short, positive decomposition of it wrt to the lattice basis.
func (l *lattice) decompose(k *big.Int) []*big.Int {
	n := len(l.inverse)

	// Calculate closest vector in lattice to <k,0,0,...> with Babai's rounding.
	c := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		c[i] = new(big.Int).Mul(k, l.inverse[i])
		round(c[i], l.det)
	}

	// Transform vectors according to c and subtract <k,0,0,...>.
	out := make([]*big.Int, n)
	temp := new(big.Int)

	for i := 0; i < n; i++ {
		out[i] = new(big.Int)

		for j := 0; j < n; j++ {
			temp.Mul(c[j], l.vectors[j][i])
			out[i].Add(out[i], temp)
		}

		out[i].Neg(out[i])
		out[i].Add(out[i], l.vectors[0][i]).Add(out[i], l.vectors[0][i])
	}
	out[0].Add(out[0], k)

	return out
}

func (l *lattice) Precompute(add func(i, j uint)) {
	n := uint(len(l.vectors))
	total := uint(1) << uint(n)

	for i := uint(0); i < n; i++ {
		for j := uint(0); j < total; j++ {
			if (j>>i)&1 == 1 {
				add(i, j)
			}
		}
	}
}

func (l *lattice) Multi(scalar *big.Int) []uint8 {
	decomp := l.decompose(scalar)

	maxLen := 0
	for _, x := range decomp {
		if x.BitLen() > maxLen {
			maxLen = x.BitLen()
		}
	}

	out := make([]uint8, maxLen)
	for j, x := range decomp {
		for i := 0; i < maxLen; i++ {
			out[i] += uint8(x.Bit(i)) << uint(j)
		}
	}

	return out
}

// round sets num to num/denom rounded to the nearest integer.
func round(num, denom *big.Int) {
	r := new(big.Int)
	num.DivMod(num, denom, r)

	if r.Cmp(half) == 1 {
		num.Add(num, big.NewInt(1))
	}
}

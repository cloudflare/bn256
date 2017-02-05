package bn256

import (
	"math/big"
)

func bigFromBase10(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

// p is a prime over which we form a basic field: 36u⁴+36u³+24u²+6u+1.
var p = bigFromBase10("65000549695646603732796438742359905742825358107623003571877145026864184071783")

var p2 = [4]uint64{0x185cac6c5e089667, 0xee5b88d120b5b59e, 0xaa6fecb86184dc21, 0x8fb501e34aa387f9}
var np = [4]uint64{0x2387f9007f17daa9, 0x734b3343ab8513c8, 0x2524282f48054c12, 0x38997ae661c3ef3c}

var rN1 = &gfP{0xcbb781e36236117d, 0xcc65f3bcec8c91b, 0x2eab68888ea1f515, 0x1fc5c0956f92f825}

// var r = &gfP{0xe7a35393a1f76999, 0x11a4772edf4a4a61, 0x559013479e7b23de, 0x704afe1cb55c7806}
var r2 = &gfP{0x9c21c3ff7e444f56, 0x409ed151b2efb0c2, 0xc6dc37b80fb1651, 0x7c36e0e62c2380b7}
var r3 = &gfP{0x2af2dfb9324a5bb8, 0x388f899054f538a4, 0xdf2ff66396b107a7, 0x24ebbbb3a2529292}

// Order is the number of elements in both G₁ and G₂: 36u⁴+36u³+18u²+6u+1.
// order-1 = (2**5) * 3 * 5743 * 280941149 * 130979359433191 * 491513138693455212421542731357 * 6518589491078791937
var Order = bigFromBase10("65000549695646603732796438742359905742570406053903786389881062969044166799969")

var xiToPMinus1Over3 = &gfP2{
	gfP{0x4f59e37c01832e57, 0xae6be39ac2bbbfe4, 0xe04ea1bb697512f8, 0x3097caa8fc40e10e},
	gfP{0xf8606916d3816f2c, 0x1e5c0d7926de927e, 0xbc45f3946d81185e, 0x80752a25aa738091},
}

var xiTo2PMinus2Over3 = &gfP2{
	gfP{0x51678e7469b3c52a, 0x4fb98f8b13319fc9, 0x29b2254db3f1df75, 0x1c044935a3d22fb2},
	gfP{0x4d2ea218872f3d2c, 0x2fcb27fc4abe7b69, 0xd31d972f0e88ced9, 0x53adc04a00a73b15},
}

var xiToPSquaredMinus1Over3 = &gfP{0x12d3cef5e1ada57d, 0xe2eca1463753babb, 0xca41e40ddccf750, 0x551337060397e04c}

var xiTo2PSquaredMinus2Over3 = &gfP{0x3642364f386c1db8, 0xe825f92d2acd661f, 0xf2aba7e846c19d14, 0x5a0bcea3dc52b7a0}

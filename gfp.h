#define gfpReduce(a1,a2,a3,a4,a5, b1,b2,b3,b4,b5) \
	\ // b = a-p
	MOVQ a1, b1 \
	MOVQ a2, b2 \
	MOVQ a3, b3 \
	MOVQ a4, b4 \
	MOVQ a5, b5 \
	\
	SUBQ ·p2+0(SB), b1 \
	SBBQ ·p2+8(SB), b2 \
	SBBQ ·p2+16(SB), b3 \
	SBBQ ·p2+24(SB), b4 \
	SBBQ $0, b5 \
	\
	\ // if b is negative then return a
	\ // else return a
	CMOVQCC b1, a1 \
	CMOVQCC b2, a2 \
	CMOVQCC b3, a3 \
	CMOVQCC b4, a4

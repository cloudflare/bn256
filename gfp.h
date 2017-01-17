#define storeBlock(a1,a2,a3,a4, r) \
	MOVQ a1,  0+r \
	MOVQ a2,  8+r \
	MOVQ a3, 16+r \
	MOVQ a4, 24+r

#define loadBlock(r, a1,a2,a3,a4) \
	MOVQ  0+r, a1 \
	MOVQ  8+r, a2 \
	MOVQ 16+r, a3 \
	MOVQ 24+r, a4

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

#define mul(ra, rb) \
	mulArb(0+ra,8+ra,16+ra,24+ra, rb)

#define mulArb(a1,a2,a3,a4, rb) \
	MOVQ a1, DX \
	MULXQ 0+rb, R8, R9 \
	MULXQ 8+rb, AX, R10 \
	ADDQ AX, R9 \
	MULXQ 16+rb, AX, R11 \
	ADCQ AX, R10 \
	MULXQ 24+rb, AX, R12 \
	ADCQ AX, R11 \
	ADCQ $0, R12 \
	MOVQ $0, R13 \
	ADCQ $0, R13 \
	\
	MOVQ a2, DX \
	MULXQ 0+rb, AX, BX \
	ADDQ AX, R9 \
	ADCQ BX, R10 \
	ADCQ $0, R11 \
	ADCQ $0, R12 \
	ADCQ $0, R13 \
	MULXQ 8+rb, AX, BX \
	ADDQ AX, R10 \
	ADCQ BX, R11 \
	ADCQ $0, R12 \
	ADCQ $0, R13 \
	MULXQ 16+rb, AX, BX \
	ADDQ AX, R11 \
	ADCQ BX, R12 \
	ADCQ $0, R13 \
	MULXQ 24+rb, AX, BX \
	ADDQ AX, R12 \
	ADCQ BX, R13 \
	MOVQ $0, R14 \
	ADCQ $0, R14 \
	\
	MOVQ a3, DX \
	MULXQ 0+rb, AX, BX \
	ADDQ AX, R10 \
	ADCQ BX, R11 \
	ADCQ $0, R12 \
	ADCQ $0, R13 \
	ADCQ $0, R14 \
	MULXQ 8+rb, AX, BX \
	ADDQ AX, R11 \
	ADCQ BX, R12 \
	ADCQ $0, R13 \
	ADCQ $0, R14 \
	MULXQ 16+rb, AX, BX \
	ADDQ AX, R12 \
	ADCQ BX, R13 \
	ADCQ $0, R14 \
	MULXQ 24+rb, AX, BX \
	ADDQ AX, R13 \
	ADCQ BX, R14 \
	MOVQ $0, R15 \
	ADCQ $0, R15 \
	\
	MOVQ a4, DX \
	MULXQ 0+rb, AX, BX \
	ADDQ AX, R11 \
	ADCQ BX, R12 \
	ADCQ $0, R13 \
	ADCQ $0, R14 \
	ADCQ $0, R15 \
	MULXQ 8+rb, AX, BX \
	ADDQ AX, R12 \
	ADCQ BX, R13 \
	ADCQ $0, R14 \
	ADCQ $0, R15 \
	MULXQ 16+rb, AX, BX \
	ADDQ AX, R13 \
	ADCQ BX, R14 \
	ADCQ $0, R15 \
	MULXQ 24+rb, AX, BX \
	ADDQ AX, R14 \
	ADCQ BX, R15

// use stack to effect.
// Pushw res, s1, flags
// popf
// pop %b
// pop %c
// be lazy: just push them again
// pushf
// push %c
// pus %b
// pushf
// op
// pus res
// pushf
// stack now has oflags, res, 2 inargs, iflags,
// push print string
// For the test output, we always push and test the 32-bit register.
#ifndef OP_SHIFTD
#ifdef OP_NOBYTE

// NOBYTE means there is no byte opcode. weird.
#define EXECSHIFT(o, size, rsize, bits, pre, res, s1, s2, flags)	\
	movw	$flags, %dx ;\
	pushw %dx ;\
	popf; \
	OPR(mov,l) $res, REG(a, e, x);	\
	PUSH(a,e) ;\
	OPR(mov,l) $s1, REG(b, e, x) ;	\
	PUSH(b,e) ;\
	OPR(o,size) REG(b,pre, rsize), REG(a,pre, rsize) ;	\
	PUSH(a,e) ;					\
	pushf ;\
	movw	$flags, %dx ;\
	cmpb $1, %al;\
	je 1f; \
	andw $0x800, %dx;		\
1:\
	pushw %dx ;\
	hlt ;\
	.byte 2; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 3; \
	.asciz #o ;							\
	.asciz "%s%s A=%08x B=%08x R=%08x CCIN=%04x CC=%04x" ;

#else
#define EXECSHIFT(o, size, rsize, bits, pre, res, s1, s2, flags) \
	movw	$flags, %dx ;\
	pushw %dx ;\
	popf; \
	OPR(mov,l) $res, REG(a, e, x);	\
	PUSH(a,e) ;\
	OPR(mov,l) $s1, REG(c, e, x) ;	\
	PUSH(c,e) ;\
	OPR(o,size) %cl, REG(a,pre, rsize) ;	\
	PUSH(a,e) ;					\
	pushf ;\
	movw	$flags, %dx ;\
	cmpb $1, %al;\
	je 1f; \
	andw $0x800, %dx;		\
1: \
	pushw %dx ;\
	hlt ;\
	.byte 2; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 3; \
	.asciz #o ;							\
	.asciz "%s%s A=%08x B=%08x R=%08x CCIN=%04x CC=%04x" ;

#endif

#else
#define EXECSHIFT(o, size, rsize, bits, pre, res, s1, s2, flags) \
	movw	$flags, %dx ;\
	pushw %dx ;\
	popf; \
	OPR(mov,l) $res, REG(a, e, x);	\
	PUSH(a,e) ;\
	OPR(mov,l) $s1, REG(c, e, x) ;	\
	PUSH(c,e) ;\
	OPR(o,size) %cl, REG(a,pre, rsize) ;	\
	 PUSH(b,e) ; \
	PUSH(a,e) ;\
	pushf ;\
	movw	$flags, %dx ;\
	cmpb $1, %al;\
	je 1f; \
	andw $0x800, %dx;		\
1: \
	pushw %dx ;\
	hlt ;\
	.byte 3; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 4; \
	.asciz #o ;							\
	.asciz "%s%s A=%08x B=%08x C=%08lx R=%08x CCIN=%04x CC=%04x" ;
#endif

#define exec_opl(o, s2, s0, s1, iflags) EXECSHIFT(o,l, x, 32, e, s0, s1, s2, iflags)
#define exec_opw(o, s2, s0, s1, iflags) EXECSHIFT(o,w, x, 16,  , s0, s1, s2, iflags)
#define exec_opb(o, s2, s0, s1, iflags) EXECSHIFT(o,b, x,  8,  , s0, s1, s2, iflags)
#if 0
void glue(test_, OP)(void)
{
    int i, n;
    n = 32;

    for(i = 0; i < n; i++)
        exec_op(0x21ad3d34, 0x12345678, i);
    for(i = 0; i < n; i++)
        exec_op(0x813f3421, 0x82345679, i);
}
#endif

#define exec_op(op,s2, s0, amt) \
    exec_opl(op, s2, s0, s1, 0); \
#ifdef OP_SHIFTD \
    exec_opw(op, s2, s0, s1, 0); \
#else \
    exec_opw(op, s2, s0, s1, 0); \
#endif \
#ifndef OP_NOBYTE \
    exec_opb(op, s2, s0, s1, 0); \
#endif \
#ifdef OP_CC \
    exec_opl(op, s2, s0, s1, CC_C); \
    exec_opw(op, s2, s0, s1, CC_C); \
    exec_opb(op, s2, s0, s1, CC_C); \

exec_op(OP, 0x21ad3d34, 0x12345678, 0);
exec_op(OP, 0x21ad3d34, 0x12345678, 1);


#undef OP
#undef OP_CC
#undef OP_SHIFTD
#undef OP_NOBYTE
#undef EXECSHIFT

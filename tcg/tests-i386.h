#ifndef OPCODES
#define REG(r,pre,sz) %pre##r##sz
#define OPR(o,sz) o##sz
#define PUSH(r,pre) push %pre##r##x
#endif
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
#define EXECOP2(o, size, rsize, bits, pre, res, s1, flags)	\
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
	pushw %dx ;\
	hlt ;\
	.byte 2; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 3; \
	.asciz #o ;							\
	.asciz "%s%s A=%08x B=%08x R=%08x CCIN=%04x CC=%04x" ;

#define EXECOP1(o, size, rsize, bits, pre, res, flags)	\
	movw	$flags, %dx ;\
	pushw %dx ;\
	popf; \
	OPR(mov,l) $res, REG(a, e, x);	\
	PUSH(a,e) ;\
	OPR(o,size) REG(a,pre, rsize) ;	\
	PUSH(a,e) ;					\
	pushf ;\
	movw	$flags, %dx ;\
	pushw %dx ;\
	hlt ;\
	.byte 2; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 2; \
	.asciz #o ;							\
	.asciz "%s%s A=%08x R=%08x CCIN=%04x CC=%04x" ;


#ifdef OP1

#define exec_opl(o,s0, s1, iflags) EXECOP1(o,l, x, 32, e, s0, iflags)
#define exec_opw(o,s0, s1, iflags) EXECOP1(o,w, x, 16,  , s0, iflags)
#define exec_opb(o,s0, s1, iflags) EXECOP1(o,b, l,  8,  , s0, iflags)

#else
#define exec_opl(o,s0, s1, iflags)  EXECOP2(o,l, x, 32, e, s0, s1, iflags)
#define exec_opw(o,s0, s1, iflags)  EXECOP2(o,w, x, 16,  , s0, s1, iflags)
#define exec_opb(o,s0, s1, iflags)  EXECOP2(o,b, l,  8,  , s0, s1, iflags)
#endif


#ifdef OP_CC 
#define exec_op(s0, s1) \
    exec_opl(OP,s0, s1, 0) \
    exec_opw(OP,s0, s1, 0) \
    exec_opb(OP,s0, s1, 0) \
    exec_opl(OP,s0, s1, CC_C) \
    exec_opw(OP,s0, s1, CC_C) \
    exec_opb(OP,s0, s1, CC_C)  
#else
#define exec_op(s0, s1) \
    exec_opl(OP,s0, s1, 0) \
    exec_opw(OP,s0, s1, 0) \
    exec_opb(OP,s0, s1, 0) 
#endif

# fuck cpp.
.code16
    exec_op(0x12345678, 0x812FADA);
    exec_op(0x12341, 0x12341);
    exec_op(0x12341, -0x12341);
    exec_op(0xffffffff, 0);
    exec_op(0xffffffff, -1);
    exec_op(0xffffffff, 1);
    exec_op(0xffffffff, 2);
    exec_op(0x7fffffff, 0);
    exec_op(0x7fffffff, 1);
    exec_op(0x7fffffff, -1);
    exec_op(0x80000000, -1);
    exec_op(0x80000000, 1);
    exec_op(0x80000000, -2);
    exec_op(0x12347fff, 0);
    exec_op(0x12347fff, 1);
    exec_op(0x12347fff, -1);
    exec_op(0x12348000, -1);
    exec_op(0x12348000, 1);
    exec_op(0x12348000, -2);
    exec_op(0x12347f7f, 0);
    exec_op(0x12347f7f, 1);
    exec_op(0x12347f7f, -1);
    exec_op(0x12348080, -1);
    exec_op(0x12348080, 1);
    exec_op(0x12348080, -2);

#undef OP
#undef OP_CC
#undef result
#undef TEST
#undef TNAME
#undef OP1
#undef exec_opl
#undef exec_opw
#undef exec_opb
#undef exec_op

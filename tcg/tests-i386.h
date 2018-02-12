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
#define EXECOP2(o, size, rsize, bits, pre, res, s1, flags)	\
	movw	$flags, %ax ;\
	pushw %ax ;\
	popf; \
	OPR(mov,size) $res, REG(c, pre, rsize);	\
	PUSH(c,pre) ;\
	OPR(mov,size) $s1, REG(b, pre, rsize) ;	\
	PUSH(b,pre) ;\
	OPR(o,size) REG(b,pre, rsize), REG(c,pre, rsize) ;	\
	PUSH(c,pre) ;					\
	pushf ;\
	hlt ;\
	.byte 2; /* number of following bytes of info */ \
	/* currently # bits per stack item, and nargs */ \
	.byte bits, 3; \
	.asciz #o ;							\
	.asciz "%-10s A=%#x B=%#x R=%#x CCIN=%04x CC=%04x" ;

#define EXECOP1(size, rsize, pre, res, flags)	\

#if 0
#define EXECOP1(size, rsize, res, flags) \
    asm (".code16\n\npush %3\n\t"\
         "popf\n\t"\
         stringify(OP) size " %" rsize "0\n\t" \
         "pushf\n\t"\
         "pop %1\n\t"\
         : "=q" (res), "=g" (flags)\
	 : "0" (res), "1" (flags));\
    result(6, "%-10s A=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",	\
           stringify(OP) size, s0, res, iflags, flags & CC_MASK, 0); \
	asm ("hlt\n")

#endif
#ifdef OP1

#define exec_opl(o,s0, s1, iflags) EXECOP1(o,l, x, 32, e, res, iflags)
#define exec_opw(o,s0, s1, iflags) EXECOP1(o,w, x, 16,  , res, iflags)
#define exec_opb(o,s0, s1, iflags) EXECOP1(o,b, l,  8,  , res, iflags)

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
    .global TNAME
TNAME:
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

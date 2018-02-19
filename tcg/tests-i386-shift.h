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


#ifndef OP_SHIFTD

#ifdef OP_NOBYTE
#define EXECSHIFT(size, rsize, res, s1, s2, flags) \
    asm ("push %4\n\t"\
         "popf\n\t"\
         stringify(OP) size " %" rsize "2, %" rsize "0\n\t" \
         "pushf\n\t"\
         "pop %1\n\t"\
         : "=g" (res), "=g" (flags)\
         : "r" (s1), "0" (res), "1" (flags));
#else
#define EXECSHIFT(size, rsize, res, s1, s2, flags) \
    asm ("push %4\n\t"\
         "popf\n\t"\
         stringify(OP) size " %%cl, %" rsize "0\n\t" \
         "pushf\n\t"\
         "pop %1\n\t"\
         : "=q" (res), "=g" (flags)\
         : "c" (s1), "0" (res), "1" (flags));
#endif

void exec_opl(long s2, long s0, long s1, long iflags)
{
    long res, flags;
    res = s0;
    flags = iflags;
    EXECSHIFT("l", "k", res, s1, s2, flags);
    /* overflow is undefined if count != 1 */
    if (s1 != 1)
      flags &= ~CC_O;
    result(7, "%-10s A=" FMTLX " B=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",
           stringify(OP) "l", s0, s1, res, iflags, flags & CC_MASK);
}

void exec_opw(long s2, long s0, long s1, long iflags)
{
    long res, flags;
    res = s0;
    flags = iflags;
    EXECSHIFT("w", "w", res, s1, s2, flags);
    /* overflow is undefined if count != 1 */
    if (s1 != 1)
      flags &= ~CC_O;
    result(7, "%-10s A=" FMTLX " B=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",
           stringify(OP) "w", s0, s1, res, iflags, flags & CC_MASK);
}

#else
#define EXECSHIFT(size, rsize, res, s1, s2, flags) \
    asm ("push %4\n\t"\
         "popf\n\t"\
         stringify(OP) size " %%cl, %" rsize "5, %" rsize "0\n\t" \
         "pushf\n\t"\
         "pop %1\n\t"\
         : "=g" (res), "=g" (flags)\
         : "c" (s1), "0" (res), "1" (flags), "r" (s2));

void exec_opl(long s2, long s0, long s1, long iflags)
{
    long res, flags;
    res = s0;
    flags = iflags;
    EXECSHIFT("l", "k", res, s1, s2, flags);
    /* overflow is undefined if count != 1 */
    if (s1 != 1)
      flags &= ~CC_O;
    result9(8, "%-10s A=" FMTLX " B=" FMTLX " C=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",
           stringify(OP) "l", s0, s2, s1, res, iflags, flags & CC_MASK);
}

void exec_opw(long s2, long s0, long s1, long iflags)
{
    long res, flags;
    res = s0;
    flags = iflags;
    EXECSHIFT("w", "w", res, s1, s2, flags);
    /* overflow is undefined if count != 1 */
    if (s1 != 1)
      flags &= ~CC_O;
    result9(8, "%-10s A=" FMTLX " B=" FMTLX " C=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",
           stringify(OP) "w", s0, s2, s1, res, iflags, flags & CC_MASK);
}

#endif

#ifndef OP_NOBYTE
void exec_opb(long s0, long s1, long iflags)
{
    long res, flags;
    res = s0;
    flags = iflags;
    EXECSHIFT("b", "b", res, s1, 0, flags);
    /* overflow is undefined if count != 1 */
    if (s1 != 1)
      flags &= ~CC_O;
    result(7, "%-10s A=" FMTLX " B=" FMTLX " R=" FMTLX " CCIN=%04lx CC=%04lx\n",
           stringify(OP) "b", s0, s1, res, iflags, flags & CC_MASK);
}
#endif

void exec_op(long s2, long s0, long s1)
{
    s2 = i2l(s2);
    s0 = i2l(s0);
    exec_opl(s2, s0, s1, 0);
#ifdef OP_SHIFTD
    exec_opw(s2, s0, s1, 0);
#else
    exec_opw(s2, s0, s1, 0);
#endif
#ifndef OP_NOBYTE
    exec_opb(s0, s1, 0);
#endif
#ifdef OP_CC
    exec_opl(s2, s0, s1, CC_C);
    exec_opw(s2, s0, s1, CC_C);
    exec_opb(s0, s1, CC_C);
#endif
}

void glue(test_, OP)(void)
{
    int i, n;
    n = 32;

    for(i = 0; i < n; i++)
        exec_op(0x21ad3d34, 0x12345678, i);
    for(i = 0; i < n; i++)
        exec_op(0x813f3421, 0x82345679, i);
}

#undef OP
#undef OP_CC
#undef OP_SHIFTD
#undef OP_NOBYTE
#undef EXECSHIFT

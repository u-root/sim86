/*
 *  x86 CPU test
 *
 *  Copyright (c) 2003 Fabrice Bellard
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program; if not, see <http://www.gnu.org/licenses/>.
 */
#define xglue(x, y) x ## y
#define glue(x, y) xglue(x, y)
#define stringify(s)	tostring(s)
#define tostring(s)	#s

#define CC_C   	0x0001
#define CC_P 	0x0004
#define CC_A	0x0010
#define CC_Z	0x0040
#define CC_S    0x0080
#define CC_O    0x0800

#undef __x86_64__
#define FMTLX "%08lx"
#define __init_call	__attribute__ ((unused,__section__ ("initcall")))

#define CC_MASK (CC_C | CC_P | CC_Z | CC_S | CC_O | CC_A)

.code16
#define OP add
.globl test_add
test_add:
#include "tests-i386.h"

#define OP sub
.global test_sub
test_sub:
#include "tests-i386.h"

#define OP xor
.global test_xor
test_xor:
#include "tests-i386.h"

#define OP and
.global test_and
test_and:
#include "tests-i386.h"

#define OP or
.global test_or
test_or:
#include "tests-i386.h"

#define OP cmp
.global test_cmp
test_cmp:
#include "tests-i386.h"

#define OP adc
.global test_adc
test_adc:
#define OP_CC
#include "tests-i386.h"

#if 0
#define OP sbb
.global test_sbb
test_sbb:
#define OP_CC
#include "tests-i386.h"
#endif

#define OP inc
.global test_inc
test_inc:
#define OP_CC
#define OP1
#include "tests-i386.h"

#define OP dec
.global test_dec
test_dec:
#define OP_CC
#define OP1
#include "tests-i386.h"

#define OP neg
.global test_neg
test_neg:
#define OP_CC
#define OP1
#include "tests-i386.h"

#define OP not
.global test_not
test_not:
#define OP_CC
#define OP1
#include "tests-i386.h"

#undef CC_MASK
#define CC_MASK (CC_C | CC_P | CC_Z | CC_S | CC_O)

#if 0
#define OP shl
.global test_shl
test_shl:
#include "tests-i386-shift.h"

#define OP shr
.global test_shr
test_shr:
#include "tests-i386-shift.h"

#define OP sar
.global test_sar
test_sar:
#include "tests-i386-shift.h"

#define OP rol
.global test_rol
test_rol:
#include "tests-i386-shift.h"

#define OP ror
.global test_ror
test_ror:
#include "tests-i386-shift.h"

#define OP rcr
.global test_rcr
test_rcr:
#define OP_CC
#include "tests-i386-shift.h"

#define OP rcl
.global test_rcl
test_rcl:
#define OP_CC
#include "tests-i386-shift.h"

#define OP shld
.global test_shld
test_shld:
#define OP_SHIFTD
#define OP_NOBYTE
#include "tests-i386-shift.h"

#define OP shrd
.global test_shrd
test_shrd:
#define OP_SHIFTD
#define OP_NOBYTE
#include "tests-i386-shift.h"

/* XXX: should be more precise ? */
#undef CC_MASK
#define CC_MASK (CC_C)

#endif
#define OP bt
.global test_bt
test_bt:
#define OP_NOBYTE
#include "tests-i386-shift.h"

#define OP bts
.global test_bts
test_bts:
#define OP_NOBYTE
#include "tests-i386-shift.h"

#define OP btr
.global test_btr
test_btr:
#define OP_NOBYTE
#include "tests-i386-shift.h"

#define OP btc
.global test_btc
test_btc:
#define OP_NOBYTE
#include "tests-i386-shift.h"
	
#if 0	
#define TEST_BSX(op, size, op0)			\
{\
    long res, val, resz;\
    val = op0;\
    asm(".code16\n\nxor %1, %1\n"\
        "mov $0x12345678, %0\n"\
        #op " %" size "2, %" size "0 ; setz %b1\n\t" \
	"hlt\n\t.asciz \"" stringify(op) "\" \n\t"	     \
        : "=&r" (res), "=&q" (resz)\
	: "r" (val));\
}

	// TODO: figure out what qemu tests should output. It's confusing.

    TEST_BSX(bsrw, "w", 0);
    TEST_BSX(bsrw, "w", 0x12340128);
    TEST_BSX(bsfw, "w", 0);
    TEST_BSX(bsfw, "w", 0x12340128);
    TEST_BSX(bsrl, "k", 0);
    TEST_BSX(bsrl, "k", 0x00340128);
    TEST_BSX(bsfl, "k", 0);
    TEST_BSX(bsfl, "k", 0x00340128);
#endif

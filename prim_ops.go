/****************************************************************************
*
*                       Realmode X86 Emulator Library
*
*               Copyright (C) 1991-2004 SciTech Software, Inc.
*                    Copyright (C) David Mosberger-Tang
*                      Copyright (C) 1999 Egbert Eich
*
*  ========================================================================
*
*  Permission to use, copy, modify, distribute, and sell this software and
*  its documentation for any purpose is hereby granted without fee,
*  provided that the above copyright notice appear in all copies and that
*  both that copyright notice and this permission notice appear in
*  supporting documentation, and that the name of the authors not be used
*  in advertising or publicity pertaining to distribution of the software
*  without specific, written prior permission.  The authors makes no
*  representations about the suitability of this software for any purpose.
*  It is provided "as is" without express or implied warranty.
*
*  THE AUTHORS DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE,
*  INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO
*  EVENT SHALL THE AUTHORS BE LIABLE FOR ANY SPECIAL, INDIRECT OR
*  CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF
*  USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
*  OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
*  PERFORMANCE OF THIS SOFTWARE.
*
*  ========================================================================
*
* Language:     ANSI C
* Environment:  Any
* Developer:    Kendall Bennett
*
* Description:  This file contains the code to implement the primitive
*               machine operations used by the emulation code in ops.c
*
* Carry Chain Calculation
*
* This represents a somewhat expensive calculation which is
* apparently required to emulate the setting of the OF and AF flag.
* The latter is not so important, but the former is.  The overflow
* flag is the XOR of the top two bits of the carry chain for an
* addition (similar for subtraction).  Since we do not want to
* simulate the addition in a bitwise manner, we try to calculate the
* carry chain given the two operands and the result.
*
* So, given the following table, which represents the addition of two
* bits, we can derive a formula for the carry chain.
*
* a   b   cin   r     cout
* 0   0   0     0     0
* 0   0   1     1     0
* 0   1   0     1     0
* 0   1   1     0     1
* 1   0   0     1     0
* 1   0   1     0     1
* 1   1   0     0     1
* 1   1   1     1     1
*
* Construction of table for cout:
*
* ab
* r  \  00   01   11  10
* |------------------
* 0  |   0    1    1   1
* 1  |   0    0    1   0
*
* By inspection, one gets:  cc = ab +  r'(a + b)
*
* That represents alot of operations, but NO CHOICE....
*
* Borrow Chain Calculation.
*
* The following table represents the subtraction of two bits, from
* which we can derive a formula for the borrow chain.
*
* a   b   bin   r     bout
* 0   0   0     0     0
* 0   0   1     1     1
* 0   1   0     1     1
* 0   1   1     0     1
* 1   0   0     1     0
* 1   0   1     0     0
* 1   1   0     0     0
* 1   1   1     1     1
*
* Construction of table for cout:
*
* ab
* r  \  00   01   11  10
* |------------------
* 0  |   0    1    0   0
* 1  |   1    1    1   0
*
* By inspection, one gets:  bc = a'b +  r(a' + b)
*
****************************************************************************/

/*------------------------- Global Variables ------------------------------*/
package main

var x86emu_parity_tab = [8]uint32{
	0x96696996,
	0x69969669,
	0x69969669,
	0x96696996,
	0x69969669,
	0x96696996,
	0x96696996,
	0x69969669,
}

func PARITY(x uint32) bool {
	return ((x86emu_parity_tab[x/32]>>(x%32))&1 == 0)
}

func XOR2(x uint32) uint32 {
	return (((x) ^ ((x) >> 1)) & 0x1)
}

/*----------------------------- Implementation ----------------------------*/

/*--------- Side effects helper functions -------*/

/****************************************************************************
REMARKS:
implements side effects for byte operations that don't overflow
****************************************************************************/

func set_parity_flag(res uint32) {
	CONDITIONAL_SET_FLAG(PARITY(res&0xFF), F_PF)
}

func set_szp_flags_8(res uint8) {
	CONDITIONAL_SET_FLAG(res&0x80, F_SF)
	CONDITIONAL_SET_FLAG(res == 0, F_ZF)
	set_parity_flag(uint32(res))
}

func set_szp_flags_16(res uint16) {
	CONDITIONAL_SET_FLAG(res&0x8000, F_SF)
	CONDITIONAL_SET_FLAG(res == 0, F_ZF)
	set_parity_flag(uint32(res))
}

func set_szp_flags_32(res uint32) {
	CONDITIONAL_SET_FLAG(res&0x80000000, F_SF)
	CONDITIONAL_SET_FLAG(res == 0, F_ZF)
	set_parity_flag(res)
}

func no_carry_byte_side_eff(res uint8) {
	CLEAR_FLAG(F_OF)
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_AF)
	set_szp_flags_8(res)
}

func no_carry_word_side_eff(res uint16) {
	CLEAR_FLAG(F_OF)
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_AF)
	set_szp_flags_16(res)
}

func no_carry_long_side_eff(res uint32) {
	CLEAR_FLAG(F_OF)
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_AF)
	set_szp_flags_32(res)
}

func calc_carry_chain(bits uint32, d uint32, s uint32, res uint32, set_carry int) {

	cc := (s & d) | ((^res) & (s | d))
	CONDITIONAL_SET_FLAG(XOR2(cc>>(bits-2)), F_OF)
	CONDITIONAL_SET_FLAG(cc&0x8, F_AF)
	if set_carry != 0 {
		CONDITIONAL_SET_FLAG(res&(1<<bits), F_CF)
	}
}

func calc_borrow_chain(bits uint32, d uint32, s uint32, res uint32, set_carry int) {

	bc := (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(XOR2(bc>>(bits-2)), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	if set_carry != 0 {
		CONDITIONAL_SET_FLAG(bc&(1<<(bits-1)), F_CF)
	}
}

/****************************************************************************
REMARKS:
Implements the AAA instruction and side effects.
****************************************************************************/
func aaa_word(d uint16) uint16 {
	if (d&0xf) > 0x9 || ACCESS_FLAG(F_AF) {
		d += 0x6
		d += 0x100
		SET_FLAG(F_AF)
		SET_FLAG(F_CF)
	} else {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_AF)
	}
	res := uint16((d & 0xFF0F))
	set_szp_flags_16(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the AAA instruction and side effects.
****************************************************************************/
func aas_word(d uint16) uint16 {
	if (d&0xf) > 0x9 || ACCESS_FLAG(F_AF) {
		d -= 0x6
		d -= 0x100
		SET_FLAG(F_AF)
		SET_FLAG(F_CF)
	} else {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_AF)
	}
	res := uint16((d & 0xFF0F))
	set_szp_flags_16(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the AAD instruction and side effects.
****************************************************************************/
func aad_word(d uint16) uint16 {
	hb := uint16((d >> 8) & 0xff)
	lb := uint16(uint8(d))
	l := uint16(((lb + 10*hb) & 0xFF))

	no_carry_byte_side_eff(uint8(l))
	return l
}

/****************************************************************************
REMARKS:
Implements the AAM instruction and side effects.
****************************************************************************/
func aam_word(d uint8) uint16 {
	var l, h uint16

	h = uint16((d / 10))
	l = uint16((d % 10))
	l |= uint16((h << 8))

	no_carry_byte_side_eff(uint8(l))
	return l
}

/****************************************************************************
REMARKS:
Implements the ADC instruction and side effects.
****************************************************************************/
func adc_byte(d uint8, s uint8) uint8 {
	var res uint32 /* all operands in native machine order */

	res = d + s
	if ACCESS_FLAG(F_CF) {
		res++
	}

	set_szp_flags_8(res)
	calc_carry_chain(8, uint32(s), uint32(d), res, 1)

	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the ADC instruction and side effects.
****************************************************************************/
func adc_word(d uint16, s uint16) uint16 {
	var res uint32 /* all operands in native machine order */

	res = d + s
	if ACCESS_FLAG(F_CF) {
		res++
	}

	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, uint32(s), uint32(d), res, 1)

	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the ADC instruction and side effects.
****************************************************************************/
func adc_long(d uint32, s uint32) uint32 {
	var lo uint32 /* all operands in native machine order */
	var hi uint32
	var res uint32

	lo = (d & 0xFFFF) + (s & 0xFFFF)
	res = d + s

	if ACCESS_FLAG(F_CF) {
		lo++
		res++
	}

	hi = (lo >> 16) + (d >> 16) + (s >> 16)

	set_szp_flags_32(res)
	calc_carry_chain(32, s, d, res, 0)

	CONDITIONAL_SET_FLAG(hi&0x10000, F_CF)

	return res
}

/****************************************************************************
REMARKS:
Implements the ADD instruction and side effects.
****************************************************************************/
func add_byte(d uint8, s uint8) uint8 {
	var res uint32 /* all operands in native machine order */

	res = d + s
	set_szp_flags_8(uint8(res))
	calc_carry_chain(8, s, d, res, 1)

	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the ADD instruction and side effects.
****************************************************************************/
func add_word(d uint16, s uint16) uint16 {
	var res uint32 /* all operands in native machine order */

	res = d + s
	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, s, d, res, 1)

	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the ADD instruction and side effects.
****************************************************************************/
func add_long(d uint32, s uint32) uint32 {
	var res uint32

	res = d + s
	set_szp_flags_32(res)
	calc_carry_chain(32, s, d, res, 0)

	CONDITIONAL_SET_FLAG(res < d || res < s, F_CF)

	return res
}

/****************************************************************************
REMARKS:
Implements the AND instruction and side effects.
****************************************************************************/
func and_byte(d uint8, s uint8) uint8 {
	var res uint8 /* all operands in native machine order */

	res = d & s

	no_carry_byte_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the AND instruction and side effects.
****************************************************************************/
func and_word(d uint16, s uint16) uint16 {
	var res uint16 /* all operands in native machine order */

	res = d & s

	no_carry_word_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the AND instruction and side effects.
****************************************************************************/
func and_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d & s
	no_carry_long_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the CMP instruction and side effects.
****************************************************************************/
func cmp_byte(d uint8, s uint8) uint8 {
	var res uint32 /* all operands in native machine order */

	res = d - s
	set_szp_flags_8(uint8(res))
	calc_borrow_chain(8, d, s, res, 1)

	return d
}

/****************************************************************************
REMARKS:
Implements the CMP instruction and side effects.
****************************************************************************/
func cmp_word(d uint16, s uint16) uint16 {
	var res uint32 /* all operands in native machine order */

	res = d - s
	set_szp_flags_16(uint16(res))
	calc_borrow_chain(16, d, s, res, 1)

	return d
}

/****************************************************************************
REMARKS:
Implements the CMP instruction and side effects.
****************************************************************************/
func cmp_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d - s
	set_szp_flags_32(res)
	calc_borrow_chain(32, d, s, res, 1)

	return d
}

/****************************************************************************
REMARKS:
Implements the DAA instruction and side effects.
****************************************************************************/
func daa_byte(d uint8) uint8 {
	var res uint32 = d
	if (d&0xf) > 9 || ACCESS_FLAG(F_AF) {
		res += 6
		SET_FLAG(F_AF)
	}
	if res > 0x9F || ACCESS_FLAG(F_CF) {
		res += 0x60
		SET_FLAG(F_CF)
	}
	set_szp_flags_8(uint8(res))
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the DAS instruction and side effects.
****************************************************************************/
func das_byte(d uint8) uint8 {
	if (d&0xf) > 9 || ACCESS_FLAG(F_AF) {
		d -= 6
		SET_FLAG(F_AF)
	}
	if d > 0x9F || ACCESS_FLAG(F_CF) {
		d -= 0x60
		SET_FLAG(F_CF)
	}
	set_szp_flags_8(d)
	return d
}

/****************************************************************************
REMARKS:
Implements the DEC instruction and side effects.
****************************************************************************/
func dec_byte(d uint8) uint8 {
	var res uint32 /* all operands in native machine order */

	res = d - 1
	set_szp_flags_8(uint8(res))
	calc_borrow_chain(8, d, 1, res, 0)

	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the DEC instruction and side effects.
****************************************************************************/
func dec_word(d uint16) uint16 {
	var res uint32 /* all operands in native machine order */

	res = d - 1
	set_szp_flags_16(uint16(res))
	calc_borrow_chain(16, d, 1, res, 0)

	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the DEC instruction and side effects.
****************************************************************************/
func dec_long(d uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d - 1

	set_szp_flags_32(res)
	calc_borrow_chain(32, d, 1, res, 0)

	return res
}

/****************************************************************************
REMARKS:
Implements the INC instruction and side effects.
****************************************************************************/
func inc_byte(d uint8) uint8 {
	var res uint32 /* all operands in native machine order */

	res = d + 1
	set_szp_flags_8(uint8(res))
	calc_carry_chain(8, d, 1, res, 0)

	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the INC instruction and side effects.
****************************************************************************/
func inc_word(d uint16) uint16 {
	var res uint32 /* all operands in native machine order */

	res = d + 1
	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, d, 1, res, 0)

	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the INC instruction and side effects.
****************************************************************************/
func inc_long(d uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d + 1
	set_szp_flags_32(res)
	calc_carry_chain(32, d, 1, res, 0)

	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func or_byte(d uint8, s uint8) uint8 {
	var res uint8 /* all operands in native machine order */

	res = d | s
	no_carry_byte_side_eff(res)

	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func or_word(d uint16, s uint16) uint16 {
	var res uint16 /* all operands in native machine order */

	res = d | s
	no_carry_word_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func or_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d | s
	no_carry_long_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func neg_byte(s uint8) uint8 {
	var res uint8

	CONDITIONAL_SET_FLAG(s != 0, F_CF)
	res = (u8) - s
	set_szp_flags_8(res)
	calc_borrow_chain(8, 0, s, res, 0)

	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func neg_word(s uint16) uint16 {
	var res uint16

	CONDITIONAL_SET_FLAG(s != 0, F_CF)
	res = uint16(-s)
	set_szp_flags_16(uint16(res))
	calc_borrow_chain(16, 0, s, res, 0)

	return res
}

/****************************************************************************
REMARKS:
Implements the OR instruction and side effects.
****************************************************************************/
func neg_long(s uint32) uint32 {
	var res uint32

	CONDITIONAL_SET_FLAG(s != 0, F_CF)
	res = (u32) - s
	set_szp_flags_32(res)
	calc_borrow_chain(32, 0, s, res, 0)

	return res
}

/****************************************************************************
REMARKS:
Implements the NOT instruction and side effects.
****************************************************************************/
func not_byte(s uint8) uint8 {
	return ^s
}

/****************************************************************************
REMARKS:
Implements the NOT instruction and side effects.
****************************************************************************/
func not_word(s uint16) uint16 {
	return ^s
}

/****************************************************************************
REMARKS:
Implements the NOT instruction and side effects.
****************************************************************************/
func not_long(s uint32) uint32 {
	return ^s
}

/****************************************************************************
REMARKS:
Implements the RCL instruction and side effects.
****************************************************************************/
func rcl_byte(d uint8, s uint8) uint8 {
	var res, int, cnt, mask, cf uint8

	/* s is the rotate distance.  It varies from 0 - 8. */
	/* have

	   CF  B_7 B_6 B_5 B_4 B_3 B_2 B_1 B_0

	   want to rotate through the carry by "s" bits.  We could
	   loop, but that's inefficient.  So the width is 9,
	   and we split into three parts:

	   The new carry flag   (was B_n)
	   the stuff in B_n-1 .. B_0
	   the stuff in B_7 .. B_n+1

	   The new rotate is done mod 9, and given this,
	   for a rotation of n bits (mod 9) the new carry flag is
	   then located n bits from the MSB.  The low part is
	   then shifted up cnt bits, and the high part is or'd
	   in.  Using CAPS for new values, and lowercase for the
	   original values, this can be expressed as:

	   IF n > 0
	   1) CF <-  b_(8-n)
	   2) B_(7) .. B_(n)  <-  b_(8-(n+1)) .. b_0
	   3) B_(n-1) <- cf
	   4) B_(n-2) .. B_0 <-  b_7 .. b_(8-(n-1))
	*/
	res = d
	cnt = s % 9
	if cnt != 0 {
		/* extract the new CARRY FLAG. */
		/* CF <-  b_(8-n)             */
		cf = (d >> (8 - cnt)) & 0x1

		/* get the low stuff which rotated
		   into the range B_7 .. B_cnt */
		/* B_(7) .. B_(n)  <-  b_(8-(n+1)) .. b_0  */
		/* note that the right hand side done by the mask */
		res = (d << cnt) & 0xff

		/* now the high stuff which rotated around
		   into the positions B_cnt-2 .. B_0 */
		/* B_(n-2) .. B_0 <-  b_7 .. b_(8-(n-1)) */
		/* shift it downward, 7-(n-2) = 9-n positions.
		   and mask off the result before or'ing in.
		*/
		mask = (1 << (cnt - 1)) - 1
		res |= (d >> (9 - cnt)) & mask

		/* if the carry flag was set, or it in.  */
		if ACCESS_FLAG(F_CF) { /* carry flag is set */
			/*  B_(n-1) <- cf */
			res |= 1 << (cnt - 1)
		}
		/* set the new carry flag, based on the variable "cf" */
		CONDITIONAL_SET_FLAG(cf, F_CF)
		/* OVERFLOW is set *IFF* cnt==1, then it is the
		   xor of CF and the most significant bit.  Blecck. */
		/* parenthesized this expression since it appears to
		   be causing OF to be missed */
		CONDITIONAL_SET_FLAG(cnt == 1 && XOR2(cf+((res>>6)&0x2)),
			F_OF)

	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the RCL instruction and side effects.
****************************************************************************/
func rcl_word(d uint16, s uint8) uint16 {
	var res, int, cnt, mask, cf uint16

	res = d
	cnt = s % 17
	if cnt != 0 {
		cf = (d >> (16 - cnt)) & 0x1
		res = (d << cnt) & 0xffff
		mask = (1 << (cnt - 1)) - 1
		res |= (d >> (17 - cnt)) & mask
		if ACCESS_FLAG(F_CF) {
			res |= 1 << (cnt - 1)
		}
		CONDITIONAL_SET_FLAG(cf, F_CF)
		CONDITIONAL_SET_FLAG(cnt == 1 && XOR2(cf+((res>>14)&0x2)),
			F_OF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the RCL instruction and side effects.
****************************************************************************/
func rcl_long(d uint32, s uint8) uint32 {
	var res, cnt, mask, cf uint32

	res = d
	cnt = s % 33
	if cnt != 0 {
		cf = (d >> (32 - cnt)) & 0x1
		res = (d << cnt) & 0xffffffff
		mask = (1 << (cnt - 1)) - 1
		res |= (d >> (33 - cnt)) & mask
		if ACCESS_FLAG(F_CF) { /* carry flag is set */
			res |= 1 << (cnt - 1)
		}
		CONDITIONAL_SET_FLAG(cf, F_CF)
		CONDITIONAL_SET_FLAG(cnt == 1 && XOR2(cf+((res>>30)&0x2)),
			F_OF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the RCR instruction and side effects.
****************************************************************************/
func rcr_byte(d uint8, s uint8) uint8 {
	var res, cnt, mask, cf, ocf uint32

	/* rotate right through carry */
	/*
	   s is the rotate distance.  It varies from 0 - 8.
	   d is the byte object rotated.

	   have

	   CF  B_7 B_6 B_5 B_4 B_3 B_2 B_1 B_0

	   The new rotate is done mod 9, and given this,
	   for a rotation of n bits (mod 9) the new carry flag is
	   then located n bits from the LSB.  The low part is
	   then shifted up cnt bits, and the high part is or'd
	   in.  Using CAPS for new values, and lowercase for the
	   original values, this can be expressed as:

	   IF n > 0
	   1) CF <-  b_(n-1)
	   2) B_(8-(n+1)) .. B_(0)  <-  b_(7) .. b_(n)
	   3) B_(8-n) <- cf
	   4) B_(7) .. B_(8-(n-1)) <-  b_(n-2) .. b_(0)
	*/
	res = d
	cnt = s % 9
	if cnt != 0 {
		/* extract the new CARRY FLAG. */
		/* CF <-  b_(n-1)              */
		if cnt == 1 {
			cf = d & 0x1
			/* note hackery here.  Access_flag(..) evaluates to either
			   0 if flag not set
			   non-zero if flag is set.
			   doing access_flag(..) != 0 casts that into either
			   0..1 in any representation of the flags register
			   (i.e. packed bit array or unpacked.)
			*/
			ocf = ACCESS_FLAG(F_CF) != 0
		} else {
			cf = (d >> (cnt - 1)) & 0x1
		}

		/* B_(8-(n+1)) .. B_(0)  <-  b_(7) .. b_n  */
		/* note that the right hand side done by the mask
		   This is effectively done by shifting the
		   object to the right.  The result must be masked,
		   in case the object came in and was treated
		   as a negative number.  Needed??? */

		mask = (1 << (8 - cnt)) - 1
		res = (d >> cnt) & mask

		/* now the high stuff which rotated around
		   into the positions B_cnt-2 .. B_0 */
		/* B_(7) .. B_(8-(n-1)) <-  b_(n-2) .. b_(0) */
		/* shift it downward, 7-(n-2) = 9-n positions.
		   and mask off the result before or'ing in.
		*/
		res |= (d << (9 - cnt))

		/* if the carry flag was set, or it in.  */
		if ACCESS_FLAG(F_CF) { /* carry flag is set */
			/*  B_(8-n) <- cf */
			res |= 1 << (8 - cnt)
		}
		/* set the new carry flag, based on the variable "cf" */
		CONDITIONAL_SET_FLAG(cf, F_CF)
		/* OVERFLOW is set *IFF* cnt==1, then it is the
		   xor of CF and the most significant bit.  Blecck. */
		/* parenthesized... */
		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(ocf+((d>>6)&0x2)),
				F_OF)
		}
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the RCR instruction and side effects.
****************************************************************************/
func rcr_word(d uint16, s uint8) uint16 {
	var res, cnt, mask, cf, ocf uint16

	/* rotate right through carry */
	res = d
	cnt = s % 17
	if cnt != 0 {
		if cnt == 1 {
			cf = d & 0x1
			ocf = ACCESS_FLAG(F_CF) != 0
		} else {
			cf = (d >> (cnt - 1)) & 0x1
		}
		mask = (1 << (16 - cnt)) - 1
		res = (d >> cnt) & mask
		res |= (d << (17 - cnt))
		if ACCESS_FLAG(F_CF) {
			res |= 1 << (16 - cnt)
		}
		CONDITIONAL_SET_FLAG(cf, F_CF)
		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(ocf+((d>>14)&0x2)),
				F_OF)
		}
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the RCR instruction and side effects.
****************************************************************************/
func rcr_long(d uint32, s uint8) uint32 {
	var res, cnt, mask, cf, ocf uint32

	/* rotate right through carry */
	res = d
	cnt = s % 33
	if cnt != 0 {
		if cnt == 1 {
			cf = d & 0x1
			ocf = ACCESS_FLAG(F_CF) != 0
		} else {
			cf = (d >> (cnt - 1)) & 0x1
		}
		mask = (1 << (32 - cnt)) - 1
		res = (d >> cnt) & mask
		if cnt != 1 {
			res |= (d << (33 - cnt))
		}
		if ACCESS_FLAG(F_CF) { /* carry flag is set */
			res |= 1 << (32 - cnt)
		}
		CONDITIONAL_SET_FLAG(cf, F_CF)
		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(ocf+((d>>30)&0x2)),
				F_OF)
		}
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the ROL instruction and side effects.
****************************************************************************/
func rol_byte(d uint8, s uint8) uint8 {
	var res, cnt, mask int

	/* rotate left */
	/*
	   s is the rotate distance.  It varies from 0 - 8.
	   d is the byte object rotated.

	   have

	   CF  B_7 ... B_0

	   The new rotate is done mod 8.
	   Much simpler than the "rcl" or "rcr" operations.

	   IF n > 0
	   1) B_(7) .. B_(n)  <-  b_(8-(n+1)) .. b_(0)
	   2) B_(n-1) .. B_(0) <-  b_(7) .. b_(8-n)
	*/
	res = d
	cnt = s % 8
	if cnt != 0 {
		/* B_(7) .. B_(n)  <-  b_(8-(n+1)) .. b_(0) */
		res = (d << cnt)

		/* B_(n-1) .. B_(0) <-  b_(7) .. b_(8-n) */
		mask = (1 << cnt) - 1
		res |= (d >> (8 - cnt)) & mask

		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
		/* OVERFLOW is set *IFF* s==1, then it is the
		   xor of CF and the most significant bit.  Blecck. */
		CONDITIONAL_SET_FLAG(s == 1 &&
			XOR2((res&0x1)+((res>>6)&0x2)),
			F_OF)
	}
	if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the ROL instruction and side effects.
****************************************************************************/
func rol_word(d uint16, s uint8) uint16 {
	var res, cnt, mask int

	res = d
	cnt = s % 16
	if cnt != 0 {
		res = (d << cnt)
		mask = (1 << cnt) - 1
		res |= (d >> (16 - cnt)) & mask
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
		CONDITIONAL_SET_FLAG(s == 1 &&
			XOR2((res&0x1)+((res>>14)&0x2)),
			F_OF)
	}
	if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the ROL instruction and side effects.
****************************************************************************/
func rol_long(d uint32, s uint8) uint32 {
	var res, cnt, mask uint32

	res = d
	if cnt = s % 32; cnt != 0 {
		res = (d << cnt)
		mask = (1 << cnt) - 1
		res |= (d >> (32 - cnt)) & mask
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
		CONDITIONAL_SET_FLAG(s == 1 &&
			XOR2((res&0x1)+((res>>30)&0x2)),
			F_OF)
	}
	if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x1, F_CF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the ROR instruction and side effects.
****************************************************************************/
func ror_byte(d uint8, s uint8) uint8 {
	var res, cnt, mask int

	/* rotate right */
	/*
	   s is the rotate distance.  It varies from 0 - 8.
	   d is the byte object rotated.

	   have

	   B_7 ... B_0

	   The rotate is done mod 8.

	   IF n > 0
	   1) B_(8-(n+1)) .. B_(0)  <-  b_(7) .. b_(n)
	   2) B_(7) .. B_(8-n) <-  b_(n-1) .. b_(0)
	*/
	res = d
	if cnt = s % 8; cnt != 0 { /* not a typo, do nada if cnt==0 */
		/* B_(7) .. B_(8-n) <-  b_(n-1) .. b_(0) */
		res = (d << (8 - cnt))

		/* B_(8-(n+1)) .. B_(0)  <-  b_(7) .. b_(n) */
		mask = (1 << (8 - cnt)) - 1
		res |= (d >> (cnt)) & mask

		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x80, F_CF)
		/* OVERFLOW is set *IFF* s==1, then it is the
		   xor of the two most significant bits.  Blecck. */
		CONDITIONAL_SET_FLAG(s == 1 && XOR2(res>>6), F_OF)
	} else if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x80, F_CF)
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the ROR instruction and side effects.
****************************************************************************/
func ror_word(d uint16, s uint8) uint16 {
	var res, cnt, mask int

	res = d
	if cnt = s % 16; cnt != 0 {
		res = (d << (16 - cnt))
		mask = (1 << (16 - cnt)) - 1
		res |= (d >> (cnt)) & mask
		CONDITIONAL_SET_FLAG(res&0x8000, F_CF)
		CONDITIONAL_SET_FLAG(s == 1 && XOR2(res>>14), F_OF)
	} else if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x8000, F_CF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the ROR instruction and side effects.
****************************************************************************/
func ror_long(d uint32, s uint8) uint32 {
	var res, cnt, mask uint32

	res = d
	if cnt = s % 32; cnt != 0 {
		res = (d << (32 - cnt))
		mask = (1 << (32 - cnt)) - 1
		res |= (d >> (cnt)) & mask
		CONDITIONAL_SET_FLAG(res&0x80000000, F_CF)
		CONDITIONAL_SET_FLAG(s == 1 && XOR2(res>>30), F_OF)
	} else if s != 0 {
		/* set the new carry flag, Note that it is the low order
		   bit of the result!!!                               */
		CONDITIONAL_SET_FLAG(res&0x80000000, F_CF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SHL instruction and side effects.
****************************************************************************/
func shl_byte(d uint8, s uint8) uint8 {
	var cnt, res, cf int

	if s < 8 {
		cnt = s % 8

		/* last bit shifted out goes into carry flag */
		if cnt > 0 {
			res = d << cnt
			cf = d & (1 << (8 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_8(uint8(res))
		} else {
			res = uint8(d)
		}

		if cnt == 1 {
			/* Needs simplification. */
			CONDITIONAL_SET_FLAG(
				(((res & 0x80) == 0x80) ^
					(ACCESS_FLAG(F_CF) != 0)),
				/* was (M.x86.R_FLG&F_CF)==F_CF)), */
				F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d<<(s-1))&0x80, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the SHL instruction and side effects.
****************************************************************************/
func shl_word(d uint16, s uint8) uint16 {
	var cnt, res, cf int

	if s < 16 {
		cnt = s % 16
		if cnt > 0 {
			res = d << cnt
			cf = d & (1 << (16 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_16(uint16(res))
		} else {
			res = uint16(d)
		}

		if cnt == 1 {
			CONDITIONAL_SET_FLAG(
				(((res & 0x8000) == 0x8000) ^
					(ACCESS_FLAG(F_CF) != 0)),
				F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d<<(s-1))&0x8000, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SHL instruction and side effects.
****************************************************************************/
func shl_long(d uint32, s uint8) uint32 {
	var cnt, res, cf int

	if s < 32 {
		cnt = s % 32
		if cnt > 0 {
			res = d << cnt
			cf = d & (1 << (32 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_32(uint32(res))
		} else {
			res = d
		}
		if cnt == 1 {
			CONDITIONAL_SET_FLAG((((res & 0x80000000) == 0x80000000) ^
				(ACCESS_FLAG(F_CF) != 0)), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d<<(s-1))&0x80000000, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SHR instruction and side effects.
****************************************************************************/
func shr_byte(d uint8, s uint8) uint8 {
	var cnt, res, cf int

	if s < 8 {
		cnt = s % 8
		if cnt > 0 {
			cf = d & (1 << (cnt - 1))
			res = d >> cnt
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_8(uint8(res))
		} else {
			res = uint8(d)
		}

		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(res>>6), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d>>(s-1))&0x1, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the SHR instruction and side effects.
****************************************************************************/
func shr_word(d uint16, s uint8) uint16 {
	var cnt, res, cf int

	if s < 16 {
		cnt = s % 16
		if cnt > 0 {
			cf = d & (1 << (cnt - 1))
			res = d >> cnt
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_16(uint16(res))
		} else {
			res = d
		}

		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(res>>14), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
		SET_FLAG(F_ZF)
		CLEAR_FLAG(F_SF)
		CLEAR_FLAG(F_PF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SHR instruction and side effects.
****************************************************************************/
func shr_long(d uint32, s uint8) uint32 {
	var cnt, res, cf int

	if s < 32 {
		cnt = s % 32
		if cnt > 0 {
			cf = d & (1 << (cnt - 1))
			res = d >> cnt
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_32(uint32(res))
		} else {
			res = d
		}
		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(res>>30), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
		SET_FLAG(F_ZF)
		CLEAR_FLAG(F_SF)
		CLEAR_FLAG(F_PF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SAR instruction and side effects.
****************************************************************************/
func sar_byte(d uint8, s uint8) uint8 {
	var cnt, res, cf, mask, sf int

	res = d
	sf = d & 0x80
	cnt = s % 8
	if cnt > 0 && cnt < 8 {
		mask = (1 << (8 - cnt)) - 1
		cf = d & (1 << (cnt - 1))
		res = (d >> cnt) & mask
		CONDITIONAL_SET_FLAG(cf, F_CF)
		if sf {
			res |= ^mask
		}
		set_szp_flags_8(uint8(res))
	} else if cnt >= 8 {
		if sf {
			res = 0xff
			SET_FLAG(F_CF)
			CLEAR_FLAG(F_ZF)
			SET_FLAG(F_SF)
			SET_FLAG(F_PF)
		} else {
			res = 0
			CLEAR_FLAG(F_CF)
			SET_FLAG(F_ZF)
			CLEAR_FLAG(F_SF)
			CLEAR_FLAG(F_PF)
		}
	}
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the SAR instruction and side effects.
****************************************************************************/
func sar_word(d uint16, s uint8) uint16 {
	var cnt, res, cf, mask, sf int

	sf = d & 0x8000
	cnt = s % 16
	res = d
	if cnt > 0 && cnt < 16 {
		mask = (1 << (16 - cnt)) - 1
		cf = d & (1 << (cnt - 1))
		res = (d >> cnt) & mask
		CONDITIONAL_SET_FLAG(cf, F_CF)
		if sf {
			res |= ^mask
		}
		set_szp_flags_16(uint16(res))
	} else if cnt >= 16 {
		if sf {
			res = 0xffff
			SET_FLAG(F_CF)
			CLEAR_FLAG(F_ZF)
			SET_FLAG(F_SF)
			SET_FLAG(F_PF)
		} else {
			res = 0
			CLEAR_FLAG(F_CF)
			SET_FLAG(F_ZF)
			CLEAR_FLAG(F_SF)
			CLEAR_FLAG(F_PF)
		}
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SAR instruction and side effects.
****************************************************************************/
func sar_long(d uint32, s uint8) uint32 {
	var cnt, res, cf, mask, sf uint32

	sf = d & 0x80000000
	cnt = s % 32
	res = d
	if cnt > 0 && cnt < 32 {
		mask = (1 << (32 - cnt)) - 1
		cf = d & (1 << (cnt - 1))
		res = (d >> cnt) & mask
		CONDITIONAL_SET_FLAG(cf, F_CF)
		if sf {
			res |= ^mask
		}
		set_szp_flags_32(res)
	} else if cnt >= 32 {
		if sf {
			res = 0xffffffff
			SET_FLAG(F_CF)
			CLEAR_FLAG(F_ZF)
			SET_FLAG(F_SF)
			SET_FLAG(F_PF)
		} else {
			res = 0
			CLEAR_FLAG(F_CF)
			SET_FLAG(F_ZF)
			CLEAR_FLAG(F_SF)
			CLEAR_FLAG(F_PF)
		}
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SHLD instruction and side effects.
****************************************************************************/
func shld_word(d uint16, fill uint16, s uint8) uint16 {
	var cnt, res, cf int

	if s < 16 {
		cnt = s % 16
		if cnt > 0 {
			res = (d << cnt) | (fill >> (16 - cnt))
			cf = d & (1 << (16 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_16(uint16(res))
		} else {
			res = d
		}
		if cnt == 1 {
			CONDITIONAL_SET_FLAG((((res & 0x8000) == 0x8000) ^
				(ACCESS_FLAG(F_CF) != 0)), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d<<(s-1))&0x8000, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SHLD instruction and side effects.
****************************************************************************/
func shld_long(d uint32, fill uint32, s uint8) uint32 {
	var cnt, res, cf int

	if s < 32 {
		cnt = s % 32
		if cnt > 0 {
			res = (d << cnt) | (fill >> (32 - cnt))
			cf = d & (1 << (32 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_32(uint32(res))
		} else {
			res = d
		}
		if cnt == 1 {
			CONDITIONAL_SET_FLAG((((res & 0x80000000) == 0x80000000) ^
				(ACCESS_FLAG(F_CF) != 0)), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CONDITIONAL_SET_FLAG((d<<(s-1))&0x80000000, F_CF)
		CLEAR_FLAG(F_OF)
		CLEAR_FLAG(F_SF)
		SET_FLAG(F_PF)
		SET_FLAG(F_ZF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SHRD instruction and side effects.
****************************************************************************/
func shrd_word(d uint16, fill uint16, s uint8) uint16 {
	var cnt, res, cf int

	if s < 16 {
		cnt = s % 16
		if cnt > 0 {
			cf = d & (1 << (cnt - 1))
			res = (d >> cnt) | (fill << (16 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_16(uint16(res))
		} else {
			res = d
		}

		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(res>>14), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
		SET_FLAG(F_ZF)
		CLEAR_FLAG(F_SF)
		CLEAR_FLAG(F_PF)
	}
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SHRD instruction and side effects.
****************************************************************************/
func shrd_long(d uint32, fill uint32, s uint8) uint32 {
	var cnt, res, cf int

	if s < 32 {
		cnt = s % 32
		if cnt > 0 {
			cf = d & (1 << (cnt - 1))
			res = (d >> cnt) | (fill << (32 - cnt))
			CONDITIONAL_SET_FLAG(cf, F_CF)
			set_szp_flags_32(uint32(res))
		} else {
			res = d
		}
		if cnt == 1 {
			CONDITIONAL_SET_FLAG(XOR2(res>>30), F_OF)
		} else {
			CLEAR_FLAG(F_OF)
		}
	} else {
		res = 0
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
		SET_FLAG(F_ZF)
		CLEAR_FLAG(F_SF)
		CLEAR_FLAG(F_PF)
	}
	return res
}

/****************************************************************************
REMARKS:
Implements the SBB instruction and side effects.
****************************************************************************/
func sbb_byte(d uint8, s uint8) uint8 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	if ACCESS_FLAG(F_CF) {
		res = d - s - 1
	} else {
		res = d - s
	}
	set_szp_flags_8(uint8(res))

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x80, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>6), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the SBB instruction and side effects.
****************************************************************************/
func sbb_word(d uint16, s uint16) uint16 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	if ACCESS_FLAG(F_CF) {
		res = d - s - 1
	} else {
		res = d - s
	}
	set_szp_flags_16(uint16(res))

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x8000, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>14), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SBB instruction and side effects.
****************************************************************************/
func sbb_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	if ACCESS_FLAG(F_CF) {
		res = d - s - 1
	} else {
		res = d - s
	}

	set_szp_flags_32(res)

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x80000000, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>30), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return res
}

/****************************************************************************
REMARKS:
Implements the SUB instruction and side effects.
****************************************************************************/
func sub_byte(d uint8, s uint8) uint8 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	res = d - s
	set_szp_flags_8(uint8(res))

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x80, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>6), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return uint8(res)
}

/****************************************************************************
REMARKS:
Implements the SUB instruction and side effects.
****************************************************************************/
func sub_word(d uint16, s uint16) uint16 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	res = d - s
	set_szp_flags_16(uint16(res))

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x8000, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>14), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return uint16(res)
}

/****************************************************************************
REMARKS:
Implements the SUB instruction and side effects.
****************************************************************************/
func sub_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */
	var bc uint32

	res = d - s
	set_szp_flags_32(res)

	/* calculate the borrow chain.  See note at top */
	bc = (res & (^d | s)) | (^d & s)
	CONDITIONAL_SET_FLAG(bc&0x80000000, F_CF)
	CONDITIONAL_SET_FLAG(XOR2(bc>>30), F_OF)
	CONDITIONAL_SET_FLAG(bc&0x8, F_AF)
	return res
}

/****************************************************************************
REMARKS:
Implements the TEST instruction and side effects.
****************************************************************************/
func test_byte(d uint8, s uint8) {
	var res uint32 /* all operands in native machine order */

	res = d & s

	CLEAR_FLAG(F_OF)
	set_szp_flags_8(uint8(res))
	/* AF == don't care */
	CLEAR_FLAG(F_CF)
}

/****************************************************************************
REMARKS:
Implements the TEST instruction and side effects.
****************************************************************************/
func test_word(d uint16, s uint16) {
	var res uint32 /* all operands in native machine order */

	res = d & s

	CLEAR_FLAG(F_OF)
	set_szp_flags_16(uint16(res))
	/* AF == don't care */
	CLEAR_FLAG(F_CF)
}

/****************************************************************************
REMARKS:
Implements the TEST instruction and side effects.
****************************************************************************/
func test_long(d uint32, s uint32) {
	var res uint32 /* all operands in native machine order */

	res = d & s

	CLEAR_FLAG(F_OF)
	set_szp_flags_32(res)
	/* AF == don't care */
	CLEAR_FLAG(F_CF)
}

/****************************************************************************
REMARKS:
Implements the XOR instruction and side effects.
****************************************************************************/
func xor_byte(d uint8, s uint8) uint8 {
	var res uint8 /* all operands in native machine order */

	res = d ^ s
	no_carry_byte_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the XOR instruction and side effects.
****************************************************************************/
func xor_word(d uint16, s uint16) uint16 {
	var res uint16 /* all operands in native machine order */

	res = d ^ s
	no_carry_word_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the XOR instruction and side effects.
****************************************************************************/
func xor_long(d uint32, s uint32) uint32 {
	var res uint32 /* all operands in native machine order */

	res = d ^ s
	no_carry_long_side_eff(res)
	return res
}

/****************************************************************************
REMARKS:
Implements the IMUL instruction and side effects.
****************************************************************************/
func imul_byte(s uint8) {
	res := uint32(int16(M.x86.A.Get8()) * int16(s))

	M.x86.A.Set(res)
	if ((M.x86.A.Get8l()&0x80) == 0 && M.x86.R_AH == 0x00) ||
		((M.x86.A.Get8l()&0x80) != 0 && M.x86.R_AH == 0xFF) {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the IMUL instruction and side effects.
****************************************************************************/
func imul_word(s uint16) {
	res := uint32(M.x86.A.Get16() * int16(s))

	M.x86.A.Set16(uint16(res))
	M.x86.D.Set16(uint16(res >> 16))
	if ((M.x86.A.Get16()&0x8000) == 0 && M.x86.D.Get16() == 0x0000) ||
		((M.x86.A.Get16()&0x8000) != 0 && M.x86.D.Get16() == 0xFFFF) {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the IMUL instruction and side effects.
****************************************************************************/
func imul_long_direct(res_lo, res_hi *uint32, d uint32, s uint32) {
	res := int64(d) * int64(s)

	*res_lo = uint32(res)
	*res_hi = (u32)(res >> 32)
}

/****************************************************************************
REMARKS:
Implements the IMUL instruction and side effects.
****************************************************************************/
func imul_long(s uint32) {
	imul_long_direct(&M.x86.gen.A.Get32(), &M.x86.gen.D.Get32(), M.x86.gen.A.Get32(), s)
	if ((M.x86.gen.A.Get32()&0x80000000) == 0 && M.x86.gen.D.Get32() == 0x00000000) ||
		((M.x86.gen.A.Get32()&0x80000000) != 0 && M.x86.gen.D.Get32() == 0xFFFFFFFF) {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the MUL instruction and side effects.
****************************************************************************/
func mul_byte(s uint8) {
	res := uint16(M.x86.A.Get8l() * s)

	M.x86.Set16(res)
	if M.x86.A.Get8h() == 0 {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the MUL instruction and side effects.
****************************************************************************/
func mul_word(s uint16) {
	var res uint32 = M.x86.A.Get16() * s

	M.x86.gen.A.Set16(uint16(res))
	M.x86.R_DX = uint16((res >> 16))
	if M.x86.D.Get16() == 0 {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the MUL instruction and side effects.
****************************************************************************/
func mul_long(s uint32) {
	res := uint64(M.x86.A.Get32()) * uint64(s)

	M.x86.A.Set32(uint32(res))
	M.x86.D.Set32((u32)(res >> 32))
	if M.x86.D.Get32() == 0 {
		CLEAR_FLAG(F_CF)
		CLEAR_FLAG(F_OF)
	} else {
		SET_FLAG(F_CF)
		SET_FLAG(F_OF)
	}
}

/****************************************************************************
REMARKS:
Implements the IDIV instruction and side effects.
****************************************************************************/
func idiv_byte(s uint8) {
	var dvd, div, mod int16

	dvd = M.x86.A.Get16()
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / int16(s)
	mod = dvd % int16(s)
	if abs(div) > 0x7f {
		x86emu_intr_raise(0)
		return
	}
	M.x86.A.Set8l(int8(div))
	M.x86.A.Set8h(int8(mod))
}

/****************************************************************************
REMARKS:
Implements the IDIV instruction and side effects.
****************************************************************************/
func idiv_word(s uint16) {
	var dvd, div, mod int32

	dvd = (int32(M.x86.D.Get16()) << 16) | M.x86.A.Get16()
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / int16(s)
	mod = dvd % int16(s)
	if abs(div) > 0x7fff {
		x86emu_intr_raise(0)
		return
	}
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_SF)
	CONDITIONAL_SET_FLAG(div == 0, F_ZF)
	set_parity_flag(mod)

	M.x86.A.Set(uint16(div))
	M.x86.D.Set(uint16(mod))
}

/****************************************************************************
REMARKS:
Implements the IDIV instruction and side effects.
****************************************************************************/
func idiv_long(s uint32) {
	var dvd, div, mod int64

	dvd = (int64(M.x86.D.Get32()) << 32) | M.x86.A.Get32()
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / int32(s)
	mod = dvd % int32(s)
	if abs(div) > 0x7fffffff {
		x86emu_intr_raise(0)
		return
	}
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_AF)
	CLEAR_FLAG(F_SF)
	SET_FLAG(F_ZF)
	set_parity_flag(mod)

	M.x86.A.Set(uint32(div))
	M.x86.D.Set(uint32(mod))
}

/****************************************************************************
REMARKS:
Implements the DIV instruction and side effects.
****************************************************************************/
func div_byte(s uint8) {
	var dvd, div, mod uint32

	dvd = M.x86.A.Get16()
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / uint8(s)
	mod = dvd % uint8(s)
	if abs(div) > 0xff {
		x86emu_intr_raise(0)
		return
	}
	M.x86.A.Set(uint8(div))
	M.x86.A.Set8h(uint8(mod))
}

/****************************************************************************
REMARKS:
Implements the DIV instruction and side effects.
****************************************************************************/
func div_word(s uint16) {
	var dvd, div, mod uint32

	dvd = (uint32(M.x86.D.Get16()) << 16) | M.x86.A.Get16()
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / uint16(s)
	mod = dvd % uint16(s)
	if abs(div) > 0xffff {
		x86emu_intr_raise(0)
		return
	}
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_SF)
	CONDITIONAL_SET_FLAG(div == 0, F_ZF)
	set_parity_flag(mod)

	M.x86.gen.A.Set16(uint16(div))
	M.x86.R_DX = uint16(mod)
}

/****************************************************************************
REMARKS:
Implements the DIV instruction and side effects.
****************************************************************************/
func div_long(s uint32) {
	var dvd, div, mod uint64

	dvd = (uint64(M.x86.D.Get32()) << 32) | uint64(M.x86.A.Get32())
	if s == 0 {
		x86emu_intr_raise(0)
		return
	}
	div = dvd / uint32(s)
	mod = dvd % uint32(s)
	if abs(div) > 0xffffffff {
		x86emu_intr_raise(0)
		return
	}
	CLEAR_FLAG(F_CF)
	CLEAR_FLAG(F_AF)
	CLEAR_FLAG(F_SF)
	SET_FLAG(F_ZF)
	set_parity_flag(mod)

	M.x86.A.Set(uint32(div))
	M.x86.D.Set(uint32(mod))
}

/****************************************************************************
REMARKS:
Implements the IN string instruction and side effects.
****************************************************************************/

func single_in(size int) {
	switch size {
	case 1:
		store_data_byte_abs(M.x86.seg.ES.Get(), M.x86.R_DI, sys_inb(M.x86.D.Get16()))
	case 2:
		store_data_word_abs(M.x86.seg.ES.Get(), M.x86.R_DI, sys_inw(M.x86.D.Get16()))
	default:
		store_data_long_abs(M.x86.seg.ES.Get(), M.x86.R_DI, sys_inl(M.x86.D.Get16()))
	}
}

func ins(size int) {
	inc := size

	if ACCESS_FLAG(F_DF) {
		inc = -size
	}
	if M.x86.mode & (SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE) {
		/* don't care whether REPE or REPNE */
		/* in until (E)CX is ZERO. */
		count := GetClrCount()
		for count > 0 {
			count--
			single_in(size)
			M.x86.DI.Add(uint32(inc))
		}
	} else {
		single_in(size)
		M.x86.DI.Add(uint32(inc))
	}
}

/****************************************************************************
REMARKS:
Implements the OUT string instruction and side effects.
****************************************************************************/

func single_out(size int) {
	switch size {
	case 1:
		sys_outb(M.x86.D.Get16(), fetch_data_byte_abs(M.x86.seg.ES.Get(), M.x86.R_SI))
	case 2:
		sys_outw(M.x86.D.Get16(), fetch_data_word_abs(M.x86.seg.ES.Get(), M.x86.R_SI))
	default:
		sys_outl(M.x86.D.Get16(), fetch_data_long_abs(M.x86.seg.ES.Get(), M.x86.R_SI))
	}
}

func outs(size int) {
	inc := uint16(size)

	if ACCESS_FLAG(F_DF) {
		inc = -size
	}
	if M.x86.mode & (SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE) {
		/* don't care whether REPE or REPNE */
		/* out until (E)CX is ZERO. */
		count := GetClrCount()
		for count > 0 {
			count--
			single_out(size)
			M.x86.SI.Add(inc)
		}
	} else {
		single_out(size)
		M.x86.SI.Add(inc)
	}
}

/****************************************************************************
PARAMETERS:
addr    - Address to fetch word from

REMARKS:
Fetches a word from emulator memory using an absolute address.
****************************************************************************/
func mem_access_word(addr uint32) uint16 {
	if CHECK_MEM_ACCESS() {
		x86emu_check_mem_access(addr)
	}
	return sys_rdw(addr)
}

/****************************************************************************
REMARKS:
Pushes a word onto the stack.

NOTE: Do not inline this, as (*sys_wrX) is already inline!
****************************************************************************/
func push_word(w uint16) {
	if CHECK_SP_ACCESS() {
		x86emu_check_sp_access()
	}
	M.x86.SP.Add(int16(-2))
	sys_wrw(uint32(M.x86.SS.Get16())<<4+M.x86.SP.Get16(), w)
}

/****************************************************************************
REMARKS:
Pushes a long onto the stack.

NOTE: Do not inline this, as (*sys_wrX) is already inline!
****************************************************************************/
func push_long(w uint32) {
	if CHECK_SP_ACCESS() {
		x86emu_check_sp_access()
	}
	M.x86.SP.Add(int16(-4))
	sys_wrl(uint32(M.x86.SS.Get16())<<4+M.x86.SP.Get16(), w)
}

/****************************************************************************
REMARKS:
Pops a word from the stack.

NOTE: Do not inline this, as (*sys_rdX) is already inline!
****************************************************************************/
func pop_word() uint16 {
	var res uint16

	if CHECK_SP_ACCESS() {
		x86emu_check_sp_access()
	}
	res = sys_rdw((uint32(M.x86.SS.Get()) << 4) + M.x86.SP.Get16())
	M.x86.SP.Add(int16(2))
	return res
}

/****************************************************************************
REMARKS:
Pops a long from the stack.

NOTE: Do not inline this, as (*sys_rdX) is already inline!
****************************************************************************/
func pop_long() uint32 {
	var res uint32

	if CHECK_SP_ACCESS() {
		x86emu_check_sp_access()
	}
	res = sys_rdl(uint32(M.x86.SS.Get())<<4 + M.x86.SP.Get16())
	M.x86.SP.Add(int16(4))
	return res
}

/****************************************************************************
REMARKS:
CPUID takes EAX/ECX as inputs, writes EAX/EBX/ECX/EDX as output
****************************************************************************/
func x86emu_cpuid() {
	feature := M.x86.A.Get32()

	switch feature {
	case 0:
		/* Regardless if we have real data from the hardware, the emulator
		 * will only support upto feature 1, which we set in register EAX.
		 * Registers EBX:EDX:ECX contain a string identifying the CPU.
		 */
		M.x86.gen.A.Set32(1)
		/* EBX:EDX:ECX = "GenuineIntel" */
		M.x86.gen.B.Set32(0x756e6547)
		M.x86.gen.D.Set32(0x49656e69)
		M.x86.gen.C.Set32(0x6c65746e)
		break
	case 1:
		/* If we don't have x86 compatible hardware, we return values from an
		 * Intel 486dx4; which was one of the first processors to have CPUID.
		 */
		M.x86.gen.A.Set32(0x00000480)
		M.x86.gen.B.Set32(0x00000000)
		M.x86.gen.C.Set32(0x00000000)
		M.x86.gen.D.Set32(0x00000002) /* VME */
		/* In the case that we have hardware CPUID instruction, we make sure
		 * that the features reported are limited to TSC and VME.
		 */
		M.x86.gen.D.Get32() &= 0x00000012
		break
	default:
		/* Finally, we don't support any additional features.  Most CPUs
		 * return all zeros when queried for invalid or unsupported feature
		 * numbers.
		 */
		M.x86.gen.A.Set32(0)
		M.x86.gen.B.Set32(0)
		M.x86.gen.C.Set32(0)
		M.x86.gen.D.Set32(0)
		break
	}
}

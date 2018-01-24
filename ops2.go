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
* Description:  This file includes subroutines to implement the decoding
*               and emulation of all the x86 extended two-byte processor
*               instructions.
*
****************************************************************************/

package main

import "fmt"

/*----------------------------- Implementation ----------------------------*/

/****************************************************************************
PARAMETERS:
op1 - Instruction op code

REMARKS:
Handles illegal opcodes.
****************************************************************************/
func x86emuOp2_illegal_op(op2 uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("ILLEGAL EXTENDED X86 OPCODE\n")
	TRACE_REGS()
	fmt.Printf("%04x:%04x: %02X ILLEGAL EXTENDED X86 OPCODE!\n",
		M.x86.seg.CS.Get(), M.x86.spc.IP.Get16()-2, op2)
	HALT_SYS()
	END_OF_INSTR()
}

/****************************************************************************
 * REMARKS:
 * Handles opcode 0x0f,0x01
 * ****************************************************************************/

func x86emuOp2_opc_01(op2 uint8) {
	const SMSW_INITIAL_VALUE = 0x10

	START_OF_INSTR()
	mod, rh, rl := fetch_decode_modrm()

	switch rh {
	case 4: // SMSW (Store Machine Status Word)
		// Decode the mod byte to find the addressing
		// Dummy implementation: Always returns 0x10 (initial value as per intel manual volume 3, figure 8-1)
		DECODE_PRINTF("SMSW\t")
		switch mod {
		case 0:
			destoffset := decode_rm00_address(rl)
			store_data_word(destoffset, SMSW_INITIAL_VALUE)
			break
		case 1:
			destoffset := decode_rm01_address(rl)
			store_data_word(destoffset, SMSW_INITIAL_VALUE)
			break
		case 2:
			destoffset := decode_rm10_address(rl)
			store_data_word(destoffset, SMSW_INITIAL_VALUE)
			break
		case 3:
			destreg := decode_rm_word_register(rl)
			destreg.Set(SMSW_INITIAL_VALUE)
			break
		}
		TRACE_AND_STEP()
		DECODE_CLEAR_SEGOVR()
		DECODE_PRINTF("\n")
		break
	default:
		DECODE_PRINTF("ILLEGAL EXTENDED X86 OPCODE IN 0F 01\n")
		TRACE_REGS()
		fmt.Printf("%04x:%04x: %02X ILLEGAL EXTENDED X86 OPCODE!\n",
			M.x86.seg.CS.Get(), M.x86.spc.IP.Get16()-2, op2)
		HALT_SYS()
		break
	}

	END_OF_INSTR()
}

/****************************************************************************
 * REMARKS:
 * Handles opcode 0x0f,0x08
 * ****************************************************************************/
func x86emuOp2_invd(op2 uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("INVD\n")
	TRACE_AND_STEP()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
 * REMARKS:
 * Handles opcode 0x0f,0x09
 * ****************************************************************************/
func x86emuOp2_wbinvd(op2 uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("WBINVD\n")
	TRACE_AND_STEP()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
 * REMARKS:
 * Handles opcode 0x0f,0x30
 * ****************************************************************************/
func x86emuOp2_wrmsr(op2 uint8) {
	/* dummy implementation, does nothing */

	START_OF_INSTR()
	DECODE_PRINTF("WRMSR\n")
	TRACE_AND_STEP()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0x31
****************************************************************************/
var counter uint64

func x86emuOp2_rdtsc(_ uint8) {
	counter += 0x10000

	/* read timestamp counter */
	/*
	 * Note that instead of actually trying to accurately measure this, we just
	 * increase the counter by a fixed amount every time we hit one of these
	 * instructions.  Feel free to come up with a better method.
	 */
	START_OF_INSTR()
	DECODE_PRINTF("RDTSC\n")
	TRACE_AND_STEP()
	M.x86.gen.A.Set32(uint32(counter))
	M.x86.gen.D.Set32(uint32(counter >> 32))
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
 * REMARKS:
 * Handles opcode 0x0f,0x32
 * ****************************************************************************/
func x86emuOp2_rdmsr(op2 uint8) {
	/* dummy implementation, always return 0 */

	START_OF_INSTR()
	DECODE_PRINTF("RDMSR\n")
	TRACE_AND_STEP()
	M.x86.gen.D.Set32(0)
	M.x86.gen.A.Set32(0)
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

func xorb(a, b bool) bool {
	if (!a && b) || (a && !b) {
		return true
	}
	return false
}
/****************************************************************************
REMARKS:
Handles opcode 0x0f,0x80-0x8F
****************************************************************************/
func x86emu_check_jump_condition(op uint8) bool {
	switch op {
	case 0x0:
		DECODE_PRINTF("JO\t")
		return ACCESS_FLAG(F_OF)
	case 0x1:
		DECODE_PRINTF("JNO\t")
		return !ACCESS_FLAG(F_OF)
		break
	case 0x2:
		DECODE_PRINTF("JB\t")
		return ACCESS_FLAG(F_CF)
		break
	case 0x3:
		DECODE_PRINTF("JNB\t")
		return !ACCESS_FLAG(F_CF)
		break
	case 0x4:
		DECODE_PRINTF("JZ\t")
		return ACCESS_FLAG(F_ZF)
		break
	case 0x5:
		DECODE_PRINTF("JNZ\t")
		return !ACCESS_FLAG(F_ZF)
		break
	case 0x6:
		DECODE_PRINTF("JBE\t")
		return ACCESS_FLAG(F_CF) || ACCESS_FLAG(F_ZF)
		break
	case 0x7:
		DECODE_PRINTF("JNBE\t")
		return !(ACCESS_FLAG(F_CF) || ACCESS_FLAG(F_ZF))
		break
	case 0x8:
		DECODE_PRINTF("JS\t")
		return ACCESS_FLAG(F_SF)
		break
	case 0x9:
		DECODE_PRINTF("JNS\t")
		return !ACCESS_FLAG(F_SF)
		break
	case 0xa:
		DECODE_PRINTF("JP\t")
		return ACCESS_FLAG(F_PF)
		break
	case 0xb:
		DECODE_PRINTF("JNP\t")
		return !ACCESS_FLAG(F_PF)
		break
	case 0xc:
		DECODE_PRINTF("JL\t")
		return xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF))
		break
	case 0xd:
		DECODE_PRINTF("JNL\t")
		return !xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF))
		break
	case 0xe:
		DECODE_PRINTF("JLE\t")
		return (xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF)) ||
			ACCESS_FLAG(F_ZF))
		break
	}
	DECODE_PRINTF("JNLE\t")
	return !(xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF)) ||
		ACCESS_FLAG(F_ZF))
}

func x86emuOp2_long_jump(op2 uint8) {

	/* conditional jump to word offset. */
	START_OF_INSTR()
	cond := x86emu_check_jump_condition(op2 & 0xF)
	target := int16(fetch_word_imm())
	target += int16(M.x86.spc.IP.Get16())
	DECODE_PRINTF2("%04x\n", target)
	TRACE_AND_STEP()
	if cond {
		M.x86.spc.IP.Set16(uint16(target))
		JMP_TRACE(M.x86.saved_cs, M.x86.saved_ip, M.x86.seg.CS.Get(), M.x86.spc.IP.Get16(), " LONG COND ")
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xC8-0xCF
****************************************************************************/
func x86emu_bswap(reg uint32) uint32 {
	// perform the byte swap
	temp := reg
	reg = (temp&0xFF000000)>>24 |
		(temp&0xFF0000)>>8 |
		(temp&0xFF00)<<8 |
		(temp&0xFF)<<24
	return reg
}

func x86emuOp2_bswap(op2 uint8) {
	/* byte swap 32 bit register */
	START_OF_INSTR()
	DECODE_PRINTF("BSWAP\t")
	switch op2 {
	case 0xc8:
		DECODE_PRINTF("EAX\n")
		M.x86.gen.A.Set32(x86emu_bswap(M.x86.gen.A.Get32()))
		break
	case 0xc9:
		DECODE_PRINTF("ECX\n")
		M.x86.gen.C.Set32(x86emu_bswap(M.x86.gen.C.Get32()))
		break
	case 0xca:
		DECODE_PRINTF("EDX\n")
		M.x86.gen.D.Set32(x86emu_bswap(M.x86.gen.D.Get32()))
		break
	case 0xcb:
		DECODE_PRINTF("EBX\n")
		M.x86.gen.B.Set32(x86emu_bswap(M.x86.gen.B.Get32()))
		break
	case 0xcc:
		DECODE_PRINTF("ESP\n")
		M.x86.spc.SP.Set32(x86emu_bswap(M.x86.spc.SP.Get32()))
		break
	case 0xcd:
		DECODE_PRINTF("EBP\n")
		M.x86.spc.BP.Set32(x86emu_bswap(M.x86.spc.BP.Get32()))
		break
	case 0xce:
		DECODE_PRINTF("ESI\n")
		M.x86.spc.SI.Set32(x86emu_bswap(M.x86.spc.SI.Get32()))
		break
	case 0xcf:
		DECODE_PRINTF("EDI\n")
		M.x86.spc.DI.Set32(x86emu_bswap(M.x86.spc.DI.Get32()))
		break
	}
	TRACE_AND_STEP()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0x90-0x9F
****************************************************************************/
func x86emuOp2_set_byte(op2 uint8) {
	var destreg register8
	var name string
	var cond bool

	START_OF_INSTR()
	switch op2 {
	case 0x90:
		name = "SETO\t"
		cond = ACCESS_FLAG(F_OF)
		break
	case 0x91:
		name = "SETNO\t"
		cond = !ACCESS_FLAG(F_OF)
		break
	case 0x92:
		name = "SETB\t"
		cond = ACCESS_FLAG(F_CF)
		break
	case 0x93:
		name = "SETNB\t"
		cond = !ACCESS_FLAG(F_CF)
		break
	case 0x94:
		name = "SETZ\t"
		cond = ACCESS_FLAG(F_ZF)
		break
	case 0x95:
		name = "SETNZ\t"
		cond = !ACCESS_FLAG(F_ZF)
		break
	case 0x96:
		name = "SETBE\t"
		cond = ACCESS_FLAG(F_CF) || ACCESS_FLAG(F_ZF)
		break
	case 0x97:
		name = "SETNBE\t"
		cond = !(ACCESS_FLAG(F_CF) || ACCESS_FLAG(F_ZF))
		break
	case 0x98:
		name = "SETS\t"
		cond = ACCESS_FLAG(F_SF)
		break
	case 0x99:
		name = "SETNS\t"
		cond = !ACCESS_FLAG(F_SF)
		break
	case 0x9a:
		name = "SETP\t"
		cond = ACCESS_FLAG(F_PF)
		break
	case 0x9b:
		name = "SETNP\t"
		cond = !ACCESS_FLAG(F_PF)
		break
	case 0x9c:
		name = "SETL\t"
		cond = xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF))
		break
	case 0x9d:
		name = "SETNL\t"
		cond = !xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF))
		break
	case 0x9e:
		name = "SETLE\t"
		cond = (xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF)) ||
			ACCESS_FLAG(F_ZF))
		break
	case 0x9f:
		name = "SETNLE\t"
		cond = !(xorb(ACCESS_FLAG(F_SF), ACCESS_FLAG(F_OF)) ||
			ACCESS_FLAG(F_ZF))
		break
	}
	DECODE_PRINTF(name)
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destoffset := decode_rmXX_address(mod, rl)
		TRACE_AND_STEP()
		if cond {
			store_data_byte(destoffset, 1)
		} else {
			store_data_byte(destoffset, 0)
		}
	} else { /* register to register */
		destreg := decode_rm_byte_register(rl)
		TRACE_AND_STEP()
		if cond {
			destreg.Set(1)
		} else {
			destreg.Set(0)
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa0
****************************************************************************/
func x86emuOp2_push_FS(_ uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("PUSH\tFS\n")
	TRACE_AND_STEP()
	push_word(M.x86.seg.FS.Get())
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa1
****************************************************************************/
func x86emuOp2_pop_FS(_ uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("POP\tFS\n")
	TRACE_AND_STEP()
	M.x86.seg.FS.Set( pop_word())
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS: CPUID takes EAX/ECX as inputs, writes EAX/EBX/ECX/EDX as output
Handles opcode 0x0f,0xa2
****************************************************************************/
func x86emuOp2_cpuid(_ uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("CPUID\n")
	TRACE_AND_STEP()
	x86emu_cpuid()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa3
****************************************************************************/
func x86emuOp2_bt_R(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BT\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		if M.x86.mode & SYSMODE_PREFIX_DATA != 0 {

			DECODE_PRINTF(",")
			shiftreg := decode_rm_long_register(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			disp := int16(shiftreg.Get()) >> 5
			srcval := fetch_data_long(srcoffset + disp)
			CONDITIONAL_SET_FLAG(srcval&(0x1<<bit), F_CF)
		} else {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			disp = int16(shiftreg.Get()) >> 4
			srcval := fetch_data_word(srcoffset + disp)
			CONDITIONAL_SET_FLAG(srcval&(0x1<<bit), F_CF)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			CONDITIONAL_SET_FLAG(*srcreg&(0x1<<bit), F_CF)
		} else {

			srcreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			CONDITIONAL_SET_FLAG(*srcreg&(0x1<<bit), F_CF)
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa4
****************************************************************************/
func x86emuOp2_shld_IMM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("SHLD\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destoffset := decode_rmXX_address(mod, rl)
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destval := fetch_data_long(destoffset)
			destval = shld_long(destval, shiftreg.Get(), shift)
			store_data_long(destoffset, destval)
		} else {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destval := fetch_data_word(destoffset)
			destval = shld_word(destval, shiftreg.Get(), shift)
			store_data_word(destoffset, destval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destreg.Set(shld_long(destreg.Get(), shiftreg.Get(), shift))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destreg.Set(shld_word(destreg.Get(), shiftreg.Get(), shift))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa5
****************************************************************************/
func x86emuOp2_shld_CL(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("SHLD\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destoffset := decode_rmXX_address(mod, rl)
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destval := fetch_data_long(destoffset)
			destval = shld_long(destval, shiftreg.Get(), M.x86.R_CL)
			store_data_long(destoffset, destval)
		} else {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destval := fetch_data_word(destoffset)
			destval = shld_word(destval, shiftreg.Get(), M.x86.R_CL)
			store_data_word(destoffset, destval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destreg.Set(shld_long(destreg.Get(), shiftreg.Get(), M.x86.R_CL))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destreg.Set(shld_word(destreg.Get(), shiftreg.Get(), M.x86.R_CL))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa8
****************************************************************************/
func x86emuOp2_push_GS(_ uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("PUSH\tGS\n")
	TRACE_AND_STEP()
	push_word(M.x86.R_GS)
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xa9
****************************************************************************/
func x86emuOp2_pop_GS(_ uint8) {
	START_OF_INSTR()
	DECODE_PRINTF("POP\tGS\n")
	TRACE_AND_STEP()
	M.x86.R_GS = pop_word()
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xab
****************************************************************************/
func x86emuOp2_bts_R(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BTS\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			disp = int16(shiftreg.Get()) >> 5
			srcval := fetch_data_long(srcoffset + disp)
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_long(srcoffset+disp, srcval|mask)
		} else {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			disp = int16(shiftreg.Get()) >> 4
			srcval := fetch_data_word(srcoffset + disp)
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_word(srcoffset+disp, srcval|mask)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg |= mask
		} else {

			srcreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg |= mask
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xac
****************************************************************************/
func x86emuOp2_shrd_IMM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("SHLD\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destoffset := decode_rmXX_address(mod, rl)
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destval := fetch_data_long(destoffset)
			destval = shrd_long(destval, shiftreg.Get(), shift)
			store_data_long(destoffset, destval)
		} else {

			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destval := fetch_data_word(destoffset)
			destval = shrd_word(destval, shiftreg.Get(), shift)
			store_data_word(destoffset, destval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destreg.Set(shrd_long(destreg.Get(), shiftreg.Get(), shift))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			shift = fetch_byte_imm()
			DECODE_PRINTF2("%d\n", shift)
			TRACE_AND_STEP()
			destreg.Set(shrd_word(destreg.Get(), shiftreg.Get(), shift))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xad
****************************************************************************/
func x86emuOp2_shrd_CL(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("SHLD\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF(",")
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destval := fetch_data_long(destoffset)
			destval = shrd_long(destval, shiftreg.Get(), M.x86.R_CL)
			store_data_long(destoffset, destval)
		} else {

			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destval := fetch_data_word(destoffset)
			destval = shrd_word(destval, shiftreg.Get(), M.x86.R_CL)
			store_data_word(destoffset, destval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destreg.Set(shrd_long(destreg.Get(), shiftreg.Get(), M.x86.R_CL))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",CL\n")
			TRACE_AND_STEP()
			destreg.Set(shrd_word(destreg.Get(), shiftreg.Get(), M.x86.R_CL))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xaf
****************************************************************************/
func x86emuOp2_imul_R_RM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("IMUL\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := fetch_data_long(srcoffset)
			TRACE_AND_STEP()
			imul_long_direct(&res_lo, &res_hi, int32(destreg.Get()), int32(srcval))
			if res_hi != 0 {
				SET_FLAG(F_CF)
				SET_FLAG(F_OF)
			} else {
				CLEAR_FLAG(F_CF)
				CLEAR_FLAG(F_OF)
			}
			destreg.Set32(uint32(res_lo))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := fetch_data_word(srcoffset)
			TRACE_AND_STEP()
			res = int16(destreg.Get()) * int16(srcval)
			if res > 0xFFFF {
				SET_FLAG(F_CF)
				SET_FLAG(F_OF)
			} else {
				CLEAR_FLAG(F_CF)
				CLEAR_FLAG(F_OF)
			}
			destreg.Set16(uint16(res))
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_LONG_REGISTER(rl)
			TRACE_AND_STEP()
			imul_long_direct(&res_lo, &res_hi, int32(destreg.Get()), int32(srcreg.Get()))
			if res_hi != 0 {
				SET_FLAG(F_CF)
				SET_FLAG(F_OF)
			} else {
				CLEAR_FLAG(F_CF)
				CLEAR_FLAG(F_OF)
			}
			destreg.Set32(uint32(res_lo))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_WORD_REGISTER(rl)
			res = int16(destreg.Get()) * int16(srcreg.Get())
			if res > 0xFFFF {
				SET_FLAG(F_CF)
				SET_FLAG(F_OF)
			} else {
				CLEAR_FLAG(F_CF)
				CLEAR_FLAG(F_OF)
			}
			destreg.Set16(uint16(res))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb2
****************************************************************************/
func x86emuOp2_lss_R_IMM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("LSS\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		dstreg := DECODE_RM_WORD_REGISTER(rh)
		DECODE_PRINTF(",")
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		dstreg.Set(fetch_data_word(srcoffset))
		M.x86.seg.SS.Set(fetch_data_word(srcoffset + 2))
	} else { /* register to register */
		/* UNDEFINED! */
		TRACE_AND_STEP()
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb3
****************************************************************************/
func x86emuOp2_btr_R(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BTR\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF(",")
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			disp = int16(shiftreg.Get()) >> 5
			srcval := fetch_data_long(srcoffset + disp)
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_long(srcoffset+disp, srcval & ^mask)
		} else {

			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			disp = int16(shiftreg.Get()) >> 4
			srcval := fetch_data_word(srcoffset + disp)
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_word(srcoffset+disp, int16(srcval & ^mask))
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg &= ^mask
		} else {

			srcreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg &= ^mask
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb4
****************************************************************************/
func x86emuOp2_lfs_R_IMM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("LFS\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		dstreg := DECODE_RM_WORD_REGISTER(rh)
		DECODE_PRINTF(",")
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		dstreg.Set(fetch_data_word(srcoffset))
		M.x86.R_FS = fetch_data_word(srcoffset + 2)
	} else { /* register to register */
		/* UNDEFINED! */
		TRACE_AND_STEP()
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb5
****************************************************************************/
func x86emuOp2_lgs_R_IMM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("LGS\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		dstreg := DECODE_RM_WORD_REGISTER(rh)
		DECODE_PRINTF(",")
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		dstreg.Set(fetch_data_word(srcoffset))
		M.x86.R_GS = fetch_data_word(srcoffset + 2)
	} else { /* register to register */
		/* UNDEFINED! */
		TRACE_AND_STEP()
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb6
****************************************************************************/
func x86emuOp2_movzx_byte_R_RM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("MOVZX\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := fetch_data_byte(srcoffset)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(srcval)
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := fetch_data_byte(srcoffset)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(srcval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_BYTE_REGISTER(rl)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(*srcreg)
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_BYTE_REGISTER(rl)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(*srcreg)
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xb7
****************************************************************************/
func x86emuOp2_movzx_word_R_RM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("MOVZX\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destreg := DECODE_RM_LONG_REGISTER(rh)
		DECODE_PRINTF(",")
		srcoffset := decode_rmXX_address(mod, rl)
		srcval := fetch_data_word(srcoffset)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		destreg.Set(srcval)
	} else { /* register to register */
		destreg := DECODE_RM_LONG_REGISTER(rh)
		DECODE_PRINTF(",")
		srcreg := DECODE_RM_WORD_REGISTER(rl)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		destreg.Set(*srcreg)
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xba
****************************************************************************/
func x86emuOp2_btX_I(_ uint8) {

	START_OF_INSTR()
	mod, rh, rl := fetch_decode_modrm()
	switch rh {
	case 4:
		DECODE_PRINTF("BT\t")
		break
	case 5:
		DECODE_PRINTF("BTS\t")
		break
	case 6:
		DECODE_PRINTF("BTR\t")
		break
	case 7:
		DECODE_PRINTF("BTC\t")
		break
	default:
		DECODE_PRINTF("ILLEGAL EXTENDED X86 OPCODE\n")
		TRACE_REGS()
		fmt.Printf("%04x:%04x: %02X%02X ILLEGAL EXTENDED X86 OPCODE EXTENSION!\n",
			M.x86.seg.CS.Get(), M.x86.spc.IP.Get16()-3, op2, (mod<<6)|(rh<<3)|rl)
		HALT_SYS()
	}
	if mod < 3 {

		srcoffset := decode_rmXX_address(mod, rl)
		shift = fetch_byte_imm()
		DECODE_PRINTF2(",%d\n", shift)
		TRACE_AND_STEP()

		if M.x86.mode & SYSMODE_PREFIX_DATA {

			bit := shift & 0x1F
			srcval := fetch_data_long(srcoffset)
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			switch rh {
			case 5:
				store_data_long(srcoffset, srcval|mask)
				break
			case 6:
				store_data_long(srcoffset, srcval & ^mask)
				break
			case 7:
				store_data_long(srcoffset, srcval^mask)
				break
			default:
				break
			}
		} else {

			bit := shift & 0xF
			srcval := fetch_data_word(srcoffset)
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			switch rh {
			case 5:
				store_data_word(srcoffset, srcval|mask)
				break
			case 6:
				store_data_word(srcoffset, srcval & ^mask)
				break
			case 7:
				store_data_word(srcoffset, srcval^mask)
				break
			default:
				break
			}
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcreg := DECODE_RM_LONG_REGISTER(rl)
			shift = fetch_byte_imm()
			DECODE_PRINTF2(",%d\n", shift)
			TRACE_AND_STEP()
			bit := shift & 0x1F
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			switch rh {
			case 5:
				*srcreg |= mask
				break
			case 6:
				*srcreg &= ^mask
				break
			case 7:
				*srcreg ^= mask
				break
			default:
				break
			}
		} else {

			srcreg := DECODE_RM_WORD_REGISTER(rl)
			shift = fetch_byte_imm()
			DECODE_PRINTF2(",%d\n", shift)
			TRACE_AND_STEP()
			bit := shift & 0xF
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			switch rh {
			case 5:
				*srcreg |= mask
				break
			case 6:
				*srcreg &= ^mask
				break
			case 7:
				*srcreg ^= mask
				break
			default:
				break
			}
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xbb
****************************************************************************/
func x86emuOp2_btc_R(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BTC\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF(",")
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			disp = int16(shiftreg.Get()) >> 5
			srcval := fetch_data_long(srcoffset + disp)
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_long(srcoffset+disp, srcval^mask)
		} else {

			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			disp = int16(shiftreg.Get()) >> 4
			srcval := fetch_data_word(srcoffset + disp)
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(srcval&mask, F_CF)
			store_data_word(srcoffset+disp, int16(srcval^mask))
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcreg := DECODE_RM_LONG_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0x1F
			mask = (0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg ^= mask
		} else {

			srcreg := DECODE_RM_WORD_REGISTER(rl)
			DECODE_PRINTF(",")
			shiftreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			bit := shiftreg.Get() & 0xF
			mask = int16(0x1 << bit)
			CONDITIONAL_SET_FLAG(*srcreg&mask, F_CF)
			*srcreg ^= mask
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xbc
****************************************************************************/
func x86emuOp2_bsf(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BSF\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF(",")
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			dstreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			srcval := fetch_data_long(srcoffset)
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(0); dstreg.Get() < 32; destreg.Set(destreg.Get() + 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		} else {

			dstreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			srcval := fetch_data_word(srcoffset)
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(0); dstreg.Get() < 16; destreg.Set(destreg.Get() + 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcval := DECODE_RM_LONG_REGISTER(rl).Get()
			DECODE_PRINTF(",")
			dstreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(0); dstreg.Get() < 32; destreg.Set(destreg.Get() + 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		} else {

			srcval := DECODE_RM_WORD_REGISTER(rl).Get()
			DECODE_PRINTF(",")
			dstreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(0); dstreg.Get() < 16; destreg.Set(destreg.Get() + 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}

		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xbd
****************************************************************************/
func x86emuOp2_bsr(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("BSR\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		srcoffset := decode_rmXX_address(mod, rl)
		DECODE_PRINTF(",")
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			dstreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			srcval := fetch_data_long(srcoffset)
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(31); dstreg.Get() > 0; destreg.Set(destreg.Get() - 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		} else {

			dstreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			srcval := fetch_data_word(srcoffset)
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(15); dstreg.Get() > 0; destreg.Set(destreg.Get() - 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			srcval := DECODE_RM_LONG_REGISTER(rl).Get()
			DECODE_PRINTF(",")
			dstreg := DECODE_RM_LONG_REGISTER(rh)
			TRACE_AND_STEP()
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(31); dstreg.Get() > 0; destreg.Set(destreg.Get() - 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		} else {

			srcval := DECODE_RM_WORD_REGISTER(rl).Get()
			DECODE_PRINTF(",")
			dstreg := DECODE_RM_WORD_REGISTER(rh)
			TRACE_AND_STEP()
			CONDITIONAL_SET_FLAG(srcval == 0, F_ZF)
			for dstreg.Set(15); dstreg.Get() > 0; destreg.Set(destreg.Get() - 1) {
				if ((srcval >> dstreg.Get()) & 1) != 0 {
					break
				}
			}
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xbe
****************************************************************************/
func x86emuOp2_movsx_byte_R_RM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("MOVSX\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := int32((int8(fetch_data_byte)(srcoffset)))
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(srcval)
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcoffset := decode_rmXX_address(mod, rl)
			srcval := int16((int8(fetch_data_byte)(srcoffset)))
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(srcval)
		}
	} else { /* register to register */
		if M.x86.mode & SYSMODE_PREFIX_DATA {

			destreg := DECODE_RM_LONG_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_BYTE_REGISTER(rl)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(int32(srcreg.Get()))
		} else {

			destreg := DECODE_RM_WORD_REGISTER(rh)
			DECODE_PRINTF(",")
			srcreg := DECODE_RM_BYTE_REGISTER(rl)
			DECODE_PRINTF("\n")
			TRACE_AND_STEP()
			destreg.Set(int16(srcreg.Get()))
		}
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/****************************************************************************
REMARKS:
Handles opcode 0x0f,0xbf
****************************************************************************/
func x86emuOp2_movsx_word_R_RM(_ uint8) {

	START_OF_INSTR()
	DECODE_PRINTF("MOVSX\t")
	mod, rh, rl := fetch_decode_modrm()
	if mod < 3 {
		destreg := DECODE_RM_LONG_REGISTER(rh)
		DECODE_PRINTF(",")
		srcoffset := decode_rmXX_address(mod, rl)
		srcval := int32((int16(fetch_data_word)(srcoffset)))
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		destreg.Set(srcval)
	} else { /* register to register */
		destreg := DECODE_RM_LONG_REGISTER(rh)
		DECODE_PRINTF(",")
		srcreg := DECODE_RM_WORD_REGISTER(rl)
		DECODE_PRINTF("\n")
		TRACE_AND_STEP()
		destreg.Set(int32(srcreg.Get()))
	}
	DECODE_CLEAR_SEGOVR()
	END_OF_INSTR()
}

/***************************************************************************
 * Double byte operation code table:
 **************************************************************************/
var x86emu_optab2 = [256]optab{
	/*  0x00 */ x86emuOp2_illegal_op, /* Group F (ring 0 PM)      */
	/*  0x01 */ x86emuOp2_opc_01, /* Group G (ring 0 PM)      */
	/*  0x02 */ x86emuOp2_illegal_op, /* lar (ring 0 PM)          */
	/*  0x03 */ x86emuOp2_illegal_op, /* lsl (ring 0 PM)          */
	/*  0x04 */ x86emuOp2_illegal_op,
	/*  0x05 */ x86emuOp2_illegal_op, /* loadall (undocumented)   */
	/*  0x06 */ x86emuOp2_illegal_op, /* clts (ring 0 PM)         */
	/*  0x07 */ x86emuOp2_illegal_op, /* loadall (undocumented)   */
	/*  0x08 */ x86emuOp2_invd, /* invd (ring 0 PM)         */
	/*  0x09 */ x86emuOp2_wbinvd, /* wbinvd (ring 0 PM)       */
	/*  0x0a */ x86emuOp2_illegal_op,
	/*  0x0b */ x86emuOp2_illegal_op,
	/*  0x0c */ x86emuOp2_illegal_op,
	/*  0x0d */ x86emuOp2_illegal_op,
	/*  0x0e */ x86emuOp2_illegal_op,
	/*  0x0f */ x86emuOp2_illegal_op,

	/*  0x10 */ x86emuOp2_illegal_op,
	/*  0x11 */ x86emuOp2_illegal_op,
	/*  0x12 */ x86emuOp2_illegal_op,
	/*  0x13 */ x86emuOp2_illegal_op,
	/*  0x14 */ x86emuOp2_illegal_op,
	/*  0x15 */ x86emuOp2_illegal_op,
	/*  0x16 */ x86emuOp2_illegal_op,
	/*  0x17 */ x86emuOp2_illegal_op,
	/*  0x18 */ x86emuOp2_illegal_op,
	/*  0x19 */ x86emuOp2_illegal_op,
	/*  0x1a */ x86emuOp2_illegal_op,
	/*  0x1b */ x86emuOp2_illegal_op,
	/*  0x1c */ x86emuOp2_illegal_op,
	/*  0x1d */ x86emuOp2_illegal_op,
	/*  0x1e */ x86emuOp2_illegal_op,
	/*  0x1f */ x86emuOp2_illegal_op,

	/*  0x20 */ x86emuOp2_illegal_op, /* mov reg32,creg (ring 0 PM) */
	/*  0x21 */ x86emuOp2_illegal_op, /* mov reg32,dreg (ring 0 PM) */
	/*  0x22 */ x86emuOp2_illegal_op, /* mov creg,reg32 (ring 0 PM) */
	/*  0x23 */ x86emuOp2_illegal_op, /* mov dreg,reg32 (ring 0 PM) */
	/*  0x24 */ x86emuOp2_illegal_op, /* mov reg32,treg (ring 0 PM) */
	/*  0x25 */ x86emuOp2_illegal_op,
	/*  0x26 */ x86emuOp2_illegal_op, /* mov treg,reg32 (ring 0 PM) */
	/*  0x27 */ x86emuOp2_illegal_op,
	/*  0x28 */ x86emuOp2_illegal_op,
	/*  0x29 */ x86emuOp2_illegal_op,
	/*  0x2a */ x86emuOp2_illegal_op,
	/*  0x2b */ x86emuOp2_illegal_op,
	/*  0x2c */ x86emuOp2_illegal_op,
	/*  0x2d */ x86emuOp2_illegal_op,
	/*  0x2e */ x86emuOp2_illegal_op,
	/*  0x2f */ x86emuOp2_illegal_op,

	/*  0x30 */ x86emuOp2_wrmsr,
	/*  0x31 */ x86emuOp2_rdtsc,
	/*  0x32 */ x86emuOp2_rdmsr,
	/*  0x33 */ x86emuOp2_illegal_op,
	/*  0x34 */ x86emuOp2_illegal_op,
	/*  0x35 */ x86emuOp2_illegal_op,
	/*  0x36 */ x86emuOp2_illegal_op,
	/*  0x37 */ x86emuOp2_illegal_op,
	/*  0x38 */ x86emuOp2_illegal_op,
	/*  0x39 */ x86emuOp2_illegal_op,
	/*  0x3a */ x86emuOp2_illegal_op,
	/*  0x3b */ x86emuOp2_illegal_op,
	/*  0x3c */ x86emuOp2_illegal_op,
	/*  0x3d */ x86emuOp2_illegal_op,
	/*  0x3e */ x86emuOp2_illegal_op,
	/*  0x3f */ x86emuOp2_illegal_op,

	/*  0x40 */ x86emuOp2_illegal_op,
	/*  0x41 */ x86emuOp2_illegal_op,
	/*  0x42 */ x86emuOp2_illegal_op,
	/*  0x43 */ x86emuOp2_illegal_op,
	/*  0x44 */ x86emuOp2_illegal_op,
	/*  0x45 */ x86emuOp2_illegal_op,
	/*  0x46 */ x86emuOp2_illegal_op,
	/*  0x47 */ x86emuOp2_illegal_op,
	/*  0x48 */ x86emuOp2_illegal_op,
	/*  0x49 */ x86emuOp2_illegal_op,
	/*  0x4a */ x86emuOp2_illegal_op,
	/*  0x4b */ x86emuOp2_illegal_op,
	/*  0x4c */ x86emuOp2_illegal_op,
	/*  0x4d */ x86emuOp2_illegal_op,
	/*  0x4e */ x86emuOp2_illegal_op,
	/*  0x4f */ x86emuOp2_illegal_op,

	/*  0x50 */ x86emuOp2_illegal_op,
	/*  0x51 */ x86emuOp2_illegal_op,
	/*  0x52 */ x86emuOp2_illegal_op,
	/*  0x53 */ x86emuOp2_illegal_op,
	/*  0x54 */ x86emuOp2_illegal_op,
	/*  0x55 */ x86emuOp2_illegal_op,
	/*  0x56 */ x86emuOp2_illegal_op,
	/*  0x57 */ x86emuOp2_illegal_op,
	/*  0x58 */ x86emuOp2_illegal_op,
	/*  0x59 */ x86emuOp2_illegal_op,
	/*  0x5a */ x86emuOp2_illegal_op,
	/*  0x5b */ x86emuOp2_illegal_op,
	/*  0x5c */ x86emuOp2_illegal_op,
	/*  0x5d */ x86emuOp2_illegal_op,
	/*  0x5e */ x86emuOp2_illegal_op,
	/*  0x5f */ x86emuOp2_illegal_op,

	/*  0x60 */ x86emuOp2_illegal_op,
	/*  0x61 */ x86emuOp2_illegal_op,
	/*  0x62 */ x86emuOp2_illegal_op,
	/*  0x63 */ x86emuOp2_illegal_op,
	/*  0x64 */ x86emuOp2_illegal_op,
	/*  0x65 */ x86emuOp2_illegal_op,
	/*  0x66 */ x86emuOp2_illegal_op,
	/*  0x67 */ x86emuOp2_illegal_op,
	/*  0x68 */ x86emuOp2_illegal_op,
	/*  0x69 */ x86emuOp2_illegal_op,
	/*  0x6a */ x86emuOp2_illegal_op,
	/*  0x6b */ x86emuOp2_illegal_op,
	/*  0x6c */ x86emuOp2_illegal_op,
	/*  0x6d */ x86emuOp2_illegal_op,
	/*  0x6e */ x86emuOp2_illegal_op,
	/*  0x6f */ x86emuOp2_illegal_op,

	/*  0x70 */ x86emuOp2_illegal_op,
	/*  0x71 */ x86emuOp2_illegal_op,
	/*  0x72 */ x86emuOp2_illegal_op,
	/*  0x73 */ x86emuOp2_illegal_op,
	/*  0x74 */ x86emuOp2_illegal_op,
	/*  0x75 */ x86emuOp2_illegal_op,
	/*  0x76 */ x86emuOp2_illegal_op,
	/*  0x77 */ x86emuOp2_illegal_op,
	/*  0x78 */ x86emuOp2_illegal_op,
	/*  0x79 */ x86emuOp2_illegal_op,
	/*  0x7a */ x86emuOp2_illegal_op,
	/*  0x7b */ x86emuOp2_illegal_op,
	/*  0x7c */ x86emuOp2_illegal_op,
	/*  0x7d */ x86emuOp2_illegal_op,
	/*  0x7e */ x86emuOp2_illegal_op,
	/*  0x7f */ x86emuOp2_illegal_op,

	/*  0x80 */ x86emuOp2_long_jump,
	/*  0x81 */ x86emuOp2_long_jump,
	/*  0x82 */ x86emuOp2_long_jump,
	/*  0x83 */ x86emuOp2_long_jump,
	/*  0x84 */ x86emuOp2_long_jump,
	/*  0x85 */ x86emuOp2_long_jump,
	/*  0x86 */ x86emuOp2_long_jump,
	/*  0x87 */ x86emuOp2_long_jump,
	/*  0x88 */ x86emuOp2_long_jump,
	/*  0x89 */ x86emuOp2_long_jump,
	/*  0x8a */ x86emuOp2_long_jump,
	/*  0x8b */ x86emuOp2_long_jump,
	/*  0x8c */ x86emuOp2_long_jump,
	/*  0x8d */ x86emuOp2_long_jump,
	/*  0x8e */ x86emuOp2_long_jump,
	/*  0x8f */ x86emuOp2_long_jump,

	/*  0x90 */ x86emuOp2_set_byte,
	/*  0x91 */ x86emuOp2_set_byte,
	/*  0x92 */ x86emuOp2_set_byte,
	/*  0x93 */ x86emuOp2_set_byte,
	/*  0x94 */ x86emuOp2_set_byte,
	/*  0x95 */ x86emuOp2_set_byte,
	/*  0x96 */ x86emuOp2_set_byte,
	/*  0x97 */ x86emuOp2_set_byte,
	/*  0x98 */ x86emuOp2_set_byte,
	/*  0x99 */ x86emuOp2_set_byte,
	/*  0x9a */ x86emuOp2_set_byte,
	/*  0x9b */ x86emuOp2_set_byte,
	/*  0x9c */ x86emuOp2_set_byte,
	/*  0x9d */ x86emuOp2_set_byte,
	/*  0x9e */ x86emuOp2_set_byte,
	/*  0x9f */ x86emuOp2_set_byte,

	/*  0xa0 */ x86emuOp2_push_FS,
	/*  0xa1 */ x86emuOp2_pop_FS,
	/*  0xa2 */ x86emuOp2_cpuid,
	/*  0xa3 */ x86emuOp2_bt_R,
	/*  0xa4 */ x86emuOp2_shld_IMM,
	/*  0xa5 */ x86emuOp2_shld_CL,
	/*  0xa6 */ x86emuOp2_illegal_op,
	/*  0xa7 */ x86emuOp2_illegal_op,
	/*  0xa8 */ x86emuOp2_push_GS,
	/*  0xa9 */ x86emuOp2_pop_GS,
	/*  0xaa */ x86emuOp2_illegal_op,
	/*  0xab */ x86emuOp2_bts_R,
	/*  0xac */ x86emuOp2_shrd_IMM,
	/*  0xad */ x86emuOp2_shrd_CL,
	/*  0xae */ x86emuOp2_illegal_op,
	/*  0xaf */ x86emuOp2_imul_R_RM,

	/*  0xb0 */ x86emuOp2_illegal_op, /* TODO: cmpxchg */
	/*  0xb1 */ x86emuOp2_illegal_op, /* TODO: cmpxchg */
	/*  0xb2 */ x86emuOp2_lss_R_IMM,
	/*  0xb3 */ x86emuOp2_btr_R,
	/*  0xb4 */ x86emuOp2_lfs_R_IMM,
	/*  0xb5 */ x86emuOp2_lgs_R_IMM,
	/*  0xb6 */ x86emuOp2_movzx_byte_R_RM,
	/*  0xb7 */ x86emuOp2_movzx_word_R_RM,
	/*  0xb8 */ x86emuOp2_illegal_op,
	/*  0xb9 */ x86emuOp2_illegal_op,
	/*  0xba */ x86emuOp2_btX_I,
	/*  0xbb */ x86emuOp2_btc_R,
	/*  0xbc */ x86emuOp2_bsf,
	/*  0xbd */ x86emuOp2_bsr,
	/*  0xbe */ x86emuOp2_movsx_byte_R_RM,
	/*  0xbf */ x86emuOp2_movsx_word_R_RM,

	/*  0xc0 */ x86emuOp2_illegal_op, /* TODO: xadd */
	/*  0xc1 */ x86emuOp2_illegal_op, /* TODO: xadd */
	/*  0xc2 */ x86emuOp2_illegal_op,
	/*  0xc3 */ x86emuOp2_illegal_op,
	/*  0xc4 */ x86emuOp2_illegal_op,
	/*  0xc5 */ x86emuOp2_illegal_op,
	/*  0xc6 */ x86emuOp2_illegal_op,
	/*  0xc7 */ x86emuOp2_illegal_op,
	/*  0xc8 */ x86emuOp2_bswap,
	/*  0xc9 */ x86emuOp2_bswap,
	/*  0xca */ x86emuOp2_bswap,
	/*  0xcb */ x86emuOp2_bswap,
	/*  0xcc */ x86emuOp2_bswap,
	/*  0xcd */ x86emuOp2_bswap,
	/*  0xce */ x86emuOp2_bswap,
	/*  0xcf */ x86emuOp2_bswap,

	/*  0xd0 */ x86emuOp2_illegal_op,
	/*  0xd1 */ x86emuOp2_illegal_op,
	/*  0xd2 */ x86emuOp2_illegal_op,
	/*  0xd3 */ x86emuOp2_illegal_op,
	/*  0xd4 */ x86emuOp2_illegal_op,
	/*  0xd5 */ x86emuOp2_illegal_op,
	/*  0xd6 */ x86emuOp2_illegal_op,
	/*  0xd7 */ x86emuOp2_illegal_op,
	/*  0xd8 */ x86emuOp2_illegal_op,
	/*  0xd9 */ x86emuOp2_illegal_op,
	/*  0xda */ x86emuOp2_illegal_op,
	/*  0xdb */ x86emuOp2_illegal_op,
	/*  0xdc */ x86emuOp2_illegal_op,
	/*  0xdd */ x86emuOp2_illegal_op,
	/*  0xde */ x86emuOp2_illegal_op,
	/*  0xdf */ x86emuOp2_illegal_op,

	/*  0xe0 */ x86emuOp2_illegal_op,
	/*  0xe1 */ x86emuOp2_illegal_op,
	/*  0xe2 */ x86emuOp2_illegal_op,
	/*  0xe3 */ x86emuOp2_illegal_op,
	/*  0xe4 */ x86emuOp2_illegal_op,
	/*  0xe5 */ x86emuOp2_illegal_op,
	/*  0xe6 */ x86emuOp2_illegal_op,
	/*  0xe7 */ x86emuOp2_illegal_op,
	/*  0xe8 */ x86emuOp2_illegal_op,
	/*  0xe9 */ x86emuOp2_illegal_op,
	/*  0xea */ x86emuOp2_illegal_op,
	/*  0xeb */ x86emuOp2_illegal_op,
	/*  0xec */ x86emuOp2_illegal_op,
	/*  0xed */ x86emuOp2_illegal_op,
	/*  0xee */ x86emuOp2_illegal_op,
	/*  0xef */ x86emuOp2_illegal_op,

	/*  0xf0 */ x86emuOp2_illegal_op,
	/*  0xf1 */ x86emuOp2_illegal_op,
	/*  0xf2 */ x86emuOp2_illegal_op,
	/*  0xf3 */ x86emuOp2_illegal_op,
	/*  0xf4 */ x86emuOp2_illegal_op,
	/*  0xf5 */ x86emuOp2_illegal_op,
	/*  0xf6 */ x86emuOp2_illegal_op,
	/*  0xf7 */ x86emuOp2_illegal_op,
	/*  0xf8 */ x86emuOp2_illegal_op,
	/*  0xf9 */ x86emuOp2_illegal_op,
	/*  0xfa */ x86emuOp2_illegal_op,
	/*  0xfb */ x86emuOp2_illegal_op,
	/*  0xfc */ x86emuOp2_illegal_op,
	/*  0xfd */ x86emuOp2_illegal_op,
	/*  0xfe */ x86emuOp2_illegal_op,
	/*  0xff */ x86emuOp2_illegal_op,
}

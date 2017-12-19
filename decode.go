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
* Description:  This file includes subroutines which are related to
*               instruction decoding and accesses of immediate data via IP.  etc.
*
****************************************************************************/

package main

import (
	"fmt"
	"os"
)

/*----------------------------- Implementation ----------------------------*/

/****************************************************************************
REMARKS:
Handles any pending asynchronous interrupts.
****************************************************************************/
func x86emu_intr_handle() {
	var intno uint8

	if M().x86.intr&INTR_SYNCH != 0 {
		intno = uint8(M().x86.intno)
		if _X86EMU_intrTab[intno] != nil {
			panic("_X86EMU_intrTab[intno](intno)")
		} else {
			push_word(uint16(M().x86.spc.FLAGS))
			CLEAR_FLAG(F_IF)
			CLEAR_FLAG(F_TF)
			push_word(M().x86.seg.CS.Get())
			M().x86.seg.CS.Set(mem_access_word(uint32(intno*4 + 2)))
			push_word(M().x86.spc.IP.Get16())
			M().x86.spc.IP.Set16(mem_access_word(uint32(intno * 4)))
			M().x86.intr = 0
		}
	}
}

/****************************************************************************
PARAMETERS:
intrnum - Interrupt number to raise

REMARKS:
Raise the specified interrupt to be handled before the execution of the
next instruction.
****************************************************************************/
func x86emu_intr_raise(intrnum uint8) {
	fmt.Printf("raising exception %x\n", intrnum)
	x86emu_dump_regs()
	M().x86.intno = intrnum
	M().x86.intr |= INTR_SYNCH
}

/****************************************************************************
REMARKS:
Main execution loop for the emulator. We return from here when the system
halts, which is normally caused by a stack fault when we return from the
original real mode call.
****************************************************************************/
func X86EMU_exec() {
	var op1 uint8

	M().x86.intr = 0
	x86emu_end_instr()

	for {
		if CHECK_IP_FETCH() {
			x86emu_check_ip_access()
		}

		/* If debugging, save the IP and CS values. */
		SAVE_IP_CS(M().x86.seg.CS.Get(), M().x86.spc.IP.Get16())
		INC_DECODED_INST_LEN(1)
		if M().x86.intr != 0 {
			if (M().x86.intr & INTR_HALTED) != 0 {
				if M().x86.spc.SP.Get16() != 0 {
					fmt.Printf("halted\n")
					X86EMU_trace_regs()
				} else {
					if M().x86.debug != 0 {
						fmt.Printf("Service completed successfully\n")
					}
				}
				return
			}
			if ((M().x86.intr & INTR_SYNCH != 0 ) && ((M().x86.intno == 0) || M().x86.intno == 2)) || !ACCESS_FLAG(F_IF) {
				x86emu_intr_handle()
			}
		}
		ip := M().x86.spc.IP.Get16()
		op1 = sys_rdb((uint32(M().x86.seg.CS.Get()) << 4 + uint32(ip)))
		M().x86.spc.IP.Set16(ip + 1)
		x86emu_optab[op1](op1)
		//if (M().x86.debug & DEBUG_EXIT) {
		//    M().x86.debug &= ~DEBUG_EXIT;
		//    return;
		//}
	}
}

/****************************************************************************
REMARKS:
Halts the system by setting the halted system flag.
****************************************************************************/
func X86EMU_halt_sys() {
	M().x86.intr |= INTR_HALTED
}

/****************************************************************************
PARAMETERS:
mod     - Mod value from decoded byte
regh    - Reg h value from decoded byte
regl    - Reg l value from decoded byte

REMARKS:
Raise the specified interrupt to be handled before the execution of the
next instruction.

NOTE: Do not inline this function, as sys_rdb is already inline!
****************************************************************************/
func fetch_decode_modrm(mod *int, regh *int, regl *int) {
	var fetched byte

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}

	ip := M().x86.spc.IP.Get16()
	fetched = sys_rdb(uint32(M().x86.seg.CS.Get()) << 4 + uint32(ip))
	M().x86.spc.IP.Set16(ip + 1)
	INC_DECODED_INST_LEN(1)
	*mod = int((fetched >> 6) & 0x03)
	*regh = int((fetched >> 3) & 0x07)
	*regl = int((fetched >> 0) & 0x07)
}

/****************************************************************************
RETURNS:
Immediate byte value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as sys_rdb is already inline!
****************************************************************************/
func fetch_byte_imm() uint32 {
	var fetched uint8

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}

	ip := M().x86.spc.IP.Get16()
	fetched = sys_rdb((uint32(M().x86.seg.CS.Get()) << 4) + uint32(ip))
	M().x86.spc.IP.Set16(ip + 1)
	INC_DECODED_INST_LEN(1)
	return uint32(fetched)
}

/****************************************************************************
RETURNS:
Immediate word value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as sys_rdw is already inline!
****************************************************************************/
func fetch_word_imm() uint16 {
	var fetched uint16

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}
	ip := M().x86.spc.IP.Get16()
	fetched = sys_rdw((uint32(M().x86.seg.CS.Get()) << 4) + uint32(ip))
	M().x86.spc.IP.Set16(ip+2)
	INC_DECODED_INST_LEN(2)
	return fetched
}

/****************************************************************************
RETURNS:
Immediate lone value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as sys_rdw is already inline!
****************************************************************************/
func fetch_long_imm() uint32 {
	var fetched uint32

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}
	ip := M().x86.spc.IP.Get16()
	fetched = sys_rdl((uint32(M().x86.seg.CS.Get()) << 4) + uint32(ip))
	M().x86.spc.IP.Set16(ip+4)
	INC_DECODED_INST_LEN(4)
	return fetched
}

/****************************************************************************
RETURNS:
Value of the default data segment

REMARKS:
Inline function that returns the default data segment for the current
instruction.

On the x86 processor, the default segment is not always DS if there is
no segment override. Address modes such as -3[BP] or 10[BP+SI] all refer to
addresses relative to SS (ie: on the stack). So, at the minimum, all
decodings of addressing modes would have to set/clear a bit describing
whether the access is relative to DS or SS.  That is the function of the
cpu-state-variable M().x86.mode. There are several potential states:

    repe prefix seen  (handled elsewhere)
    repne prefix seen  (ditto)

    cs segment override
    ds segment override
    es segment override
    fs segment override
    gs segment override
    ss segment override

    ds/ss select (in absence of override)

Each of the above 7 items are handled with a bit in the mode field.
****************************************************************************/
func get_data_segment() uint16 {

	switch M().x86.mode & SYSMODE_SEGMASK {
	case 0: /* default case: use ds register */
	case SYSMODE_SEGOVR_DS:
	case SYSMODE_SEGOVR_DS | SYSMODE_SEG_DS_SS:
		return M().x86.seg.DS.Get()
	case SYSMODE_SEG_DS_SS: /* non-overridden, use ss register */
		return M().x86.seg.SS.Get()
	case SYSMODE_SEGOVR_CS:
	case SYSMODE_SEGOVR_CS | SYSMODE_SEG_DS_SS:
		return M().x86.seg.CS.Get()
	case SYSMODE_SEGOVR_ES:
	case SYSMODE_SEGOVR_ES | SYSMODE_SEG_DS_SS:
		return M().x86.seg.ES.Get()
	case SYSMODE_SEGOVR_FS:
	case SYSMODE_SEGOVR_FS | SYSMODE_SEG_DS_SS:
		return M().x86.seg.FS.Get()
	case SYSMODE_SEGOVR_GS:
	case SYSMODE_SEGOVR_GS | SYSMODE_SEG_DS_SS:
		return M().x86.seg.GS.Get()
	case SYSMODE_SEGOVR_SS:
	case SYSMODE_SEGOVR_SS | SYSMODE_SEG_DS_SS:
		return M().x86.seg.SS.Get()
	default:

		HALT_SYS()
		return 0
	}
	return 0
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Byte value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_byte(offset uint) uint8 {

	panic("fix me")

	return 0 // sys_rdb((get_data_segment() << 4) + offset);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Word value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_word(offset uint) uint16 {

	panic("fix me")
	return 0 // return sys_rdw((get_data_segment() << 4) + offset);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Long value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_long(offset uint16) uint32 {

	return sys_rdl(uint32(get_data_segment() << 4) + uint32(offset))
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Byte value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_byte_abs(segment uint16, offset uint16) uint8 {

	return sys_rdb(uint32(segment << 4 + offset))
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Word value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_word_abs(segment uint16, offset uint16) uint16 {

	panic("fix me")
	return 0
	//return sys_rdw((uint32(segment) << 4) + offset);
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Long value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_long_abs(segment uint, offset uint) uint32 {

	panic("fix me")

	return 0 //sys_rdl((uint32(segment) << 4) + offset);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a word value to an segmented memory location. The segment used is
the current 'default' segment, which may have been overridden.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_byte(offset uint, val uint8) byte {

	panic("fix me")
	return 0 //(*sys_wrb)((get_data_segment() << 4) + offset, val);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a word value to an segmented memory location. The segment used is
the current 'default' segment, which may have been overridden.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_word(offset uint, val uint16) uint16 {

	panic("fix me")
	return 0 // (*sys_wrw)((get_data_segment() << 4) + offset, val);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a long value to an segmented memory location. The segment used is
the current 'default' segment, which may have been overridden.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_long(offset uint, val uint32) uint32 {

	panic("fix me")

	return 0 // (*sys_wrl)((get_data_segment() << 4) + offset, val);
}

/****************************************************************************
PARAMETERS:
segment - Segment to store data at
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a byte value to an absolute memory location.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_byte_abs(segment uint, offset uint, val uint8) {

	panic("fix me")
}

/****************************************************************************
PARAMETERS:
segment - Segment to store data at
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a word value to an absolute memory location.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_word_abs(segment uint, offset uint, val uint16) {

	panic("fix me")

	//(*sys_wrw)((uint32(segment)<<4)+offset, val)
}

/****************************************************************************
PARAMETERS:
segment - Segment to store data at
offset  - Offset to store data at
val     - Value to store

REMARKS:
Writes a long value to an absolute memory location.

NOTE: Do not inline this function as (*sys_wrX) is already inline!
****************************************************************************/
func store_data_long_abs(segment uint, offset uint, val uint32) {

	panic("fix me")

	//(*sys_wrl)((uint32(segment) << 4) + offset, val);
}

/****************************************************************************
PARAMETERS:
reg - Register to decode

RETURNS:
Pointer to the appropriate register

REMARKS:
Return a pointer to the register given by the R/RM field of the
modrm byte, for byte operands. Also enables the decoding of instructions.
****************************************************************************/
func decode_rm_byte_register(reg int) register {
	switch reg {
	case 0:
		DECODE_PRINTF("AL")
		return M().x86.gen.A
	case 1:
		DECODE_PRINTF("CL")
		return M().x86.gen.C
	case 2:
		DECODE_PRINTF("DL")
		return M().x86.gen.D
	case 3:
		DECODE_PRINTF("BL")
		return M().x86.gen.B
	case 4:
		DECODE_PRINTF("AH")
		return M().x86.gen.A
	case 5:
		DECODE_PRINTF("CH")
		return M().x86.gen.C
	case 6:
		DECODE_PRINTF("DH")
		return M().x86.gen.D
	case 7:
		DECODE_PRINTF("BH")
		return M().x86.gen.B
	}
	HALT_SYS()
	return nil /* NOT REACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
reg - Register to decode

RETURNS:
Pointer to the appropriate register

REMARKS:
Return a pointer to the register given by the R/RM field of the
modrm byte, for word operands.  Also enables the decoding of instructions.
****************************************************************************/
func decode_rm_word_register(reg int) register {
	switch reg {
	case 0:
		DECODE_PRINTF("AX")
		return &M().x86.gen.A
	case 1:
		DECODE_PRINTF("CX")
		return &M().x86.gen.C
	case 2:
		DECODE_PRINTF("DX")
		return &M().x86.gen.D
	case 3:
		DECODE_PRINTF("BX")
		return &M().x86.gen.B
	case 4:
		DECODE_PRINTF("SP")
		return &M().x86.spc.SP
	case 5:
		DECODE_PRINTF("BP")
		return &M().x86.spc.BP
	case 6:
		DECODE_PRINTF("SI")
		return &M().x86.spc.SI
	case 7:
		DECODE_PRINTF("DI")
		return &M().x86.spc.DI
	}
	HALT_SYS()
	return nil /* NOTREACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
reg - Register to decode

RETURNS:
Pointer to the appropriate register

REMARKS:
Return a pointer to the register given by the R/RM field of the
modrm byte, for dword operands.  Also enables the decoding of instructions.
****************************************************************************/
func decode_rm_long_register(reg int) register {
	switch reg {
	case 0:
		DECODE_PRINTF("EAX")
		return &M().x86.gen.A
	case 1:
		DECODE_PRINTF("ECX")
		return &M().x86.gen.C
	case 2:
		DECODE_PRINTF("EDX")
		return &M().x86.gen.D
	case 3:
		DECODE_PRINTF("EBX")
		return &M().x86.gen.B
	case 4:
		DECODE_PRINTF("ESP")
		return &M().x86.spc.SP
	case 5:
		DECODE_PRINTF("EBP")
		return &M().x86.spc.BP
	case 6:
		DECODE_PRINTF("ESI")
		return &M().x86.spc.SI
	case 7:
		DECODE_PRINTF("EDI")
		return &M().x86.spc.DI
	}
	HALT_SYS()
	return nil /* NOTREACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
reg - Register to decode

RETURNS:
Pointer to the appropriate register

REMARKS:
Return a pointer to the register given by the R/RM field of the
modrm byte, for word operands, modified from above for the weirdo
special case of segreg operands.  Also enables the decoding of instructions.
****************************************************************************/
func decode_rm_seg_register(reg int) register16 {
	switch reg {
	case 0:
		DECODE_PRINTF("ES")
		return M().x86.seg.ES
	case 1:
		DECODE_PRINTF("CS")
		return M().x86.seg.CS
	case 2:
		DECODE_PRINTF("SS")
		return M().x86.seg.SS
	case 3:
		DECODE_PRINTF("DS")
		return M().x86.seg.DS
	case 4:
		DECODE_PRINTF("FS")
		return M().x86.seg.FS
	case 5:
		DECODE_PRINTF("GS")
		return M().x86.seg.GS
	case 6:
	case 7:
		DECODE_PRINTF("ILLEGAL SEGREG")
		break
	}
	HALT_SYS()
	return nil /* NOT REACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
scale - scale value of SIB byte
index - index value of SIB byte

RETURNS:
Value of scale * index

REMARKS:
Decodes scale/index of SIB byte and returns relevant offset part of
effective address.
****************************************************************************/
func decode_sib_si(scale uint32, index uint32) uint32 {
	scale = 1 << scale
	if scale > 1 {
		DECODE_PRINTF2("[%d*", scale)
	} else {
		DECODE_PRINTF("[")
	}
	switch index {
	case 0:
		DECODE_PRINTF("EAX]")
		return M().x86.gen.A.Get32() * index
	case 1:
		DECODE_PRINTF("ECX]")
		return M().x86.gen.C.Get32() * index
	case 2:
		DECODE_PRINTF("EDX]")
		return M().x86.gen.D.Get32() * index
	case 3:
		DECODE_PRINTF("EBX]")
		return M().x86.gen.B.Get32() * index
	case 4:
		DECODE_PRINTF("0]")
		return 0
	case 5:
		DECODE_PRINTF("EBP]")
		return M().x86.spc.BP.Get32() * index
	case 6:
		DECODE_PRINTF("ESI]")
		return M().x86.spc.SI.Get32() * index
	case 7:
		DECODE_PRINTF("EDI]")
		return M().x86.spc.DI.Get32() * index
	}
	HALT_SYS()
	return 0 /* NOT REACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
mod - MOD value of preceding ModR/M byte

RETURNS:
Offset in memory for the address decoding

REMARKS:
Decodes SIB addressing byte and returns calculated effective address.
****************************************************************************/
func decode_sib_address(mod uint32) uint32 {
	var sib = uint32(fetch_byte_imm())
	var ss = uint32((sib >> 6) & 0x03)
	var index = uint32((sib >> 3) & 0x07)
	var base = uint32(sib & 0x07)
	var offset uint32
	var displacement uint32

	switch base {
	case 0:
		DECODE_PRINTF("[EAX]")
		offset = M().x86.gen.A.Get32()
		break
	case 1:
		DECODE_PRINTF("[ECX]")
		offset = M().x86.gen.C.Get32()
		break
	case 2:
		DECODE_PRINTF("[EDX]")
		offset = M().x86.gen.D.Get32()
		break
	case 3:
		DECODE_PRINTF("[EBX]")
		offset = M().x86.gen.B.Get32()
		break
	case 4:
		DECODE_PRINTF("[ESP]")
		offset = M().x86.spc.SP.Get32()
		break
	case 5:
		switch mod {
		case 0:
			displacement = fetch_long_imm()
			DECODE_PRINTF2("[%d]", displacement)
			offset = displacement
			break
		case 1:
			displacement = fetch_byte_imm()
			DECODE_PRINTF2("[%d][EBP]", displacement)
			offset = M().x86.spc.BP.Get32() + displacement
			break
		case 2:
			displacement = fetch_long_imm()
			DECODE_PRINTF2("[%d][EBP]", displacement)
			offset = M().x86.spc.BP.Get32() + displacement
			break
		default:
			HALT_SYS()
		}
		DECODE_PRINTF("[EAX]")
		offset = M().x86.gen.A.Get32()
		break
	case 6:
		DECODE_PRINTF("[ESI]")
		offset = M().x86.spc.SI.Get32()
		break
	case 7:
		DECODE_PRINTF("[EDI]")
		offset = M().x86.spc.DI.Get32()
		break
	default:
		HALT_SYS()
	}
	offset += decode_sib_si(ss, index)
	return offset
}

/****************************************************************************
PARAMETERS:
rm  - RM value to decode

RETURNS:
Offset in memory for the address decoding

REMARKS:
Return the offset given by mod=00 addressing.  Also enables the
decoding of instructions.

NOTE:   The code which specifies the corresponding segment (ds vs ss)
        below in the case of [BP+..].  The assumption here is that at the
        point that this subroutine is called, the bit corresponding to
        SYSMODE_SEG_DS_SS will be zero.  After every instruction
        except the segment override instructions, this bit (as well
        as any bits indicating segment overrides) will be clear.  So
        if a SS access is needed, set this bit.  Otherwise, DS access
        occurs (unless any of the segment override bits are set).
****************************************************************************/
func decode_rm00_address(rm uint32) uint32 {
	var offset uint32

	if M().x86.mode & SYSMODE_PREFIX_ADDR != 0 {
		/* 32-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF("[EAX]")
			return M().x86.gen.A.Get32()
		case 1:
			DECODE_PRINTF("[ECX]")
			return M().x86.gen.C.Get32()
		case 2:
			DECODE_PRINTF("[EDX]")
			return M().x86.gen.D.Get32()
		case 3:
			DECODE_PRINTF("[EBX]")
			return M().x86.gen.B.Get32()
		case 4:
			return decode_sib_address(0)
		case 5:
			offset = fetch_long_imm()
			DECODE_PRINTF2("[%08x]", offset)
			return offset
		case 6:
			DECODE_PRINTF("[ESI]")
			return M().x86.spc.SI.Get32()
		case 7:
			DECODE_PRINTF("[EDI]")
			return M().x86.spc.DI.Get32()
		}
	} else {
		/* 16-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF("[BX+SI]")
			return uint32(M().x86.gen.B.Get16() + M().x86.spc.SI.Get16())
		case 1:
			DECODE_PRINTF("[BX+DI]")
			return uint32(M().x86.gen.B.Get16() + M().x86.spc.DI.Get16())
		case 2:
			DECODE_PRINTF("[BP+SI]")
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32(M().x86.spc.BP.Get16() + M().x86.spc.SI.Get16())
		case 3:
			DECODE_PRINTF("[BP+DI]")
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32(M().x86.spc.BP.Get16() + M().x86.spc.DI.Get16())
		case 4:
			DECODE_PRINTF("[SI]")
			return uint32(M().x86.spc.SI.Get16())
		case 5:
			DECODE_PRINTF("[DI]")
			return uint32(M().x86.spc.DI.Get16())
		case 6:
			offset = uint32(fetch_word_imm())
			DECODE_PRINTF2("[%04x]", offset)
			return offset
		case 7:
			DECODE_PRINTF("[BX]")
			return uint32(M().x86.gen.B.Get16())
		}
	}
	HALT_SYS()
	return 0
}

/****************************************************************************
PARAMETERS:
rm  - RM value to decode

RETURNS:
Offset in memory for the address decoding

REMARKS:
Return the offset given by mod=01 addressing.  Also enables the
decoding of instructions.
****************************************************************************/
func decode_rm01_address(rm uint32) uint32 {
	var displacement uint32

	if M().x86.mode & SYSMODE_PREFIX_ADDR != 0{
		/* 32-bit addressing */
		if rm != 4 {
			displacement = uint32(fetch_byte_imm())
		}

		switch rm {
		case 0:
			DECODE_PRINTF2("%d[EAX]", displacement)
			return M().x86.gen.A.Get32() + displacement
		case 1:
			DECODE_PRINTF2("%d[ECX]", displacement)
			return M().x86.gen.C.Get32() + displacement
		case 2:
			DECODE_PRINTF2("%d[EDX]", displacement)
			return M().x86.gen.D.Get32() + displacement
		case 3:
			DECODE_PRINTF2("%d[EBX]", displacement)
			return M().x86.gen.B.Get32() + displacement
		case 4:
			{
				var offset = uint32(decode_sib_address(1))
				displacement = uint32(fetch_byte_imm())
				DECODE_PRINTF2("[%d]", displacement)
				return offset + displacement
			}
		case 5:
			DECODE_PRINTF2("%d[EBP]", displacement)
			return M().x86.spc.BP.Get32() + displacement
		case 6:
			DECODE_PRINTF2("%d[ESI]", displacement)
			return M().x86.spc.SI.Get32() + displacement
		case 7:
			DECODE_PRINTF2("%d[EDI]", displacement)
			return M().x86.spc.DI.Get32() + displacement
		}
	} else {
		/* 16-bit addressing */
		d16 := uint16(fetch_byte_imm())
		switch rm {
		case 0:
			DECODE_PRINTF2("%d[BX+SI]", d16)
			return uint32((M().x86.gen.B.Get16() + M().x86.spc.SI.Get16() + d16))
		case 1:
			DECODE_PRINTF2("%d[BX+DI]", d16)
			return uint32((M().x86.gen.B.Get16() + M().x86.spc.DI.Get16() + d16))
		case 2:
			DECODE_PRINTF2("%d[BP+SI]", d16)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + M().x86.spc.SI.Get16() + d16))
		case 3:
			DECODE_PRINTF2("%d[BP+DI]", d16)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + M().x86.spc.DI.Get16() + d16))
		case 4:
			DECODE_PRINTF2("%d[SI]", d16)
			return uint32((M().x86.spc.SI.Get16() + d16))
		case 5:
			DECODE_PRINTF2("%d[DI]", d16)
			return uint32((M().x86.spc.DI.Get16() + d16))
		case 6:
			DECODE_PRINTF2("%d[BP]", d16)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + d16))
		case 7:
			DECODE_PRINTF2("%d[BX]", d16)
			return uint32((M().x86.gen.B.Get16() + d16))
		}
	}
	HALT_SYS()
	return 0 /* SHOULD NOT HAPPEN */
}

/****************************************************************************
PARAMETERS:
rm  - RM value to decode

RETURNS:
Offset in memory for the address decoding

REMARKS:
Return the offset given by mod=10 addressing.  Also enables the
decoding of instructions.
****************************************************************************/
func decode_rm10_address(rm uint32) uint32 {
	if M().x86.mode & SYSMODE_PREFIX_ADDR != 0{
		var displacement uint32

		/* 32-bit addressing */
		if rm != 4 {
			displacement = fetch_long_imm()
		}

		switch rm {
		case 0:
			DECODE_PRINTF2("%d[EAX]", displacement)
			return M().x86.gen.A.Get32() + displacement
		case 1:
			DECODE_PRINTF2("%d[ECX]", displacement)
			return M().x86.gen.C.Get32() + displacement
		case 2:
			DECODE_PRINTF2("%d[EDX]", displacement)
			return M().x86.gen.D.Get32() + displacement
		case 3:
			DECODE_PRINTF2("%d[EBX]", displacement)
			return M().x86.gen.B.Get32() + displacement
		case 4:
			{
				var offset = decode_sib_address(2)
				displacement = fetch_long_imm()
				DECODE_PRINTF2("[%d]", displacement)
				return offset + displacement
			}
		case 5:
			DECODE_PRINTF2("%d[EBP]", displacement)
			return M().x86.spc.BP.Get32() + displacement
		case 6:
			DECODE_PRINTF2("%d[ESI]", displacement)
			return M().x86.spc.SI.Get32() + displacement
		case 7:
			DECODE_PRINTF2("%d[EDI]", displacement)
			return M().x86.spc.DI.Get32() + displacement
		}
	} else {
		var displacement = uint16(fetch_word_imm())

		/* 16-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF2("%d[BX+SI]", displacement)
			return uint32((M().x86.gen.B.Get16() + M().x86.spc.SI.Get16() + displacement))
		case 1:
			DECODE_PRINTF2("%d[BX+DI]", displacement)
			return uint32((M().x86.gen.B.Get16() + M().x86.spc.DI.Get16() + displacement))
		case 2:
			DECODE_PRINTF2("%d[BP+SI]", displacement)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + M().x86.spc.SI.Get16() + displacement))
		case 3:
			DECODE_PRINTF2("%d[BP+DI]", displacement)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + M().x86.spc.DI.Get16() + displacement))
		case 4:
			DECODE_PRINTF2("%d[SI]", displacement)
			return uint32((M().x86.spc.SI.Get16() + displacement))
		case 5:
			DECODE_PRINTF2("%d[DI]", displacement)
			return uint32((M().x86.spc.DI.Get16() + displacement))
		case 6:
			DECODE_PRINTF2("%d[BP]", displacement)
			M().x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((M().x86.spc.BP.Get16() + displacement))
		case 7:
			DECODE_PRINTF2("%d[BX]", displacement)
			return uint32((M().x86.gen.B.Get16() + displacement))
		}
	}
	HALT_SYS()
	return 0 /* SHOULD NOT HAPPEN */
}

/****************************************************************************
PARAMETERS:
mod - modifier
rm  - RM value to decode

RETURNS:
Offset in memory for the address decoding, multiplexing calls to
the decode_rmXX_address functions

REMARKS:
Return the offset given by "mod" addressing.
****************************************************************************/

func decode_rmXX_address(mod uint32, rm uint32) uint32 {
	if mod == 0 {
		return decode_rm00_address(rm)
	}
	if mod == 1 {
		return decode_rm01_address(rm)
	}
	return decode_rm10_address(rm)
}

func push_word(w uint16) {
	panic("fix me")
}

func mem_access_word(addr uint32) uint16 {
	panic("fix me")
}

func INC_DECODED_INST_LEN(amt uint16) {
	if (DEBUG_DECODE())  	       {
		x86emu_inc_decoded_inst_len(uint32(amt))
	}
}
func sys_rdw(add uint32) uint16 {
	panic("fix me")
}
func sys_rdl(add uint32) uint32 {
	panic("fix me")
}

func HALT_SYS() {
	os.Exit(0)
}

func DECODE_PRINTF(x string, y...interface{}) {
     	if (DEBUG_DECODE()) {
		x86emu_decode_printf(x, y...)
	}
}

func DECODE_PRINTF2(x string, y...interface{}) {
     	if (DEBUG_DECODE()) {
		x86emu_decode_printf(x, y...)
	}
}

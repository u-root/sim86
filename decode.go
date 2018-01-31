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
	"log"
	"os"
)

/*----------------------------- Implementation ----------------------------*/

/****************************************************************************
REMARKS:
Handles any pending asynchronous interrupts.
****************************************************************************/
func x86emu_intr_handle() {
	var intno uint8

	if M.x86.intr&INTR_SYNCH != 0 {
		intno = uint8(M.x86.intno)
		if _X86EMU_intrTab[intno] != nil {
			panic("_X86EMU_intrTab[intno](intno)")
		} else {
			push_word(uint16(G(FLAGS)))
			CLEAR_FLAG(F_IF)
			CLEAR_FLAG(F_TF)
			push_word(G16(CS))
			S16(CS, mem_access_word(uint32(intno*4 + 2)))
			push_word(G16(IP))
			S16(IP, mem_access_word(uint32(intno * 4)))
			M.x86.intr = 0
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
	M.x86.intno = intrnum
	M.x86.intr |= INTR_SYNCH
}

/****************************************************************************
REMARKS:
Main execution loop for the emulator. We return from here when the system
halts, which is normally caused by a stack fault when we return from the
original real mode call.
****************************************************************************/
func X86EMU_exec() {
	var op1 uint8

	M.x86.intr = 0
	x86emu_end_instr()

	for {
		if CHECK_IP_FETCH() {
			x86emu_check_ip_access()
		}

		/* If debugging, save the IP and CS values. */
		SAVE_IP_CS(G16(CS), G16(IP))
		INC_DECODED_INST_LEN(1)
		if M.x86.intr != 0 {
			if (M.x86.intr & INTR_HALTED) != 0 {
				if G16(SP) != 0 {
					fmt.Printf("halted\n")
					X86EMU_trace_regs()
				} else {
					if M.x86.debug != 0 {
						fmt.Printf("Service completed successfully\n")
					}
				}
				return
			}
			if ((M.x86.intr&INTR_SYNCH != 0) && ((M.x86.intno == 0) || M.x86.intno == 2)) || !ACCESS_FLAG(F_IF) {
				x86emu_intr_handle()
			}
		}
		ip := G16(IP)
		op1 = sys_rdb((uint32(G16(CS))<<4 + uint32(ip)))
		fmt.Printf("Set ip to %d\n", ip)
		S16(IP, ip + 1)
		ip = G16(IP)
		fmt.Printf("Set ip to %d\n", ip)
		x86emu_dump_regs()
		x86emu_optab[op1](op1)
		if M.x86.exit {
			M.x86.exit = false
			return
		}
		x86emu_dump_regs()
	}
}

/****************************************************************************
REMARKS:
Halts the system by setting the halted system flag.
****************************************************************************/
func X86EMU_halt_sys() {
	M.x86.intr |= INTR_HALTED
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
func fetch_decode_modrm() (uint32, uint32, uint32) {
	var fetched byte

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}

	ip := G16(IP)
	fetched = sys_rdb(uint32(G16(CS))<<4 + uint32(ip))
	S16(IP, ip + 1)
	INC_DECODED_INST_LEN(1)
	mod := uint32((fetched >> 6) & 0x03)
	regh := uint32((fetched >> 3) & 0x07)
	regl := uint32((fetched >> 0) & 0x07)
	return mod, regh, regl
}

/****************************************************************************
RETURNS:
Immediate byte value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as sys_rdb is already inline!
****************************************************************************/
func fetch_byte_imm() uint8 {
	var fetched uint8

	if CHECK_IP_FETCH() {
		x86emu_check_ip_access()
	}

	ip := G16(IP)
	fetched = sys_rdb((uint32(G16(CS)) << 4) + uint32(ip))
	S16(IP, ip + 1)
	INC_DECODED_INST_LEN(1)
	return fetched
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
	ip := G16(IP)
	fetched = sys_rdw((uint32(G16(CS)) << 4) + uint32(ip))
	S16(IP, ip + 2)
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
	ip := G16(IP)
	fetched = sys_rdl((uint32(G16(CS)) << 4) + uint32(ip))
	S16(IP, ip + 4)
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
cpu-state-variable M.x86.mode. There are several potential states:

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

	switch M.x86.mode & SYSMODE_SEGMASK {
	case 0: /* default case: use ds register */
	case SYSMODE_SEGOVR_DS:
	case SYSMODE_SEGOVR_DS | SYSMODE_SEG_DS_SS:
		return G16(DS)
	case SYSMODE_SEG_DS_SS: /* non-overridden, use ss register */
		return G16(SS)
	case SYSMODE_SEGOVR_CS:
	case SYSMODE_SEGOVR_CS | SYSMODE_SEG_DS_SS:
		return G16(CS)
	case SYSMODE_SEGOVR_ES:
	case SYSMODE_SEGOVR_ES | SYSMODE_SEG_DS_SS:
		return G16(ES)
	case SYSMODE_SEGOVR_FS:
	case SYSMODE_SEGOVR_FS | SYSMODE_SEG_DS_SS:
		return G16(FS)
	case SYSMODE_SEGOVR_GS:
	case SYSMODE_SEGOVR_GS | SYSMODE_SEG_DS_SS:
		return G16(GS)
	case SYSMODE_SEGOVR_SS:
	case SYSMODE_SEGOVR_SS | SYSMODE_SEG_DS_SS:
		return G16(SS)
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
func fetch_data_byte(offset uint32) uint8 {

	return sys_rdb(uint32(get_data_segment()<<4) + offset)
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Word value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_word(offset uint32) uint16 {
	return sys_rdw(uint32(get_data_segment()<<4) + uint32(offset))
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Long value read from the absolute memory location.

NOTE: Do not inline this function as sys_rdX is already inline!
****************************************************************************/
func fetch_data_long(offset uint32) uint32 {

	return sys_rdl(uint32(get_data_segment()<<4) + offset)
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
	return sys_rdb(uint32(segment<<4 + offset))
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

	panic("fix mdecoe")
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
func fetch_data_long_abs(segment uint16, offset uint16) uint32 {
	var i uint32
	sysr((uint32(segment)<<4)+uint32(offset), &i)
	return i
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
func store_data_byte(offset uint32, val uint8) {
	sysw(uint32(get_data_segment()<<4)+offset, val)
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
func store_data_word(offset uint32, val uint16) {
	sysw(uint32(get_data_segment()<<4)+offset, val)
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
func store_data_long(offset uint32, val uint32) {
	sysw(uint32(get_data_segment()<<4)+offset, val)
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
func store_data_byte_abs(segment uint16, offset uint16, val uint8) {
	sysw((uint32(segment)<<4 + uint32(offset)), val)
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
func store_data_word_abs(segment uint16, offset uint16, val uint16) {
	sysw((uint32(segment)<<4 + uint32(offset)), val)
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
func store_data_long_abs(segment uint16, offset uint16, val uint32) {

	sysw((uint32(segment)<<4)+uint32(offset), val)

}
type regmap struct {
	n string
	r regtype
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
func decode_rm_byte_register(reg uint32) regtype {
	switch reg {
	case 0:
		DECODE_PRINTF("AL")
		return AL
	case 1:
		DECODE_PRINTF("CL")
		return CL
	case 2:
		DECODE_PRINTF("DL")
		return DL
	case 3:
		DECODE_PRINTF("BL")
		return BL
	case 4:
		DECODE_PRINTF("AH")
		return AH
	case 5:
		DECODE_PRINTF("CH")
		return CH
	case 6:
		DECODE_PRINTF("DH")
		return DH
	case 7:
		DECODE_PRINTF("BH")
		return BH
	}
	HALT_SYS()
	log.Panicf("bad register in decode_rm_byte_register: %04x", reg)
	return 0
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
func decode_rm_word_register(reg uint32) regtype {
	var rmap = [...]regmap {
		{"AX", AX},
		{"CX", CX},
		{"DX", DX},
		{"BX", BX},
		{"SP", SP},
		{"BP", BP},
		{"SI", SI},
		{"DI", DI},
	}
	if int(reg) < len(rmap) {
		DECODE_PRINTF(rmap[reg].n)
		return rmap[reg].r
	}
	HALT_SYS()
	log.Panicf("decode_rm_word_register bad register %d\n", reg)
	return 0
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
func decode_rm_long_register(reg uint32) regtype {
	var rmap = [...]regmap {
		{"EAX", EAX},
		{"ECX", ECX},
		{"EDX", EDX},
		{"EBX", EBX},
		{"ESP", ESP},
		{"EBP", EBP},
		{"ESI", ESI},
		{"EDI", EDI},
	}
	if int(reg) < len(rmap) {
		DECODE_PRINTF(rmap[reg].n)
		return rmap[reg].r
	}
	HALT_SYS()
	log.Panicf("decode_rm_long_register bad register %d\n", reg)
	return 0
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
func decode_rm_seg_register(reg uint32) regtype {
	var rmap = [...]regmap {
		{"ES", ES},
		{"CS", CS},
		{"SS", SS},
		{"DS", DS},
		{"FS", FS},
		{"GS", GS},
	}
	if int(reg) < len(rmap) {
		DECODE_PRINTF(rmap[reg].n)
		return rmap[reg].r
	}
	HALT_SYS()
	log.Panicf("decode_rm_seg_register bad register %d\n", reg)
	return 0
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
		return G32(EAX) * index
	case 1:
		DECODE_PRINTF("ECX]")
		return G32(ECX) * index
	case 2:
		DECODE_PRINTF("EDX]")
		return G32(EDX) * index
	case 3:
		DECODE_PRINTF("EBX]")
		return G32(EBX) * index
	case 4:
		DECODE_PRINTF("0]")
		return 0
	case 5:
		DECODE_PRINTF("EBP]")
		return G(BP) * index
	case 6:
		DECODE_PRINTF("ESI]")
		return G(SI) * index
	case 7:
		DECODE_PRINTF("EDI]")
		return G(DI) * index
	}
	HALT_SYS()
	log.Panicf("decode_sib_si bad index register %d\n", index)
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
		offset = G32(EAX)
		break
	case 1:
		DECODE_PRINTF("[ECX]")
		offset = G32(ECX)
		break
	case 2:
		DECODE_PRINTF("[EDX]")
		offset = G32(EDX)
		break
	case 3:
		DECODE_PRINTF("[EBX]")
		offset = G32(EBX)
		break
	case 4:
		DECODE_PRINTF("[ESP]")
		offset = G(SP)
		break
	case 5:
		switch mod {
		case 0:
			displacement = fetch_long_imm()
			DECODE_PRINTF2("[%d]", displacement)
			offset = displacement
			break
		case 1:
			displacement = uint32(fetch_byte_imm())
			DECODE_PRINTF2("[%d][EBP]", displacement)
			offset = G(BP) + displacement
			break
		case 2:
			displacement = fetch_long_imm()
			DECODE_PRINTF2("[%d][EBP]", displacement)
			offset = G(BP) + displacement
			break
		default:
			HALT_SYS()
		}
		DECODE_PRINTF("[EAX]")
		offset = G32(EAX)
		break
	case 6:
		DECODE_PRINTF("[ESI]")
		offset = G(SI)
		break
	case 7:
		DECODE_PRINTF("[EDI]")
		offset = G(DI)
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

	if M.x86.mode&SYSMODE_PREFIX_ADDR != 0 {
		/* 32-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF("[EAX]")
			return G32(EAX)
		case 1:
			DECODE_PRINTF("[ECX]")
			return G32(ECX)
		case 2:
			DECODE_PRINTF("[EDX]")
			return G32(EDX)
		case 3:
			DECODE_PRINTF("[EBX]")
			return G32(EBX)
		case 4:
			return decode_sib_address(0)
		case 5:
			offset = fetch_long_imm()
			DECODE_PRINTF2("[%08x]", offset)
			return offset
		case 6:
			DECODE_PRINTF("[ESI]")
			return G(SI)
		case 7:
			DECODE_PRINTF("[EDI]")
			return G(DI)
		}
	} else {
		/* 16-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF("[BX+SI]")
			return uint32(G16(BX) + G16(SI))
		case 1:
			DECODE_PRINTF("[BX+DI]")
			return uint32(G16(BX) + G16(DI))
		case 2:
			DECODE_PRINTF("[BP+SI]")
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32(G16(BP) + G16(SI))
		case 3:
			DECODE_PRINTF("[BP+DI]")
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32(G16(BP) + G16(DI))
		case 4:
			DECODE_PRINTF("[SI]")
			return uint32(G16(SI))
		case 5:
			DECODE_PRINTF("[DI]")
			return uint32(G16(DI))
		case 6:
			offset = uint32(fetch_word_imm())
			DECODE_PRINTF2("[%04x]", offset)
			return offset
		case 7:
			DECODE_PRINTF("[BX]")
			return uint32(G16(BX))
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

	if M.x86.mode&SYSMODE_PREFIX_ADDR != 0 {
		/* 32-bit addressing */
		if rm != 4 {
			displacement = uint32(fetch_byte_imm())
		}

		switch rm {
		case 0:
			DECODE_PRINTF2("%d[EAX]", displacement)
			return G32(EAX) + displacement
		case 1:
			DECODE_PRINTF2("%d[ECX]", displacement)
			return G32(ECX) + displacement
		case 2:
			DECODE_PRINTF2("%d[EDX]", displacement)
			return G32(EDX) + displacement
		case 3:
			DECODE_PRINTF2("%d[EBX]", displacement)
			return G32(EBX) + displacement
		case 4:
			{
				var offset = uint32(decode_sib_address(1))
				displacement = uint32(fetch_byte_imm())
				DECODE_PRINTF2("[%d]", displacement)
				return offset + displacement
			}
		case 5:
			DECODE_PRINTF2("%d[EBP]", displacement)
			return G(BP) + displacement
		case 6:
			DECODE_PRINTF2("%d[ESI]", displacement)
			return G(SI) + displacement
		case 7:
			DECODE_PRINTF2("%d[EDI]", displacement)
			return G(DI) + displacement
		}
	} else {
		/* 16-bit addressing */
		d16 := uint16(fetch_byte_imm())
		switch rm {
		case 0:
			DECODE_PRINTF2("%d[BX+SI]", d16)
			return uint32((G16(BX) + G16(SI) + d16))
		case 1:
			DECODE_PRINTF2("%d[BX+DI]", d16)
			return uint32((G16(BX) + G16(DI) + d16))
		case 2:
			DECODE_PRINTF2("%d[BP+SI]", d16)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + G16(SI) + d16))
		case 3:
			DECODE_PRINTF2("%d[BP+DI]", d16)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + G16(DI) + d16))
		case 4:
			DECODE_PRINTF2("%d[SI]", d16)
			return uint32((G16(SI) + d16))
		case 5:
			DECODE_PRINTF2("%d[DI]", d16)
			return uint32((G16(DI) + d16))
		case 6:
			DECODE_PRINTF2("%d[BP]", d16)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + d16))
		case 7:
			DECODE_PRINTF2("%d[BX]", d16)
			return uint32((G16(BX) + d16))
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
	if M.x86.mode&SYSMODE_PREFIX_ADDR != 0 {
		var displacement uint32

		/* 32-bit addressing */
		if rm != 4 {
			displacement = fetch_long_imm()
		}

		switch rm {
		case 0:
			DECODE_PRINTF2("%d[EAX]", displacement)
			return G32(EAX) + displacement
		case 1:
			DECODE_PRINTF2("%d[ECX]", displacement)
			return G32(ECX) + displacement
		case 2:
			DECODE_PRINTF2("%d[EDX]", displacement)
			return G32(EDX) + displacement
		case 3:
			DECODE_PRINTF2("%d[EBX]", displacement)
			return G32(EBX) + displacement
		case 4:
			{
				var offset = decode_sib_address(2)
				displacement = fetch_long_imm()
				DECODE_PRINTF2("[%d]", displacement)
				return offset + displacement
			}
		case 5:
			DECODE_PRINTF2("%d[EBP]", displacement)
			return G(BP) + displacement
		case 6:
			DECODE_PRINTF2("%d[ESI]", displacement)
			return G(SI) + displacement
		case 7:
			DECODE_PRINTF2("%d[EDI]", displacement)
			return G(DI) + displacement
		}
	} else {
		var displacement = uint16(fetch_word_imm())

		/* 16-bit addressing */
		switch rm {
		case 0:
			DECODE_PRINTF2("%d[BX+SI]", displacement)
			return uint32((G16(BX) + G16(SI) + displacement))
		case 1:
			DECODE_PRINTF2("%d[BX+DI]", displacement)
			return uint32((G16(BX) + G16(DI) + displacement))
		case 2:
			DECODE_PRINTF2("%d[BP+SI]", displacement)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + G16(SI) + displacement))
		case 3:
			DECODE_PRINTF2("%d[BP+DI]", displacement)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + G16(DI) + displacement))
		case 4:
			DECODE_PRINTF2("%d[SI]", displacement)
			return uint32((G16(SI) + displacement))
		case 5:
			DECODE_PRINTF2("%d[DI]", displacement)
			return uint32((G16(DI) + displacement))
		case 6:
			DECODE_PRINTF2("%d[BP]", displacement)
			M.x86.mode |= SYSMODE_SEG_DS_SS
			return uint32((G16(BP) + displacement))
		case 7:
			DECODE_PRINTF2("%d[BX]", displacement)
			return uint32((G16(BX) + displacement))
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

func INC_DECODED_INST_LEN(amt uint16) {
	if DEBUG_DECODE() {
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

func DECODE_PRINTF(x string, y ...interface{}) {
	if DEBUG_DECODE() {
		x86emu_decode_printf(x, y...)
	}
}

func DECODE_PRINTF2(x string, y ...interface{}) {
	if DEBUG_DECODE() {
		x86emu_decode_printf(x, y...)
	}
}

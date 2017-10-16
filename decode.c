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

#include "x86emui.h"

/*----------------------------- Implementation ----------------------------*/

/****************************************************************************
REMARKS:
Handles any pending asynchronous interrupts.
****************************************************************************/
static void x86emu_intr_handle(void)
{
    u8  intno;

    if (intr & INTR_SYNCH) {
        intno = intno;
        if (_X86EMU_intrTab[intno]) {
		loggy("intr %d\n", intno);
        } else {
            push_word((u16)FLG);
            CLEAR_FLAG(F_IF);
            CLEAR_FLAG(F_TF);
            push_word(CS);
            CS = mem_access_word(intno * 4 + 2);
            push_word(IP);
            IP = mem_access_word(intno * 4);
            intr = 0;
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
void x86emu_intr_raise(
    u8 intrnum)
{
    loggy("%s, raising exception %x\n", __func__, intrnum);
    x86emu_dump_regs();
    intno = intrnum;
    intr |= INTR_SYNCH;
}

/****************************************************************************
REMARKS:
Main execution loop for the emulator. We return from here when the system
halts, which is normally caused by a stack fault when we return from the
original real mode call.
****************************************************************************/
void X86EMU_exec(void)
{
    u8 op1;

    intr = 0;
    DB(x86emu_end_instr();)

    for (;;) {
DB(     if (CHECK_IP_FETCH())
            x86emu_check_ip_access();)
        /* If debugging, save the IP and CS values. */
        SAVE_IP_CS(CS, IP);
        INC_DECODED_INST_LEN(1);
        if (intr) {
            if (intr & INTR_HALTED) {
DB(             if (SP != 0) {
                    loggy("halted\n");
                    X86EMU_trace_regs();
                    }
                else {
                    if (debug)
                        loggy("Service completed successfully\n");
                    })
                return;
            }
            if (((intr & INTR_SYNCH) && (intno == 0 || intno == 2)) ||
                !ACCESS_FLAG(F_IF)) {
                x86emu_intr_handle();
            }
        }
        op1 = (*sys_rdb)(((u32)CS << 4) + (IP++));
	x86_byte_dispatch(op1);
        //if (debug & DEBUG_EXIT) {
        //    debug &= ~DEBUG_EXIT;
        //    return;
        //}
    }
}

/****************************************************************************
REMARKS:
Halts the system by setting the halted system flag.
****************************************************************************/
void X86EMU_halt_sys(void)
{
    intr |= INTR_HALTED;
}

/****************************************************************************
PARAMETERS:
mod     - Mod value from decoded byte
regh    - Reg h value from decoded byte
regl    - Reg l value from decoded byte

REMARKS:
Raise the specified interrupt to be handled before the execution of the
next instruction.

NOTE: Do not inline this function, as (*sys_rdb) is already inline!
****************************************************************************/
void fetch_decode_modrm(
    int *mod,
    int *regh,
    int *regl)
{
    int fetched;

DB( if (CHECK_IP_FETCH())
        x86emu_check_ip_access();)
    fetched = (*sys_rdb)(((u32)CS << 4) + (IP++));
    INC_DECODED_INST_LEN(1);
    *mod  = (fetched >> 6) & 0x03;
    *regh = (fetched >> 3) & 0x07;
    *regl = (fetched >> 0) & 0x07;
}

/****************************************************************************
RETURNS:
Immediate byte value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as (*sys_rdb) is already inline!
****************************************************************************/
u8 fetch_byte_imm(void)
{
    u8 fetched;

DB( if (CHECK_IP_FETCH())
        x86emu_check_ip_access();)
    fetched = (*sys_rdb)(((u32)CS << 4) + (IP++));
    INC_DECODED_INST_LEN(1);
    return fetched;
}

/****************************************************************************
RETURNS:
Immediate word value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as (*sys_rdw) is already inline!
****************************************************************************/
u16 fetch_word_imm(void)
{
    u16 fetched;

DB( if (CHECK_IP_FETCH())
        x86emu_check_ip_access();)
    fetched = (*sys_rdw)(((u32)CS << 4) + (IP));
    IP += 2;
    INC_DECODED_INST_LEN(2);
    return fetched;
}

/****************************************************************************
RETURNS:
Immediate lone value read from instruction queue

REMARKS:
This function returns the immediate byte from the instruction queue, and
moves the instruction pointer to the next value.

NOTE: Do not inline this function, as (*sys_rdw) is already inline!
****************************************************************************/
u32 fetch_long_imm(void)
{
    u32 fetched;

DB( if (CHECK_IP_FETCH())
        x86emu_check_ip_access();)
    fetched = (*sys_rdl)(((u32)CS << 4) + (IP));
    IP += 4;
    INC_DECODED_INST_LEN(4);
    return fetched;
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
cpu-state-variable mode. There are several potential states:

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
_INLINE u32 get_data_segment(void)
{
#define GET_SEGMENT(segment)
    switch (mode & SYSMODE_SEGMASK) {
      case 0:                   /* default case: use ds register */
      case SYSMODE_SEGOVR_DS:
      case SYSMODE_SEGOVR_DS | SYSMODE_SEG_DS_SS:
        return  DS;
      case SYSMODE_SEG_DS_SS:   /* non-overridden, use ss register */
        return  SS;
      case SYSMODE_SEGOVR_CS:
      case SYSMODE_SEGOVR_CS | SYSMODE_SEG_DS_SS:
        return  CS;
      case SYSMODE_SEGOVR_ES:
      case SYSMODE_SEGOVR_ES | SYSMODE_SEG_DS_SS:
        return  ES;
      case SYSMODE_SEGOVR_FS:
      case SYSMODE_SEGOVR_FS | SYSMODE_SEG_DS_SS:
        return  FS;
      case SYSMODE_SEGOVR_GS:
      case SYSMODE_SEGOVR_GS | SYSMODE_SEG_DS_SS:
        return  GS;
      case SYSMODE_SEGOVR_SS:
      case SYSMODE_SEGOVR_SS | SYSMODE_SEG_DS_SS:
        return  SS;
      default:
#ifdef  DEBUG
        loggy("error: should not happen:  multiple overrides.\n");
#endif
        HALT_SYS();
        return 0;
    }
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Byte value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u8 fetch_data_byte(
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    return (*sys_rdb)((get_data_segment() << 4) + offset);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Word value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u16 fetch_data_word(
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    return (*sys_rdw)((get_data_segment() << 4) + offset);
}

/****************************************************************************
PARAMETERS:
offset  - Offset to load data from

RETURNS:
Long value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u32 fetch_data_long(
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    return (*sys_rdl)((get_data_segment() << 4) + offset);
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Byte value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u8 fetch_data_byte_abs(
    uint segment,
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    return (*sys_rdb)(((u32)segment << 4) + offset);
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Word value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u16 fetch_data_word_abs(
    uint segment,
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    return (*sys_rdw)(((u32)segment << 4) + offset);
}

/****************************************************************************
PARAMETERS:
segment - Segment to load data from
offset  - Offset to load data from

RETURNS:
Long value read from the absolute memory location.

NOTE: Do not inline this function as (*sys_rdX) is already inline!
****************************************************************************/
u32 fetch_data_long_abs(
    uint segment,
    uint offset)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    return (*sys_rdl)(((u32)segment << 4) + offset);
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
void store_data_byte(
    uint offset,
    u8 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    (*sys_wrb)((get_data_segment() << 4) + offset, val);
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
void store_data_word(
    uint offset,
    u16 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    (*sys_wrw)((get_data_segment() << 4) + offset, val);
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
void store_data_long(
    uint offset,
    u32 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access((u16)get_data_segment(), offset);
#endif
    (*sys_wrl)((get_data_segment() << 4) + offset, val);
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
void store_data_byte_abs(
    uint segment,
    uint offset,
    u8 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    (*sys_wrb)(((u32)segment << 4) + offset, val);
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
void store_data_word_abs(
    uint segment,
    uint offset,
    u16 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    (*sys_wrw)(((u32)segment << 4) + offset, val);
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
void store_data_long_abs(
    uint segment,
    uint offset,
    u32 val)
{
#ifdef DEBUG
    if (CHECK_DATA_ACCESS())
        x86emu_check_data_access(segment, offset);
#endif
    (*sys_wrl)(((u32)segment << 4) + offset, val);
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
u8* decode_rm_byte_register(
    int reg)
{
    switch (reg) {
      case 0:
        DECODE_PRINTF("AL");
        return &AL;
      case 1:
        DECODE_PRINTF("CL");
        return &CL;
      case 2:
        DECODE_PRINTF("DL");
        return &DL;
      case 3:
        DECODE_PRINTF("BL");
        return &BL;
      case 4:
        DECODE_PRINTF("AH");
        return &AH;
      case 5:
        DECODE_PRINTF("CH");
        return &CH;
      case 6:
        DECODE_PRINTF("DH");
        return &DH;
      case 7:
        DECODE_PRINTF("BH");
        return &BH;
    }
    HALT_SYS();
    return NULL;                /* NOT REACHED OR REACHED ON ERROR */
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
u16* decode_rm_word_register(
    int reg)
{
    switch (reg) {
      case 0:
        DECODE_PRINTF("AX");
        return &AX;
      case 1:
        DECODE_PRINTF("CX");
        return &CX;
      case 2:
        DECODE_PRINTF("DX");
        return &DX;
      case 3:
        DECODE_PRINTF("BX");
        return &BX;
      case 4:
        DECODE_PRINTF("SP");
        return &SP;
      case 5:
        DECODE_PRINTF("BP");
        return &BP;
      case 6:
        DECODE_PRINTF("SI");
        return &SI;
      case 7:
        DECODE_PRINTF("DI");
        return &DI;
    }
    HALT_SYS();
    return NULL;                /* NOTREACHED OR REACHED ON ERROR */
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
u32* decode_rm_long_register(
    int reg)
{
    switch (reg) {
      case 0:
        DECODE_PRINTF("EAX");
        return &EAX;
      case 1:
        DECODE_PRINTF("ECX");
        return &ECX;
      case 2:
        DECODE_PRINTF("EDX");
        return &EDX;
      case 3:
        DECODE_PRINTF("EBX");
        return &EBX;
      case 4:
        DECODE_PRINTF("ESP");
        return &ESP;
      case 5:
        DECODE_PRINTF("EBP");
        return &EBP;
      case 6:
        DECODE_PRINTF("ESI");
        return &ESI;
      case 7:
        DECODE_PRINTF("EDI");
        return &EDI;
    }
    HALT_SYS();
    return NULL;                /* NOTREACHED OR REACHED ON ERROR */
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
u16* decode_rm_seg_register(
    int reg)
{
    switch (reg) {
      case 0:
        DECODE_PRINTF("ES");
        return &ES;
      case 1:
        DECODE_PRINTF("CS");
        return &CS;
      case 2:
        DECODE_PRINTF("SS");
        return &SS;
      case 3:
        DECODE_PRINTF("DS");
        return &DS;
      case 4:
        DECODE_PRINTF("FS");
        return &FS;
      case 5:
        DECODE_PRINTF("GS");
        return &GS;
      case 6:
      case 7:
        DECODE_PRINTF("ILLEGAL SEGREG");
        break;
    }
    HALT_SYS();
    return NULL;                /* NOT REACHED OR REACHED ON ERROR */
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
static unsigned decode_sib_si(
    int scale,
    int index)
{
    scale = 1 << scale;
    if (scale > 1) {
        DECODE_PRINTF2("[%d*", scale);
    } else {
        DECODE_PRINTF("[");
    }
    switch (index) {
      case 0:
        DECODE_PRINTF("EAX]");
        return EAX * index;
      case 1:
        DECODE_PRINTF("ECX]");
        return ECX * index;
      case 2:
        DECODE_PRINTF("EDX]");
        return EDX * index;
      case 3:
        DECODE_PRINTF("EBX]");
        return EBX * index;
      case 4:
        DECODE_PRINTF("0]");
        return 0;
      case 5:
        DECODE_PRINTF("EBP]");
        return EBP * index;
      case 6:
        DECODE_PRINTF("ESI]");
        return ESI * index;
      case 7:
        DECODE_PRINTF("EDI]");
        return EDI * index;
    }
    HALT_SYS();
    return 0;                   /* NOT REACHED OR REACHED ON ERROR */
}

/****************************************************************************
PARAMETERS:
mod - MOD value of preceding ModR/M byte

RETURNS:
Offset in memory for the address decoding

REMARKS:
Decodes SIB addressing byte and returns calculated effective address.
****************************************************************************/
static unsigned decode_sib_address(
    int mod)
{
    int sib   = fetch_byte_imm();
    int ss    = (sib >> 6) & 0x03;
    int index = (sib >> 3) & 0x07;
    int base  = sib & 0x07;
    int offset = 0;
    int displacement;

    switch (base) {
      case 0:
        DECODE_PRINTF("[EAX]");
        offset = EAX;
        break;
      case 1:
        DECODE_PRINTF("[ECX]");
        offset = ECX;
        break;
      case 2:
        DECODE_PRINTF("[EDX]");
        offset = EDX;
        break;
      case 3:
        DECODE_PRINTF("[EBX]");
        offset = EBX;
        break;
      case 4:
        DECODE_PRINTF("[ESP]");
        offset = ESP;
        break;
      case 5:
        switch (mod) {
          case 0:
            displacement = (s32)fetch_long_imm();
            DECODE_PRINTF2("[%d]", displacement);
            offset = displacement;
            break;
          case 1:
            displacement = (s8)fetch_byte_imm();
            DECODE_PRINTF2("[%d][EBP]", displacement);
            offset = EBP + displacement;
            break;
          case 2:
            displacement = (s32)fetch_long_imm();
            DECODE_PRINTF2("[%d][EBP]", displacement);
            offset = EBP + displacement;
            break;
          default:
            HALT_SYS();
        }
        DECODE_PRINTF("[EAX]");
        offset = EAX;
        break;
      case 6:
        DECODE_PRINTF("[ESI]");
        offset = ESI;
        break;
      case 7:
        DECODE_PRINTF("[EDI]");
        offset = EDI;
        break;
      default:
        HALT_SYS();
    }
    offset += decode_sib_si(ss, index);
    return offset;
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
unsigned decode_rm00_address(
    int rm)
{
    unsigned offset;

    if (mode & SYSMODE_PREFIX_ADDR) {
        /* 32-bit addressing */
        switch (rm) {
          case 0:
            DECODE_PRINTF("[EAX]");
            return EAX;
          case 1:
            DECODE_PRINTF("[ECX]");
            return ECX;
          case 2:
            DECODE_PRINTF("[EDX]");
            return EDX;
          case 3:
            DECODE_PRINTF("[EBX]");
            return EBX;
          case 4:
            return decode_sib_address(0);
          case 5:
            offset = fetch_long_imm();
            DECODE_PRINTF2("[%08x]", offset);
            return offset;
          case 6:
            DECODE_PRINTF("[ESI]");
            return ESI;
          case 7:
            DECODE_PRINTF("[EDI]");
            return EDI;
        }
    } else {
        /* 16-bit addressing */
        switch (rm) {
          case 0:
            DECODE_PRINTF("[BX+SI]");
            return (BX + SI) & 0xffff;
          case 1:
            DECODE_PRINTF("[BX+DI]");
            return (BX + DI) & 0xffff;
          case 2:
            DECODE_PRINTF("[BP+SI]");
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + SI) & 0xffff;
          case 3:
            DECODE_PRINTF("[BP+DI]");
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + DI) & 0xffff;
          case 4:
            DECODE_PRINTF("[SI]");
            return SI;
          case 5:
            DECODE_PRINTF("[DI]");
            return DI;
          case 6:
            offset = fetch_word_imm();
            DECODE_PRINTF2("[%04x]", offset);
            return offset;
          case 7:
            DECODE_PRINTF("[BX]");
            return BX;
        }
    }
    HALT_SYS();
    return 0;
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
unsigned decode_rm01_address(
    int rm)
{
    int displacement;

    if (mode & SYSMODE_PREFIX_ADDR) {
        /* 32-bit addressing */
        if (rm != 4)
            displacement = (s8)fetch_byte_imm();
        else
            displacement = 0;

        switch (rm) {
          case 0:
            DECODE_PRINTF2("%d[EAX]", displacement);
            return EAX + displacement;
          case 1:
            DECODE_PRINTF2("%d[ECX]", displacement);
            return ECX + displacement;
          case 2:
            DECODE_PRINTF2("%d[EDX]", displacement);
            return EDX + displacement;
          case 3:
            DECODE_PRINTF2("%d[EBX]", displacement);
            return EBX + displacement;
          case 4: {
            int offset = decode_sib_address(1);
            displacement = (s8)fetch_byte_imm();
            DECODE_PRINTF2("[%d]", displacement);
            return offset + displacement;
          }
          case 5:
            DECODE_PRINTF2("%d[EBP]", displacement);
            return EBP + displacement;
          case 6:
            DECODE_PRINTF2("%d[ESI]", displacement);
            return ESI + displacement;
          case 7:
            DECODE_PRINTF2("%d[EDI]", displacement);
            return EDI + displacement;
        }
    } else {
        /* 16-bit addressing */
        displacement = (s8)fetch_byte_imm();
        switch (rm) {
          case 0:
            DECODE_PRINTF2("%d[BX+SI]", displacement);
            return (BX + SI + displacement) & 0xffff;
          case 1:
            DECODE_PRINTF2("%d[BX+DI]", displacement);
            return (BX + DI + displacement) & 0xffff;
          case 2:
            DECODE_PRINTF2("%d[BP+SI]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + SI + displacement) & 0xffff;
          case 3:
            DECODE_PRINTF2("%d[BP+DI]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + DI + displacement) & 0xffff;
          case 4:
            DECODE_PRINTF2("%d[SI]", displacement);
            return (SI + displacement) & 0xffff;
          case 5:
            DECODE_PRINTF2("%d[DI]", displacement);
            return (DI + displacement) & 0xffff;
          case 6:
            DECODE_PRINTF2("%d[BP]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + displacement) & 0xffff;
          case 7:
            DECODE_PRINTF2("%d[BX]", displacement);
            return (BX + displacement) & 0xffff;
        }
    }
    HALT_SYS();
    return 0;                   /* SHOULD NOT HAPPEN */
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
unsigned decode_rm10_address(
    int rm)
{
    if (mode & SYSMODE_PREFIX_ADDR) {
        int displacement;

        /* 32-bit addressing */
        if (rm != 4)
            displacement = (s32)fetch_long_imm();
        else
            displacement = 0;

        switch (rm) {
          case 0:
            DECODE_PRINTF2("%d[EAX]", displacement);
            return EAX + displacement;
          case 1:
            DECODE_PRINTF2("%d[ECX]", displacement);
            return ECX + displacement;
          case 2:
            DECODE_PRINTF2("%d[EDX]", displacement);
            return EDX + displacement;
          case 3:
            DECODE_PRINTF2("%d[EBX]", displacement);
            return EBX + displacement;
          case 4: {
            int offset = decode_sib_address(2);
            displacement = (s32)fetch_long_imm();
            DECODE_PRINTF2("[%d]", displacement);
            return offset + displacement;
          }
          case 5:
            DECODE_PRINTF2("%d[EBP]", displacement);
            return EBP + displacement;
          case 6:
            DECODE_PRINTF2("%d[ESI]", displacement);
            return ESI + displacement;
          case 7:
            DECODE_PRINTF2("%d[EDI]", displacement);
            return EDI + displacement;
        }
    } else {
        int displacement = (s16)fetch_word_imm();

        /* 16-bit addressing */
        switch (rm) {
          case 0:
            DECODE_PRINTF2("%d[BX+SI]", displacement);
            return (BX + SI + displacement) & 0xffff;
          case 1:
            DECODE_PRINTF2("%d[BX+DI]", displacement);
            return (BX + DI + displacement) & 0xffff;
          case 2:
            DECODE_PRINTF2("%d[BP+SI]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + SI + displacement) & 0xffff;
          case 3:
            DECODE_PRINTF2("%d[BP+DI]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + DI + displacement) & 0xffff;
          case 4:
            DECODE_PRINTF2("%d[SI]", displacement);
            return (SI + displacement) & 0xffff;
          case 5:
            DECODE_PRINTF2("%d[DI]", displacement);
            return (DI + displacement) & 0xffff;
          case 6:
            DECODE_PRINTF2("%d[BP]", displacement);
            mode |= SYSMODE_SEG_DS_SS;
            return (BP + displacement) & 0xffff;
          case 7:
            DECODE_PRINTF2("%d[BX]", displacement);
            return (BX + displacement) & 0xffff;
        }
    }
    HALT_SYS();
    return 0;                   /* SHOULD NOT HAPPEN */
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

unsigned decode_rmXX_address(int mod, int rm)
{
  if (mod == 0)
    return decode_rm00_address(rm);
  if (mod == 1)
    return decode_rm01_address(rm);
  return decode_rm10_address(rm);
}

/****************************************************************************
*
*						Realmode X86 Emulator Library
*
*            	Copyright (C) 1996-1999 SciTech Software, Inc.
* 				     Copyright (C) David Mosberger-Tang
* 					   Copyright (C) 1999 Egbert Eich
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
* Language:		ANSI C
* Environment:	Any
* Developer:    Kendall Bennett
*
* Description:  Header file for x86 register definitions.
*
****************************************************************************/

#ifndef __X86EMU_REGS_H
#define __X86EMU_REGS_H

/*---------------------- Macros and type definitions ----------------------*/

#pragma pack(1)

/*
 * General EAX, EBX, ECX, EDX type registers.  Note that for
 * portability, and speed, the issue of byte swapping is not addressed
 * in the registers.
 */

/* flag conditions   */
#define FB_CF 0x0001            /* CARRY flag  */
#define FB_PF 0x0004            /* PARITY flag */
#define FB_AF 0x0010            /* AUX  flag   */
#define FB_ZF 0x0040            /* ZERO flag   */
#define FB_SF 0x0080            /* SIGN flag   */
#define FB_TF 0x0100            /* TRAP flag   */
#define FB_IF 0x0200            /* INTERRUPT ENABLE flag */
#define FB_DF 0x0400            /* DIR flag    */
#define FB_OF 0x0800            /* OVERFLOW flag */

/* 80286 and above always have bit#1 set */
#define F_ALWAYS_ON  (0x0002)   /* flag bits always on */

/*
 * Define a mask for only those flag bits we will ever pass back
 * (via PUSHF)
 */
#define F_MSK (FB_CF|FB_PF|FB_AF|FB_ZF|FB_SF|FB_TF|FB_IF|FB_DF|FB_OF)

/* following bits masked in to a 16bit quantity */

#define F_CF 0x0001             /* CARRY flag  */
#define F_PF 0x0004             /* PARITY flag */
#define F_AF 0x0010             /* AUX  flag   */
#define F_ZF 0x0040             /* ZERO flag   */
#define F_SF 0x0080             /* SIGN flag   */
#define F_TF 0x0100             /* TRAP flag   */
#define F_IF 0x0200             /* INTERRUPT ENABLE flag */
#define F_DF 0x0400             /* DIR flag    */
#define F_OF 0x0800             /* OVERFLOW flag */

#define TOGGLE_FLAG(flag)     	(FLG ^= (flag))
#define SET_FLAG(flag)        	(FLG |= (flag))
#define CLEAR_FLAG(flag)      	(FLG &= ~(flag))
#define ACCESS_FLAG(flag)     	(FLG & (flag))
#define CLEARALL_FLAG(m)    	(FLG = 0)

#define CONDITIONAL_SET_FLAG(COND,FLAG) \
  if (COND) SET_FLAG(FLAG); else CLEAR_FLAG(FLAG)

#define F_PF_CALC 0x010000      /* PARITY flag has been calced    */
#define F_ZF_CALC 0x020000      /* ZERO flag has been calced      */
#define F_SF_CALC 0x040000      /* SIGN flag has been calced      */

#define F_ALL_CALC      0xff0000        /* All have been calced   */

/*
 * Emulator machine state.
 * Segment usage control.
 */
#define SYSMODE_SEG_DS_SS       0x00000001
#define SYSMODE_SEGOVR_CS       0x00000002
#define SYSMODE_SEGOVR_DS       0x00000004
#define SYSMODE_SEGOVR_ES       0x00000008
#define SYSMODE_SEGOVR_FS       0x00000010
#define SYSMODE_SEGOVR_GS       0x00000020
#define SYSMODE_SEGOVR_SS       0x00000040
#define SYSMODE_PREFIX_REPE     0x00000080
#define SYSMODE_PREFIX_REPNE    0x00000100
#define SYSMODE_PREFIX_DATA     0x00000200
#define SYSMODE_PREFIX_ADDR     0x00000400
//phueper: for REP(E|NE) Instructions, we need to decide whether it should be
//using the 32bit ECX register as or the 16bit CX register as count register
#define SYSMODE_32BIT_REP       0x00000800
#define SYSMODE_INTR_PENDING    0x10000000
#define SYSMODE_EXTRN_INTR      0x20000000
#define SYSMODE_HALTED          0x40000000

#define SYSMODE_SEGMASK (SYSMODE_SEG_DS_SS      | \
						 SYSMODE_SEGOVR_CS      | \
						 SYSMODE_SEGOVR_DS      | \
						 SYSMODE_SEGOVR_ES      | \
						 SYSMODE_SEGOVR_FS      | \
						 SYSMODE_SEGOVR_GS      | \
						 SYSMODE_SEGOVR_SS)
#define SYSMODE_CLRMASK (SYSMODE_SEG_DS_SS      | \
						 SYSMODE_SEGOVR_CS      | \
						 SYSMODE_SEGOVR_DS      | \
						 SYSMODE_SEGOVR_ES      | \
						 SYSMODE_SEGOVR_FS      | \
						 SYSMODE_SEGOVR_GS      | \
						 SYSMODE_SEGOVR_SS      | \
						 SYSMODE_PREFIX_DATA    | \
						 SYSMODE_PREFIX_ADDR    | \
						 SYSMODE_32BIT_REP)

#define  INTR_SYNCH           0x1
#define  INTR_ASYNCH          0x2
#define  INTR_HALTED          0x4

/****************************************************************************
REMARKS:
Structure maintaining the emulator machine state.

MEMBERS:
mem_base		- Base real mode memory for the emulator
abseg			- Base for the absegment
mem_size		- Size of the real mode memory block for the emulator
private			- private data pointer
x86			- X86 registers
****************************************************************************/
	unsigned long	mem_base;
	unsigned long	mem_size;
	unsigned long	abseg;
	void*        	private;
	/*
	 * MODE contains information on:
	 *  REPE prefix             2 bits  repe,repne
	 *  SEGMENT overrides       5 bits  normal,DS,SS,CS,ES
	 *  Delayed flag set        3 bits  (zero, signed, parity)
	 *  reserved                6 bits
	 *  interrupt #             8 bits  instruction raised interrupt
	 *  BIOS video segregs      4 bits
	 *  Interrupt Pending       1 bits
	 *  Extern interrupt        1 bits
	 *  Halted                  1 bits
	 */
	u32                         mode;
	volatile int                intr;   /* mask of pending interrupts */
	volatile int                         debug;
	u32 A, B, C, D;
		u32 SP, BP, SI, DI, IP;
	u32 FLAGS;
	u16 CS, DS, SS, ES, FS, GS;
	int                         check;
	u16                         saved_ip;
	u16                         saved_cs;
	int                         enc_pos;
	int                         enc_str_pos;
	char                        decode_buf[32]; /* encoded byte stream  */
	char                        decoded_buf[256]; /* disassembled strings */
	u8                          intno;
	u8                          __pad[3];

#pragma pack()

/*----------------------------- Global Variables --------------------------*/

/* Global emulator machine state.
 *
 * We keep it global to avoid pointer dereferences in the code for speed.
 */

#define X86_EAX EAX
#define X86_EBX EBX
#define X86_ECX ECX
#define X86_EDX EDX
#define X86_ESI ESI
#define X86_EDI EDI
#define X86_EBP EBP
#define X86_EIP EIP
#define X86_ESP ESP
#define X86_EFLAGS EFLG

#define X86_FLAGS FLG
#define X86_AX AX
#define X86_BX BX
#define X86_CX CX
#define X86_DX DX
#define X86_SI SI
#define X86_DI DI
#define X86_BP BP
#define X86_IP IP
#define X86_SP SP
#define X86_CS CS
#define X86_DS DS
#define X86_ES ES
#define X86_SS SS
#define X86_FS FS
#define X86_GS GS

#define X86_AL AL
#define X86_BL BL
#define X86_CL CL
#define X86_DL DL

#define X86_AH AH
#define X86_BH BH
#define X86_CH CH
#define X86_DH DH

#endif /* __X86EMU_REGS_H */

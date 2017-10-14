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
* Description:  Header file for public specific functions.
*               Any application linking against us should only
*               include this header
*
****************************************************************************/

#ifndef __X86EMU_X86EMU_H
#define __X86EMU_X86EMU_H

#include "types.h"
#define	X86API
#define	X86APIP	*
#include "regs.h"

/*---------------------- Macros and type definitions ----------------------*/

#pragma	pack(1)

/****************************************************************************
  Here are the default memory read and write
  function in case they are needed as fallbacks.
***************************************************************************/
extern u8 X86API rdb(u32 addr);
extern u16 X86API rdw(u32 addr);
extern u32 X86API rdl(u32 addr);
extern void X86API wrb(u32 addr, u8 val);
extern void X86API wrw(u32 addr, u16 val);
extern void X86API wrl(u32 addr, u32 val);

#pragma	pack()

/*--------------------- type definitions -----------------------------------*/

//typedef void (X86APIP X86EMU_intrFuncs)(int num);
typedef u32 X86EMU_intrFuncs;
extern X86EMU_intrFuncs _X86EMU_intrTab[256];

/*-------------------------- Function Prototypes --------------------------*/

void 	X86EMU_prepareForInt(int num);

void X86EMU_setMemBase(void *base, int size);

/* decode.c */

void 	X86EMU_exec(void);
void 	X86EMU_halt_sys(void);

#define	HALT_SYS()	\
	loggy("halt_sys: in %s\n", __func__);	\
	X86EMU_halt_sys();
/* Debug options */

#define DEBUG_DECODE_F          0x000001 /* print decoded instruction  */
#define DEBUG_TRACE_F           0x000002 /* dump regs before/after execution */
#define DEBUG_STEP_F            0x000004
#define DEBUG_DISASSEMBLE_F     0x000008
#define DEBUG_BREAK_F           0x000010
#define DEBUG_SVC_F             0x000020
#define DEBUG_FS_F              0x000080
#define DEBUG_PROC_F            0x000100
#define DEBUG_SYSINT_F          0x000200 /* bios system interrupts. */
#define DEBUG_TRACECALL_F       0x000400
#define DEBUG_INSTRUMENT_F      0x000800
#define DEBUG_MEM_TRACE_F       0x001000
#define DEBUG_IO_TRACE_F        0x002000
#define DEBUG_TRACECALL_REGS_F  0x004000
#define DEBUG_DECODE_NOPRINT_F  0x008000
#define DEBUG_SAVE_IP_CS_F      0x010000
#define DEBUG_TRACEJMP_F        0x020000
#define DEBUG_TRACEJMP_REGS_F   0x040000
#define DEBUG_SYS_F             (DEBUG_SVC_F|DEBUG_FS_F|DEBUG_PROC_F)

void 	X86EMU_trace_regs(void);
void 	X86EMU_trace_xregs(void);
void 	X86EMU_dump_memory(u16 seg, u16 off, u32 amt);
int 	X86EMU_trace_on(void);
int 	X86EMU_trace_off(void);

void loggy(char *, ...);

#endif /* __X86EMU_X86EMU_H */

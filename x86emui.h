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
* Description:  Header file for system specific functions. These functions
*				are always compiled and linked in the OS dependent libraries,
*				and never in a binary portable driver.
*
****************************************************************************/

#ifndef __X86EMU_X86EMUI_H
#define __X86EMU_X86EMUI_H

/* If we are compiling in C++ mode, we can compile some functions as
 * inline to increase performance (however the code size increases quite
 * dramatically in this case).
 */

#define	_INLINE static

/* Get rid of unused parameters in C++ compilation mode */

#define	X86EMU_UNUSED(v) _

#include <stdio.h>
#include "x86emu.h"
#include "regs.h"
#include "debug.h"
#include "decode.h"
#include "ops.h"
#include "prim_ops.h"
#include "fpu.h"
#include "fpu_regs.h"

#include <string.h>
/*--------------------------- Inline Functions ----------------------------*/

extern u8  	(X86APIP sys_rdb)(u32 addr);
extern u16 	(X86APIP sys_rdw)(u32 addr);
extern u32 	(X86APIP sys_rdl)(u32 addr);
extern void (X86APIP sys_wrb)(u32 addr,u8 val);
extern void (X86APIP sys_wrw)(u32 addr,u16 val);
extern void (X86APIP sys_wrl)(u32 addr,u32 val);

extern u8  	(X86APIP sys_inb)(X86EMU_pioAddr addr);
extern u16 	(X86APIP sys_inw)(X86EMU_pioAddr addr);
extern u32 	(X86APIP sys_inl)(X86EMU_pioAddr addr);
extern void (X86APIP sys_outb)(X86EMU_pioAddr addr,u8 val);
extern void (X86APIP sys_outw)(X86EMU_pioAddr addr,u16 val);
extern void	(X86APIP sys_outl)(X86EMU_pioAddr addr,u32 val);
void panic(char *, ...);
#endif /* __X86EMU_X86EMUI_H */

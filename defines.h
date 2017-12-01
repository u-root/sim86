#define s8 char
#define u8 unsigned char
#define u16 unsigned short
#define s16 short
#define u32 unsigned int
#define s32 int
#define u64 unsigned long
#define size_t int
#define uint unsigned int
#define sint int;
#define X86EMU_pioAddr unsigned short
int printf(const char *format, ...);
#define DEBUG

static int CHECK_IP_FETCH_F = 0x1;
static int CHECK_SP_ACCESS_F = 0x2;
static int CHECK_MEM_ACCESS_F = 0x4; /*using regular linear pointer */
static int CHECK_DATA_ACCESS_F = 0x8; /*using segment:offset*/

# define CHECK_IP_FETCH()              	(M.x86.check & CHECK_IP_FETCH_F)
# define CHECK_SP_ACCESS()             	(M.x86.check & CHECK_SP_ACCESS_F)
# define CHECK_MEM_ACCESS()            	(M.x86.check & CHECK_MEM_ACCESS_F)
# define CHECK_DATA_ACCESS()           	(M.x86.check & CHECK_DATA_ACCESS_F)

# define DEBUG_INSTRUMENT()    	(M.x86.debug & DEBUG_INSTRUMENT_F)
# define DEBUG_DECODE()        	(M.x86.debug & DEBUG_DECODE_F)
# define DEBUG_TRACE()         	(M.x86.debug & DEBUG_TRACE_F)
# define DEBUG_STEP()          	(M.x86.debug & DEBUG_STEP_F)
# define DEBUG_DISASSEMBLE()   	(M.x86.debug & DEBUG_DISASSEMBLE_F)
# define DEBUG_BREAK()         	(M.x86.debug & DEBUG_BREAK_F)
# define DEBUG_SVC()           	(M.x86.debug & DEBUG_SVC_F)
# define DEBUG_SAVE_IP_CS()     (M.x86.debug & DEBUG_SAVE_IP_CS_F)

# define DEBUG_FS()            	(M.x86.debug & DEBUG_FS_F)
# define DEBUG_PROC()          	(M.x86.debug & DEBUG_PROC_F)
# define DEBUG_SYSINT()        	(M.x86.debug & DEBUG_SYSINT_F)
# define DEBUG_TRACECALL()     	(M.x86.debug & DEBUG_TRACECALL_F)
# define DEBUG_TRACECALLREGS() 	(M.x86.debug & DEBUG_TRACECALL_REGS_F)
# define DEBUG_TRACEJMP()       (M.x86.debug & DEBUG_TRACEJMP_F)
# define DEBUG_TRACEJMPREGS()   (M.x86.debug & DEBUG_TRACEJMP_REGS_F)
# define DEBUG_SYS()           	(M.x86.debug & DEBUG_SYS_F)
# define DEBUG_MEM_TRACE()     	(M.x86.debug & DEBUG_MEM_TRACE_F)
# define DEBUG_IO_TRACE()      	(M.x86.debug & DEBUG_IO_TRACE_F)
# define DEBUG_DECODE_NOPRINT() (M.x86.debug & DEBUG_DECODE_NOPRINT_F)

# define DECODE_PRINTF(x)     	if (DEBUG_DECODE()) \
									x86emu_decode_printf(x)
# define DECODE_PRINTF2(x,y)  	if (DEBUG_DECODE()) \
									x86emu_decode_printf2(x,y)

/*
 * The following allow us to look at the bytes of an instruction.  The
 * first INCR_INSTRN_LEN, is called every time bytes are consumed in
 * the decoding process.  The SAVE_IP_CS is called initially when the
 * major opcode of the instruction is accessed.
 */
#define INC_DECODED_INST_LEN(x)                    	\
	if (DEBUG_DECODE())  	                       	\
		x86emu_inc_decoded_inst_len(x)

#define SAVE_IP_CS(x,y)                               			\
	if (DEBUG_DECODE() | DEBUG_TRACECALL() | DEBUG_BREAK() \
              | DEBUG_IO_TRACE() | DEBUG_SAVE_IP_CS()) { \
		M.x86.saved_cs = x;                          			\
		M.x86.saved_ip = y;                          			\
	}
#define TRACE_REGS()                                   		\
	if (DEBUG_DISASSEMBLE()) {                         		\
		x86emu_just_disassemble();                        	\
		goto EndOfTheInstructionProcedure;             		\
	}                                                   	\
	if (DEBUG_TRACE() || DEBUG_DECODE()) X86EMU_trace_regs()

# define SINGLE_STEP()		if (DEBUG_STEP()) x86emu_single_step()

#define TRACE_AND_STEP()	\
	TRACE_REGS();			\
	SINGLE_STEP()

# define START_OF_INSTR()
# define END_OF_INSTR()		EndOfTheInstructionProcedure: x86emu_end_instr();
# define END_OF_INSTR_NO_TRACE()	x86emu_end_instr();

# define  CALL_TRACE(u,v,w,x,s)                                 \
	if (DEBUG_TRACECALLREGS())									\
		x86emu_dump_regs();                                     \
	if (DEBUG_TRACECALL())                                     	\
		printf("%04x:%04x: CALL %s%04x:%04x\n", u , v, s, w, x);
# define RETURN_TRACE(u,v,w,x,s)                                    \
	if (DEBUG_TRACECALLREGS())									\
		x86emu_dump_regs();                                     \
	if (DEBUG_TRACECALL())                                     	\
		printf("%04x:%04x: RET %s %04x:%04x\n",u,v,s,w,x);
# define  JMP_TRACE(u,v,w,x,s)                                 \
   if (DEBUG_TRACEJMPREGS()) \
      x86emu_dump_regs(); \
   if (DEBUG_TRACEJMP()) \
      printf("%04x:%04x: JMP %s%04x:%04x\n", u , v, s, w, x);

#define	DB(x)	x

#define X86EMU_DEBUG_ONLY(x) x
#define X86EMU_UNUSED(x) x


u32		get_flags_asm(void);
u16     aaa_word_asm(u32 *flags,u16 d);
u16     aas_word_asm(u32 *flags,u16 d);
u16     aad_word_asm(u32 *flags,u16 d);
u16     aam_word_asm(u32 *flags,u8 d);
u8      adc_byte_asm(u32 *flags,u8 d, u8 s);
u16     adc_word_asm(u32 *flags,u16 d, u16 s);
u32     adc_long_asm(u32 *flags,u32 d, u32 s);
u8      add_byte_asm(u32 *flags,u8 d, u8 s);
u16     add_word_asm(u32 *flags,u16 d, u16 s);
u32     add_long_asm(u32 *flags,u32 d, u32 s);
u8      and_byte_asm(u32 *flags,u8 d, u8 s);
u16     and_word_asm(u32 *flags,u16 d, u16 s);
u32     and_long_asm(u32 *flags,u32 d, u32 s);
u8      cmp_byte_asm(u32 *flags,u8 d, u8 s);
u16     cmp_word_asm(u32 *flags,u16 d, u16 s);
u32     cmp_long_asm(u32 *flags,u32 d, u32 s);
u8      daa_byte_asm(u32 *flags,u8 d);
u8      das_byte_asm(u32 *flags,u8 d);
u8      dec_byte_asm(u32 *flags,u8 d);
u16     dec_word_asm(u32 *flags,u16 d);
u32     dec_long_asm(u32 *flags,u32 d);
u8      inc_byte_asm(u32 *flags,u8 d);
u16     inc_word_asm(u32 *flags,u16 d);
u32     inc_long_asm(u32 *flags,u32 d);
u8      or_byte_asm(u32 *flags,u8 d, u8 s);
u16     or_word_asm(u32 *flags,u16 d, u16 s);
u32     or_long_asm(u32 *flags,u32 d, u32 s);
u8      neg_byte_asm(u32 *flags,u8 d);
u16     neg_word_asm(u32 *flags,u16 d);
u32     neg_long_asm(u32 *flags,u32 d);
u8      not_byte_asm(u32 *flags,u8 d);
u16     not_word_asm(u32 *flags,u16 d);
u32     not_long_asm(u32 *flags,u32 d);
u8      rcl_byte_asm(u32 *flags,u8 d, u8 s);
u16     rcl_word_asm(u32 *flags,u16 d, u8 s);
u32     rcl_long_asm(u32 *flags,u32 d, u8 s);
u8      rcr_byte_asm(u32 *flags,u8 d, u8 s);
u16     rcr_word_asm(u32 *flags,u16 d, u8 s);
u32     rcr_long_asm(u32 *flags,u32 d, u8 s);
u8      rol_byte_asm(u32 *flags,u8 d, u8 s);
u16     rol_word_asm(u32 *flags,u16 d, u8 s);
u32     rol_long_asm(u32 *flags,u32 d, u8 s);
u8      ror_byte_asm(u32 *flags,u8 d, u8 s);
u16     ror_word_asm(u32 *flags,u16 d, u8 s);
u32     ror_long_asm(u32 *flags,u32 d, u8 s);
u8      shl_byte_asm(u32 *flags,u8 d, u8 s);
u16     shl_word_asm(u32 *flags,u16 d, u8 s);
u32     shl_long_asm(u32 *flags,u32 d, u8 s);
u8      shr_byte_asm(u32 *flags,u8 d, u8 s);
u16     shr_word_asm(u32 *flags,u16 d, u8 s);
u32     shr_long_asm(u32 *flags,u32 d, u8 s);
u8      sar_byte_asm(u32 *flags,u8 d, u8 s);
u16     sar_word_asm(u32 *flags,u16 d, u8 s);
u32     sar_long_asm(u32 *flags,u32 d, u8 s);
u16		shld_word_asm(u32 *flags,u16 d, u16 fill, u8 s);
u32     shld_long_asm(u32 *flags,u32 d, u32 fill, u8 s);
u16		shrd_word_asm(u32 *flags,u16 d, u16 fill, u8 s);
u32     shrd_long_asm(u32 *flags,u32 d, u32 fill, u8 s);
u8      sbb_byte_asm(u32 *flags,u8 d, u8 s);
u16     sbb_word_asm(u32 *flags,u16 d, u16 s);
u32     sbb_long_asm(u32 *flags,u32 d, u32 s);
u8      sub_byte_asm(u32 *flags,u8 d, u8 s);
u16     sub_word_asm(u32 *flags,u16 d, u16 s);
u32     sub_long_asm(u32 *flags,u32 d, u32 s);
void	test_byte_asm(u32 *flags,u8 d, u8 s);
void	test_word_asm(u32 *flags,u16 d, u16 s);
void	test_long_asm(u32 *flags,u32 d, u32 s);
u8      xor_byte_asm(u32 *flags,u8 d, u8 s);
u16     xor_word_asm(u32 *flags,u16 d, u16 s);
u32     xor_long_asm(u32 *flags,u32 d, u32 s);
void    imul_byte_asm(u32 *flags,u16 *ax,u8 d,u8 s);
void    imul_word_asm(u32 *flags,u16 *ax,u16 *dx,u16 d,u16 s);
void    imul_long_asm(u32 *flags,u32 *eax,u32 *edx,u32 d,u32 s);
void    mul_byte_asm(u32 *flags,u16 *ax,u8 d,u8 s);
void    mul_word_asm(u32 *flags,u16 *ax,u16 *dx,u16 d,u16 s);
void    mul_long_asm(u32 *flags,u32 *eax,u32 *edx,u32 d,u32 s);
void	idiv_byte_asm(u32 *flags,u8 *al,u8 *ah,u16 d,u8 s);
void	idiv_word_asm(u32 *flags,u16 *ax,u16 *dx,u16 dlo,u16 dhi,u16 s);
void	idiv_long_asm(u32 *flags,u32 *eax,u32 *edx,u32 dlo,u32 dhi,u32 s);
void	div_byte_asm(u32 *flags,u8 *al,u8 *ah,u16 d,u8 s);
void	div_word_asm(u32 *flags,u16 *ax,u16 *dx,u16 dlo,u16 dhi,u16 s);
void	div_long_asm(u32 *flags,u32 *eax,u32 *edx,u32 dlo,u32 dhi,u32 s);
u16     aaa_word (u16 d);
u16     aas_word (u16 d);
u16     aad_word (u16 d);
u16     aam_word (u8 d);
u8      adc_byte (u8 d, u8 s);
u16     adc_word (u16 d, u16 s);
u32     adc_long (u32 d, u32 s);
u8      add_byte (u8 d, u8 s);
u16     add_word (u16 d, u16 s);
u32     add_long (u32 d, u32 s);
u8      and_byte (u8 d, u8 s);
u16     and_word (u16 d, u16 s);
u32     and_long (u32 d, u32 s);
u8      cmp_byte (u8 d, u8 s);
u16     cmp_word (u16 d, u16 s);
u32     cmp_long (u32 d, u32 s);
u8      daa_byte (u8 d);
u8      das_byte (u8 d);
u8      dec_byte (u8 d);
u16     dec_word (u16 d);
u32     dec_long (u32 d);
u8      inc_byte (u8 d);
u16     inc_word (u16 d);
u32     inc_long (u32 d);
u8      or_byte (u8 d, u8 s);
u16     or_word (u16 d, u16 s);
u32     or_long (u32 d, u32 s);
u8      neg_byte (u8 s);
u16     neg_word (u16 s);
u32     neg_long (u32 s);
u8      not_byte (u8 s);
u16     not_word (u16 s);
u32     not_long (u32 s);
u8      rcl_byte (u8 d, u8 s);
u16     rcl_word (u16 d, u8 s);
u32     rcl_long (u32 d, u8 s);
u8      rcr_byte (u8 d, u8 s);
u16     rcr_word (u16 d, u8 s);
u32     rcr_long (u32 d, u8 s);
u8      rol_byte (u8 d, u8 s);
u16     rol_word (u16 d, u8 s);
u32     rol_long (u32 d, u8 s);
u8      ror_byte (u8 d, u8 s);
u16     ror_word (u16 d, u8 s);
u32     ror_long (u32 d, u8 s);
u8      shl_byte (u8 d, u8 s);
u16     shl_word (u16 d, u8 s);
u32     shl_long (u32 d, u8 s);
u8      shr_byte (u8 d, u8 s);
u16     shr_word (u16 d, u8 s);
u32     shr_long (u32 d, u8 s);
u8      sar_byte (u8 d, u8 s);
u16     sar_word (u16 d, u8 s);
u32     sar_long (u32 d, u8 s);
u16     shld_word (u16 d, u16 fill, u8 s);
u32     shld_long (u32 d, u32 fill, u8 s);
u16     shrd_word (u16 d, u16 fill, u8 s);
u32     shrd_long (u32 d, u32 fill, u8 s);
u8      sbb_byte (u8 d, u8 s);
u16     sbb_word (u16 d, u16 s);
u32     sbb_long (u32 d, u32 s);
u8      sub_byte (u8 d, u8 s);
u16     sub_word (u16 d, u16 s);
u32     sub_long (u32 d, u32 s);
void    test_byte (u8 d, u8 s);
void    test_word (u16 d, u16 s);
void    test_long (u32 d, u32 s);
u8      xor_byte (u8 d, u8 s);
u16     xor_word (u16 d, u16 s);
u32     xor_long (u32 d, u32 s);
void    imul_byte (u8 s);
void    imul_word (u16 s);
void    imul_long (u32 s);
void 	imul_long_direct(u32 *res_lo, u32* res_hi,u32 d, u32 s);
void    mul_byte (u8 s);
void    mul_word (u16 s);
void    mul_long (u32 s);
void    idiv_byte (u8 s);
void    idiv_word (u16 s);
void    idiv_long (u32 s);
void    div_byte (u8 s);
void    div_word (u16 s);
void    div_long (u32 s);
void    ins (int size);
void    outs (int size);
u16     mem_access_word (int addr);
void    push_word (u16 w);
void    push_long (u32 w);
u16     pop_word (void);
u32	pop_long (void);
void	x86emu_cpuid (void);
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

/*---------------------- Macros and type definitions ----------------------*/

/*typedef*/ struct I32_reg_t{
	u32 e_reg;
	} ;

/*typedef*/ struct I16_reg_t{
	u16 x_reg;
	} ;

/*typedef*/ struct I8_reg_t{
	u8 l_reg, h_reg;
	} ;

struct i386_general_register{
	struct I32_reg_t   I32_reg;
	struct I16_reg_t   I16_reg;
	struct I8_reg_t    I8_reg;
	} ;

struct i386_general_regs {
	struct i386_general_register A, B, C, D;
	};

/*typedef*/ struct i386_general_regs Gen_reg_t;

struct i386_special_regs {
	struct i386_general_register SP, BP, SI, DI, IP;
	u32 FLAGS;
	};

/*
 * Segment registers here represent the 16 bit quantities
 * CS, DS, ES, SS.
 */

struct i386_segment_regs {
	u16 CS, DS, SS, ES, FS, GS;
	};

/* 8 bit registers */
#define R_AH  gen.A.I8_reg.h_reg
#define R_AL  gen.A.I8_reg.l_reg
#define R_BH  gen.B.I8_reg.h_reg
#define R_BL  gen.B.I8_reg.l_reg
#define R_CH  gen.C.I8_reg.h_reg
#define R_CL  gen.C.I8_reg.l_reg
#define R_DH  gen.D.I8_reg.h_reg
#define R_DL  gen.D.I8_reg.l_reg

/* 16 bit registers */
#define R_AX  gen.A.I16_reg.x_reg
#define R_BX  gen.B.I16_reg.x_reg
#define R_CX  gen.C.I16_reg.x_reg
#define R_DX  gen.D.I16_reg.x_reg

/* 32 bit extended registers */
#define R_EAX  gen.A.I32_reg.e_reg
#define R_EBX  gen.B.I32_reg.e_reg
#define R_ECX  gen.C.I32_reg.e_reg
#define R_EDX  gen.D.I32_reg.e_reg

/* special registers */
#define R_SP  spc.SP.I16_reg.x_reg
#define R_BP  spc.BP.I16_reg.x_reg
#define R_SI  spc.SI.I16_reg.x_reg
#define R_DI  spc.DI.I16_reg.x_reg
#define R_IP  spc.IP.I16_reg.x_reg
#define R_FLG spc.FLAGS

/* special registers */
#define R_SP  spc.SP.I16_reg.x_reg
#define R_BP  spc.BP.I16_reg.x_reg
#define R_SI  spc.SI.I16_reg.x_reg
#define R_DI  spc.DI.I16_reg.x_reg
#define R_IP  spc.IP.I16_reg.x_reg
#define R_FLG spc.FLAGS

/* special registers */
#define R_ESP  spc.SP.I32_reg.e_reg
#define R_EBP  spc.BP.I32_reg.e_reg
#define R_ESI  spc.SI.I32_reg.e_reg
#define R_EDI  spc.DI.I32_reg.e_reg
#define R_EIP  spc.IP.I32_reg.e_reg
#define R_EFLG spc.FLAGS

/* segment registers */
#define R_CS  seg.CS
#define R_DS  seg.DS
#define R_SS  seg.SS
#define R_ES  seg.ES
#define R_FS  seg.FS
#define R_GS  seg.GS

/* flag conditions   */
static int FB_CF = 0x0001; /* CARRY flag */
static int FB_PF = 0x0004; /* PARITY flag */
static int FB_AF = 0x0010; /* AUX flag */
static int FB_ZF = 0x0040; /* ZERO flag */
static int FB_SF = 0x0080; /* SIGN flag */
static int FB_TF = 0x0100; /* TRAP flag */
static int FB_IF = 0x0200; /* INTERRUPT ENABLE flag */
static int FB_DF = 0x0400; /* DIR flag */
static int FB_OF = 0x0800; /* OVERFLOW flag */

/* 80286 and above always have bit#1 set */
#define F_ALWAYS_ON (0x0002) /* flag bits always on */

/*
 * Define a mask for only those flag bits we will ever pass back
 * (via PUSHF)
 */
#define F_MSK (FB_CF|FB_PF|FB_AF|FB_ZF|FB_SF|FB_TF|FB_IF|FB_DF|FB_OF)

/* following bits masked in to a 16bit quantity */

enum {
 F_CF = 0x0001, /* CARRY flag */
 F_PF = 0x0004, /* PARITY flag */
 F_AF = 0x0010, /* AUX flag */
 F_ZF = 0x0040, /* ZERO flag */
 F_SF = 0x0080, /* SIGN flag */
 F_TF = 0x0100, /* TRAP flag */
 F_IF = 0x0200, /* INTERRUPT ENABLE flag */
 F_DF = 0x0400, /* DIR flag */
 F_OF = 0x0800 /* OVERFLOW flag */
};

#define TOGGLE_FLAG(flag)     	(M.x86.R_FLG ^= (flag))
#define SET_FLAG(flag)        	(M.x86.R_FLG |= (flag))
#define CLEAR_FLAG(flag)      	(M.x86.R_FLG &= ~(flag))
#define ACCESS_FLAG(flag)     	(M.x86.R_FLG & (flag))
#define CLEARALL_FLAG(m)    	(M.x86.R_FLG = 0)

#define CONDITIONAL_SET_FLAG(COND,FLAG) \
  if (COND) SET_FLAG(FLAG); else CLEAR_FLAG(FLAG)

static int F_PF_CALC = 0x010000; /* PARITY flag has been calced */
static int F_ZF_CALC = 0x020000; /* ZERO flag has been calced */
static int F_SF_CALC = 0x040000; /* SIGN flag has been calced */

static int F_ALL_CALC = 0xff0000; /* All have been calced */

/*
 * Emulator machine state.
 * Segment usage control.
 */
enum {
 SYSMODE_SEG_DS_SS = 0x00000001,
 SYSMODE_SEGOVR_CS = 0x00000002,
 SYSMODE_SEGOVR_DS = 0x00000004,
 SYSMODE_SEGOVR_ES = 0x00000008,
 SYSMODE_SEGOVR_FS = 0x00000010,
 SYSMODE_SEGOVR_GS = 0x00000020,
 SYSMODE_SEGOVR_SS = 0x00000040,
 SYSMODE_PREFIX_REPE = 0x00000080,
 SYSMODE_PREFIX_REPNE = 0x00000100,
 SYSMODE_PREFIX_DATA = 0x00000200,
 SYSMODE_PREFIX_ADDR = 0x00000400,
//phueper: for REP(E|NE) Instructions, we need to decide whether it should be
//using the 32bit ECX register as or the 16bit CX register as count register
 SYSMODE_32BIT_REP = 0x00000800,
 SYSMODE_INTR_PENDING = 0x10000000,
 SYSMODE_EXTRN_INTR = 0x20000000,
 SYSMODE_HALTED = 0x40000000
 };

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

static int INTR_SYNCH = 0x1;
static unsigned int INTR_ASYNCH = 0x2;
static unsigned int INTR_HALTED = 0x4;

/*typedef*/ struct X86EMU_regs{
	struct i386_general_regs    gen;
	struct i386_special_regs    spc;
	struct i386_segment_regs    seg;
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
	volatile u32                         debug;

	int                         check;
	u16                         saved_ip;
	u16                         saved_cs;
	int                         enc_pos;
	int                         enc_str_pos;
	char                        decode_buf[32]; /* encoded byte stream  */
	char                        decoded_buf[256]; /* disassembled strings */

	u8                          intno;
	u8                          __pad[3];
	} ;

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
/*typedef*/ struct X86EMU_sysEnv{
	unsigned long	mem_base;
	unsigned long	mem_size;
	unsigned long	abseg;
	void*        	private;
	struct X86EMU_regs		x86;
	} ;


/*----------------------------- Global Variables --------------------------*/


/* Global emulator machine state.
 *
 * We keep it global to avoid pointer dereferences in the code for speed.
 */

extern    struct X86EMU_sysEnv	_X86EMU_env;
#define   M             _X86EMU_env

#define X86_EAX M.x86.R_EAX
#define X86_EBX M.x86.R_EBX
#define X86_ECX M.x86.R_ECX
#define X86_EDX M.x86.R_EDX
#define X86_ESI M.x86.R_ESI
#define X86_EDI M.x86.R_EDI
#define X86_EBP M.x86.R_EBP
#define X86_EIP M.x86.R_EIP
#define X86_ESP M.x86.R_ESP
#define X86_EFLAGS M.x86.R_EFLG

#define X86_FLAGS M.x86.R_FLG
#define X86_AX M.x86.R_AX
#define X86_BX M.x86.R_BX
#define X86_CX M.x86.R_CX
#define X86_DX M.x86.R_DX
#define X86_SI M.x86.R_SI
#define X86_DI M.x86.R_DI
#define X86_BP M.x86.R_BP
#define X86_IP M.x86.R_IP
#define X86_SP M.x86.R_SP
#define X86_CS M.x86.R_CS
#define X86_DS M.x86.R_DS
#define X86_ES M.x86.R_ES
#define X86_SS M.x86.R_SS
#define X86_FS M.x86.R_FS
#define X86_GS M.x86.R_GS

#define X86_AL M.x86.R_AL
#define X86_BL M.x86.R_BL
#define X86_CL M.x86.R_CL
#define X86_DL M.x86.R_DL

#define X86_AH M.x86.R_AH
#define X86_BH M.x86.R_BH
#define X86_CH M.x86.R_CH
#define X86_DH M.x86.R_DH


/****************************************************************************
REMARKS:
Data structure containing pointers to programmed I/O functions used by the
emulator. This is used so that the user program can hook all programmed
I/O for the emulator to handled as necessary by the user program. By
default the emulator contains simple functions that do not do access the
hardware in any way. To allow the emulator access the hardware, you will
need to override the programmed I/O functions using the X86EMU_setupPioFuncs
function.

HEADER:
x86emu.h

MEMBERS:
inb		- Function to read a byte from an I/O port
inw		- Function to read a word from an I/O port
inl     - Function to read a dword from an I/O port
outb	- Function to write a byte to an I/O port
outw    - Function to write a word to an I/O port
outl    - Function to write a dword to an I/O port
****************************************************************************/
/*typedef*/ struct X86EMU_pioFuncs{
	u8  	(* inb)(X86EMU_pioAddr addr);
	u16 	(* inw)(X86EMU_pioAddr addr);
	u32 	(* inl)(X86EMU_pioAddr addr);
	void 	(* outb)(X86EMU_pioAddr addr, u8 val);
	void 	(* outw)(X86EMU_pioAddr addr, u16 val);
	void 	(* outl)(X86EMU_pioAddr addr, u32 val);
	} ;

/****************************************************************************
REMARKS:
Data structure containing pointers to memory access functions used by the
emulator. This is used so that the user program can hook all memory
access functions as necessary for the emulator. By default the emulator
contains simple functions that only access the internal memory of the
emulator. If you need specialized functions to handle access to different
types of memory (ie: hardware framebuffer accesses and BIOS memory access
etc), you will need to override this using the X86EMU_setupMemFuncs
function.

HEADER:
x86emu.h

MEMBERS:
rdb		- Function to read a byte from an address
rdw		- Function to read a word from an address
rdl     - Function to read a dword from an address
wrb		- Function to write a byte to an address
wrw    	- Function to write a word to an address
wrl    	- Function to write a dword to an address
****************************************************************************/
/*typedef*/ struct X86EMU_memFuncs{
	u8  	(* rdb)(u32 addr);
	u16 	(* rdw)(u32 addr);
	u32 	(* rdl)(u32 addr);
	void 	(* wrb)(u32 addr, u8 val);
	void 	(* wrw)(u32 addr, u16 val);
	void	(* wrl)(u32 addr, u32 val);
	} ;

/****************************************************************************
  Here are the default memory read and write
  function in case they are needed as fallbacks.
***************************************************************************/
extern u8  rdb(u32 addr);
extern u16  rdw(u32 addr);
extern u32  rdl(u32 addr);
extern void  wrb(u32 addr, u8 val);
extern void  wrw(u32 addr, u16 val);
extern void  wrl(u32 addr, u32 val);



/*--------------------- type definitions -----------------------------------*/

typedef void (* X86EMU_intrFuncs)(int num);
struct X86EMU_intrFuncs {
       void *f;
};
extern struct X86EMU_intrFuncs * _X86EMU_intrTab[256];

/*-------------------------- Function Prototypes --------------------------*/


void 	X86EMU_setupMemFuncs(struct X86EMU_memFuncs *funcs);
void 	X86EMU_setupPioFuncs(struct X86EMU_pioFuncs *funcs);
void 	X86EMU_setupIntrFuncs(struct X86EMU_intrFuncs funcs[]);
void 	X86EMU_prepareForInt(int num);

void X86EMU_setMemBase(void *base, size_t size);

/* decode.c */

void 	X86EMU_exec(void);
void 	X86EMU_halt_sys(void);

#define	HALT_SYS()	\
	printf("halt_sys: in %s\n", __func__);	\
	X86EMU_halt_sys();

/* Debug options */

static unsigned int DEBUG_DECODE_F = 0x000001; /* print decoded instruction */
static unsigned int DEBUG_TRACE_F = 0x000002; /* dump regs before/after execution */
static unsigned int DEBUG_STEP_F = 0x000004;
static unsigned int DEBUG_DISASSEMBLE_F = 0x000008;
static unsigned int DEBUG_BREAK_F = 0x000010;
static unsigned int DEBUG_SVC_F = 0x000020;
static unsigned int DEBUG_FS_F = 0x000080;
static unsigned int DEBUG_PROC_F = 0x000100;
static unsigned int DEBUG_SYSINT_F = 0x000200; /* bios system interrupts. */
static unsigned int DEBUG_TRACECALL_F = 0x000400;
static unsigned int DEBUG_INSTRUMENT_F = 0x000800;
static unsigned int DEBUG_MEM_TRACE_F = 0x001000;
static unsigned int DEBUG_IO_TRACE_F = 0x002000;
static unsigned int DEBUG_TRACECALL_REGS_F = 0x004000;
static unsigned int DEBUG_DECODE_NOPRINT_F = 0x008000;
static unsigned int DEBUG_SAVE_IP_CS_F = 0x010000;
static unsigned int DEBUG_TRACEJMP_F = 0x020000;
static unsigned int DEBUG_TRACEJMP_REGS_F = 0x040000;
static unsigned int DEBUG_SYS_F;
static void initDEBUG_SYS_F() {
	DEBUG_SYS_F  =            (DEBUG_SVC_F|DEBUG_FS_F|DEBUG_PROC_F);
}

void 	X86EMU_trace_regs(void);
void 	X86EMU_trace_xregs(void);
void 	X86EMU_dump_memory(u16 seg, u16 off, u32 amt);
int 	X86EMU_trace_on(void);
int 	X86EMU_trace_off(void);

void x86emu_inc_decoded_inst_len (int x);
void x86emu_decode_printf (const char *x);
void x86emu_decode_printf2 (const char *x, int y);
void x86emu_just_disassemble (void);
void x86emu_single_step (void);
void x86emu_end_instr (void);
void x86emu_dump_regs (void);
void x86emu_dump_xregs (void);
void x86emu_print_int_vect (u16 iv);
void x86emu_instrument_instruction (void);
void x86emu_check_ip_access (void);
void x86emu_check_sp_access (void);
void x86emu_check_mem_access (u32 p);
void x86emu_check_data_access (uint s, uint o);

void disassemble_forward (u16 seg, u16 off, int n);
// decode.h
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
* Description:  Header file for instruction decoding logic.
*
****************************************************************************/


/*---------------------- Macros and type definitions ----------------------*/

/* Instruction Decoding Stuff */

#define FETCH_DECODE_MODRM(mod,rh,rl) 	fetch_decode_modrm(&mod,&rh,&rl)
#define DECODE_RM_BYTE_REGISTER(r)    	decode_rm_byte_register(r)
#define DECODE_RM_WORD_REGISTER(r)    	decode_rm_word_register(r)
#define DECODE_RM_LONG_REGISTER(r)    	decode_rm_long_register(r)
#define DECODE_CLEAR_SEGOVR()         	M.x86.mode &= ~SYSMODE_CLRMASK

/*-------------------------- Function Prototypes --------------------------*/


void 	x86emu_intr_raise (u8 type);
void    fetch_decode_modrm (int *mod,int *regh,int *regl);
u8      fetch_byte_imm (void);
u16     fetch_word_imm (void);
u32     fetch_long_imm (void);
u8      fetch_data_byte (uint offset);
u8      fetch_data_byte_abs (uint segment, uint offset);
u16     fetch_data_word (uint offset);
u16     fetch_data_word_abs (uint segment, uint offset);
u32     fetch_data_long (uint offset);
u32     fetch_data_long_abs (uint segment, uint offset);
void    store_data_byte (uint offset, u8 val);
void    store_data_byte_abs (uint segment, uint offset, u8 val);
void    store_data_word (uint offset, u16 val);
void    store_data_word_abs (uint segment, uint offset, u16 val);
void    store_data_long (uint offset, u32 val);
void    store_data_long_abs (uint segment, uint offset, u32 val);
u8* 	decode_rm_byte_register(int reg);
u16* 	decode_rm_word_register(int reg);
u32* 	decode_rm_long_register(int reg);
u16* 	decode_rm_seg_register(int reg);
unsigned decode_rm00_address(int rm);
unsigned decode_rm01_address(int rm);
unsigned decode_rm10_address(int rm);
unsigned decode_rmXX_address(int mod, int rm);


extern u8  	(* sys_rdb)(u32 addr);
extern u16 	(* sys_rdw)(u32 addr);
extern u32 	(* sys_rdl)(u32 addr);
extern void (* sys_wrb)(u32 addr,u8 val);
extern void (* sys_wrw)(u32 addr,u16 val);
extern void (* sys_wrl)(u32 addr,u32 val);

extern u8  	(* sys_inb)(X86EMU_pioAddr addr);
extern u16 	(* sys_inw)(X86EMU_pioAddr addr);
extern u32 	(* sys_inl)(X86EMU_pioAddr addr);
extern void (* sys_outb)(X86EMU_pioAddr addr,u8 val);
extern void (* sys_outw)(X86EMU_pioAddr addr,u16 val);
extern void	(* sys_outl)(X86EMU_pioAddr addr,u32 val);

       
// ops.h
extern void (*x86emu_optab[0x100])(u8 op1);
extern void (*x86emu_optab2[0x100])(u8 op2);
int x86emu_check_jump_condition(u8 op);

// fpu.h

extern void x86emuOp_esc_coprocess_d8 (u8 op1);
extern void x86emuOp_esc_coprocess_d9 (u8 op1);
extern void x86emuOp_esc_coprocess_da (u8 op1);
extern void x86emuOp_esc_coprocess_db (u8 op1);
extern void x86emuOp_esc_coprocess_dc (u8 op1);
extern void x86emuOp_esc_coprocess_dd (u8 op1);
extern void x86emuOp_esc_coprocess_de (u8 op1);
extern void x86emuOp_esc_coprocess_df (u8 op1);

#define NULL ((void *)0)
#define _INLINE 

package main

import (
	"fmt"
	"log"
)

var notyet = `
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


#define TOGGLE_FLAG(flag)     	M.x86.R_FLG ^= (flag)
#define SET_FLAG(flag)        	M.x86.R_FLG |= (flag)
#define CLEAR_FLAG(flag)      	M.x86.R_FLG &= ~(flag)
#define ACCESS_FLAG(flag)     	M.x86.R_FLG & (flag)
#define CLEARALL_FLAG(m)    	M.x86.R_FLG = 0

#define CONDITIONAL_SET_FLAG(COND,FLAG) \
  if (COND) SET_FLAG(FLAG); else CLEAR_FLAG(FLAG)

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

/* Instruction Decoding Stuff */

#define FETCH_DECODE_MODRM(mod,rh,rl) 	fetch_decode_modrm(&mod,&rh,&rl)
#define DECODE_RM_BYTE_REGISTER(r)    	decode_rm_byte_register(r)
#define DECODE_RM_WORD_REGISTER(r)    	decode_rm_word_register(r)
#define DECODE_RM_LONG_REGISTER(r)    	decode_rm_long_register(r)
#define DECODE_CLEAR_SEGOVR()         	M.x86.mode &= ~SYSMODE_CLRMASK

/*-------------------------- Function Prototypes --------------------------*/

`

var M = &_X86EMU_env
var ro = uint32(0)

func TOGGLE_FLAG(flag uint32) {
	f := G32(EFLAGS)
	f ^= flag
	S(EFLAGS, f)
}
func SET_FLAG(flag uint32) {
	f := G32(EFLAGS)
	f |= flag
	S(EFLAGS, f)
}
func CLEAR_FLAG(flag uint32) {
	f := G32(EFLAGS)
	f &= ^flag
	S(EFLAGS, f)
}
func ACCESS_FLAG(flag uint32) bool {
	return G32(EFLAGS)&flag != 0
}
func CLEARALL_FLAG(_ uint32) {
	S(EFLAGS, 0)
}

// :.,$s/func \(.*\) {^M\(.*\)/func \1() {\2}/^M}
func CHECK_IP_FETCH() bool {
	return (M.x86.check & CHECK_IP_FETCH_F) != 0
}
func CHECK_SP_ACCESS() bool {
	return (M.x86.check & CHECK_SP_ACCESS_F) != 0
}
func CHECK_MEM_ACCESS() bool {
	return (M.x86.check & CHECK_MEM_ACCESS_F) != 0
}
func CHECK_DATA_ACCESS() bool {
	return (M.x86.check & CHECK_DATA_ACCESS_F) != 0
}

func DEBUG_INSTRUMENT() bool {
	return (M.x86.debug & DEBUG_INSTRUMENT_F) != 0
}
func DEBUG_DECODE() bool {
	return (M.x86.debug & DEBUG_DECODE_F) != 0
}
func DEBUG_TRACE() bool {
	return (M.x86.debug & DEBUG_TRACE_F) != 0
}
func DEBUG_STEP() bool {
	return (M.x86.debug & DEBUG_STEP_F) != 0
}
func DEBUG_DISASSEMBLE() bool {
	return (M.x86.debug & DEBUG_DISASSEMBLE_F) != 0
}
func DEBUG_BREAK() bool {
	return (M.x86.debug & DEBUG_BREAK_F) != 0
}
func DEBUG_SVC() bool {
	return (M.x86.debug & DEBUG_SVC_F) != 0
}
func DEBUG_SAVE_IP_CS() bool {
	return (M.x86.debug & DEBUG_SAVE_IP_CS_F) != 0
}

func DEBUG_FS() bool {
	return (M.x86.debug & DEBUG_FS_F) != 0
}
func DEBUG_PROC() bool {
	return (M.x86.debug & DEBUG_PROC_F) != 0
}
func DEBUG_SYSINT() bool {
	return (M.x86.debug & DEBUG_SYSINT_F) != 0
}
func DEBUG_TRACECALL() bool {
	return (M.x86.debug & DEBUG_TRACECALL_F) != 0
}
func DEBUG_TRACECALLREGS() bool {
	return (M.x86.debug & DEBUG_TRACECALL_REGS_F) != 0
}
func DEBUG_TRACEJMP() bool {
	return (M.x86.debug & DEBUG_TRACEJMP_F) != 0
}
func DEBUG_TRACEJMPREGS() bool {
	return (M.x86.debug & DEBUG_TRACEJMP_REGS_F) != 0
}
func DEBUG_SYS() bool {
	return (M.x86.debug & DEBUG_SYS_F) != 0
}
func DEBUG_MEM_TRACE() bool {
	return (M.x86.debug & DEBUG_MEM_TRACE_F) != 0
}
func DEBUG_IO_TRACE() bool {
	return (M.x86.debug & DEBUG_IO_TRACE_F) != 0
}
func DEBUG_DECODE_NOPRINT() bool {
	return (M.x86.debug & DEBUG_DECODE_NOPRINT_F) != 0
}
func initDEBUG_SYS_F() {
	DEBUG_SYS_F = (DEBUG_SVC_F | DEBUG_FS_F | DEBUG_PROC_F)
}
func sys_rdb(ip uint32) byte {
	return memory[ip]
}
func SAVE_IP_CS(cs, ip uint16) {
	if DEBUG_DECODE() || DEBUG_TRACECALL() || DEBUG_BREAK() || DEBUG_IO_TRACE() || DEBUG_SAVE_IP_CS() {
		M.x86.saved_cs = cs
		M.x86.saved_ip = ip
	}
}

func labs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func START_OF_INSTR() {
	fmt.Printf("Start instruction\n")
}
func END_OF_INSTR() {
	fmt.Printf("End instruction\n")
	x86emu_end_instr()
}

// trace_regs traces. it also returns true if the caller should immediately
// return
func TRACE_REGS() bool {
	if DEBUG_DISASSEMBLE() {
		x86emu_just_disassemble()
		return true
	}
	if DEBUG_TRACE() || DEBUG_DECODE() {
		X86EMU_trace_regs()
	}
	return false
}
func SINGLE_STEP() bool {
	if DEBUG_STEP() {
		x86emu_single_step()
		return true
	}
	return false
}
func TRACE_AND_STEP() bool {
	TRACE_REGS()
	return SINGLE_STEP()
}

func DECODE_CLEAR_SEGOVR() {
	M.x86.mode &= ^SYSMODE_CLRMASK
}
func JMP_TRACE(u, v, w, x uint16, s string) {
	if DEBUG_TRACEJMPREGS() {
		x86emu_dump_regs()
	}
	if DEBUG_TRACEJMP() {
		fmt.Printf("%04x:%04x: JMP %s%04x:%04x\n", u, v, s, w, x)
	}
}

func CALL_TRACE(u, v, w, x uint16, s string) {
	if DEBUG_TRACECALLREGS() {
		x86emu_dump_regs()
	}
	if DEBUG_TRACECALL() {
		fmt.Printf("%04x:%04x: CALL %s%04x:%04x\n", u, v, s, w, x)
	}
}

func halted() bool {
	return M.x86.intr&INTR_HALTED != 0
}

func RETURN_TRACE(u, v, w, x uint16, s string) {
	if DEBUG_TRACECALLREGS() {
		x86emu_dump_regs()
	}
	if DEBUG_TRACECALL() {
		fmt.Printf("%04x:%04x: RET %s %04x:%04x\n", u, v, s, w, x)
	}
}
func sys_inb(i uint16) uint8 {
	panic("io")
}
func sys_inw(i uint16) uint16 {
	panic("io")
}
func sys_inl(i uint16) uint32 {
	panic("io")
}

func sys_outb(i uint16, v uint8) {
	panic("io")
}
func sys_outw(i uint16, v uint16) {
	panic("io")
}
func sys_outl(i uint16, v uint32) {
	panic("io")
}
func sysw(addr uint32, i interface{}) {
	if addr > uint32(len(memory))-4 {
		log.Panicf("sysw: address %#x out of range; max is %#x", addr, len(memory))
	}
	if addr < ro {
		fx86emu_dump_regs(printer)
		log.Panicf("sysw: address %#x in ro %#x", addr, ro)
	}
	switch v := i.(type) {
	case uint32:
		memory[addr+0] = uint8(v)
		memory[addr+1] = uint8(v >> 8)
		memory[addr+2] = uint8(v >> 16)
		memory[addr+3] = uint8(v >> 24)
	case uint16:
		memory[addr+0] = uint8(v)
		memory[addr+1] = uint8(v >> 8)
	case uint8:
		memory[addr+0] = v
	default:
		log.Panicf("sysw: %T", v)
	}
}
func sysr(addr uint32, i interface{}) {
	switch v := i.(type) {
	case *uint32:
		*v = uint32(memory[addr+3])<<24 | uint32(memory[addr+2])<<16 | uint32(memory[addr+1])<<8 | uint32(memory[addr+0])
	case *uint16:
		*v = uint16(memory[addr+1])<<8 | uint16(memory[addr+0])
	case *uint8:
		*v = memory[addr]
	default:
		log.Panicf("sysr: %T", v)
	}
}
func CONDITIONAL_SET_FLAG_BOOL(cond bool, flag uint32) {
	if cond {
		SET_FLAG(flag)
		return
	}
	CLEAR_FLAG(flag)
}

func CONDITIONAL_SET_FLAG(cond interface{}, flag uint32) {
	switch v := cond.(type) {
	case bool:
		CONDITIONAL_SET_FLAG_BOOL(v, flag)
	case uint32:
		CONDITIONAL_SET_FLAG_BOOL(v != 0, flag)
	case uint16:
		CONDITIONAL_SET_FLAG_BOOL(v != 0, flag)
	case uint8:
		CONDITIONAL_SET_FLAG_BOOL(v != 0, flag)
	default:
		log.Panicf("CSF: %T", cond)
	}
}

// Mode returns true if the bits are *all* set in the mode.
func Mode(m uint32) bool {
	return M.x86.mode&m == m
}

func PC() uint32 {
	return uint32(G16(CS))<<4 + uint32(G16(IP))
}

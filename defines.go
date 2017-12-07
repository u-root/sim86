// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (VarDecl): 791: function pointers are not supported
// Warning (VarDecl): 792: function pointers are not supported

package main

type __int128_t int64
type __uint128_t uint64
type __builtin_ms_va_list []byte

var CHECK_IP_FETCH_F int = 1
var CHECK_SP_ACCESS_F int = 2
var CHECK_MEM_ACCESS_F int = 4
var CHECK_DATA_ACCESS_F int = 8

type I32_reg_t struct {
	e_reg uint32
}
type I16_reg_t struct {
	x_reg uint16
}
type I8_reg_t struct {
	l_reg uint8
	h_reg uint8
}
type i386_general_register struct {
	I32_reg I32_reg_t
	I16_reg I16_reg_t
	I8_reg  I8_reg_t
}
type i386_general_regs struct {
	A i386_general_register
	B i386_general_register
	C i386_general_register
	D i386_general_register
}

var Gen_reg_t i386_general_regs

type i386_special_regs struct {
	SP    i386_general_register
	BP    i386_general_register
	SI    i386_general_register
	DI    i386_general_register
	IP    i386_general_register
	FLAGS uint32
}
type i386_segment_regs struct {
	CS uint16
	DS uint16
	SS uint16
	ES uint16
	FS uint16
	GS uint16
}

var FB_CF int = 1
var FB_PF int = 4
var FB_AF int = 16
var FB_ZF int = 64
var FB_SF int = 128
var FB_TF int = 256
var FB_IF int = 512
var FB_DF int = 1024
var FB_OF int = 2048

const F_CF int = 1
const F_PF int = 4
const F_AF int = 16
const F_ZF int = 64
const F_SF int = 128
const F_TF int = 256
const F_IF int = 512
const F_DF int = 1024
const F_OF int = 2048

var F_PF_CALC int = 65536
var F_ZF_CALC int = 131072
var F_SF_CALC int = 262144
var F_ALL_CALC int = 16711680

const SYSMODE_SEG_DS_SS int = 1
const SYSMODE_SEGOVR_CS int = 2
const SYSMODE_SEGOVR_DS int = 4
const SYSMODE_SEGOVR_ES int = 8
const SYSMODE_SEGOVR_FS int = 16
const SYSMODE_SEGOVR_GS int = 32
const SYSMODE_SEGOVR_SS int = 64
const SYSMODE_PREFIX_REPE int = 128
const SYSMODE_PREFIX_REPNE int = 256
const SYSMODE_PREFIX_DATA int = 512
const SYSMODE_PREFIX_ADDR int = 1024
const SYSMODE_32BIT_REP int = 2048
const SYSMODE_INTR_PENDING int = 268435456
const SYSMODE_EXTRN_INTR int = 536870912
const SYSMODE_HALTED int = 1073741824

var INTR_SYNCH int = 1
var INTR_ASYNCH uint32 = uint32(2)
var INTR_HALTED uint32 = uint32(4)

type X86EMU_regs struct {
	gen         i386_general_regs
	spc         i386_special_regs
	seg         i386_segment_regs
	mode        uint32
	intr        int
	debug       uint32
	check       int
	saved_ip    uint16
	saved_cs    uint16
	enc_pos     int
	enc_str_pos int
	decode_buf  []byte
	decoded_buf []byte
	intno       uint8
	__pad       []uint8
}
type X86EMU_sysEnv struct {
	mem_base uint32
	mem_size uint32
	abseg    uint32
	private  []byte
	x86      X86EMU_regs
}

var _X86EMU_env X86EMU_sysEnv

type X86EMU_intrFuncs struct {
	f uint32
}

var _X86EMU_intrTab [][]byte
var DEBUG_DECODE_F uint32 = uint32(1)
var DEBUG_TRACE_F uint32 = uint32(2)
var DEBUG_STEP_F uint32 = uint32(4)
var DEBUG_DISASSEMBLE_F uint32 = uint32(8)
var DEBUG_BREAK_F uint32 = uint32(16)
var DEBUG_SVC_F uint32 = uint32(32)
var DEBUG_FS_F uint32 = uint32(128)
var DEBUG_PROC_F uint32 = uint32(256)
var DEBUG_SYSINT_F uint32 = uint32(512)
var DEBUG_TRACECALL_F uint32 = uint32(1024)
var DEBUG_INSTRUMENT_F uint32 = uint32(2048)
var DEBUG_MEM_TRACE_F uint32 = uint32(4096)
var DEBUG_IO_TRACE_F uint32 = uint32(8192)
var DEBUG_TRACECALL_REGS_F uint32 = uint32(16384)
var DEBUG_DECODE_NOPRINT_F uint32 = uint32(32768)
var DEBUG_SAVE_IP_CS_F uint32 = uint32(65536)
var DEBUG_TRACEJMP_F uint32 = uint32(131072)
var DEBUG_TRACEJMP_REGS_F uint32 = uint32(262144)
var DEBUG_SYS_F uint32

func initDEBUG_SYS_F() {
	DEBUG_SYS_F = (DEBUG_SVC_F | DEBUG_FS_F | DEBUG_PROC_F)
}


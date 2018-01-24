// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (VarDecl): 791: function pointers are not supported
// Warning (VarDecl): 792: function pointers are not supported

package main

const F_CF uint32 = 1
const F_PF uint32 = 4
const F_AF uint32 = 16
const F_ZF uint32 = 64
const F_SF uint32 = 128
const F_TF uint32 = 256
const F_IF uint32 = 512
const F_DF uint32 = 1024
const F_OF uint32 = 2048

const SYSMODE_SEG_DS_SS uint32 = 1
const SYSMODE_SEGOVR_CS uint32 = 2
const SYSMODE_SEGOVR_DS uint32 = 4
const SYSMODE_SEGOVR_ES uint32 = 8
const SYSMODE_SEGOVR_FS uint32 = 16
const SYSMODE_SEGOVR_GS uint32 = 32
const SYSMODE_SEGOVR_SS uint32 = 64

const SYSMODE_SEGMASK uint32 = (SYSMODE_SEG_DS_SS |
	SYSMODE_SEGOVR_CS |
	SYSMODE_SEGOVR_DS |
	SYSMODE_SEGOVR_ES |
	SYSMODE_SEGOVR_FS |
	SYSMODE_SEGOVR_GS |
	SYSMODE_SEGOVR_SS)
const SYSMODE_CLRMASK uint32 = (SYSMODE_SEG_DS_SS |
	SYSMODE_SEGOVR_CS |
	SYSMODE_SEGOVR_DS |
	SYSMODE_SEGOVR_ES |
	SYSMODE_SEGOVR_FS |
	SYSMODE_SEGOVR_GS |
	SYSMODE_SEGOVR_SS |
	SYSMODE_PREFIX_DATA |
	SYSMODE_PREFIX_ADDR |
	SYSMODE_32BIT_REP)

const SYSMODE_PREFIX_REPE uint32 = 128
const SYSMODE_PREFIX_REPNE uint32 = 256
const SYSMODE_PREFIX_DATA uint32 = 512
const SYSMODE_PREFIX_ADDR uint32 = 1024
const SYSMODE_32BIT_REP uint32 = 2048
const SYSMODE_INTR_PENDING uint32 = 268435456
const SYSMODE_EXTRN_INTR uint32 = 536870912
const SYSMODE_HALTED uint32 = 1073741824

// These can almost certainly become const at some point.
var CHECK_IP_FETCH_F int = 1
var CHECK_SP_ACCESS_F int = 2
var CHECK_MEM_ACCESS_F int = 4
var CHECK_DATA_ACCESS_F int = 8
var FB_CF uint32 = 1
var FB_PF uint32 = 4
var FB_AF uint32 = 16
var FB_ZF uint32 = 64
var FB_SF uint32 = 128
var FB_TF uint32 = 256
var FB_IF uint32 = 512
var FB_DF uint32 = 1024
var FB_OF uint32 = 2048
var F_MSK uint32 = (FB_CF | FB_PF | FB_AF | FB_ZF | FB_SF | FB_TF | FB_IF | FB_DF | FB_OF)

/* 80286 and above always have bit#1 set */
var F_ALWAYS_ON uint32 = (0x0002) /* flag bits always on */
var F_PF_CALC int = 65536
var F_ZF_CALC int = 131072
var F_SF_CALC int = 262144
var F_ALL_CALC int = 16711680
var INTR_SYNCH uint32 = 1
var INTR_ASYNCH uint32 = uint32(2)
var INTR_HALTED uint32 = uint32(4)
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

// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (FieldDecl): 300: function pointers are not supported
// Warning (FieldDecl): 301: function pointers are not supported
// Warning (FieldDecl): 302: function pointers are not supported
// Warning (FieldDecl): 303: function pointers are not supported
// Warning (FieldDecl): 304: function pointers are not supported
// Warning (FieldDecl): 305: function pointers are not supported
// Warning (FieldDecl): 308: function pointers are not supported
// Warning (FieldDecl): 309: function pointers are not supported
// Warning (FieldDecl): 310: function pointers are not supported
// Warning (FieldDecl): 311: function pointers are not supported
// Warning (FieldDecl): 312: function pointers are not supported
// Warning (FieldDecl): 313: function pointers are not supported
// Warning (TypedefDecl): 331: function pointers are not supported
// Warning (VarDecl): 430: function pointers are not supported
// Warning (VarDecl): 431: function pointers are not supported
// Warning (VarDecl): 432: function pointers are not supported
// Warning (VarDecl): 433: function pointers are not supported
// Warning (VarDecl): 434: function pointers are not supported
// Warning (VarDecl): 435: function pointers are not supported
// Warning (VarDecl): 437: function pointers are not supported
// Warning (VarDecl): 438: function pointers are not supported
// Warning (VarDecl): 439: function pointers are not supported
// Warning (VarDecl): 440: function pointers are not supported
// Warning (VarDecl): 441: function pointers are not supported
// Warning (VarDecl): 442: function pointers are not supported
// Warning (VarDecl): 446: function pointers are not supported
// Warning (VarDecl): 447: function pointers are not supported
// Warning (BinaryOperator): 993: unsigned char
// Warning (BinaryOperator): 1009: unsigned short
// Warning (BinaryOperator): 1025: unsigned int
// Warning (BinaryOperator): 2174: unsigned int
// Warning (BinaryOperator): 2178: unsigned int
// Warning (BinaryOperator): 2187: unsigned int
// Warning (CallExpr): 2509: unknown function: sys_inb
// Warning (CallExpr): 2509: probably an incorrect type translation 1
// Warning (CallExpr): 2511: unknown function: sys_inw
// Warning (CallExpr): 2511: probably an incorrect type translation 1
// Warning (CallExpr): 2513: unknown function: sys_inl
// Warning (CallExpr): 2513: probably an incorrect type translation 1
// Warning (BinaryOperator): 2521: int
// Warning (CallExpr): 2551: unknown function: sys_outb
// Warning (CallExpr): 2553: unknown function: sys_outw
// Warning (CallExpr): 2555: unknown function: sys_outl
// Warning (BinaryOperator): 2563: int
// Warning (CallExpr): 2588: unknown function: sys_rdw
// Warning (ReturnStmt): 2588: probably an incorrect type translation 1
// Warning (CallExpr): 2602: unknown function: sys_wrw
// Warning (CallExpr): 2616: unknown function: sys_wrl
// Warning (CallExpr): 2631: unknown function: sys_rdw
// Warning (BinaryOperator): 2631: probably an incorrect type translation 1
// Warning (CallExpr): 2648: unknown function: sys_rdl
// Warning (BinaryOperator): 2648: probably an incorrect type translation 1

package main

import "unsafe"
import "github.com/elliotchance/c2go/noarch"

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

type X86EMU_pioFuncs struct {
	inb interface {
	}
	inw interface {
	}
	inl interface {
	}
	outb interface {
	}
	outw interface {
	}
	outl interface {
	}
}
type X86EMU_memFuncs struct {
	rdb interface {
	}
	rdw interface {
	}
	rdl interface {
	}
	wrb interface {
	}
	wrw interface {
	}
	wrl interface {
	}
}
type X86EMU_intrFuncs interface {
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

var sys_rdb interface {
}
var sys_rdw interface {
}
var sys_rdl interface {
}
var sys_wrb interface {
}
var sys_wrw interface {
}
var sys_wrl interface {
}
var sys_inb interface {
}
var sys_inw interface {
}
var sys_inl interface {
}
var sys_outb interface {
}
var sys_outw interface {
}
var sys_outl interface {
}
var x86emu_optab interface {
}
var x86emu_optab2 interface {
}
var x86emu_parity_tab []uint32 = []uint32{2523490710, 1771476585, 1771476585, 2523490710, 1771476585, 2523490710, 2523490710, 1771476585}

func set_parity_flag(res uint32) {
	if ((x86emu_parity_tab[(res&255)/32] >> uint64(((res & 255) % 32))) & 1) == uint32(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func set_szp_flags_8(res uint8) {
	if res&128 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if res == uint8(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_parity_flag(uint32(res))
}
func set_szp_flags_16(res uint16) {
	if res&32768 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if res == uint16(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_parity_flag(uint32(res))
}
func set_szp_flags_32(res uint32) {
	if res&2147483648 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if res == uint32(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_parity_flag(res)
}
func no_carry_byte_side_eff(res uint8) {
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_8(res)
}
func no_carry_word_side_eff(res uint16) {
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_16(res)
}
func no_carry_long_side_eff(res uint32) {
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_32(res)
}
func calc_carry_chain(bits int, d uint32, s uint32, res uint32, set_carry int) {
	var cc uint32
	cc = (s&d)|((^res)&(s|d))
	if (((cc >> uint64((bits - 2))) ^ ((cc >> uint64((bits - 2))) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if cc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if set_carry != 0 {
		if res&(1<<uint64(bits)) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
}
func calc_borrow_chain(bits int, d uint32, s uint32, res uint32, set_carry int) {
	var bc uint32
	bc = (res&(^d|s))|(^d&s)
	if (((bc >> uint64((bits - 2))) ^ ((bc >> uint64((bits - 2))) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if set_carry != 0 {
		if bc&(1<<uint64((bits-1))) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
}
func aaa_word(d uint16) uint16 {
	var res uint16
	if (d&15) > 9 || (_X86EMU_env.x86.spc.FLAGS&(F_AF)) != 0 {
		d += 6
		d += 256
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	res = (d & 65295)
	set_szp_flags_16(res)
	return res
}
func aas_word(d uint16) uint16 {
	var res uint16
	if (d&15) > 9 || (_X86EMU_env.x86.spc.FLAGS&(F_AF)) != 0 {
		d -= 6
		d -= 256
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	res = (d & 65295)
	set_szp_flags_16(res)
	return res
}
func aad_word(d uint16) uint16 {
	var l uint16
	var hb uint8
	var lb uint8
	hb = uint8(((d >> uint64(8)) & 255))
	lb = uint8((d & 255))
	l = uint16(((lb + 10*hb) & 255))
	no_carry_byte_side_eff(uint8(l & 255))
	return l
}
func aam_word(d uint8) uint16 {
	var h uint16
	var l uint16
	h = uint16((d / 10))
	l = uint16((d % 10))
	l |= (h << uint64(8))
	no_carry_byte_side_eff(uint8(l & 255))
	return l
}
func adc_byte(d uint8, s uint8) uint8 {
	var res uint32
	res = uint32(d+s)
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		func() uint32 {
			res += 1
			return res
		}()
	}
	set_szp_flags_8(uint8(res))
	calc_carry_chain(8, uint32(s), uint32(d), res, 1)
	return uint8(res)
}
func adc_word(d uint16, s uint16) uint16 {
	var res uint32
	res = uint32(d+s)
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		func() uint32 {
			res += 1
			return res
		}()
	}
	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, uint32(s), uint32(d), res, 1)
	return uint16(res)
}
func adc_long(d uint32, s uint32) uint32 {
	var lo uint32
	var hi uint32
	var res uint32
	lo = (d&65535)+(s&65535)
	res = d+s
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		func() uint32 {
			lo += 1
			return lo
		}()
		func() uint32 {
			res += 1
			return res
		}()
	}
	hi = (lo>>uint64(16))+(d>>uint64(16))+(s>>uint64(16))
	set_szp_flags_32(res)
	calc_carry_chain(32, s, d, res, 0)
	if hi&65536 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func add_byte(d uint8, s uint8) uint8 {
	var res uint32
	res = uint32(d+s)
	set_szp_flags_8(uint8(res))
	calc_carry_chain(8, uint32(s), uint32(d), res, 1)
	return uint8(res)
}
func add_word(d uint16, s uint16) uint16 {
	var res uint32
	res = uint32(d+s)
	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, uint32(s), uint32(d), res, 1)
	return uint16(res)
}
func add_long(d uint32, s uint32) uint32 {
	var res uint32
	res = d+s
	set_szp_flags_32(res)
	calc_carry_chain(32, s, d, res, 0)
	if res < d || res < s {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func and_byte(d uint8, s uint8) uint8 {
	var res uint8
	res = d&s
	no_carry_byte_side_eff(res)
	return res
}
func and_word(d uint16, s uint16) uint16 {
	var res uint16
	res = d&s
	no_carry_word_side_eff(res)
	return res
}
func and_long(d uint32, s uint32) uint32 {
	var res uint32
	res = d&s
	no_carry_long_side_eff(res)
	return res
}
func cmp_byte(d uint8, s uint8) uint8 {
	var res uint32
	res = uint32(d-s)
	set_szp_flags_8(uint8(res))
	calc_borrow_chain(8, uint32(d), uint32(s), res, 1)
	return d
}
func cmp_word(d uint16, s uint16) uint16 {
	var res uint32
	res = uint32(d-s)
	set_szp_flags_16(uint16(res))
	calc_borrow_chain(16, uint32(d), uint32(s), res, 1)
	return d
}
func cmp_long(d uint32, s uint32) uint32 {
	var res uint32
	res = d-s
	set_szp_flags_32(res)
	calc_borrow_chain(32, d, s, res, 1)
	return d
}
func daa_byte(d uint8) uint8 {
	var res uint32 = uint32(d)
	if (d&15) > 9 || (_X86EMU_env.x86.spc.FLAGS&(F_AF)) != 0 {
		res += 6
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if res > 159 || (_X86EMU_env.x86.spc.FLAGS&(F_CF)) != 0 {
		res += 96
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_szp_flags_8(uint8(res))
	return uint8(res)
}
func das_byte(d uint8) uint8 {
	if (d&15) > 9 || (_X86EMU_env.x86.spc.FLAGS&(F_AF)) != 0 {
		d -= 6
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if d > 159 || (_X86EMU_env.x86.spc.FLAGS&(F_CF)) != 0 {
		d -= 96
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_szp_flags_8(d)
	return d
}
func dec_byte(d uint8) uint8 {
	var res uint32
	res = uint32(d-1)
	set_szp_flags_8(uint8(res))
	calc_borrow_chain(8, uint32(d), uint32(1), res, 0)
	return uint8(res)
}
func dec_word(d uint16) uint16 {
	var res uint32
	res = uint32(d-1)
	set_szp_flags_16(uint16(res))
	calc_borrow_chain(16, uint32(d), uint32(1), res, 0)
	return uint16(res)
}
func dec_long(d uint32) uint32 {
	var res uint32
	res = d-1
	set_szp_flags_32(res)
	calc_borrow_chain(32, d, uint32(1), res, 0)
	return res
}
func inc_byte(d uint8) uint8 {
	var res uint32
	res = uint32(d+1)
	set_szp_flags_8(uint8(res))
	calc_carry_chain(8, uint32(d), uint32(1), res, 0)
	return uint8(res)
}
func inc_word(d uint16) uint16 {
	var res uint32
	res = uint32(d+1)
	set_szp_flags_16(uint16(res))
	calc_carry_chain(16, uint32(d), uint32(1), res, 0)
	return uint16(res)
}
func inc_long(d uint32) uint32 {
	var res uint32
	res = d+1
	set_szp_flags_32(res)
	calc_carry_chain(32, d, uint32(1), res, 0)
	return res
}
func or_byte(d uint8, s uint8) uint8 {
	var res uint8
	res = d|s
	no_carry_byte_side_eff(res)
	return res
}
func or_word(d uint16, s uint16) uint16 {
	var res uint16
	res = d|s
	no_carry_word_side_eff(res)
	return res
}
func or_long(d uint32, s uint32) uint32 {
	var res uint32
	res = d|s
	no_carry_long_side_eff(res)
	return res
}
func neg_byte(s uint8) uint8 {
	var res uint8
	if s != uint8(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	res = -s
	set_szp_flags_8(res)
	calc_borrow_chain(8, uint32(0), uint32(s), uint32(res), 0)
	return res
}
func neg_word(s uint16) uint16 {
	var res uint16
	if s != uint16(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	res = -s
	set_szp_flags_16(res)
	calc_borrow_chain(16, uint32(0), uint32(s), uint32(res), 0)
	return res
}
func neg_long(s uint32) uint32 {
	var res uint32
	if s != uint32(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	res = -s
	set_szp_flags_32(res)
	calc_borrow_chain(32, uint32(0), s, res, 0)
	return res
}
func not_byte(s uint8) uint8 {
	return ^s
}
func not_word(s uint16) uint16 {
	return ^s
}
func not_long(s uint32) uint32 {
	return ^s
}
func rcl_byte(d uint8, s uint8) uint8 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%9)
		return cnt
	}()) != uint32(0) {
		cf = uint32((d>>uint64((8-cnt)))&1)
		res = uint32((d<<uint64(cnt))&255)
		mask = uint32((1<<uint64((cnt-1)))-1)
		res |= (d>>uint64((9-cnt)))&mask
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((cnt-1))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) && (((cf+((res>>uint64(6))&2))^((cf+((res>>uint64(6))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return uint8(res)
}
func rcl_word(d uint16, s uint8) uint16 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%17)
		return cnt
	}()) != uint32(0) {
		cf = uint32((d>>uint64((16-cnt)))&1)
		res = uint32((d<<uint64(cnt))&65535)
		mask = uint32((1<<uint64((cnt-1)))-1)
		res |= (d>>uint64((17-cnt)))&mask
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((cnt-1))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) && (((cf+((res>>uint64(14))&2))^((cf+((res>>uint64(14))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return uint16(res)
}
func rcl_long(d uint32, s uint8) uint32 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	res = d
	if (func() uint32 {
		cnt = uint32(s%33)
		return cnt
	}()) != uint32(0) {
		cf = (d>>uint64((32-cnt)))&1
		res = (d<<uint64(cnt))&4294967295
		mask = uint32((1<<uint64((cnt-1)))-1)
		res |= (d>>uint64((33-cnt)))&mask
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((cnt-1))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) && (((cf+((res>>uint64(30))&2))^((cf+((res>>uint64(30))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return res
}
func rcr_byte(d uint8, s uint8) uint8 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	var ocf uint32 = uint32(0)
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%9)
		return cnt
	}()) != uint32(0) {
		if cnt == uint32(1) {
			cf = uint32(d&1)
			ocf = noarch.BoolToUint32((_X86EMU_env.x86.spc.FLAGS&(F_CF)) != uint32(0))
		} else {
			cf = uint32((d>>uint64((cnt-1)))&1)
		}
		mask = uint32((1<<uint64((8-cnt)))-1)
		res = uint32((d>>uint64(cnt))&mask)
		res |= (d << uint64((9 - cnt)))
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((8-cnt))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) {
			if (((ocf + ((d >> uint64(6)) & 2)) ^ ((ocf + ((d >> uint64(6)) & 2)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint8(res)
}
func rcr_word(d uint16, s uint8) uint16 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	var ocf uint32 = uint32(0)
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%17)
		return cnt
	}()) != uint32(0) {
		if cnt == uint32(1) {
			cf = uint32(d&1)
			ocf = noarch.BoolToUint32((_X86EMU_env.x86.spc.FLAGS&(F_CF)) != uint32(0))
		} else {
			cf = uint32((d>>uint64((cnt-1)))&1)
		}
		mask = uint32((1<<uint64((16-cnt)))-1)
		res = uint32((d>>uint64(cnt))&mask)
		res |= (d << uint64((17 - cnt)))
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((16-cnt))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) {
			if (((ocf + ((d >> uint64(14)) & 2)) ^ ((ocf + ((d >> uint64(14)) & 2)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint16(res)
}
func rcr_long(d uint32, s uint8) uint32 {
	var res uint32
	var cnt uint32
	var mask uint32
	var cf uint32
	var ocf uint32 = uint32(0)
	res = d
	if (func() uint32 {
		cnt = uint32(s%33)
		return cnt
	}()) != uint32(0) {
		if cnt == uint32(1) {
			cf = d&1
			ocf = noarch.BoolToUint32((_X86EMU_env.x86.spc.FLAGS&(F_CF)) != uint32(0))
		} else {
			cf = (d>>uint64((cnt-1)))&1
		}
		mask = uint32((1<<uint64((32-cnt)))-1)
		res = (d>>uint64(cnt))&mask
		if cnt != uint32(1) {
			res |= (d << uint64((33 - cnt)))
		}
		if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
			res |= 1<<uint64((32-cnt))
		}
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if cnt == uint32(1) {
			if (((ocf + ((d >> uint64(30)) & 2)) ^ ((ocf + ((d >> uint64(30)) & 2)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return res
}
func rol_byte(d uint8, s uint8) uint8 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%8)
		return cnt
	}()) != uint32(0) {
		res = uint32((d << uint64(cnt)))
		mask = uint32((1<<uint64(cnt))-1)
		res |= (d>>uint64((8-cnt)))&mask
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && ((((res&1)+((res>>uint64(6))&2))^(((res&1)+((res>>uint64(6))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	if s != uint8(0) {
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return uint8(res)
}
func rol_word(d uint16, s uint8) uint16 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%16)
		return cnt
	}()) != uint32(0) {
		res = uint32((d << uint64(cnt)))
		mask = uint32((1<<uint64(cnt))-1)
		res |= (d>>uint64((16-cnt)))&mask
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && ((((res&1)+((res>>uint64(14))&2))^(((res&1)+((res>>uint64(14))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	if s != uint8(0) {
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return uint16(res)
}
func rol_long(d uint32, s uint8) uint32 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = d
	if (func() uint32 {
		cnt = uint32(s%32)
		return cnt
	}()) != uint32(0) {
		res = (d << uint64(cnt))
		mask = uint32((1<<uint64(cnt))-1)
		res |= (d>>uint64((32-cnt)))&mask
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && ((((res&1)+((res>>uint64(30))&2))^(((res&1)+((res>>uint64(30))&2))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	if s != uint8(0) {
		if res&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	}
	return res
}
func ror_byte(d uint8, s uint8) uint8 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%8)
		return cnt
	}()) != uint32(0) {
		res = uint32((d << uint64((8 - cnt))))
		mask = uint32((1<<uint64((8-cnt)))-1)
		res |= (d>>uint64((cnt)))&mask
		if res&128 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && (((res>>uint64(6))^((res>>uint64(6))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		if s != uint8(0) {
			if res&128 != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint8(res)
}
func ror_word(d uint16, s uint8) uint16 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = uint32(d)
	if (func() uint32 {
		cnt = uint32(s%16)
		return cnt
	}()) != uint32(0) {
		res = uint32((d << uint64((16 - cnt))))
		mask = uint32((1<<uint64((16-cnt)))-1)
		res |= (d>>uint64((cnt)))&mask
		if res&32768 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && (((res>>uint64(14))^((res>>uint64(14))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		if s != uint8(0) {
			if res&32768 != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint16(res)
}
func ror_long(d uint32, s uint8) uint32 {
	var res uint32
	var cnt uint32
	var mask uint32
	res = d
	if (func() uint32 {
		cnt = uint32(s%32)
		return cnt
	}()) != uint32(0) {
		res = (d << uint64((32 - cnt)))
		mask = uint32((1<<uint64((32-cnt)))-1)
		res |= (d>>uint64((cnt)))&mask
		if res&2147483648 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if s == uint8(1) && (((res>>uint64(30))^((res>>uint64(30))>>uint64(1)))&1) != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		if s != uint8(0) {
			if res&2147483648 != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return res
}
func shl_byte(d uint8, s uint8) uint8 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 8 {
		cnt = uint32(s%8)
		if cnt > 0 {
			res = uint32(d<<uint64(cnt))
			cf = uint32(d&(1<<uint64((8-cnt))))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_8(uint8(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if ((res & 128) == uint32(128)) ^ ((_X86EMU_env.x86.spc.FLAGS & (F_CF)) != uint32(0)) {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d<<uint64((s-1)))&128 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint8(res)
}
func shl_word(d uint16, s uint8) uint16 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 16 {
		cnt = uint32(s%16)
		if cnt > 0 {
			res = uint32(d<<uint64(cnt))
			cf = uint32(d&(1<<uint64((16-cnt))))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_16(uint16(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if ((res & 32768) == uint32(32768)) ^ ((_X86EMU_env.x86.spc.FLAGS & (F_CF)) != uint32(0)) {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d<<uint64((s-1)))&32768 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func shl_long(d uint32, s uint8) uint32 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 32 {
		cnt = uint32(s%32)
		if cnt > 0 {
			res = d<<uint64(cnt)
			cf = d&(1<<uint64((32-cnt)))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_32(res)
		} else {
			res = d
		}
		if cnt == uint32(1) {
			if ((res & 2147483648) == uint32(2147483648)) ^ ((_X86EMU_env.x86.spc.FLAGS & (F_CF)) != uint32(0)) {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d<<uint64((s-1)))&2147483648 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func shr_byte(d uint8, s uint8) uint8 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 8 {
		cnt = uint32(s%8)
		if cnt > 0 {
			cf = uint32(d&(1<<uint64((cnt-1))))
			res = uint32(d>>uint64(cnt))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_8(uint8(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if (((res >> uint64(6)) ^ ((res >> uint64(6)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d>>uint64((s-1)))&1 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint8(res)
}
func shr_word(d uint16, s uint8) uint16 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 16 {
		cnt = uint32(s%16)
		if cnt > 0 {
			cf = uint32(d&(1<<uint64((cnt-1))))
			res = uint32(d>>uint64(cnt))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_16(uint16(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if (((res >> uint64(14)) ^ ((res >> uint64(14)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func shr_long(d uint32, s uint8) uint32 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 32 {
		cnt = uint32(s%32)
		if cnt > 0 {
			cf = d&(1<<uint64((cnt-1)))
			res = d>>uint64(cnt)
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_32(res)
		} else {
			res = d
		}
		if cnt == uint32(1) {
			if (((res >> uint64(30)) ^ ((res >> uint64(30)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func sar_byte(d uint8, s uint8) uint8 {
	var cnt uint32
	var res uint32
	var cf uint32
	var mask uint32
	var sf uint32
	res = uint32(d)
	sf = uint32(d&128)
	cnt = uint32(s%8)
	if cnt > 0 && cnt < 8 {
		mask = uint32((1<<uint64((8-cnt)))-1)
		cf = uint32(d&(1<<uint64((cnt-1))))
		res = uint32((d>>uint64(cnt))&mask)
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if sf != 0 {
			res |= ^mask
		}
		set_szp_flags_8(uint8(res))
	} else {
		if cnt >= 8 {
			if sf != 0 {
				res = uint32(255)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				res = uint32(0)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint8(res)
}
func sar_word(d uint16, s uint8) uint16 {
	var cnt uint32
	var res uint32
	var cf uint32
	var mask uint32
	var sf uint32
	sf = uint32(d&32768)
	cnt = uint32(s%16)
	res = uint32(d)
	if cnt > 0 && cnt < 16 {
		mask = uint32((1<<uint64((16-cnt)))-1)
		cf = uint32(d&(1<<uint64((cnt-1))))
		res = uint32((d>>uint64(cnt))&mask)
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if sf != 0 {
			res |= ^mask
		}
		set_szp_flags_16(uint16(res))
	} else {
		if cnt >= 16 {
			if sf != 0 {
				res = uint32(65535)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				res = uint32(0)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return uint16(res)
}
func sar_long(d uint32, s uint8) uint32 {
	var cnt uint32
	var res uint32
	var cf uint32
	var mask uint32
	var sf uint32
	sf = d&2147483648
	cnt = uint32(s%32)
	res = d
	if cnt > 0 && cnt < 32 {
		mask = uint32((1<<uint64((32-cnt)))-1)
		cf = d&(1<<uint64((cnt-1)))
		res = (d>>uint64(cnt))&mask
		if cf != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		if sf != 0 {
			res |= ^mask
		}
		set_szp_flags_32(res)
	} else {
		if cnt >= 32 {
			if sf != 0 {
				res = uint32(4294967295)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				res = uint32(0)
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		}
	}
	return res
}
func shld_word(d uint16, fill uint16, s uint8) uint16 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 16 {
		cnt = uint32(s%16)
		if cnt > 0 {
			res = uint32((d<<uint64(cnt))|(fill>>uint64((16-cnt))))
			cf = uint32(d&(1<<uint64((16-cnt))))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_16(uint16(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if ((res & 32768) == uint32(32768)) ^ ((_X86EMU_env.x86.spc.FLAGS & (F_CF)) != uint32(0)) {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d<<uint64((s-1)))&32768 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func shld_long(d uint32, fill uint32, s uint8) uint32 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 32 {
		cnt = uint32(s%32)
		if cnt > 0 {
			res = (d<<uint64(cnt))|(fill>>uint64((32-cnt)))
			cf = d&(1<<uint64((32-cnt)))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_32(res)
		} else {
			res = d
		}
		if cnt == uint32(1) {
			if ((res & 2147483648) == uint32(2147483648)) ^ ((_X86EMU_env.x86.spc.FLAGS & (F_CF)) != uint32(0)) {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		if (d<<uint64((s-1)))&2147483648 != 0 {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS |= (F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func shrd_word(d uint16, fill uint16, s uint8) uint16 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 16 {
		cnt = uint32(s%16)
		if cnt > 0 {
			cf = uint32(d&(1<<uint64((cnt-1))))
			res = uint32((d>>uint64(cnt))|(fill<<uint64((16-cnt))))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_16(uint16(res))
		} else {
			res = uint32(d)
		}
		if cnt == uint32(1) {
			if (((res >> uint64(14)) ^ ((res >> uint64(14)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func shrd_long(d uint32, fill uint32, s uint8) uint32 {
	var cnt uint32
	var res uint32
	var cf uint32
	if s < 32 {
		cnt = uint32(s%32)
		if cnt > 0 {
			cf = d&(1<<uint64((cnt-1)))
			res = (d>>uint64(cnt))|(fill<<uint64((32-cnt)))
			if cf != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
			set_szp_flags_32(res)
		} else {
			res = d
		}
		if cnt == uint32(1) {
			if (((res >> uint64(30)) ^ ((res >> uint64(30)) >> uint64(1))) & 1) != 0 {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS |= (F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			} else {
				(func() uint32 {
					_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
					return _X86EMU_env.x86.spc.FLAGS
				}())
			}
		} else {
			(func() uint32 {
				_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
				return _X86EMU_env.x86.spc.FLAGS
			}())
		}
	} else {
		res = uint32(0)
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_PF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func sbb_byte(d uint8, s uint8) uint8 {
	var res uint32
	var bc uint32
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		res = uint32(d-s-1)
	} else {
		res = uint32(d-s)
	}
	set_szp_flags_8(uint8(res))
	bc = (res&(^d|s))|(^d&s)
	if bc&128 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(6)) ^ ((bc >> uint64(6)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint8(res)
}
func sbb_word(d uint16, s uint16) uint16 {
	var res uint32
	var bc uint32
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		res = uint32(d-s-1)
	} else {
		res = uint32(d-s)
	}
	set_szp_flags_16(uint16(res))
	bc = (res&(^d|s))|(^d&s)
	if bc&32768 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(14)) ^ ((bc >> uint64(14)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func sbb_long(d uint32, s uint32) uint32 {
	var res uint32
	var bc uint32
	if (_X86EMU_env.x86.spc.FLAGS & (F_CF)) != 0 {
		res = d-s-1
	} else {
		res = d-s
	}
	set_szp_flags_32(res)
	bc = (res&(^d|s))|(^d&s)
	if bc&2147483648 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(30)) ^ ((bc >> uint64(30)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func sub_byte(d uint8, s uint8) uint8 {
	var res uint32
	var bc uint32
	res = uint32(d-s)
	set_szp_flags_8(uint8(res))
	bc = (res&(^d|s))|(^d&s)
	if bc&128 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(6)) ^ ((bc >> uint64(6)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint8(res)
}
func sub_word(d uint16, s uint16) uint16 {
	var res uint32
	var bc uint32
	res = uint32(d-s)
	set_szp_flags_16(uint16(res))
	bc = (res&(^d|s))|(^d&s)
	if bc&32768 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(14)) ^ ((bc >> uint64(14)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return uint16(res)
}
func sub_long(d uint32, s uint32) uint32 {
	var res uint32
	var bc uint32
	res = d-s
	set_szp_flags_32(res)
	bc = (res&(^d|s))|(^d&s)
	if bc&2147483648 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if (((bc >> uint64(30)) ^ ((bc >> uint64(30)) >> uint64(1))) & 1) != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	if bc&8 != 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	return res
}
func test_byte(d uint8, s uint8) {
	var res uint32
	res = uint32(d&s)
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_8(uint8(res))
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
}
func test_word(d uint16, s uint16) {
	var res uint32
	res = uint32(d&s)
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_16(uint16(res))
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
}
func test_long(d uint32, s uint32) {
	var res uint32
	res = d&s
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_szp_flags_32(res)
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
}
func xor_byte(d uint8, s uint8) uint8 {
	var res uint8
	res = d^s
	no_carry_byte_side_eff(res)
	return res
}
func xor_word(d uint16, s uint16) uint16 {
	var res uint16
	res = d^s
	no_carry_word_side_eff(res)
	return res
}
func xor_long(d uint32, s uint32) uint32 {
	var res uint32
	res = d^s
	no_carry_long_side_eff(res)
	return res
}
func imul_byte(s uint8) {
	var res int16 = int16((_X86EMU_env.x86.gen.A.I8_reg.l_reg * s))
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = uint16(res)
	if ((_X86EMU_env.x86.gen.A.I8_reg.l_reg&128) == uint8(0) && _X86EMU_env.x86.gen.A.I8_reg.h_reg == uint8(0)) || ((_X86EMU_env.x86.gen.A.I8_reg.l_reg&128) != uint8(0) && _X86EMU_env.x86.gen.A.I8_reg.h_reg == uint8(255)) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func imul_word(s uint16) {
	var res int = int(_X86EMU_env.x86.gen.A.I16_reg.x_reg * s)
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = uint16(res)
	_X86EMU_env.x86.gen.D.I16_reg.x_reg = uint16((res >> uint64(16)))
	if ((_X86EMU_env.x86.gen.A.I16_reg.x_reg&32768) == uint16(0) && _X86EMU_env.x86.gen.D.I16_reg.x_reg == uint16(0)) || ((_X86EMU_env.x86.gen.A.I16_reg.x_reg&32768) != uint16(0) && _X86EMU_env.x86.gen.D.I16_reg.x_reg == uint16(65535)) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func imul_long_direct(res_lo []uint32, res_hi []uint32, d uint32, s uint32) {
	var d_lo uint32
	var d_hi uint32
	var d_sign uint32
	var s_lo uint32
	var s_hi uint32
	var s_sign uint32
	var rlo_lo uint32
	var rlo_hi uint32
	var rhi_lo uint32
	if (func() uint32 {
		d_sign = d&2147483648
		return d_sign
	}()) != uint32(0) {
		d = -d
	}
	d_lo = d&65535
	d_hi = d>>uint64(16)
	if (func() uint32 {
		s_sign = s&2147483648
		return s_sign
	}()) != uint32(0) {
		s = -s
	}
	s_lo = s&65535
	s_hi = s>>uint64(16)
	rlo_lo = d_lo*s_lo
	rlo_hi = (d_hi*s_lo+d_lo*s_hi)+(rlo_lo>>uint64(16))
	rhi_lo = d_hi*s_hi+(rlo_hi>>uint64(16))
	res_lo[0] = (rlo_hi<<uint64(16))|(rlo_lo&65535)
	res_hi[0] = rhi_lo
	if d_sign != s_sign {
		d = ^res_lo[0]
		s = (((d&65535)+1)>>uint64(16))+(d>>uint64(16))
		res_lo[0] = ^res_lo[0]+1
		res_hi[0] = ^res_hi[0]+(s>>uint64(16))
	}
}
func imul_long(s uint32) {
	imul_long_direct((*[1]uint32)(unsafe.Pointer(&_X86EMU_env.x86.gen.A.I32_reg.e_reg))[:], (*[1]uint32)(unsafe.Pointer(&_X86EMU_env.x86.gen.D.I32_reg.e_reg))[:], _X86EMU_env.x86.gen.A.I32_reg.e_reg, s)
	if ((_X86EMU_env.x86.gen.A.I32_reg.e_reg&2147483648) == uint32(0) && _X86EMU_env.x86.gen.D.I32_reg.e_reg == uint32(0)) || ((_X86EMU_env.x86.gen.A.I32_reg.e_reg&2147483648) != uint32(0) && _X86EMU_env.x86.gen.D.I32_reg.e_reg == uint32(4294967295)) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func mul_byte(s uint8) {
	var res uint16 = uint16((_X86EMU_env.x86.gen.A.I8_reg.l_reg * s))
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = res
	if _X86EMU_env.x86.gen.A.I8_reg.h_reg == uint8(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func mul_word(s uint16) {
	var res uint32 = uint32(_X86EMU_env.x86.gen.A.I16_reg.x_reg * s)
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = uint16(res)
	_X86EMU_env.x86.gen.D.I16_reg.x_reg = uint16((res >> uint64(16)))
	if _X86EMU_env.x86.gen.D.I16_reg.x_reg == uint16(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func mul_long(s uint32) {
	var a uint32
	var a_lo uint32
	var a_hi uint32
	var s_lo uint32
	var s_hi uint32
	var rlo_lo uint32
	var rlo_hi uint32
	var rhi_lo uint32
	a = _X86EMU_env.x86.gen.A.I32_reg.e_reg
	a_lo = a&65535
	a_hi = a>>uint64(16)
	s_lo = s&65535
	s_hi = s>>uint64(16)
	rlo_lo = a_lo*s_lo
	rlo_hi = (a_hi*s_lo+a_lo*s_hi)+(rlo_lo>>uint64(16))
	rhi_lo = a_hi*s_hi+(rlo_hi>>uint64(16))
	_X86EMU_env.x86.gen.A.I32_reg.e_reg = (rlo_hi<<uint64(16))|(rlo_lo&65535)
	_X86EMU_env.x86.gen.D.I32_reg.e_reg = rhi_lo
	if _X86EMU_env.x86.gen.D.I32_reg.e_reg == uint32(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_CF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_OF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
}
func idiv_byte(s uint8) {
	var dvd int
	var div int
	var mod int
	dvd = int(_X86EMU_env.x86.gen.A.I16_reg.x_reg)
	if s == uint8(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	div = dvd/s
	mod = dvd%s
	if func() int {
		var __x int = (div)
		return func() int {
			if __x < 0 {
				return -__x
			} else {
				return __x
			}
		}()
	}() > 127 {
		x86emu_intr_raise(uint8(0))
		return
	}
	_X86EMU_env.x86.gen.A.I8_reg.l_reg = uint8(div)
	_X86EMU_env.x86.gen.A.I8_reg.h_reg = uint8(mod)
}
func idiv_word(s uint16) {
	var dvd int
	var div int
	var mod int
	dvd = int(((_X86EMU_env.x86.gen.D.I16_reg.x_reg)<<uint64(16))|_X86EMU_env.x86.gen.A.I16_reg.x_reg)
	if s == uint16(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	div = dvd/s
	mod = dvd%s
	if func() int {
		var __x int = (div)
		return func() int {
			if __x < 0 {
				return -__x
			} else {
				return __x
			}
		}()
	}() > 32767 {
		x86emu_intr_raise(uint8(0))
		return
	}
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	if div == 0 {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_parity_flag(uint32(mod))
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = uint16(div)
	_X86EMU_env.x86.gen.D.I16_reg.x_reg = uint16(mod)
}
func idiv_long(s uint32) {
	var div int = 0
	var mod int
	var h_dvd int = int(_X86EMU_env.x86.gen.D.I32_reg.e_reg)
	var l_dvd uint32 = _X86EMU_env.x86.gen.A.I32_reg.e_reg
	var abs_s uint32 = s & 2147483647
	var abs_h_dvd uint32 = uint32(h_dvd & 2147483647)
	var h_s uint32 = abs_s >> uint64(1)
	var l_s uint32 = abs_s << uint64(31)
	var counter int = 31
	var carry int
	if s == uint32(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	for {
		div <<= uint64(1)
		carry = func() int {
			if l_dvd >= l_s {
				return 0
			} else {
				return 1
			}
		}()
		if abs_h_dvd < (h_s + carry) {
			h_s >>= uint64(1)
			l_s = abs_s<<uint64((func() int {
				counter -= 1
				return counter
			}()))
			continue
		} else {
			abs_h_dvd -= (h_s + carry)
			l_dvd = func() uint32 {
				if carry != 0 {
					return uint32(((4294967295 - l_s) + l_dvd + 1))
				} else {
					return (l_dvd - l_s)
				}
			}()
			h_s >>= uint64(1)
			l_s = abs_s<<uint64((func() int {
				counter -= 1
				return counter
			}()))
			div |= 1
			continue
		}
		if !(counter > -1) {
			break
		}
	}
	if abs_h_dvd != 0 || (l_dvd > abs_s) {
		x86emu_intr_raise(uint8(0))
		return
	}
	div |= ((h_dvd & 268435456) ^ (s & 268435456))
	mod = int(l_dvd)
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_parity_flag(uint32(mod))
	_X86EMU_env.x86.gen.A.I32_reg.e_reg = uint32(div)
	_X86EMU_env.x86.gen.D.I32_reg.e_reg = uint32(mod)
}
func div_byte(s uint8) {
	var dvd uint32
	var div uint32
	var mod uint32
	dvd = uint32(_X86EMU_env.x86.gen.A.I16_reg.x_reg)
	if s == uint8(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	div = dvd/s
	mod = dvd%s
	if func() int {
		var __x int = int((div))
		return func() int {
			if __x < 0 {
				return -__x
			} else {
				return __x
			}
		}()
	}() > 255 {
		x86emu_intr_raise(uint8(0))
		return
	}
	_X86EMU_env.x86.gen.A.I8_reg.l_reg = uint8(div)
	_X86EMU_env.x86.gen.A.I8_reg.h_reg = uint8(mod)
}
func div_word(s uint16) {
	var dvd uint32
	var div uint32
	var mod uint32
	dvd = uint32(((_X86EMU_env.x86.gen.D.I16_reg.x_reg)<<uint64(16))|_X86EMU_env.x86.gen.A.I16_reg.x_reg)
	if s == uint16(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	div = dvd/s
	mod = dvd%s
	if func() int {
		var __x int = int((div))
		return func() int {
			if __x < 0 {
				return -__x
			} else {
				return __x
			}
		}()
	}() > 65535 {
		x86emu_intr_raise(uint8(0))
		return
	}
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	if div == uint32(0) {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	} else {
		(func() uint32 {
			_X86EMU_env.x86.spc.FLAGS &= ^(F_ZF)
			return _X86EMU_env.x86.spc.FLAGS
		}())
	}
	set_parity_flag(mod)
	_X86EMU_env.x86.gen.A.I16_reg.x_reg = uint16(div)
	_X86EMU_env.x86.gen.D.I16_reg.x_reg = uint16(mod)
}
func div_long(s uint32) {
	var div int = 0
	var mod int
	var h_dvd int = int(_X86EMU_env.x86.gen.D.I32_reg.e_reg)
	var l_dvd uint32 = _X86EMU_env.x86.gen.A.I32_reg.e_reg
	var h_s uint32 = s
	var l_s uint32 = uint32(0)
	var counter int = 32
	var carry int
	if s == uint32(0) {
		x86emu_intr_raise(uint8(0))
		return
	}
	for {
		div <<= uint64(1)
		carry = func() int {
			if l_dvd >= l_s {
				return 0
			} else {
				return 1
			}
		}()
		if h_dvd < (h_s + carry) {
			h_s >>= uint64(1)
			l_s = s<<uint64((func() int {
				counter -= 1
				return counter
			}()))
			continue
		} else {
			h_dvd -= (h_s + carry)
			l_dvd = func() uint32 {
				if carry != 0 {
					return uint32(((4294967295 - l_s) + l_dvd + 1))
				} else {
					return (l_dvd - l_s)
				}
			}()
			h_s >>= uint64(1)
			l_s = s<<uint64((func() int {
				counter -= 1
				return counter
			}()))
			div |= 1
			continue
		}
		if !(counter > -1) {
			break
		}
	}
	if h_dvd != 0 || (l_dvd > s) {
		x86emu_intr_raise(uint8(0))
		return
	}
	mod = int(l_dvd)
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_CF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_AF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS &= ^(F_SF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	(func() uint32 {
		_X86EMU_env.x86.spc.FLAGS |= (F_ZF)
		return _X86EMU_env.x86.spc.FLAGS
	}())
	set_parity_flag(uint32(mod))
	_X86EMU_env.x86.gen.A.I32_reg.e_reg = uint32(div)
	_X86EMU_env.x86.gen.D.I32_reg.e_reg = uint32(mod)
}
func single_in(size int) {
	if size == 1 {
		store_data_byte_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.DI.I16_reg.x_reg), nil)
	} else {
		if size == 2 {
			store_data_word_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.DI.I16_reg.x_reg), nil)
		} else {
			store_data_long_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.DI.I16_reg.x_reg), nil)
		}
	}
}
func ins(size int) {
	var inc int = size
	if (_X86EMU_env.x86.spc.FLAGS & (F_DF)) != 0 {
		inc = -size
	}
	if _X86EMU_env.x86.mode&(SYSMODE_PREFIX_REPE|SYSMODE_PREFIX_REPNE) != 0 {
		var count uint32 = uint32((func() uint16 {
			if (_X86EMU_env.x86.mode & SYSMODE_32BIT_REP) != 0 {
				return uint16(_X86EMU_env.x86.gen.C.I32_reg.e_reg)
			} else {
				return _X86EMU_env.x86.gen.C.I16_reg.x_reg
			}
		}()))
		for func() uint32 {
			count -= 1
			return count
		}() != 0 {
			single_in(size)
			_X86EMU_env.x86.spc.DI.I16_reg.x_reg += inc
		}
		_X86EMU_env.x86.gen.C.I16_reg.x_reg = uint16(0)
		if _X86EMU_env.x86.mode&SYSMODE_32BIT_REP != 0 {
			_X86EMU_env.x86.gen.C.I32_reg.e_reg = uint32(0)
		}
		_X86EMU_env.x86.mode &= ^(SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE)
	} else {
		single_in(size)
		_X86EMU_env.x86.spc.DI.I16_reg.x_reg += inc
	}
}
func single_out(size int) {
	if size == 1 {
		sys_outb(_X86EMU_env.x86.gen.D.I16_reg.x_reg, fetch_data_byte_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.SI.I16_reg.x_reg)))
	} else {
		if size == 2 {
			sys_outw(_X86EMU_env.x86.gen.D.I16_reg.x_reg, fetch_data_word_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.SI.I16_reg.x_reg)))
		} else {
			sys_outl(_X86EMU_env.x86.gen.D.I16_reg.x_reg, fetch_data_long_abs(uint32(_X86EMU_env.x86.seg.ES), uint32(_X86EMU_env.x86.spc.SI.I16_reg.x_reg)))
		}
	}
}
func outs(size int) {
	var inc int = size
	if (_X86EMU_env.x86.spc.FLAGS & (F_DF)) != 0 {
		inc = -size
	}
	if _X86EMU_env.x86.mode&(SYSMODE_PREFIX_REPE|SYSMODE_PREFIX_REPNE) != 0 {
		var count uint32 = uint32((func() uint16 {
			if (_X86EMU_env.x86.mode & SYSMODE_32BIT_REP) != 0 {
				return uint16(_X86EMU_env.x86.gen.C.I32_reg.e_reg)
			} else {
				return _X86EMU_env.x86.gen.C.I16_reg.x_reg
			}
		}()))
		for func() uint32 {
			count -= 1
			return count
		}() != 0 {
			single_out(size)
			_X86EMU_env.x86.spc.SI.I16_reg.x_reg += inc
		}
		_X86EMU_env.x86.gen.C.I16_reg.x_reg = uint16(0)
		if _X86EMU_env.x86.mode&SYSMODE_32BIT_REP != 0 {
			_X86EMU_env.x86.gen.C.I32_reg.e_reg = uint32(0)
		}
		_X86EMU_env.x86.mode &= ^(SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE)
	} else {
		single_out(size)
		_X86EMU_env.x86.spc.SI.I16_reg.x_reg += inc
	}
}
func mem_access_word(addr int) uint16 {
	if (_X86EMU_env.x86.check & CHECK_MEM_ACCESS_F) != 0 {
		x86emu_check_mem_access(uint32(addr))
	}
	return nil
}
func push_word(w uint16) {
	if (_X86EMU_env.x86.check & CHECK_SP_ACCESS_F) != 0 {
		x86emu_check_sp_access()
	}
	_X86EMU_env.x86.spc.SP.I16_reg.x_reg -= 2
	sys_wrw((_X86EMU_env.x86.seg.SS<<uint64(4))+_X86EMU_env.x86.spc.SP.I16_reg.x_reg, w)
}
func push_long(w uint32) {
	if (_X86EMU_env.x86.check & CHECK_SP_ACCESS_F) != 0 {
		x86emu_check_sp_access()
	}
	_X86EMU_env.x86.spc.SP.I16_reg.x_reg -= 4
	sys_wrl((_X86EMU_env.x86.seg.SS<<uint64(4))+_X86EMU_env.x86.spc.SP.I16_reg.x_reg, w)
}
func pop_word() uint16 {
	var res uint16
	if (_X86EMU_env.x86.check & CHECK_SP_ACCESS_F) != 0 {
		x86emu_check_sp_access()
	}
	res = sys_rdw((_X86EMU_env.x86.seg.SS<<uint64(4))+_X86EMU_env.x86.spc.SP.I16_reg.x_reg)
	_X86EMU_env.x86.spc.SP.I16_reg.x_reg += 2
	return res
}
func pop_long() uint32 {
	var res uint32
	if (_X86EMU_env.x86.check & CHECK_SP_ACCESS_F) != 0 {
		x86emu_check_sp_access()
	}
	res = sys_rdl((_X86EMU_env.x86.seg.SS<<uint64(4))+_X86EMU_env.x86.spc.SP.I16_reg.x_reg)
	_X86EMU_env.x86.spc.SP.I16_reg.x_reg += 4
	return res
}
func x86emu_cpuid() {
	var feature uint32 = _X86EMU_env.x86.gen.A.I32_reg.e_reg
	switch feature {
	case 0:
		_X86EMU_env.x86.gen.A.I32_reg.e_reg = uint32(1)
		_X86EMU_env.x86.gen.B.I32_reg.e_reg = uint32(1970169159)
		_X86EMU_env.x86.gen.D.I32_reg.e_reg = uint32(1231384169)
		_X86EMU_env.x86.gen.C.I32_reg.e_reg = uint32(1818588270)
	case 1:
		_X86EMU_env.x86.gen.A.I32_reg.e_reg = uint32(1152)
		_X86EMU_env.x86.gen.B.I32_reg.e_reg = uint32(0)
		_X86EMU_env.x86.gen.C.I32_reg.e_reg = uint32(0)
		_X86EMU_env.x86.gen.D.I32_reg.e_reg = uint32(2)
		_X86EMU_env.x86.gen.D.I32_reg.e_reg &= 18
	default:
		_X86EMU_env.x86.gen.A.I32_reg.e_reg = uint32(0)
		_X86EMU_env.x86.gen.B.I32_reg.e_reg = uint32(0)
		_X86EMU_env.x86.gen.C.I32_reg.e_reg = uint32(0)
		_X86EMU_env.x86.gen.D.I32_reg.e_reg = uint32(0)
	}
}
func __init() {
}

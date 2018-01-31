// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (VarDecl): 791: function pointers are not supported
// Warning (VarDecl): 792: function pointers are not supported

package main

import "log"

type regtype uint16

// A regtype encodes width in low byte and shift amount in high byte
const (
	d  regtype = 32
	w          = 16
	bl         = 8
	bh         = 0x808
)

func R(r regtype) (regtype, regtype, regtype) {
	reg, shift, size := r>>8, (r>>4)&0xf, (r&0xf)*8
	if reg < 0 || reg > 15 {
		log.Panicf("R %x: bogus register # %02x", r, reg)
	}
	if shift != 0 && shift != 8 {
		log.Panicf("R %x: bogus register shift", r, shift)
	}
	if size != 1 && size != 2 && size != 4 {
		log.Panicf("R %x: bogus register size", r, size)
	}
	return reg, shift, size

}
func S(r regtype, val interface{}) {
	reg, shift, size := R(r)
	switch v := val.(type) {
	case uint32:
		switch size {
		case 4:
			if M.x86.mode&SYSMODE_32BIT_REP != 0 {
				M.x86.regs[reg] = v
			}
			M.x86.regs[reg] = M.x86.regs[reg]&0xffff0000 | uint32(v)
		default:
			log.Panicf("R %x: Can't assign 32 bits to %d bits", size)
		}
	case uint16:
		switch size {
		case 4, 2:
			M.x86.regs[reg] = M.x86.regs[reg]&0xffff0000 | uint32(v)
		default:
			log.Panicf("R %x: Can't assign 16 bits to %d bits", size)
		}
	case uint8:
		mask := 0xff << shift
		M.x86.regs[reg] = (M.x86.regs[reg] &^ mask) | v<<shift
	default:
		log.Panicf("Can't assign type %T to register", val)
	}
}

// Get gets the register as uint32. The amount of data depends on the SYSMODE.
// Note you can't just return the u32, always, in the none 32-bit case you have to
// return the low 16 bits, upper 16 0.
func G(r regtype) uint32 {
	reg, shift, size := R(r)
	v := M.x86.regs[reg]
	switch {
	case size == 32:
		if M.x86.mode&SYSMODE_32BIT_REP != 0 {
			return v
		}
		return uint32(uint16(v))
	case size == 16:
		return uint32(uint16(v))
	case size == 8 && shift == 0:
		return uint32(uint8(v))
	default:
		return uint32(uint8(v >> 8))
	}
}

// Changes takes a variable and adds it. It can be negative.
// In this case, due to the mode, we use the ability to override
// the number of bits in the register.
func Change(r regtype, i int) {
	S(r, G(r)+uint32(i))
}

func Dec(r regtype) {
	Change(r, -1)
}

func Inc(r regtype) {
	Change(r, 1)
}

// Simple encoding
// reg size is low nibl (#bytes)
// reg shift is next nibl
// reg # is next byte
const (
	AL     regtype = 0x0001
	AH             = 0x0081
	AX             = 0x0002
	EAX            = 0x0004
	BL             = 0x0101
	BH             = 0x0181
	BX             = 0x0102
	EBX            = 0x0104
	CL             = 0x0201
	CH             = 0x0281
	CX             = 0x0202
	ECX            = 0x0204
	DL             = 0x0301
	DH             = 0x0381
	DX             = 0x0302
	EDX            = 0x3004
	SP             = 0x0402
	ESP            = 0x0404
	BP             = 0x0502
	EBP            = 0x0504
	SI             = 0x0602
	ESI            = 0x0604
	DI             = 0x0702
	EDI            = 0x0704
	IP             = 0x0802
	EIP            = 0x0804
	FLAGS          = 0x0902
	EFLAGS         = 0x0904
	CS             = 0x00A02
	DS             = 0x00B02
	SS             = 0x00C02
	ES             = 0x00D02
	FS             = 0x00E02
	GS             = 0x00F02
)

type X86EMU_regs struct {
	regmem      [64]uint32
	mode        uint32
	intr        uint32
	debug       uint32
	check       int
	exit        bool
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

type X86EMU_intrFuncs struct {
	f uint32
}

type __int128_t int64
type __uint128_t uint64
type __builtin_ms_va_list []byte

type optab func(uint8)
type intrtab func(uint8)

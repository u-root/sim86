// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (VarDecl): 791: function pointers are not supported
// Warning (VarDecl): 792: function pointers are not supported

package main

import "log"

type register interface {
	Set32(uint32)
	Get32() uint32
	Set16(uint16)
	Get16() uint16
	Seth8(uint8)
	Geth8() uint8
	Setl8(uint8)
	Getl8() uint8
	Set(v interface{})
	Get() uint32
}

type register16 interface {
	Set(uint16)
	Get() uint16
	Set16(uint16)
	Get16() uint16
}

type register8 interface {
	Set(uint8)
	Get() uint8
	Set8(uint8)
	Get8() uint8
	Seth8(uint8)
	Geth8() uint8
}

type reg32 struct {
	reg uint32
}

type reg16 struct {
	reg uint32
}

type reg8 struct {
	reg uint32
}

//func (r reg32) Get() uint32 {
//	return r.reg32
//}

func (r reg32) Set32(i uint32) {
	r.reg = i
}
func (r reg32) Get32() uint32 {
	return r.reg32
}

func (r reg32) Set16(i uint16) {
	r.reg = uint32(i)
}
func (r reg32) Get16() uint16 {
	return uint16(r.reg)
}

func (r reg16) Set(i uint16) {
	r.reg = (r.reg & 0xffff0000) | uint32(i)
}
func (r reg16) Get() uint16 {
	return uint16(r.reg)
}
func (r reg16) Set16(i uint16) {
	r.Set(i)
}
func (r reg16) Get16() uint16 {
	return r.Get()
}

func (r reg8) Set(i uint8) {
	r.reg = (r.reg & 0xffffff00) | uint32(i)
}
func (r reg8) Get() uint8 {
	return uint8(r.reg)
}
func (r reg8) Set8(i uint8) {
	r.Set(i)
}
func (r reg8) Get8() uint8 {
	return r.Get()
}
func (r reg8) Seth8(i uint8) {
	r.reg = (r.reg & 0xfff00ff) | uint32(i)<<8
}
func (r reg8) Geth8() uint8 {
	return uint8(r.reg >> 8)
}

func (r reg32) Seth8(i uint8) {
	r.reg = (r.reg & 0x0000ff00) | uint32(i)<<8
}
func (r reg32) Geth8() uint8 {
	return uint8(r.reg >> 8)
}
func (r reg32) Setl8(i uint8) {
	r.reg = (r.reg & 0xffffff00) | uint32(i)
}
func (r reg32) Getl8() uint8 {
	return uint8(r.reg >> 8)
}

func (r reg32) Set(v interface{}) {
	switch i := v.(type) {
	case uint32:
		r.Set32(i)
	case uint16:
		r.Set16(i)
	case uint8:
		r.Setl8(i)
	default:
		log.Fatalf("Can't set register with %v", v)
	}
}

func (r reg32) Add(v interface{}) {
	switch i := v.(type) {
	case uint32:
		r.Set32(r.Get32() + i)
	case uint16:
		r.Set16(r.Get16() + i)
	case uint8:
		r.Setl8(r.Getl8() + i)
	default:
		log.Fatalf("Can't add register with %v", v)
	}
}

// Get gets the register as uint32. The amount of data depends on the SYSMODE.
func (r reg32) Get() uint32 {
	if M.x86.mode&SYSMODE_32BIT_REP != 0 {
		return r.Get32()
	}
	return r.Get16()
}

// Changes takes a variable and adds it. It can be negative.
// In this case, due to the mode, we use the ability to override
// the number of bits in the register.
func (r reg32) Change(i int) {
	if M.x86.mode&SYSMODE_32BIT_REP != 0 {
		r.Set32(r.Get32() + uint32(i))
	} else {
		r.Set16(r.Get16() + uint16(i))
	}
}

func (r reg32) Dec() {
	r.Change(-1)
}

func (r reg32) Inc() {
	r.Change(1)
}

type i386_general_regs struct {
	A reg32
	B reg32
	C reg32
	D reg32
}

type i386_special_regs struct {
	SP    reg32
	BP    reg32
	SI    reg32
	DI    reg32
	IP    reg32
	FLAGS uint32
}
type i386_segment_regs struct {
	CS reg16
	DS reg16
	SS reg16
	ES reg16
	FS reg16
	GS reg16
}

type X86EMU_regs struct {
	gen         i386_general_regs
	spc         i386_special_regs
	seg         i386_segment_regs
	mode        uint32
	intr        uint32
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

type X86EMU_intrFuncs struct {
	f uint32
}

type __int128_t int64
type __uint128_t uint64
type __builtin_ms_va_list []byte

type optab func(uint8)
type intrtab func(uint8)

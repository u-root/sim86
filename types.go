// Warning (TypedefDecl): 0: I couldn't find an appropriate Go type for the C type '__NSConstantString_tag'.
// Warning (VarDecl): 791: function pointers are not supported
// Warning (VarDecl): 792: function pointers are not supported

package main

type register interface {
	Set32(uint32)
	Get32() uint32
	Set16(uint16)
	Get16() uint16
	Seth8(uint8)
	Geth8() uint8
	Setl8(uint8)
	Getl8() uint8
}

type register16 interface {
	Set16(uint16)
	Get16(uint16)
}

type reg struct {
	reg uint32
}

type reg16 struct {
	reg uint16
}

func (r reg) Set32(i uint32) {
	r.reg = i
}
func (r reg) Get32() uint32 {
	return r.reg
}

func (r reg) Set16(i uint16) {
	r.reg = uint32(i)
}
func (r reg) Get16() uint16 {
	return uint16(r.reg)
}

func (r reg) Seth8(i uint8) {
	r.reg = (r.reg & 0x0000ff00) | uint32(i)<<8
}
func (r reg) Get8h() uint8 {
	return uint8(r.reg >> 8)
}
func (r reg) Setl8(i uint8) {
	r.reg = (r.reg & 0xffffff00) | uint32(i)
}
func (r reg) Get8l() uint8 {
	return uint8(r.reg >> 8)
}

func (r reg16) Set(i uint16) {
	r.reg = i
}
func (r reg16) Get() uint16 {
	return r.reg
}

type i386_general_regs struct {
	A reg
	B reg
	C reg
	D reg
}

type i386_special_regs struct {
	SP    reg
	BP    reg
	SI    reg
	DI    reg
	IP    reg
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

type X86EMU_intrFuncs struct {
	f uint32
}

type __int128_t int64
type __uint128_t uint64
type __builtin_ms_va_list []byte

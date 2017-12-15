package mem

var memory []uint32

func (a u32address) Write(val uint32) {
	memory[a] = val
}

func (a u32address) Read() uint32 {
	return memory[a]
}

func NewLongMemReader(a u32address) LongReader {
	return LongReader(a)
}

func NewLongMemWriter(a u32address) LongWriter{
	return LongWriter(a)
}

func (a u16address) Write(val uint16) {
	mem := memory[a>>1]
	v := uint32(val)
	if a & 1 == 0 {
		mem = ((mem>>16)<<16) | v
	} else {
		mem = (mem & 0xffff) | (v<<16)
	}
	memory[a>>1] = mem
}

func (a u16address) Read() uint16 {
	mem := memory[a>>1]
	if a & 1 == 0 {
		return uint16(mem)
	}
	return uint16(mem>>16)
}

func NewWordMemReader(a u16address) WordReader {
	return WordReader(a)
}

func NewWordMemWriter(a u16address) WordWriter{
	return WordWriter(a)
}

func (a u8address) Write(val uint8) {
	ff := 0xff
	b := (a&3)*8
	v := memory[a>>2]
	v &= ^uint32(ff<<b)
	v |= uint32(val)<<b
	memory[a>>2] = v
}

func (a u8address) Read() uint8 {
	b := (a&3)*8
	v := memory[a]
	return uint8(v>>b)
}

func NewByteMemReader(a u8address) ByteReader {
	return ByteReader(a)
}

func NewByteMemWriter(a u8address) ByteWriter{
	return ByteWriter(a)
}



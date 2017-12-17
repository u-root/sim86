package mem

var memory []byte

func (a u32address) Write(val uint32) {
	memory[a] = byte(val)
	memory[a+1] = byte(val>>8)
	memory[a+2] = byte(val>>16)
	memory[a+3] = byte(val>>24)
}

func (a u32address) Read() uint32 {
	return uint32(memory[a]) | uint32(memory[a+1]) << 8 | uint32(memory[a+2]) << 16 | uint32(memory[a+3]) << 24
}

func NewLongMemReader(a u32address) LongReader {
	return LongReader(a)
}

func NewLongMemWriter(a u32address) LongWriter{
	return LongWriter(a)
}

func (a u16address) Write(val uint16) {
	memory[a] = byte(val)
	memory[a+1] = byte(val>>8)
}

func (a u16address) Read() uint16 {
	return uint16(memory[a]) | uint16(memory[a+1]) << 8
}

func NewWordMemReader(a u16address) WordReader {
	return WordReader(a)
}

func NewWordMemWriter(a u16address) WordWriter{
	return WordWriter(a)
}

func (a u8address) Write(val uint8) {
	memory[a] = val
}

func (a u8address) Read() uint8 {
	return memory[a]
}

func NewByteMemReader(a u8address) ByteReader {
	return ByteReader(a)
}

func NewByteMemWriter(a u8address) ByteWriter{
	return ByteWriter(a)
}



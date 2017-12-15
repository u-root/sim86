package mem

type (
	u32address uint32
	u16address uint32
	u8address uint32
)

type LongWriter interface {
	Write(uint32)
}

type LongReader interface {
	Read() uint32
}

type WordWriter interface {
	Write(uint16)
}

type WordReader interface {
	Read() uint16
}

type ByteWriter interface {
	Write(uint8)
}

type ByteReader interface {
	Read() uint8
}


package main

import (
	"io/ioutil"
	"testing"
)

func TestIP(t *testing.T) {
	ip := uint16(0x1234)
	S(IP, ip)
	if G16(IP) != ip {
		t.Errorf("ip: got %04x, want %04x", G(IP), ip)
		x86emu_dump_regs()
	}
}

func TestEAX(t *testing.T) {
	eax := uint32(0xdeadbeef)
	S(EAX, eax)
	S(AL, uint8(0xaa))
	S(AH, uint8(0x44))
	if G(EAX) != 0x44aa {
		t.Errorf("EAX: got %08x, want %08x", G(EAX), 0x44aa)
		x86emu_dump_xregs()
	}
	// TODO: change mode, check G again
}

func TestBinary(t *testing.T) {
	b, err := ioutil.ReadFile("test.bin")
	if err != nil {
		t.Fatal(err)
	}
	copy(memory[:], b)
	S(CS, uint16(0))
	S(IP, uint16(0))
	M.x86.debug |= DEBUG_DISASSEMBLE_F | DEBUG_DECODE_F | DEBUG_TRACE_F
	X86EMU_exec()
}

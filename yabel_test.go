package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const (
	StackSeg     uint16 = 0xa000
	StackPointer uint16 = 0xfff0
	TOS          uint32 = 0xafff0
)

func TestIP(t *testing.T) {
	ip := uint16(0x1234)
	S(IP, ip)
	if G16(IP) != ip {
		t.Errorf("ip: got %04x, want %04x", G(IP), ip)
		fx86emu_dump_regs(t.Logf)
	}
}

func TestEAX(t *testing.T) {
	eax := uint32(0xdeadbeef)
	M.x86.mode |= SYSMODE_PREFIX_DATA
	S(EAX, eax)
	S(AL, uint8(0xaa))
	S(AH, uint8(0x44))
	if G32(EAX) != 0xdead44aa {
		t.Errorf("EAX: got %08x, want %08x", G32(EAX), 0xdead44aa)
		fx86emu_dump_xregs(t.Logf)
	}
	// TODO: change mode, check G again
}

type regval struct {
	r regtype
	v uint32
}

type check struct {
	n string
	r []regval
}

func TestBinary(t *testing.T) {
	var checks = []check{
		{n: "Halt", r: []regval{{IP, 1}, {SP, 0x2000}}},
		{n: "seg", r: []regval{{AX, 0x23}, {SS, 0x20}, {ES, 0x21}, {FS, 0x22}, {IP, 0x13}}},
		{n: "jmpcsip", r: []regval{{CS, 0x2}, {IP, 0x1}}},
		{n: "pushpop", r: []regval{{EBX, 0x12345678}, {CX, 0x5678}, {EDX, 0x12345678}}},
		{n: "qemu-test-i386-1", r: []regval{{CS, 0x2}, {IP, 0x16}, {EAX, 1}}},
		{n: "qemu-test-i386-2", r: []regval{{CS, 0x2}, {IP, 0x28}, {EBX, 0x12345678}, {ECX, 0x2}}},
		{n: "qemu-test-i386-3", r: []regval{{CS, 0x0}, {IP, 0x76}}},
		{n: "qemu-test-i386-4", r: []regval{{CS, 0x0}, {IP, 0x76}, {AX, 0x39}}},
	}

	b, err := ioutil.ReadFile("test.bin")
	if err != nil {
		t.Fatal(err)
	}

	// Fill memory with hlt.
	for i := range memory {
		memory[i] = 0xf4
	}

	copy(memory[:], b)
	S(CS, uint16(0))
	S(IP, uint16(0x0))
	S(SP, uint16(0xfff0))
	M.x86.debug |= DEBUG_DISASSEMBLE_F | DEBUG_DECODE_F | DEBUG_TRACE_F
	for _, c := range checks {
		S(SS, uint16(0))
		t.Logf("Start Test %s", c.n)
		X86EMU_exec()
		for i, r := range c.r {
			M.x86.mode &= ^SYSMODE_PREFIX_DATA
			if _, _, size := R(r.r); size == 4 { // simulate prefix
				M.x86.mode |= SYSMODE_PREFIX_DATA
			}
			if G32(r.r) != r.v {
				t.Errorf("%v: %d'th test fails: reg %s got %08x, want %08x", c.n, i, r.r.String(), G32(r.r), r.v)
			}
			if PC() > uint32(len(b)) {
				t.Fatalf("PC %08x: ran off the end of the test", PC())
			}
		}
		t.Logf("Done Test %s", c.n)
	}

	// test the reset vector
	S16(CS, 0xf000)
	S16(IP, 0xfff0)
	X86EMU_exec()
	if G16(CS) != 0xf000 || G16(IP) != 0x0001 {
		t.Fatalf("reset vector test: CS:IP is %04x:%04x, want 0xf000:0x0001", G16(CS), G16(IP))
	}

	// qemu tests
	b, err = ioutil.ReadFile("tcg/tests.bin")
	if err != nil {
		t.Fatal(err)
	}
	copy(memory[0:], b)

	S16(CS, 0)
	S16(IP, 0)
	for int(G16(IP)) < len(b) {
		S(SS, StackSeg)
		S(SP, StackPointer)
		t.Logf("Start code at %#04x:%#04x", G16(CS), G16(IP))
		X86EMU_exec()
		t.Logf("End code at %#04x:%#04x", G16(CS), G16(IP))
		fx86emu_dump_xregs(t.Logf)
		TestOutput := uint32(G16(IP))
		sp := int(G16(SS))<<4 + int(G16(SP))
		t.Logf("TestOutput at %#x; sp at %#x; vars %02x", TestOutput, sp, memory[TestOutput:TestOutput+0x10])
		dsz := sys_rdb(TestOutput)
		bits := sys_rdb(TestOutput + 1)
		nargs := int(sys_rdb(TestOutput + 2))
		t.Logf("dsz %d bits %d nargs %d", dsz, bits, nargs)
		// Can't scan for null. Damn.
		opx := TestOutput + uint32(dsz) + 1
		t.Logf("opx is %#x", opx)
		// Can't quite work out null terminators in strings library
		// but as loves them. So ...
		var opcode string
		for {
			b := sys_rdb(opx)
			if b == 0 {
				break
			}
			opcode = opcode + string([]byte{b})
			opx++
		}
		opx++
		var f string
		for {
			b := sys_rdb(opx)
			if b == 0 {
				break
			}
			f = f + string([]byte{b})
			opx++
		}
		t.Logf("f is %s and o is %s", f, opcode)
		args := []interface{}{opcode, map[byte]string{8:"b", 16:"w", 32: "l"}[bits]}
		t.Logf("Process %d args ", nargs)
		tos := TOS
		for i := 0; i < nargs; i++ {
			switch bits {
			case 16, 8:
				t.Logf("word at tos %#x is %02x", tos, memory[tos-2:tos])
				args = append(args, uint16(memory[tos-2])|uint16(memory[tos-1])<<8)
				tos -= 2
			case 32:
				t.Logf("long at tos %#x is %02x", tos, memory[tos-4:tos])
				args = append(args, uint32(memory[tos-4])|
					uint32(memory[tos-3])<<8|
					uint32(memory[tos-2])<<16|
					uint32(memory[tos-1])<<24)
				tos -= 4
			default:
				t.Fatalf("Bogus bit size: %d", dsz)
			}
		}
		// And, the iflags and flags are always there and always 16 bits
		args = append(args, uint16(memory[tos-2]))
		args = append(args, uint16(memory[tos-4]))
		out := fmt.Sprintf(f, args...)
		if _, ok := testout[out]; ! ok {
			t.Errorf("%s: can't find it in output", out)
		}

		t.Logf(f, args...)
		// opx is at the null after the fmt string. Bump IP to start again.
		S16(IP, uint16(opx+1))
	}

}

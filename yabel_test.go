package main

import (
	"debug/elf"
	"io"
	"io/ioutil"
	"strings"
	"testing"
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
	S(SP, uint16(0x2000))
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
	f, err := elf.Open("tcg/a.out")
	if err != nil {
		t.Fatalf("elf")
	}
	s := f.Section("initcall")
	if s == nil {
		t.Fatal(err)
	}
	for _, p := range f.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}
		if p.Vaddr > uint64(len(memory)) {
			t.Fatalf("p.Vaddr (%#x) > len(memory) %#x", p.Vaddr, len(memory))
		}
		if p.Vaddr+p.Filesz > uint64(len(memory)) {
			t.Fatalf("p.Vaddr (%#x) + p.Filesz %#x  len(memory) %#x", p.Vaddr, p.Filesz, len(memory))
		}
		t.Logf("Read in %d bytes at %#x", p.Filesz, p.Vaddr)
		b := make([]byte, p.Filesz)
		n, err := p.ReadAt(b, 0)
		// The elf package is strangely disfunctional
		if n < len(b) || (err != nil && err != io.EOF) {
			t.Fatalf("got %d bytes, err %v: wanted %d, nil", n, err, p.Filesz)
		}
		//t.Logf("b is %02x:", b)
		copy(memory[p.Vaddr:], b)
	}

	syms, err := f.Symbols()
	if err != nil {
		t.Fatal(err)
	}
	var addrs []elf.Symbol
	var TestOutput uint32
	for _, s := range syms {
		t.Logf("Check %v", s)
		if s.Name == "TestOutput" {
			TestOutput = uint32(s.Value)
			continue
		}
		if len(s.Name) < 5 || s.Name[:5] != "test_" {
			continue
		}
		addrs = append(addrs, s)
	}

	t.Logf("Now run a.out")

	for _, s := range addrs {
		S16(CS, 0)
		ip := uint16(s.Value)
		t.Logf("Start %s at %04x:%04x", s.Name, 0, ip)
		S16(IP, ip)
		X86EMU_exec()
		t.Logf("Finished")
		fx86emu_dump_xregs(t.Logf)
		narg := uint32(sys_rdw(TestOutput))
		t.Fatalf("nargs at %#x is %d", TestOutput, narg)
		if narg < 1 {
			continue
		}
		f := strings.TrimSpace(string(memory[sys_rdw(TestOutput+2):]))
		o := strings.TrimSpace(string(memory[sys_rdw(TestOutput+4):]))
		args := []interface{}{o}
		for i := uint32(0); i < narg - 2; i++ {
			args = append(args, uint32(sys_rdw(TestOutput+6+i)))
		}
		t.Logf("f is %s and o is %s", f, o)
		t.Logf(f, args...)

	}
}

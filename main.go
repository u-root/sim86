package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	disassemble = flag.Bool("d", false, "Disassemble")
	step        = flag.Bool("s", false, "Single step")
	trace       = flag.Bool("t", false, "Trace")
	cmds        = bufio.NewReader(os.Stdin)
)

var memory [1 << 20]byte

func cpuInit() {
	// due to init loop
	copy(x86emu_optab[:], _x86emu_optab[:])
}

func main() {
	log.Printf("x86 emulator")
	flag.Parse()
	if *disassemble {
		M.x86.debug |= DEBUG_DISASSEMBLE_F | DEBUG_DECODE_F
	}
	if *step {
		M.x86.debug |= DEBUG_STEP_F
	}
	if *trace {
		M.x86.debug |= DEBUG_TRACE_F
	}
	cpuInit()
	X86EMU_exec()
	log.Printf("Done")
}

package main

var (
	x86emu_optab  = make(map[uint8]func(uint8))
	x86emu_optab2 = make(map[uint8]func(uint8))
)

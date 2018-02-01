.code16
.globl Halt
Halt:hlt
Segs:
	movw $0x20, %AX
	movw %AX, %SS
	hlt



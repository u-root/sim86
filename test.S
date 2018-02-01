.code16
.globl Halt
Halt:hlt
Segs:
	movw $0x20, %AX
	movw %AX, %SS
	addb $1, %al
	movw %AX, %ES
	addb $1, %al
	movw %AX, %FS
	addb $1, %al
	movw %AX, %GS
	hlt



package main

func incamount(scale int) int {
	if ACCESS_FLAG(F_DF) { /* down */
		return scale * -1
	}
	return scale * 1
}

// Count returns cx/x depending on mode
func Count(mode uint32) uint32 {
	if Mode(mode) {
		return G32(ECX)
	}
	return uint32(G16(CX))
}
func ClrCount(mode uint32) {
	if Mode(mode) {
		S(ECX, 0)
		return
	}
	S(CX, 0)
}

// DecCount decrements count, depending on the mode.
func DecCount() {
	if Mode(SYSMODE_PREFIX_ADDR) {
		S(CX, G16(CX)-1)
		return
	}
	S(ECX, G32(ECX)-1)
}

// GetClrCount gets the c/cx register and clears it, as well as
// clearing the REPE/REPNE bits from mode.
func GetClrCount() uint32 {
	M.x86.mode &= ^(SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE)
	var count uint32
	if M.x86.mode&SYSMODE_32BIT_REP == 0 {
		count = uint32(G16(CX))
		S(CX, 0)
	} else {
		count = G32(ECX)
		S(ECX, 0)
	}

	return count
}

// Halted returns 1 if we are halted
func Halted() bool {
	return M.x86.intr&INTR_HALTED != 0
}

func Counting() bool {
	return M.x86.mode&(SYSMODE_PREFIX_REPE|SYSMODE_PREFIX_REPNE) != 0
}

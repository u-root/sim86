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
		return M.x86.gen.C.Get32()
	}
	return uint32(M.x86.gen.C.Get16())
}
func ClrCount(mode uint32) {
	if Mode(mode) {
		M.x86.gen.C.Set32(0)
		return
	}
	M.x86.gen.C.Set16(0)
}

// DecCount decrements count, depending on the mode.
func DecCount() {
	if Mode(SYSMODE_PREFIX_ADDR) {
		M.x86.gen.C.Set16(M.x86.gen.C.Get16() - 1)
		return
	}
	M.x86.gen.C.Set32(M.x86.gen.C.Get32() - 1)
}

// GetClrCount gets the c/cx register and clears it, as well as
// clearing the REPE/REPNE bits from mode.
func GetClrCount() uint32 {
	M.x86.mode &= ^(SYSMODE_PREFIX_REPE | SYSMODE_PREFIX_REPNE)
	var count uint32
	if M.x86.mode&SYSMODE_32BIT_REP == 0 {
		count = uint32(M.x86.gen.C.Get16())
		M.x86.gen.C.Set16(0)
	} else {
		count = M.x86.gen.C.Get32()
		M.x86.gen.C.Set32(0)
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

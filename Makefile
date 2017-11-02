all: cppops.c
	c2go transpile cppops.c

cppops.c: ops.c defines.h
	cpp --include defines.h ops.c cppops.c
	sed -i 's/= \(genop.*\)\[\(.*\)\](\(.*\),\(.*\))/= call4("\1", "\2", "\3", "\4")/' cppops.c
	sed -i 's/= (\*\(genop.*\)\[\(.*\)\]) *(\(.*\),\(.*\))/= call4("\1", "\2", "\3", "\4")/' cppops.c
	sed -i 's/(.x86emu_optab2.op2.)(op2)/call4("x86emu_optab2.op2", "op2", "op2", "")/' cppops.c
	sed -i 's/(._X86EMU_intrTab.\(.*\)\])(\(.*\))/call4("_X86EMU_intrTab", "\1", "\2", "")/' cppops.c
	sed -i 's/= (\*\(opcD.*\)\[\(.*\)\]) *(\(.*\),\(.*\))/= call4("\1", "\2", "\3", "\4")/' cppops.c
	sed -i 's/= (\*\(opcD.*\)\[\(.*\)\]) *(\(.*\),\(.*\))/= call4("\1", "\2", "\3", "\4")/' cppops.c
	sed -i '/^#/d' cppops.c

all:
	cpp --include defines.h ops.c cppops.c
	c2go transpile cppops.c

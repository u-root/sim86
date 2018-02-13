include config-host.mak

all:
	cc -E -I . test.c > t.c
	cc -m32 -S -I . test.c > t.s
	cc -m32 t.c
	cc -E -I . tests.S > tests.s
	as tests.s
	objcopy -O binary a.out tests.bin
	ndisasm tests.bin > tests.asm

clean:
	rm -f t.c a.out


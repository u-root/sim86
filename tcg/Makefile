include config-host.mak

all:
	go run gen.go
old:
	cc -E -I . test.c > t.c
	cc -m32 -S -I . test.c > t.s
	cc -m32 t.c
	cc -E -I . tests.S > testsout.S
	sed -i 's/^#.*//' testsout.S
	sed -i 's/;/\n/g' testsout.S
	sed -i 's/#/\n#/g' testsout.S
	as -a testsout.S > testsout.asm
	objcopy -O binary a.out tests.bin

clean:
	rm -f t.c a.out


include config-host.mak

all:
	cc -E -I . test.c > t.c
	cc -m32 t.c

clean:
	rm -f t.c a.out


all:
	c2go transpile sys.c
	go build sys.go
	c2go transpile ops.c
	go build ops.go
	c2go transpile ops2.c
	c2go transpile debug.c
	c2go transpile decode.c
	c2go transpile prim_ops.c

package main

import (
	"fmt"
	"text/template"
)

type test struct {
	O string
	S string
	X string
	A string
	B string
	F string
	R string
	E string
}

var (
	tests = []test {
		{"add", },
	}
	execop2=`
	movw	$,{{.F}}, %dx 
	pushw %dx 
	popf 
        mov{{.S}}  ${{.R}}, {{.E}}{{.A}}{{.X}}
        pushl %e{{.A}}x
        mov{{.S}}  ${{.R}}, {{.E}}{{.B}}{{.X}}
        pushl %e{{.B}}x
        {{.O}}{{.S}} {{.E}}{{.B}}{{.S}}, {{.E}}{{.A}}{{.S}}
        pushl %e{{.A}}x
	pushf 
	movw	${{.F}}, %dx 
	pushw %dx 
	hlt 
	.byte 2 /* number of following bytes of info */ 
	/* currently # bits per stack item, and nargs */ 
	.byte bits, 3 
	.asciz #o 							
	.asciz "%s%s A=%08x B=%08x R=%08x CCIN=%04x CC=%04x" 
`
	
execop1=`
	movw	$flags, %dx 
	pushw %dx 
	popf 
	OPR(mov,l) $res, REG(a, e, x)	
	PUSH(a,e) 
	OPR(o,size) REG(a,pre, rsize) 	
	PUSH(a,e) 					
	pushf 
	movw	$flags, %dx 
	pushw %dx 
	hlt
	.byte 2 /* number of following bytes of info */ 
	/* currently # bits per stack item, and nargs */ 
	.byte bits, 2 
	.asciz #o 							
	.asciz "%s%s A=%08x R=%08x CCIN=%04x CC=%04x" 
`
	code = map[string]string {
		"add": execop2,
	}
)


func gen(t test) {
	c := code[t.O]
	
}
func main() {
	for _, t := range tests {
		gen(t)
	}
}

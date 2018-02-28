package main

import (
	"log"
	"os"
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
		{O: "add", },
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
	ops = []*template.Template{
		template.Must(template.New("op1").Parse(execop1)),
		template.Must(template.New("op2").Parse(execop2)),

	}
	s = []string{"b", "w", "l"}
	b = []int{8, 16, 32}
	x = []string{"x", "x", "l"}
	e = []string{"", "", "e"}
)


func gen(t test) {
	for _, a := range ops {
		for i := 0; i < 3; i++ {
			var tt = test{
				O:t.O,
				S: s[i],
				X: x[i],
				A: "a",
				B: "b",
				F: t.F,
				R: "a",
				E: e[i],
			}

			if err := a.Execute(os.Stdout, tt); err != nil {
				log.Print(err)
			}
		}
	}
}

func main() {
	for _, t := range tests {
		gen(t)
	}
}

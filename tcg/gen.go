package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const CC_C = "1"

type test struct {
	O           string
	S           string
	X           string
	A           string
	B           string
	F           string
	Arg0        string
	Arg1        string
	Bits        string
	RegOpSuffix string
	E           string
}

type op2 struct {
	A string
	B string
}

var (
	outFile *os.File
	outName = flag.String("o", "testsout.S", "Output file for assembly")
	tests2  = []test{
		{O: "adc", F: "0", S: " "}, {O: "adc", F: CC_C, S: " "},
		{O: "sbb", F: "0"}, {O: "sbb", F: CC_C},
		{O: "add", F: "0"},
		{O: "sub", F: "0"},
		{O: "xor", F: "0"},
		{O: "and", F: "0"},
		{O: "or", F: "0"},
		{O: "cmp", F: "0"},
	}
	tests1 = []test{
		{O: "inc", F: "0"}, {O: "inc", F: CC_C},
		{O: "dec", F: "0"}, {O: "dec", F: CC_C},
		{O: "neg", F: "0"}, {O: "neg", F: CC_C},
		{O: "not", F: "0"}, {O: "not", F: CC_C},
	}
	shifts = []test {
		{O: "rcr", F: "0"}, {O: "rcr", F: CC_C},
		{O: "rcl", F: "0"}, {O: "rcl", F: CC_C},
	}
	execop2 = `
	movw	${{.F}}, %dx 
	pushw %dx 
	popf 
        movl  ${{.Arg0}}, %e{{.A}}x
        push %e{{.A}}x
        movl  ${{.Arg1}}, %e{{.B}}x
        push %e{{.B}}x
        {{.O}}{{.S}} %{{.E}}{{.B}}{{.X}}, %{{.E}}{{.A}}{{.X}}
        push %e{{.A}}x
	pushf 
	movw	${{.F}}, %dx 
	pushw %dx 
	hlt 
	.byte 2 /* number of following bytes of info */ 
	/* currently # bits per stack item, and nargs */ 
	.byte {{.Bits}}, 3 
	.asciz "{{.O}}"
	.asciz "%s%s A=%08x B=%08x R=%08x CCIN=%04x CC=%04x" 
`

	execop1 = `
	movw	${{.F}}, %dx 
	pushw %dx 
	popf 
        movl  ${{.Arg0}}, %e{{.A}}x
        push %e{{.A}}x
        {{.O}}{{.S}}  %{{.E}}{{.A}}{{.X}}
        push %e{{.A}}x
	pushf 
	movw	${{.F}}, %dx 
	pushw %dx 
	hlt
	.byte 2 /* number of following bytes of info */ 
	/* currently # bits per stack item, and nargs */ 
	.byte {{.Bits}}, 2 
	.asciz "{{.O}}"
	.asciz "%s%s A=%08x R=%08x CCIN=%04x CC=%04x" 
`
	code = map[string]string{
		"add": execop2,
	}
	ops = []*template.Template{
		template.Must(template.New("op1").Parse(execop1)),
		template.Must(template.New("op2").Parse(execop2)),
	}
	s    = []string{"b", "w", "l"}
	b    = []int{8, 16, 32}
	x    = []string{"x", "x", "l"}
	e    = []string{"", "", "e"}
	ops1 = []string{
		"0x12345678",
		"0x12341",
		"0xffffffff",
		"0x7fffffff",
		"0x80000000",
		"0x12347fff",
		"0x12348000",
		"0x12347f7f",
		"0x12348080",
	}
	ops2 = []op2{
		{A: "0x12345678", B: "0x812FADA"},
		{A: "0x12341", B: "0x12341"},
		{A: "0x12341", B: "-0x12341"},
		{A: "0xffffffff", B: "0"},
		{A: "0xffffffff", B: "-1"},
		{A: "0xffffffff", B: "1"},
		{A: "0xffffffff", B: "2"},
		{A: "0x7fffffff", B: "0"},
		{A: "0x7fffffff", B: "1"},
		{A: "0x7fffffff", B: "-1"},
		{A: "0x80000000", B: "-1"},
		{A: "0x80000000", B: "1"},
		{A: "0x80000000", B: "-2"},
		{A: "0x12347fff", B: "0"},
		{A: "0x12347fff", B: "1"},
		{A: "0x12347fff", B: "-1"},
		{A: "0x12348000", B: "-1"},
		{A: "0x12348000", B: "1"},
		{A: "0x12348000", B: "-2"},
		{A: "0x12347f7f", B: "0"},
		{A: "0x12347f7f", B: "1"},
		{A: "0x12347f7f", B: "-1"},
		{A: "0x12348080", B: "-1"},
		{A: "0x12348080", B: "1"},
		{A: "0x12348080", B: "-2"},
	}
)

func gen1(t test, operands []string) {
	for _, o := range operands {
		for i := 0; i < 3; i++ {
			bits := "8"
			lxx := "l"
			switch i {
			case 0:
			case 1:
				bits = "16"
				lxx = "x"
			case 2:
				bits = "32"
				lxx = "x"
			default:
				log.Panic("fix me")
			}
			var tt = test{
				O:    t.O,
				X:    lxx,
				A:    "a",
				F:    t.F,
				Arg0: o,
				E:    e[i],
				Bits: bits,
			}
			if t.S == "" {
				tt.S = s[i]
			}

			if err := ops[0].Execute(outFile, tt); err != nil {
				log.Print(err)
			}
		}
	}
}

func gen2(t test, operands []op2) {
	for _, o := range operands {
		for i := 0; i < 3; i++ {
			bits := "8"
			lxx := "l"
			switch i {
			case 0:
			case 1:
				bits = "16"
				lxx = "x"
			case 2:
				bits = "32"
				lxx = "x"
			default:
				log.Panic("fix me")
			}
			var tt = test{
				O:    t.O,
				X:    lxx,
				A:    "a",
				B:    "b",
				F:    t.F,
				Arg0: o.A,
				Arg1: o.B,
				E:    e[i],
				Bits: bits,
			}
			if t.S == "" {
				tt.S = s[i]
			}

			if err := ops[1].Execute(outFile, tt); err != nil {
				log.Print(err)
			}
		}
	}
}

func main() {
	var err error
	flag.Parse()
	outFile, err = os.OpenFile(*outName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(outFile, ".code16\n")
	for _, t := range tests1 {
		gen1(t, ops1)
	}
	for _, t := range tests2 {
		gen2(t, ops2)
	}
	c := exec.Command("as", []string{"-a", "testsout.S"}...)
	c.Stdout, err = os.OpenFile("testsout.asm", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}

	c = exec.Command("objcopy", []string{"-O", "binary", "a.out", "tests.bin"}...)
	if o, err := c.CombinedOutput(); err != nil {
		log.Fatalf("%v %v", string(o), err)
	}
}

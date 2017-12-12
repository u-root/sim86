/****************************************************************************
*
*                       Realmode X86 Emulator Library
*
*               Copyright (C) 1991-2004 SciTech Software, Inc.
*                    Copyright (C) David Mosberger-Tang
*                      Copyright (C) 1999 Egbert Eich
*
*  ========================================================================
*
*  Permission to use, copy, modify, distribute, and sell this software and
*  its documentation for any purpose is hereby granted without fee,
*  provided that the above copyright notice appear in all copies and that
*  both that copyright notice and this permission notice appear in
*  supporting documentation, and that the name of the authors not be used
*  in advertising or publicity pertaining to distribution of the software
*  without specific, written prior permission.  The authors makes no
*  representations about the suitability of this software for any purpose.
*  It is provided "as is" without express or implied warranty.
*
*  THE AUTHORS DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE,
*  INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO
*  EVENT SHALL THE AUTHORS BE LIABLE FOR ANY SPECIAL, INDIRECT OR
*  CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF
*  USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
*  OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
*  PERFORMANCE OF THIS SOFTWARE.
*
*  ========================================================================
*
* Language:     ANSI C
* Environment:  Any
* Developer:    Kendall Bennett
*
* Description:  This file contains the code to handle debugging of the
*               emulator.
*
****************************************************************************/

package main

import (
	"fmt"
)

/*----------------------------- Implementation ----------------------------*/

/* should look something like debug's output. */
func X86EMU_trace_regs() {
	if DEBUG_TRACE() {
		if (M().x86.mode & uint32(SYSMODE_PREFIX_DATA | SYSMODE_PREFIX_ADDR)) != 0  {
			x86emu_dump_xregs()
		} else {
			x86emu_dump_regs()
		}
	}
	if DEBUG_DECODE() && !DEBUG_DECODE_NOPRINT() {
		fmt.Printf("%04x:%04x ", M().x86.saved_cs, M().x86.saved_ip)
		print_encoded_bytes(M().x86.saved_cs, M().x86.saved_ip)
		print_decoded_instruction()
	}
}

func X86EMU_trace_xregs() {
	if DEBUG_TRACE() {
		x86emu_dump_xregs()
	}
}

func x86emu_just_disassemble() {
	/*
	 * This routine called if the flag DEBUG_DISASSEMBLE is set kind
	 * of a hack!
	 */
	fmt.Printf("%04x:%04x ", M().x86.saved_cs, M().x86.saved_ip)
	print_encoded_bytes(M().x86.saved_cs, M().x86.saved_ip)
	print_decoded_instruction()
}

func disassemble_forward(seg uint16, off uint16, n int) {
	var (
		tregs *X86EMU_sysEnv
		i    int
		op1  uint8
	)
	/*
	 * hack, hack, hack.  What we do is use the exact machinery set up
	 * for execution, except that now there is an additional state
	 * flag associated with the "execution", and we are using a copy
	 * of the register struct.  All the major opcodes, once fully
	 * decoded, have the following two steps: TRACE_REGS(r,m);
	 * SINGLE_STEP(r,m); which disappear if DEBUG is not defined to
	 * the preprocessor.  The TRACE_REGS macro expands to:
	 *
	 * if (debug&DEBUG_DISASSEMBLE)
	 *     {just_disassemble(); goto EndOfInstruction;}
	 *     if (debug&DEBUG_TRACE) trace_regs(r,m);
	 *
	 * ......  and at the last line of the routine.
	 *
	 * EndOfInstruction: end_instr();
	 *
	 * Up to the point where TRACE_REG is expanded, NO modifications
	 * are done to any register EXCEPT the IP register, for fetch and
	 * decoding purposes.
	 *
	 * This was done for an entirely different reason, but makes a
	 * nice way to get the system to help debug codes.
	 */
	tregs = M()
	tregs.x86.spc.IP.Set16(off)
	tregs.x86.seg.CS.Set(seg)

	/* reset the decoding buffers */
	tregs.x86.enc_str_pos = 0
	tregs.x86.enc_pos = 0

	/* turn on the "disassemble only, no execute" flag */
	tregs.x86.debug |= DEBUG_DISASSEMBLE_F

	/* DUMP NEXT n instructions to screen in straight_line fashion */
	/*
	 * This looks like the regular instruction fetch stream, except
	 * that when this occurs, each fetched opcode, upon seeing the
	 * DEBUG_DISASSEMBLE flag set, exits immediately after decoding
	 * the instruction.  XXX --- CHECK THAT MEM IS NOT AFFECTED!!!
	 * Note the use of a copy of the register structure...
	 */
	for i := 0; i < n; i += 1 {
		ip := M().x86.spc.IP.Get16()
		op1 = sys_rdb(uint32(M().x86.seg.CS.Get()) << 4 + uint32(ip))
		M().x86.spc.IP.Set16(ip+1)
		x86emu_optab[op1](op1)
	}
	/* end major hack mode. */
}

func x86emu_check_ip_access() {
	/* NULL as of now */
}

func x86emu_check_sp_access() {
}

func x86emu_check_mem_access(_ uint32) {
	/*  check bounds, etc */
}

func x86emu_check_data_access(_, _ uint) {
	/*  check bounds, etc */
}

func x86emu_inc_decoded_inst_len(x int) {
	M().x86.enc_pos += x
}

func x86emu_decode_printf(x string, y ...interface{}) {
	M().x86.decoded_buf = []byte(string(M().x86.decoded_buf) + fmt.Sprintf(x, y...))
}

func x86emu_decode_printf2(x string, y int) {
	x86emu_decode_printf(x , y)
}

func x86emu_end_instr() {
	M().x86.decoded_buf = []byte{}
}

func print_encoded_bytes(s uint16, o uint16) {
	var notyet = `
    for (i:=0; i< M().x86.enc_pos; i++) {
	    snfmt.Printf(buf1+2*i, 64 - 2 * i, "%02x", fetch_data_byte_abs(s,o+i));
    }
    fmt.Printf("%-20s ",buf1);`
}

func print_decoded_instruction() {
	fmt.Printf("%s", M().x86.decoded_buf)
}

func x86emu_print_int_vect(iv uint16) {
	var seg, off uint16

	if iv > 256 {
		return
	}
	seg = fetch_data_word_abs(0, iv*4)
	off = fetch_data_word_abs(0, iv*4+2)
	fmt.Printf("%04x:%04x ", seg, off)
}

func X86EMU_dump_memory(seg uint16, o uint16, amt uint32) {
	var (
		off = uint32(o)
		start = uint32(off) & 0xfffffff0
		end   = uint32(off+16) & 0xfffffff0
		i     uint32
	)

	for end <= off+amt {
		fmt.Printf("%04x:%04x ", seg, start)
		for i = start; i < off; i++ {
			fmt.Printf("   ")
		}
		for i < end {
			fmt.Printf("%02x ", fetch_data_byte_abs(seg, uint16(i)))
			fmt.Printf("\n")
			start = end
			end = start + 16
			i++
		}
	}
}

func x86emu_single_step() {
	var notyet = `
    char s[1024];
    int ps[10];
    int ntok;
    int cmd;
    int done;
        int segment;
    int offset;
    static int breakpoint;
    static int noDecode = 1;

        if (DEBUG_BREAK()) {
                if (M().x86.saved_ip != breakpoint) {
                        return;
                } else {
              M().x86.debug &= ~DEBUG_DECODE_NOPRINT_F;
                        M().x86.debug |= DEBUG_TRACE_F;
                        M().x86.debug &= ~DEBUG_BREAK_F;
                        print_decoded_instruction ();
                        X86EMU_trace_regs();
                }
        }
    done=0;
    offset = M().x86.saved_ip;
    while (!done) {
        fmt.Printf("-");
        (void)fgets(s, 1023, stdin);
        cmd = parse_line(s, ps, &ntok);
        switch(cmd) {
          case 'u':
            disassemble_forward(M().x86.saved_cs,uint16(offset),10);
            break;
          case 'd':
                            if (ntok == 2) {
                                    segment = M().x86.saved_cs;
                                    offset = ps[1];
                                    X86EMU_dump_memory(segment,uint16(offset),16);
                                    offset += 16;
                            } else if (ntok == 3) {
                                    segment = ps[1];
                                    offset = ps[2];
                                    X86EMU_dump_memory(segment,uint16(offset),16);
                                    offset += 16;
                            } else {
                                    segment = M().x86.saved_cs;
                                    X86EMU_dump_memory(segment,uint16(offset),16);
                                    offset += 16;
                            }
            break;
          case 'c':
            M().x86.debug ^= DEBUG_TRACECALL_F;
            break;
          case 's':
            M().x86.debug ^= DEBUG_SVC_F | DEBUG_SYS_F | DEBUG_SYSINT_F;
            break;
          case 'r':
            X86EMU_trace_regs();
            break;
          case 'x':
            X86EMU_trace_xregs();
            break;
          case 'g':
            if (ntok == 2) {
                breakpoint = ps[1];
        if (noDecode) {
                        M().x86.debug |= DEBUG_DECODE_NOPRINT_F;
        } else {
                        M().x86.debug &= ~DEBUG_DECODE_NOPRINT_F;
        }
        M().x86.debug &= ~DEBUG_TRACE_F;
        M().x86.debug |= DEBUG_BREAK_F;
        done = 1;
            }
            break;
          case 'q':
          M().x86.debug |= DEBUG_EXIT;
          return;
      case 'P':
          noDecode = (noDecode)?0:1;
          fmt.Printf("Toggled decoding to %s\n",(noDecode)?"FALSE":"TRUE");
          break;
          case 't':
      case 0:
            done = 1;
            break;
        }
    }
`
}

func X86EMU_trace_on() int {
	M().x86.debug = M().x86.debug | DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F
	return M().x86.debug
}

func X86EMU_trace_off() int {
	M().x86.debug = M().x86.debug & ^(DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F)
	return M().x86.debug
}

func parse_line(s string, ps *int, n *int) error {
	var notyet = `
    int cmd;

    *n = 0;
    while(*s == ' ' || *s == '\t') s++;
    ps[*n] = *s;
    switch (*s) {
      case '\n':
        *n += 1;
        return 0;
      default:
        cmd = *s;
        *n += 1;
    }

    while (1) {
        while (*s != ' ' && *s != '\t' && *s != '\n')  s++;

        if (*s == '\n')
            return cmd;

        while(*s == ' ' || *s == '\t') s++;

        sscanf(s,"%x",&ps[*n]);
        *n += 1;
    }
    `
	return errors.New("not yet")
}
func x86emu_dump_regs() {
	fmt.Printf("\tAX=%04x  ", M().x86.R_AX)
	fmt.Printf("BX=%04x  ", M().x86.R_BX)
	fmt.Printf("CX=%04x  ", M().x86.R_CX)
	fmt.Printf("DX=%04x  ", M().x86.R_DX)
	fmt.Printf("SP=%04x  ", M().x86.R_SP)
	fmt.Printf("BP=%04x  ", M().x86.R_BP)
	fmt.Printf("SI=%04x  ", M().x86.R_SI)
	fmt.Printf("DI=%04x\n", M().x86.R_DI)
	fmt.Printf("\tDS=%04x  ", M().x86.R_DS)
	fmt.Printf("ES=%04x  ", M().x86.R_ES)
	fmt.Printf("SS=%04x  ", M().x86.R_SS)
	fmt.Printf("CS=%04x  ", M().x86.R_CS)
	fmt.Printf("IP=%04x   ", M().x86.R_IP)
	/* CHECKED... */
	if ACCESS_FLAG(F_OF) {
		fmt.Printf("OV ")
	} else {
		fmt.Printf("NV ")
	}
	if ACCESS_FLAG(F_DF) {
		fmt.Printf("DN ")
	} else {
		fmt.Printf("UP ")
	}
	if ACCESS_FLAG(F_IF) {
		fmt.Printf("EI ")
	} else {
		fmt.Printf("DI ")
	}
	if ACCESS_FLAG(F_SF) {
		fmt.Printf("NG ")
	} else {
		fmt.Printf("PL ")
	}
	if ACCESS_FLAG(F_ZF) {
		fmt.Printf("ZR ")
	} else {
		fmt.Printf("NZ ")
	}
	if ACCESS_FLAG(F_AF) {
		fmt.Printf("AC ")
	} else {
		fmt.Printf("NA ")
	}
	if ACCESS_FLAG(F_PF) {
		fmt.Printf("PE ")
	} else {
		fmt.Printf("PO ")
	}
	if ACCESS_FLAG(F_CF) {
		fmt.Printf("CY ")
	} else {
		fmt.Printf("NC ")
	}
	fmt.Printf("\n")
}

func x86emu_dump_xregs() {
	fmt.Printf("\tEAX=%08x  ", M().x86.R_EAX)
	fmt.Printf("EBX=%08x  ", M().x86.R_EBX)
	fmt.Printf("ECX=%08x  ", M().x86.R_ECX)
	fmt.Printf("EDX=%08x\n", M().x86.R_EDX)
	fmt.Printf("\tESP=%08x  ", M().x86.R_ESP)
	fmt.Printf("EBP=%08x  ", M().x86.R_EBP)
	fmt.Printf("ESI=%08x  ", M().x86.R_ESI)
	fmt.Printf("EDI=%08x\n", M().x86.R_EDI)
	fmt.Printf("\tDS=%04x  ", M().x86.R_DS)
	fmt.Printf("ES=%04x  ", M().x86.R_ES)
	fmt.Printf("SS=%04x  ", M().x86.R_SS)
	fmt.Printf("CS=%04x  ", M().x86.R_CS)
	fmt.Printf("EIP=%08x\n\t", M().x86.R_EIP)

	/* CHECKED... */
	if ACCESS_FLAG(F_OF) {
		fmt.Printf("OV ")
	} else {
		fmt.Printf("NV ")
	}
	if ACCESS_FLAG(F_DF) {
		fmt.Printf("DN ")
	} else {
		fmt.Printf("UP ")
	}
	if ACCESS_FLAG(F_IF) {
		fmt.Printf("EI ")
	} else {
		fmt.Printf("DI ")
	}
	if ACCESS_FLAG(F_SF) {
		fmt.Printf("NG ")
	} else {
		fmt.Printf("PL ")
	}
	if ACCESS_FLAG(F_ZF) {
		fmt.Printf("ZR ")
	} else {
		fmt.Printf("NZ ")
	}
	if ACCESS_FLAG(F_AF) {
		fmt.Printf("AC ")
	} else {
		fmt.Printf("NA ")
	}
	if ACCESS_FLAG(F_PF) {
		fmt.Printf("PE ")
	} else {
		fmt.Printf("PO ")
	}
	if ACCESS_FLAG(F_CF) {
		fmt.Printf("CY ")
	} else {
		fmt.Printf("NC ")
	}
	fmt.Printf("\n")
}

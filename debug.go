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
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*----------------------------- Implementation ----------------------------*/

/* should look something like debug's output. */
func X86EMU_trace_regs() {
	if DEBUG_TRACE() {
		if (M.x86.mode & uint32(SYSMODE_PREFIX_DATA|SYSMODE_PREFIX_ADDR)) != 0 {
			x86emu_dump_xregs()
		} else {
			x86emu_dump_regs()
		}
	}
	fmt.Printf("trace regs DD %v DDD %v\n", DEBUG_DECODE(),DEBUG_DECODE_NOPRINT())
	if DEBUG_DECODE() && !DEBUG_DECODE_NOPRINT() {
		fmt.Printf("%04x:%04x ", M.x86.saved_cs, M.x86.saved_ip)
		print_encoded_bytes(M.x86.saved_cs, M.x86.saved_ip)
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
	fmt.Printf("%04x:%04x ", M.x86.saved_cs, M.x86.saved_ip)
	print_encoded_bytes(M.x86.saved_cs, M.x86.saved_ip)
	print_decoded_instruction()
}

func disassemble_forward(seg uint16, off uint16, n int) {
	var (
		tregs X86EMU_sysEnv
		op1   uint8
	)
	fmt.Printf("DAF %04x:%04x %d\n", seg, off, n)
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
	tregs = *M
	M.x86.spc.IP.Set16(off)
	M.x86.seg.CS.Set(seg)

	/* reset the decoding buffers */
	M.x86.enc_str_pos = 0
	M.x86.enc_pos = 0
	M.x86.decoded_buf = []byte{}

	/* turn on the "disassemble only, no execute" flag */
	M.x86.debug |= DEBUG_DISASSEMBLE_F

	/* DUMP NEXT n instructions to screen in straight_line fashion */
	/*
	 * This looks like the regular instruction fetch stream, except
	 * that when this occurs, each fetched opcode, upon seeing the
	 * DEBUG_DISASSEMBLE flag set, exits immediately after decoding
	 * the instruction.  XXX --- CHECK THAT MEM IS NOT AFFECTED!!!
	 * Note the use of a copy of the register structure...
	 */
	for i := 0; i < n; i += 1 {
		ip := M.x86.spc.IP.Get16()
		op1 = sys_rdb(uint32(M.x86.seg.CS.Get())<<4 + uint32(ip))
		M.x86.spc.IP.Set16(ip + 1)
		x86emu_optab[op1](op1)
	}
	*M = tregs
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

func x86emu_inc_decoded_inst_len(x uint32) {
	M.x86.enc_pos += int(x)
}

func x86emu_decode_printf(x string, y ...interface{}) {
	fmt.Printf(x, y...)
//M.x86.decoded_buf = []byte(string(M.x86.decoded_buf) + fmt.Sprintf(x, y...))
}

func x86emu_decode_printf2(x string, y int) {
	fmt.Printf(x, y)
//	x86emu_decode_printf(x, y)
}

func x86emu_end_instr() {
	M.x86.decoded_buf = []byte{}
	M.x86.enc_pos = 0
}

func print_encoded_bytes(s uint16, o uint16) {
	fmt.Printf("print encode bytes s %x o %x enc-pos %x\n", s, o, M.x86.enc_pos)
	for i := uint16(0); i < uint16(M.x86.enc_pos); i++ {
		fmt.Printf("%02x", fetch_data_byte_abs(s, o+i))
	}
}

func print_decoded_instruction() {
	fmt.Printf("%s", M.x86.decoded_buf)
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
		off   = uint32(o)
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
			i++
		}
		fmt.Printf("\n")
		start = end
		end = start + 16
	}
}

var (
	breakpoint uint16
	noDecode   = true
)

func x86emu_single_step() error {
	var (
		segment, offset uint16
	)

	if DEBUG_BREAK() {
		if M.x86.saved_ip != breakpoint {
			return nil
		} else {
			M.x86.debug &= ^DEBUG_DECODE_NOPRINT_F
			M.x86.debug |= DEBUG_TRACE_F
			M.x86.debug &= ^DEBUG_BREAK_F
			print_decoded_instruction()
			X86EMU_trace_regs()
		}
	}
	var done bool
	offset = M.x86.saved_ip
	for !done {
		fmt.Printf("-")
		cmd, ps, err := parse_line(cmds)
		log.Printf("parse_line: %v %v %v", cmd, ps, err)
		if err != nil {
			return err
		}
		switch cmd {
		case "u":
			disassemble_forward(M.x86.saved_cs, uint16(offset), 10)
		case "d":
			if len(ps) == 1 {
				segment = M.x86.saved_cs
				offset = ps[0]
				X86EMU_dump_memory(segment, uint16(offset), 16)
				offset += 16
			} else if len(ps) == 2 {
				segment = ps[0]
				offset = ps[1]
				X86EMU_dump_memory(segment, uint16(offset), 16)
				offset += 16
			} else {
				segment = M.x86.saved_cs
				X86EMU_dump_memory(segment, uint16(offset), 16)
				offset += 16
			}
		case "c":
			M.x86.debug ^= DEBUG_TRACECALL_F
		case "s":
			M.x86.debug ^= DEBUG_SVC_F | DEBUG_SYS_F | DEBUG_SYSINT_F
		case "r":
			X86EMU_trace_regs()
		case "x":
			X86EMU_trace_xregs()
		case "g":
			if len(ps) == 1 {
				breakpoint = ps[0]
				if noDecode {
					M.x86.debug |= DEBUG_DECODE_NOPRINT_F
				} else {
					M.x86.debug &= ^DEBUG_DECODE_NOPRINT_F
				}
				M.x86.debug &= ^DEBUG_TRACE_F
				M.x86.debug |= DEBUG_BREAK_F
				done = true
			}
		case "q":
			M.x86.exit = true
		case "P":
			noDecode = !noDecode
			fmt.Printf("Toggled decoding to %v\n", noDecode)
		case "t":
		case "":
			done = true
		}
	}
	return nil
}

func X86EMU_trace_on() uint32 {
	log.Printf("ton")
	M.x86.debug = M.x86.debug | DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F
	return M.x86.debug
}

func X86EMU_trace_off() uint32 {
	log.Printf("toff")
	M.x86.debug = M.x86.debug & ^(DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F)
	return M.x86.debug
}

func parse_line(r *bufio.Reader) (string, []uint16, error) {
	data, err := r.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	fields := strings.Fields(data)
	log.Printf("felds %v", fields)
	if len(fields) == 0 {
		return "", nil, nil
	}

	cmd := fields[0]
	var vals []uint16
	for i := range fields[1:] {
		v, err := strconv.ParseUint(fields[i], 0, 16)
		if err != nil {
			return cmd, vals, err
		}
		vals = append(vals, uint16(v))
	}
	return cmd, vals, nil
}

func x86emu_dump_regs() {
	fmt.Printf("\tAX=%04x  ", M.x86.gen.A.Get16())
	fmt.Printf("BX=%04x  ", M.x86.gen.B.Get16())
	fmt.Printf("CX=%04x  ", M.x86.gen.C.Get16())
	fmt.Printf("DX=%04x  ", M.x86.gen.D.Get16())
	fmt.Printf("SP=%04x  ", M.x86.spc.SP.Get16())
	fmt.Printf("BP=%04x  ", M.x86.spc.BP.Get16())
	fmt.Printf("SI=%04x  ", M.x86.spc.SI.Get16())
	fmt.Printf("DI=%04x\n", M.x86.spc.DI.Get16())
	fmt.Printf("\tDS=%04x  ", M.x86.seg.DS.Get())
	fmt.Printf("ES=%04x  ", M.x86.seg.ES.Get())
	fmt.Printf("SS=%04x  ", M.x86.seg.SS.Get())
	fmt.Printf("CS=%04x  ", M.x86.seg.CS.Get())
	fmt.Printf("IP=%04x   ", M.x86.spc.IP.Get16())
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
	fmt.Printf("\tAX=%08x  ", M.x86.gen.A.Get32())
	fmt.Printf("BX=%08x  ", M.x86.gen.B.Get32())
	fmt.Printf("CX=%08x  ", M.x86.gen.C.Get32())
	fmt.Printf("DX=%08x  ", M.x86.gen.D.Get32())
	fmt.Printf("SP=%08x  ", M.x86.spc.SP.Get32())
	fmt.Printf("BP=%08x  ", M.x86.spc.BP.Get32())
	fmt.Printf("SI=%08x  ", M.x86.spc.SI.Get32())
	fmt.Printf("DI=%08x\n", M.x86.spc.DI.Get32())
	fmt.Printf("\tDS=%04x  ", M.x86.seg.DS.Get())
	fmt.Printf("ES=%04x  ", M.x86.seg.ES.Get())
	fmt.Printf("SS=%04x  ", M.x86.seg.SS.Get())
	fmt.Printf("CS=%04x  ", M.x86.seg.CS.Get())
	fmt.Printf("IP=%08x   ", M.x86.spc.IP.Get32())

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

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

#include "x86emui.h"

/*----------------------------- Implementation ----------------------------*/

static void     print_encoded_bytes (u16 s, u16 o);
static void     print_decoded_instruction (void);
int      parse_line (char *s, int *ps, int *n);

/* should look something like debug's output. */
void X86EMU_trace_regs (void)
{
    if (DEBUG_TRACE()) {
	if (mode & (SYSMODE_PREFIX_DATA | SYSMODE_PREFIX_ADDR)) {
	        x86emu_dump_xregs();
	} else {
	        x86emu_dump_regs();
	}
    }
    if (DEBUG_DECODE() && ! DEBUG_DECODE_NOPRINT()) {
        loggy("%04x:%04x ",saved_cs, saved_ip);
        print_encoded_bytes( saved_cs, saved_ip);
        print_decoded_instruction();
    }
}

void X86EMU_trace_xregs (void)
{
    if (DEBUG_TRACE()) {
        x86emu_dump_xregs();
    }
}

void x86emu_just_disassemble (void)
{
    /*
     * This routine called if the flag DEBUG_DISASSEMBLE is set kind
     * of a hack!
     */
    loggy("%04x:%04x ",saved_cs, saved_ip);
    print_encoded_bytes( saved_cs, saved_ip);
    print_decoded_instruction();
}

void disassemble_forward (u16 seg, u16 off, int n)
{
    X86EMU_sysEnv tregs;
    int i;
    u8 op1;
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
    tregs = M;
    tregs.x86.R_IP = off;
    tregs.x86.R_CS = seg;

    /* reset the decoding buffers */
    tregs.x86.enc_str_pos = 0;
    tregs.x86.enc_pos = 0;

    /* turn on the "disassemble only, no execute" flag */
    tregs.x86.debug |= DEBUG_DISASSEMBLE_F;

    /* DUMP NEXT n instructions to screen in straight_line fashion */
    /*
     * This looks like the regular instruction fetch stream, except
     * that when this occurs, each fetched opcode, upon seeing the
     * DEBUG_DISASSEMBLE flag set, exits immediately after decoding
     * the instruction.  XXX --- CHECK THAT MEM IS NOT AFFECTED!!!
     * Note the use of a copy of the register structure...
     */
    for (i=0; i<n; i++) {
        op1 = (*sys_rdb)(((u32)CS<<4) + (IP++));
        x86_byte_dispatch(op1);
    }
    /* end major hack mode. */
}

void x86emu_check_ip_access (void)
{
    /* NULL as of now */
}

void x86emu_check_sp_access (void)
{
}

void x86emu_check_mem_access (u32 dummy)
{
    /*  check bounds, etc */
}

void x86emu_check_data_access (uint dummy1, uint dummy2)
{
    /*  check bounds, etc */
}

void x86emu_inc_decoded_inst_len (int x)
{
    enc_pos += x;
}

void x86emu_decode_log (const char *x)
{
    strcpy(decoded_buf+enc_str_pos,x);
    enc_str_pos += strlen(x);
}

void x86emu_decode_log2 (const char *x, int y)
{
    char temp[100];
    snloggy(temp, sizeof (temp), x,y);
    strcpy(decoded_buf+enc_str_pos,temp);
    enc_str_pos += strlen(temp);
}

void x86emu_end_instr (void)
{
    enc_str_pos = 0;
    enc_pos = 0;
}

static void print_encoded_bytes (u16 s, u16 o)
{
    int i;
    char buf1[64];
    for (i=0; i< enc_pos; i++) {
	    snloggy(buf1+2*i, 64 - 2 * i, "%02x", fetch_data_byte_abs(s,o+i));
    }
    loggy("%-20s ",buf1);
}

static void print_decoded_instruction (void)
{
    loggy("%s", decoded_buf);
}

void x86emu_print_int_vect (u16 iv)
{
    u16 seg,off;

    if (iv > 256) return;
    seg   = fetch_data_word_abs(0,iv*4);
    off   = fetch_data_word_abs(0,iv*4+2);
    loggy("%04x:%04x ", seg, off);
}

void X86EMU_dump_memory (u16 seg, u16 off, u32 amt)
{
    u32 start = off & 0xfffffff0;
    u32 end  = (off+16) & 0xfffffff0;
    u32 i;

    while (end <= off + amt) {
        loggy("%04x:%04x ", seg, start);
        for (i=start; i< off; i++)
          loggy("   ");
        for (       ; i< end; i++)
          loggy("%02x ", fetch_data_byte_abs(seg,i));
        loggy("\n");
        start = end;
        end = start + 16;
    }
}

void x86emu_single_step (void)
{
#if 0
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
                if (saved_ip != breakpoint) {
                        return;
                } else {
              debug &= ~DEBUG_DECODE_NOPRINT_F;
                        debug |= DEBUG_TRACE_F;
                        debug &= ~DEBUG_BREAK_F;
                        print_decoded_instruction ();
                        X86EMU_trace_regs();
                }
        }
    done=0;
    offset = saved_ip;
    while (!done) {
        loggy("-");
        (void)fgets(s, 1023, stdin);
        cmd = parse_line(s, ps, &ntok);
        switch(cmd) {
          case 'u':
            disassemble_forward(saved_cs,(u16)offset,10);
            break;
          case 'd':
                            if (ntok == 2) {
                                    segment = saved_cs;
                                    offset = ps[1];
                                    X86EMU_dump_memory(segment,(u16)offset,16);
                                    offset += 16;
                            } else if (ntok == 3) {
                                    segment = ps[1];
                                    offset = ps[2];
                                    X86EMU_dump_memory(segment,(u16)offset,16);
                                    offset += 16;
                            } else {
                                    segment = saved_cs;
                                    X86EMU_dump_memory(segment,(u16)offset,16);
                                    offset += 16;
                            }
            break;
          case 'c':
            debug ^= DEBUG_TRACECALL_F;
            break;
          case 's':
            debug ^= DEBUG_SVC_F | DEBUG_SYS_F | DEBUG_SYSINT_F;
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
                        debug |= DEBUG_DECODE_NOPRINT_F;
        } else {
                        debug &= ~DEBUG_DECODE_NOPRINT_F;
        }
        debug &= ~DEBUG_TRACE_F;
        debug |= DEBUG_BREAK_F;
        done = 1;
            }
            break;
          case 'q':
          debug |= DEBUG_EXIT;
          return;
      case 'P':
          noDecode = (noDecode)?0:1;
          loggy("Toggled decoding to %s\n",(noDecode)?"FALSE":"TRUE");
          break;
          case 't':
      case 0:
            done = 1;
            break;
        }
    }
#endif
}

int X86EMU_trace_on(void)
{
    return debug |= DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F;
}

int X86EMU_trace_off(void)
{
    return debug &= ~(DEBUG_STEP_F | DEBUG_DECODE_F | DEBUG_TRACE_F);
}

int parse_line (char *s, int *ps, int *n)
{
#if 0
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
#else
    return 0;
#endif
}

void x86emu_dump_regs (void)
{
    loggy("\tAX=%04x  ", AX );
    loggy("BX=%04x  ", BX );
    loggy("CX=%04x  ", CX );
    loggy("DX=%04x  ", DX );
    loggy("SP=%04x  ", SP );
    loggy("BP=%04x  ", BP );
    loggy("SI=%04x  ", SI );
    loggy("DI=%04x\n", DI );
    loggy("\tDS=%04x  ", DS );
    loggy("ES=%04x  ", ES );
    loggy("SS=%04x  ", SS );
    loggy("CS=%04x  ", CS );
    loggy("IP=%04x   ", IP );
    if (ACCESS_FLAG(F_OF))    loggy("OV ");     /* CHECKED... */
    else                        loggy("NV ");
    if (ACCESS_FLAG(F_DF))    loggy("DN ");
    else                        loggy("UP ");
    if (ACCESS_FLAG(F_IF))    loggy("EI ");
    else                        loggy("DI ");
    if (ACCESS_FLAG(F_SF))    loggy("NG ");
    else                        loggy("PL ");
    if (ACCESS_FLAG(F_ZF))    loggy("ZR ");
    else                        loggy("NZ ");
    if (ACCESS_FLAG(F_AF))    loggy("AC ");
    else                        loggy("NA ");
    if (ACCESS_FLAG(F_PF))    loggy("PE ");
    else                        loggy("PO ");
    if (ACCESS_FLAG(F_CF))    loggy("CY ");
    else                        loggy("NC ");
    loggy("\n");
}

void x86emu_dump_xregs (void)
{
    loggy("\tEAX=%08x  ", EAX );
    loggy("EBX=%08x  ", EBX );
    loggy("ECX=%08x  ", ECX );
    loggy("EDX=%08x\n", EDX );
    loggy("\tESP=%08x  ", ESP );
    loggy("EBP=%08x  ", EBP );
    loggy("ESI=%08x  ", ESI );
    loggy("EDI=%08x\n", EDI );
    loggy("\tDS=%04x  ", DS );
    loggy("ES=%04x  ", ES );
    loggy("SS=%04x  ", SS );
    loggy("CS=%04x  ", CS );
    loggy("EIP=%08x\n\t", EIP );
    if (ACCESS_FLAG(F_OF))    loggy("OV ");     /* CHECKED... */
    else                        loggy("NV ");
    if (ACCESS_FLAG(F_DF))    loggy("DN ");
    else                        loggy("UP ");
    if (ACCESS_FLAG(F_IF))    loggy("EI ");
    else                        loggy("DI ");
    if (ACCESS_FLAG(F_SF))    loggy("NG ");
    else                        loggy("PL ");
    if (ACCESS_FLAG(F_ZF))    loggy("ZR ");
    else                        loggy("NZ ");
    if (ACCESS_FLAG(F_AF))    loggy("AC ");
    else                        loggy("NA ");
    if (ACCESS_FLAG(F_PF))    loggy("PE ");
    else                        loggy("PO ");
    if (ACCESS_FLAG(F_CF))    loggy("CY ");
    else                        loggy("NC ");
    loggy("\n");
}

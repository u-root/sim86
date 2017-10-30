#define CHECK_IP_FETCH_F                0x1
#define CHECK_SP_ACCESS_F               0x2
#define CHECK_MEM_ACCESS_F              0x4 /*using regular linear pointer */
#define CHECK_DATA_ACCESS_F             0x8 /*using segment:offset*/

# define CHECK_IP_FETCH()              	(M.x86.check & CHECK_IP_FETCH_F)
# define CHECK_SP_ACCESS()             	(M.x86.check & CHECK_SP_ACCESS_F)
# define CHECK_MEM_ACCESS()            	(M.x86.check & CHECK_MEM_ACCESS_F)
# define CHECK_DATA_ACCESS()           	(M.x86.check & CHECK_DATA_ACCESS_F)

# define DEBUG_INSTRUMENT()    	(M.x86.debug & DEBUG_INSTRUMENT_F)
# define DEBUG_DECODE()        	(M.x86.debug & DEBUG_DECODE_F)
# define DEBUG_TRACE()         	(M.x86.debug & DEBUG_TRACE_F)
# define DEBUG_STEP()          	(M.x86.debug & DEBUG_STEP_F)
# define DEBUG_DISASSEMBLE()   	(M.x86.debug & DEBUG_DISASSEMBLE_F)
# define DEBUG_BREAK()         	(M.x86.debug & DEBUG_BREAK_F)
# define DEBUG_SVC()           	(M.x86.debug & DEBUG_SVC_F)
# define DEBUG_SAVE_IP_CS()     (M.x86.debug & DEBUG_SAVE_IP_CS_F)

# define DEBUG_FS()            	(M.x86.debug & DEBUG_FS_F)
# define DEBUG_PROC()          	(M.x86.debug & DEBUG_PROC_F)
# define DEBUG_SYSINT()        	(M.x86.debug & DEBUG_SYSINT_F)
# define DEBUG_TRACECALL()     	(M.x86.debug & DEBUG_TRACECALL_F)
# define DEBUG_TRACECALLREGS() 	(M.x86.debug & DEBUG_TRACECALL_REGS_F)
# define DEBUG_TRACEJMP()       (M.x86.debug & DEBUG_TRACEJMP_F)
# define DEBUG_TRACEJMPREGS()   (M.x86.debug & DEBUG_TRACEJMP_REGS_F)
# define DEBUG_SYS()           	(M.x86.debug & DEBUG_SYS_F)
# define DEBUG_MEM_TRACE()     	(M.x86.debug & DEBUG_MEM_TRACE_F)
# define DEBUG_IO_TRACE()      	(M.x86.debug & DEBUG_IO_TRACE_F)
# define DEBUG_DECODE_NOPRINT() (M.x86.debug & DEBUG_DECODE_NOPRINT_F)

# define DECODE_PRINTF(x)     	if (DEBUG_DECODE()) \
									x86emu_decode_printf(x)
# define DECODE_PRINTF2(x,y)  	if (DEBUG_DECODE()) \
									x86emu_decode_printf2(x,y)

/*
 * The following allow us to look at the bytes of an instruction.  The
 * first INCR_INSTRN_LEN, is called every time bytes are consumed in
 * the decoding process.  The SAVE_IP_CS is called initially when the
 * major opcode of the instruction is accessed.
 */
#define INC_DECODED_INST_LEN(x)                    	\
	if (DEBUG_DECODE())  	                       	\
		x86emu_inc_decoded_inst_len(x)

#define SAVE_IP_CS(x,y)                               			\
	if (DEBUG_DECODE() | DEBUG_TRACECALL() | DEBUG_BREAK() \
              | DEBUG_IO_TRACE() | DEBUG_SAVE_IP_CS()) { \
		M.x86.saved_cs = x;                          			\
		M.x86.saved_ip = y;                          			\
	}
#define TRACE_REGS()                                   		\
	if (DEBUG_DISASSEMBLE()) {                         		\
		x86emu_just_disassemble();                        	\
		goto EndOfTheInstructionProcedure;             		\
	}                                                   	\
	if (DEBUG_TRACE() || DEBUG_DECODE()) X86EMU_trace_regs()

# define SINGLE_STEP()		if (DEBUG_STEP()) x86emu_single_step()

#define TRACE_AND_STEP()	\
	TRACE_REGS();			\
	SINGLE_STEP()

# define START_OF_INSTR()
# define END_OF_INSTR()		EndOfTheInstructionProcedure: x86emu_end_instr();
# define END_OF_INSTR_NO_TRACE()	x86emu_end_instr();

# define  CALL_TRACE(u,v,w,x,s)                                 \
	if (DEBUG_TRACECALLREGS())									\
		x86emu_dump_regs();                                     \
	if (DEBUG_TRACECALL())                                     	\
		printf("%04x:%04x: CALL %s%04x:%04x\n", u , v, s, w, x);
# define RETURN_TRACE(u,v,w,x,s)                                    \
	if (DEBUG_TRACECALLREGS())									\
		x86emu_dump_regs();                                     \
	if (DEBUG_TRACECALL())                                     	\
		printf("%04x:%04x: RET %s %04x:%04x\n",u,v,s,w,x);
# define  JMP_TRACE(u,v,w,x,s)                                 \
   if (DEBUG_TRACEJMPREGS()) \
      x86emu_dump_regs(); \
   if (DEBUG_TRACEJMP()) \
      printf("%04x:%04x: JMP %s%04x:%04x\n", u , v, s, w, x);

#define	DB(x)	x

#define X86EMU_DEBUG_ONLY(x) x

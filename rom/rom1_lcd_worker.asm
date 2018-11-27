LCD_FLAG		EQU	0xE800
LCD_DATA		EQU	0xE801
LCD_LOCAL		EQU	0x1000
LAMP_SRC		EQU	0xEA00
LAMP_DEST		EQU	0x9500
SIO_CMD_NULL		EQU	0x00
SIO_CMD_TX_ABORT	EQU	0x08
SIO_CMD_RST_STATUS	EQU	0x10
SIO_CMD_CH_RESET	EQU	0x18
SIO_CMD_EI_RX		EQU	0x20
SIO_CMD_TX_IR_RST	EQU	0x28
SIO_CMD_ERROR_RST	EQU	0x30
SIO_CMD_RETI		EQU	0x38
ORG 0x0

START:
	DI
	IM	2
	LD	A,INTERRUPT_VECTORS / 0x0100
	LD	I,A
	LD	SP,0E7FFh	; Stack pointer
	CALL	INIT1		; Sets up IO F4, F0 and F1, has a big delay in the end
	JR	INIT2
	; This region had what looks like a second start for the CPU?

	; -- L0019 --
INIT2:
	CALL	INIT_EI_RETI
	CALL	INIT_CTC
	CALL	INIT_SIO
	CALL	INIT_PIO
	; CALL	INIT_4B ; Unknown, possibly unpopulated IO?
	CALL	SHARED_MEM_INIT
	; CALL	INIT_44_45
	; CALL	INIT_48_49
	CALL	LCD_INIT
	; CALL	L03F7 ; Zeroes 9416h, not ported
	; CALL	L0466 ; Sets 8003h - 8008h memory up
	; CALL	L0486 ; Sets up 920Dh - 9212h
	; CALL	L0494 ; Sets up 930Dh - 9312h
	; CALL	L0476 ; Sets up 8908h - 890Eh
	LD	A,0FFh
	LD	(9311h),A
	LD	HL,0FFFFh	; address or value?
	LD	(890Ch),HL
	; EI
	JP	MAIN_LOOP

	SEEK 0x70
	ORG 0x70
	; Interrupts
INTERRUPT_VECTORS:
	DW	SIO_B_TX_EMPTY
	DW	SIO_B_STATUS_CHANGE
	DW	SIO_B_RX_AVAILABLE
	DW	SIO_B_ERROR
	DW	SIO_A_TX_EMPTY
	DW	SIO_A_STATUS_CHANGE
	DW	SIO_A_RX_AVAILABLE
	DW	SIO_A_ERROR
CTC_VECTORS:
	DW	CTC_CH1
	DW	CTC_CH2
	DW	CTC_CH3
	DW	CTC_CH4
PIO_VECTORS:
	DW	PIO_A
	DW	PIO_B

	SEEK 0x90
	ORG 0x90
	; -- L00C8 --
INIT1: PROC
	LD	A,00h
	OUT	(0F4h),A
	CALL	INIT_F0_F1
	LD	HL,012Ch
loop1:	LD	B,00h
loop2:	DJNZ	loop2
	DEC	HL
	LD	A,H
	OR	L
	JR	NZ,loop1
	RET
ENDP

MAIN_LOOP: PROC
	LD	HL,LCD_SPLASH
	CALL	LCD_UPDATE
	CALL	LAMP_COPY	; Copy 1Bh bytes from LAMP_SRC to LAMP_DEST
	; CALL	L5783		; Goes to the big jump table
	; CALL	L00E5		; Send ping request, zero memory flags
loop:	; LD	A,(9417h)	; Check flag 9417h, if not zero zero it and execute
	; AND	A
	; JR	Z,loopend
	; LD	A,00h
	; LD	(9417h),A
	CALL	LAMP_COPY	; Copy 1Bh bytes from LAMP_SRC to LAMP_DEST
	; CALL	L05D3		; Goes to the big jump table
	CALL	LAMP_UPDATE	; Output 1Bh bytes from LAMP_DEST + 0x1A down to ports 00h - 03h
	; CALL	L0409		; Conditionally sends PING request
loopend:
	; CALL	L57BB		; Huge conditional tree
	CALL	LCD_COPY	; Update LCD from shared memory
	JP	loop
ENDP

	; --- L5783 ---
	; Shifts out 0x1B bytes to the various lamps
LAMP_UPDATE: PROC
	LD	HL,LAMP_DEST + 0x1A
	LD	B,1Bh
loop:	LD	A,(HL)
	OUT	(01h),A
	OUT	(00h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	OUT	(03h),A
	DEC	HL
	DJNZ	loop
	OUT	(02h),A
	RET
ENDP

	; --- L57A3 ---
LAMP_COPY:
	LD	HL,LAMP_SRC
	LD	DE,LAMP_DEST
	LD	BC,001Bh
	LDIR
	RET

	; --- L0304 ---
	; Unknown, possibly unpopulated devices
INIT_48_49: PROC
	LD	A,0FFh
	OUT	(49h),A ; 0xFF
	LD	D,00h
	LD	B,08h
	LD	A,D
	OUT	(48h),A ; 0x00
	LD	A,0FCh
	OUT	(49h),A ; 0xFC
	AND	0FBh
	OUT	(49h),A ; 0xF8
	LD	A,0FFh
	OUT	(49h),A ; 0xFF
loop1:	LD	A,0FFh
	OUT	(48h),A ; 0xFF
	LD	A,0FEh
	OUT	(49h),A ; 0xFE
	AND	0FBh
	OUT	(49h),A ; 0xFB
	LD	A,0FFh
	OUT	(49h),A ; 0xFF
	INC	D
	DJNZ	loop1	; 8 times
	LD	D,00h
	LD	HL,INIT_48_DATA
	LD	B,08h
loop2:	LD	C,03h
	LD	A,D
	OUT	(48h),A	; 0x00 - 0x08
	LD	A,0FCh
	OUT	(49h),A	; 0xFC
	AND	0FBh
	OUT	(49h),A ; 0xFB
	LD	A,0FFh
	OUT	(49h),A ; 0xFF
loop3:	LD	A,(HL)
	OUT	(48h),A ; Data from table
	LD	A,0FDh
	OUT	(49h),A ; 0xFD
	AND	0FBh
	OUT	(49h),A ; 0xFB
	LD	A,0FFh
	OUT	(49h),A ; 0xFF
	INC	HL
	DEC	C
	JP	NZ,loop3	; 3 times
	INC	D
	DJNZ	loop2		; 8 times
	RET
ENDP

INIT_48_DATA:
	DB	03h
	DB	05h
	DB	14h
	DB	3Fh
	DB	00h
	DB	00h
	DB	00h
	DB	3Fh
	DB	00h
	DB	3Fh
	DB	3Fh
	DB	00h
	DB	00h
	DB	00h
	DB	3Fh
	DB	3Fh
	DB	00h
	DB	3Fh
	DB	00h
	DB	3Fh
	DB	3Fh
	DB	3Fh
	DB	3Fh
	DB	3Fh

	; --- L02E2 ---
	; Unknown device init? writes incrementing bytes to 44 and data to 45
INIT_44_45: PROC
	LD	D,00h
	LD	HL,INIT_45_DATA
	LD	B,10h
loop:	LD	A,D
	OUT	(44h),A
	LD	A,(HL)
	OUT	(45h),A
	INC	D
	INC	HL
	DJNZ	loop
	RET
ENDP

INIT_45_DATA:
	DB	64h
	DB	50h
	DB	52h
	DB	2Ch
	DB	1Ch
	DB	01h
	DB	19h
	DB	1Ah
	DB	50h
	DB	0Fh
	DB	20h
	DB	0Fh
	DB	0F8h
	DB	00h
	DB	0F8h
	DB	00h

	; --- L03E2 ---
	; Not entirely sure about this, seems strange...
SHARED_MEM_INIT: PROC
	LD	HL,7D33h
loop:	DEC	HL
	LD	A,H
	OR	L
	JR	NZ,loop
	LD	HL,0F000h
	LD	A,(HL)
	LD	HL,0F800h
	LD	A,(HL)
	LD	HL,0E800h
	LD	A,(HL)
	RET
ENDP

	; --- L02D7 ---
INIT_4B:
	LD	HL,INIT_4B_DATA
	LD	C,4Bh
	LD	B,01h
	OTIR
	RET

INIT_4B_DATA:
	DB	80h

	; --- L0136 ---
INIT_PIO:
	LD	HL,PIO_A_DATA
	LD	C,1Dh
	LD	B,03h
	OTIR
	LD	HL,PIO_B_DATA
	LD	C,1Fh
	LD	B,03h
	OTIR
	RET

PIO_A_DATA:
	DB	0CFh	; Mode 3
	DB	00h	; All pins outputs
	DB	07h	; No interrupts
PIO_B_DATA:
	DB	0CFh	; Mode 3
	DB	00h	; All pins outputs
	DB	07h	; No interrupts

	; --- L015A ---
INIT_SIO:
	LD	HL,SIO_A_DATA
	LD	C,19h
	LD	B,0Dh
	OTIR
	LD	HL,SIO_B_DATA
	LD	C,1Bh
	LD	B,0Fh
	OTIR
	LD	HL,0000h
	LD	(9411h),HL
	LD	A,01h
	LD	(9413h),A
	RET

SIO_A_DATA:
	DB	00h	; Write register 0
	DB	00h
	DB	18h	; CMD error reset
	DB	10h	; CMD reset status interrupts
	DB	04h	; Write register 4
	DB	45h	; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
	DB	03h	; Write register 3
	DB	0C1h	; RX enable, 8 bits per character
	DB	05h	; Write register 5
	DB	0EAh	; TX crc disable, RTS low, TX enable, DTR low, external sync
	DB	10h	; CMD reset status interrupts
	DB	01h	; Write register 1
	DB	12h	; Transmit interrupt enable, receive interrupt on all characters
SIO_B_DATA:
	DB	00h	; Write register 0
	DB	00h
	DB	18h	; CMD error reset
	DB	10h	; CMD reset status interrupts
	DB	04h	; Write register 4
	DB	45h	; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
	DB	03h	; Write register 3
	DB	0C1h	; RX enable, 8 bits per character
	DB	05h	; Write register 5
	DB	0EAh	; TX crc disable, RTS low, TX enable, DTR low, external sync
	DB	02h	; Write register 2
	DB	INTERRUPT_VECTORS % 0x0100
	; DB	70h
	DB	10h	; CMD reset status interrupts
	DB	01h	; Write register 1
	DB	16h	; Transmit interrupt enable, receive interrupt on all characters, status affects vector


	; --- L0194 ---
INIT_CTC:
	LD	HL,CTC_1_DATA
	LD	C,10h
	LD	B,03h
	OTIR
	LD	HL,CTC_2_DATA	; address or value?
	LD	C,11h
	LD	B,02h
	OTIR
	LD	HL,CTC_3_DATA	; address or value?
	LD	C,12h
	LD	B,02h
	OTIR
	LD	HL,CTC_4_DATA	; address or value?
	LD	C,13h
	LD	B,02h
	OTIR
	RET

CTC_1_DATA:
	DB	CTC_VECTORS % 0x0100	; Interrupt vector location
	DB	0DDh	; Interrupt enabled, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	0Ch	; Time constant
CTC_2_DATA:
	DB	5Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	01h	; Time constant
CTC_3_DATA:
	DB	0A7h	; Interrupt enabled, timer, prescaler 16, falling edge, time constant follow, software reset
	DB	0A3h	; Time constant
CTC_4_DATA:
	DB	5Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	01h	; Time constant


	; --- L00DC ---
INIT_F0_F1:
	LD	A,5Bh
	OUT	(0F0h),A
	LD	A,0B1h
	OUT	(0F1h),A
	RET

	; --- L00B8 ---
INIT_EI_RETI:
	LD	A,0FBh
	LD	(8000h),A
	LD	A,0EDh
	LD	(8001h),A
	LD	A,4Dh
	LD	(8002h),A
	RET


	; --- L0376 ---
	; Init both LCD controllers
	; Set data length 8, lines 2, display on, cursor off, blink off,
	; auto-increment writes, no display shift
LCD_INIT:
	LD	A,30h
	CALL	LCD1_REGISTER
	CALL	LCD_DELAY3
	LD	A,30h
	CALL	LCD1_REGISTER
	CALL	LCD_DELAY3
	LD	A,30h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY3
	LD	A,30h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY3
	LD	A,30h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY3
	LD	A,30h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY3
	LD	A,38h
	CALL	LCD1_REGISTER
	CALL	LCD_DELAY2
	LD	A,38h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY2
	LD	A,0Ch
	CALL	LCD1_REGISTER
	CALL	LCD_DELAY2
	LD	A,0Ch
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY2
	LD	A,06h
	CALL	LCD1_REGISTER
	CALL	LCD_DELAY2
	LD	A,06h
	CALL	LCD2_REGISTER
	CALL	LCD_DELAY2
	LD	HL,9450h
	LD	(94F0h),HL
	LD	A,00h
	LD	(94F2h),A
	RET

	; --- L0DE0 ---
	; Writes A to LCD controller 1 registers
LCD1_REGISTER:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	OR	01h
	OUT	(1Eh),A
	LD	A,D
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	CALL	LCD_DELAY2
	POP	DE
	RET

	; --- L0DFE ---
	; Writes A to LCD controller 2 registers
LCD2_REGISTER:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	OR	02h
	OUT	(1Eh),A
	LD	A,D
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	CALL	LCD_DELAY2
	POP	DE
	RET

	; Reads 160 bytes starting from HL and writes to the LCD
LCD_UPDATE: PROC
	LD	A,00h
	CALL	SET_LCD1_ADDR
	LD	B,28h
loop1:	LD	A,(HL)
	CALL	WRITE_LCD1
	INC	HL
	DJNZ	loop1
	LD	A,40h	; 64 address for the second row
	CALL	SET_LCD1_ADDR
	LD	B,28h
loop2:	LD	A,(HL)
	CALL	WRITE_LCD1
	INC	HL
	DJNZ	loop2
	LD	A,00h
	CALL	SET_LCD2_ADDR
	LD	B,28h	; '('
loop3:	LD	A,(HL)
	CALL	WRITE_LCD2
	INC	HL
	DJNZ	loop3
	LD	A,78h	; 'x'
	SUB	50h	; 'P'
	ADD	A,18h
	CALL	SET_LCD2_ADDR
	LD	B,28h	; '('
loop4:	LD	A,(HL)
	CALL	WRITE_LCD2
	INC	HL
	DJNZ	loop4
	RET
ENDP

; Set LCD1 DDRAM address to A
SET_LCD1_ADDR:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	OR	01h
	OUT	(1Eh),A
	LD	A,D
	AND	7Fh
	OR	80h
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	CALL	LCD_DELAY1
	POP	DE
	RET

; SET LCD2 DDRAM address to A
SET_LCD2_ADDR:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	OR	02h
	OUT	(1Eh),A
	LD	A,D
	AND	7Fh
	OR	80h
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	00h
	OUT	(1Eh),A
	CALL	LCD_DELAY1
	POP	DE
	RET

; Write A to LCD1 DDRAM
WRITE_LCD1:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	04h
	OUT	(1Eh),A
	OR	01h
	OUT	(1Eh),A
	LD	A,D
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	04h
	OUT	(1Eh),A
	CALL	LCD_DELAY1
	POP	DE
	RET

; Write A to LCD2 DDRAM
WRITE_LCD2:
	PUSH	DE
	LD	D,A
	IN	A,(1Eh)
	AND	0F0h
	OR	04h
	OUT	(1Eh),A
	OR	02h
	OUT	(1Eh),A
	LD	A,D
	OUT	(1Ch),A
	IN	A,(1Eh)
	AND	0F0h
	OR	04h
	OUT	(1Eh),A
	CALL	LCD_DELAY1
	POP	DE
	RET

; Delay loop for the LCD control
LCD_DELAY1: PROC
	PUSH	BC
	LD	B,19h
loop:	DJNZ	loop
	POP	BC
	RET
ENDP

	; --- L0F23 ---
LCD_DELAY2: PROC
	PUSH	BC
	LD	B,00h
loop:	DJNZ	loop
	POP	BC
	RET
ENDP

	; --- L0F0B ---
LCD_DELAY3: PROC
	PUSH	BC
	LD	B,00h
loop1:	DJNZ	loop1
loop2:	DJNZ	loop2
loop3:	DJNZ	loop3
loop4:	DJNZ	loop4
loop5:	DJNZ	loop5
	POP	BC
	RET
ENDP

	; Check LCD_FLAG location for non-zero value
	; If so zero it and copy the new LCD data down
	; to LCD_LOCAL and update the LCD
LCD_COPY: PROC
	LD	HL, LCD_FLAG
	LD	A,(HL)
	AND	A
	JR	Z, return
	LD	A, 0x00
	LD	(HL), A
	LD	HL,LCD_DATA
	LD	DE,LCD_LOCAL
	LD	BC,0xA0		; 160 bytes
	LDIR
	LD	HL,LCD_LOCAL
	CALL	LCD_UPDATE
return:	RET
ENDP

	; -- INTERRUPT HANDLERS --
SIO_B_TX_EMPTY:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_B_STATUS_CHANGE:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_B_RX_AVAILABLE:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_B_ERROR:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	LD	A, SIO_CMD_ERROR_RST
	OUT (1Bh), A
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_A_TX_EMPTY:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_A_STATUS_CHANGE:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_A_RX_AVAILABLE:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

SIO_A_ERROR:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	LD	A, SIO_CMD_ERROR_RST
	OUT (18h), A
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

CTC_CH1:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

CTC_CH2:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

CTC_CH3:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

CTC_CH4:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

PIO_A:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

PIO_B:
	DI
	PUSH AF
	PUSH BC
	PUSH DE
	PUSH HL
	POP HL
	POP DE
	POP BC
	POP AF
	JP 0x8000

LCD_SPLASH:
	DB "****************************************"
	DB "*	VIDEO PRODUCTION SYSTEM VPS-300P   *"
	DB "*	   HACKED BY HELSINKI HACKLAB	   *"
	DB "****************************************"

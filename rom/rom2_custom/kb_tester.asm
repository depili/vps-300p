; SIO commands
SIO_CMD_NULL			EQU	0x00
SIO_CMD_TX_ABORT		EQU	0x08
SIO_CMD_RST_STATUS		EQU	0x10
SIO_CMD_CH_RESET		EQU	0x18
SIO_CMD_EI_RX			EQU	0x20
SIO_CMD_TX_IR_RST		EQU	0x28
SIO_CMD_ERROR_RST		EQU	0x30
SIO_CMD_RETI			EQU	0x38

; Serial receive
RX_COUNTER			EQU	0x9000
RX_POINTER			EQU	0x9002  ; Write pointer
RX_TYPE				EQU	0x9004  ; first byte of the message
RX_MAX_BYTES			EQU	0xFF - 0x06

; Serial transmit, 0xFF byte ring buffer
TX_COUNTER			EQU	0x9100	; Bytes in the TX buffer
TX_PTR_WRITE			EQU	0x9102	; Write pointer
TX_PTR_READ			EQU	0x9104	; Read pointer
TX_BUFFER_START			EQU	0x9200	; Start of the TX buffer, needs to be 0xXX00 as we use 8 bit compare

; LCD
LCD_FLAG			EQU	0xF800
LCD_SRC				EQU	0xC000  ; Source for LCD data
LCD_DEST			EQU	0xF801  ; Shared memory copy destination for lcd data
LCD_BYTES			EQU	0xA0    ; 40 bytes per line, 4 lines = 160 = 0xA0 bytes
LCD_POINTER			EQU	0x8100  ; Pointer for writing bytes to the LCD_SRC buffer
LCD_WRITE_DEST			EQU	LCD_DEST

; Lamps
LAMP_SRC			EQU	0xEA00  ; Source for lamp data in shared memory
LAMP_DEST			EQU	0x9500  ; Local memory destination for lamp data
LAMP_BYTES			EQU	0x1B    ; 27 bytes for lamp data

; Keyboard commands
KEYB_CMD_END			EQU	0xE0
KEYB_CMD_WRITE_DISPLAY  	EQU	0x90    ; Auto increment from address 0
KEYB_CMD_READ_DATA		EQU	0x50    ; Auto increment from address 0
; Keyboard memory
KEYB_CMD_BYTE			EQU	0xB000	; Command byte to send with encoded key event
KEYB_DATA_BYTE			EQU	0xB001	; Data byte goes here
KEYB_1_DEST			EQU	0xB100	; New data goes here
KEYB_1_OLD			EQU	0xB108	; Old data to compare
KEYB_1_DISP			EQU	0xB110
PATTERN_0			EQU	0xB110	; Pattern number, least significant digit
PATTERN_1			EQU	0xB111
PATTERN_2			EQU	0xB112	; Pattern number, most significant digit
FRAMERATE_0			EQU	0xB115	; Frame rate, least significant digit
FRAMERATE_1			EQU	0xB114
FRAMERATE_3			EQU	0xB113	; Frame rate, most significant digit
KEYB_2_DEST			EQU	0xB200
KEYB_2_OLD			EQU	0xB208
KEYB_2_DISP			EQU	0xB210
KEYB_3_DEST			EQU	0xB300
KEYB_3_OLD			EQU	0xB308
KEYB_3_DISP			EQU	0xB310

; ADC
ADC_CMD_TRIGGER			EQU	0xFF
ADC_0_DEST			EQU	0xA000	; Hue
ADC_1_DEST			EQU	0xA002	; Joystick Y (up-down)
ADC_2_DEST			EQU	0xA004	; Clip
ADC_3_DEST			EQU	0xA006	; T-bar
ADC_4_DEST			EQU	0xA008	; Gain
ADC_5_DEST			EQU	0xA00A	; Joystick Z (rotate)
ADC_6_DEST			EQU	0xA00C	; Joystick X (left-right)
ADC_READ_COUNTER		EQU	0xA00E

; LCD numbers and letters
SEG_OFF				EQU	0xFF	; All off
SEG_0				EQU	0x40	; 0
SEG_1				EQU	0x79	; 1
SEG_2				EQU	0x24	; 2
SEG_3				EQU	0x30	; 3
SEG_4				EQU	0x19	; 4
SEG_5				EQU	0x12	; 5
SEG_6				EQU	0x02	; 6
SEG_7				EQU	0x78	; 7
SEG_8				EQU	0x80	; 8
SEG_9				EQU	0x10	; 9
SEG_A				EQU	0x08	; A
SEG_B				EQU	0x03	; B
SEG_C				EQU	0x46	; C
SEG_D				EQU	0x21	; d
SEG_E				EQU	0x06	; E
SEG_F				EQU	0x0E	; F

; IO PORTS
KEYB_1_DATA			EQU	0x00
KEYB_1_CMD			EQU	0x01
KEYB_2_DATA			EQU	0x02
KEYB_2_CMD			EQU	0x03
KEYB_3_DATA			EQU	0x04
KEYB_3_CMD			EQU	0x05
ADC_MUX				EQU	0x08
ADC_IO				EQU	0x09
CTC_CH1				EQU	0x10
CTC_CH2				EQU	0x11
CTC_CH3				EQU	0x12
CTC_CH4				EQU	0x13
SIO_A_DATA			EQU	0x18
SIO_A_CMD			EQU	0x19
SIO_B_DATA			EQU	0x1A
SIO_B_CMD			EQU	0x1B
PIO_A_DATA			EQU	0x1C
PIO_A_CMD			EQU	0x1D
PIO_B_DATA			EQU	0x1E
PIO_B_CMD			EQU	0x1F

TX_A_B(data) MACRO
	LD	A, data
	CALL	TX_BUF_WRITE
ENDM

ORG 0x000
START:
	DI
	IM	2
	LD	A,00h
	LD	I,A
	LD	SP,0F7FFh	; Stack
	CALL	TX_INIT
	TX_A_B	"H"
	TX_A_B	"e"
	TX_A_B	"l"
	TX_A_B	"l"
	TX_A_B	"o"
	LD	A, (TX_COUNTER)
	CALL	DIG2
	LD	A, "\n"
	CALL	OUTCHAR

	CALL	TX_BUF_READ
	CALL	OUTCHAR

	CALL	TX_BUF_READ
	CALL	OUTCHAR

	CALL	TX_BUF_READ
	CALL	OUTCHAR

	CALL	TX_BUF_READ
	CALL	OUTCHAR

	CALL	TX_BUF_READ
	CALL	OUTCHAR

	LD	A, "\n"
	CALL	OUTCHAR

	ZFILL	0xB000, 0x0FFF

	LD	IX, KEYB_1_DEST
	LD	(IX), 0x00
	LD	(IX+1), 0x00
	LD	(IX+2), 0x00
	LD	(IX+3), 0x00
	LD	(IX+4), 0x00
	LD	(IX+5), 0x00
	LD	(IX+6), 0x00
	LD	(IX+7), 0x80
	LD	IY, KEYB_1_OLD
	LD	A, 0x81
	LD	(KEYB_CMD_BYTE), A
	CALL	KB_PROCESS

	CALL	TX_BUF_READ
	CALL	DIG2

	CALL	TX_BUF_READ
	CALL	DIG2

	LD	A, (TX_COUNTER)
	CALL	DIG2

	DI
	HALT

	; IX = new data
	; IY = old data
	; C = Counter for the checked byte, later encoded into the keyboard event byte
	; Messes up AF
KB_PROCESS: PROC
	LD	C, 0x00		; Byte being checked
loop:	LD	B, (IX)
	LD	A, (IY)
	XOR	B
	CALL	NZ, byte_diff
	INC	C
	LD	A, C
	CP	0x08
	RET	Z
	INC	IX
	INC	IY
	JR	loop
ENDP

	; A - old reading
	; B - new new reading
	; E - Encoded bit number and key up/down bit
	; Messes up AF
	; Do not corrupt: IX, IY, C
BYTE_DIFF: PROC
	BIT	0, A
	JR	NZ, bit_0
	BIT	1, A
	JR	NZ, bit_1
	BIT	2, A
	JR	NZ, bit_2
	BIT	3, A
	JR	NZ, bit_3
	BIT	4, A
	JR	NZ, bit_4
	BIT	5, A
	JR	NZ, bit_5
	BIT	6, A
	JR	NZ, bit_6
	BIT	7, A
	JR	NZ, bit_7
	RET
bit_0:
	KB_BIT	0
bit_1:
	KB_BIT	1
bit_2:
	KB_BIT	2
bit_3:
	KB_BIT	3
bit_4:
	KB_BIT	4
bit_5:
	KB_BIT	5
bit_6:
	KB_BIT	6
bit_7:
	KB_BIT	7
ENDP

KB_BIT(chk_bit) MACRO
	LD	E, 0 + (chk_bit << 3)		; Bit code
	BIT	chk_bit, B
	JR	Z, send
	SET	6, E
send:
	CALL	SEND_KB
	RES	chk_bit, A
	JP	BYTE_DIFF
ENDM

SEND_KB:
	PUSH	AF
	PUSH	BC
	LD	A, (KEYB_CMD_BYTE)
	CALL	TX_BUF_WRITE
	LD	A, C
	OR	E
	CALL	TX_BUF_WRITE
	POP	BC
	POP	AF
	RET

TX_INIT:
	LD	A, 0x00
	LD	(TX_COUNTER), A
	LD	HL, TX_BUFFER_START
	LD	(TX_PTR_WRITE), HL
	LD	(TX_PTR_READ), HL
	RET

	; Write A to the TX ring buffer
	; Destroys AF, D, HL
TX_BUF_WRITE: PROC
	DI
	LD	D, A			; Tuck away the byte to be written
	LD	A, (TX_COUNTER)
	CP	0xFF
	JR	Z, return		; Buffer full
	LD	A, D			; Restore the byte to be written
	LD	HL, (TX_PTR_WRITE)
	LD	(HL), A
	INC	L			; Increment the pointer, will loop on the last byte
	LD	(TX_PTR_WRITE), HL
	LD	HL, TX_COUNTER
	INC	(HL)
	CALL	SIO_A_TX_EI		; Enable the transmit buffer empty interrupt
return: EI
	RET
ENDP

	; Read from the TX buffer
	; If the buffer becomes empty disable the TX interrupt
	; Destroys: AF, HL
TX_BUF_READ:
	LD	HL, (TX_PTR_READ)
	LD	A, (HL)
	INC	L
	LD	(TX_PTR_READ), HL
	LD	HL, TX_COUNTER
	DEC	(HL)
	CALL	Z, SIO_A_TX_DI		; Zero bytes in the buffer, disable TX interrupt
	RET

SIO_A_TX_EI:
	LD	A, 0x01
	OUT	(SIO_A_CMD), A		; Write register 1
	LD	A, 0x12
	OUT	(SIO_A_CMD), A		; RX interrupt on all characters, TX interrupt on
	RET

SIO_A_TX_DI:
	PUSH	AF
	LD	A, 0x01
	OUT	(SIO_A_CMD), A		; Write register 1
	LD	A, 0x10
	OUT	(SIO_A_CMD), A		; RX interrupt on all characters
	POP	AF
	RET


DIG2:
	PUSH	AF		; Nibble-Paar ausdrucken
	RRCA			; unterstes Nibble retten
	RRCA			; oberes Nibble rechtsbuendig
	RRCA			; positionieren
	RRCA
	AND	00001111b
	ADD	A,90H		; binaer in ASCII (hex)
	DAA
	ADC	A,40H
	DAA
	CALL	OUTCHAR		; Zeichen ausgeben
	POP	AF		; jetzt unteres Nibble verarbeiten
	AND	00001111b	; Nibble maskieren
	ADD	A,90H		; binaer in ASCII (hex)
	DAA
	ADC	A,40H
	DAA
	CALL	OUTCHAR
	RET

OUTCHAR:			; Zeichen auf Console ausgeben
	OUT	(1),A
	RET

; Macro to zero fill memory
ZFILL(ptr,len) MACRO
	LD  HL,ptr
	LD  DE,ptr+1
	LD  BC,len-1
	LD  (HL),0
	LDIR
ENDM

; Memory copy
MCOPY(src,dest,len) MACRO
	LD HL,src
	LD DE,dest
	LD BC, len
	LDIR
ENDM

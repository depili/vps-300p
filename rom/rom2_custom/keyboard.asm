	; IX = new data
	; IY = old data
	; C = Counter for the checked byte, later encoded into the keyboard event byte
	; Messes up AF
KEYB_PROCESS: PROC
	LD	C, 0x00		; Byte being checked
loop:	LD	B, (IX)
	LD	A, (IY)
	XOR	B
	CALL	NZ, KEYB_BYTE_DIFF
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
KEYB_BYTE_DIFF: PROC
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

	; Send a keyboard event via TX buffer
	; C - changed byte
	; E - changed bit and key up/down
KEYB_SEND:
	PUSH	AF
	PUSH	BC
	LD	A, (KEYB_CMD_BYTE)
	CALL	SIO_A_TX_BLOCKING
	LD	A, C
	OR	E
	CALL	SIO_A_TX_BLOCKING
	POP	BC
	POP	AF
	RET

	; --- START PROC L035F ---
INIT_KEYB:
	LD	HL,KEYB_1_INIT_DATA
	LD	C,KEYB_1_CMD
	LD	B,05h
	OTIR
	CALL	KEYB_DELAY
	LD	HL,KEYB_2_INIT_DATA
	LD	C,KEYB_2_CMD
	LD	B,05h
	OTIR
	CALL	KEYB_DELAY
	LD	HL,KEYB_3_INIT_DATA
	LD	C,KEYB_3_CMD
	LD	B,05h
	OTIR
	CALL	KEYB_DELAY
	LD	A,KEYB_CMD_WRITE_DISPLAY	; Fill keyboard 2 with 0x00
	OUT	(KEYB_2_CMD),A
	LD	A,00h
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	OUT	(KEYB_2_DATA),A
	ZFILL	0xB100, 0x300			; Zero keyboard memories
	RET

	; --- START PROC L0CCF ---
KEYB_READ_1:
	IN	A,(PIO_B_DATA)
	BIT	7,A
	RET	Z
	LD	A,KEYB_CMD_READ_DATA
	OUT	(KEYB_1_CMD),A
	LD	HL,KEYB_1_DEST
	LD	B,08h
	LD	C,KEYB_1_DATA
	INIR
	LD	A,KEYB_CMD_END
	OUT	(KEYB_1_CMD),A
	; TX_A	"1"
	; TX_A_KB	KEYB_1_DEST
	LD	A, TX_CMD_KB_1
	LD	(KEYB_CMD_BYTE), A
	LD	IX, KEYB_1_DEST
	LD	IY, KEYB_1_OLD
	CALL	KEYB_PROCESS
	MCOPY	KEYB_1_DEST, KEYB_1_OLD, 0x08
	RET

	; --- START PROC L0CEF ---
KEYB_READ_2:
	IN	A,(PIO_B_DATA)
	BIT	6,A
	RET	Z
	LD	A,KEYB_CMD_READ_DATA
	OUT	(KEYB_2_CMD),A
	LD	HL,KEYB_2_DEST
	LD	B,08h
	LD	C,KEYB_2_DATA
	INIR
	LD	A,KEYB_CMD_END
	OUT	(KEYB_2_CMD),A
	; TX_A	"2"
	; TX_A_KB	KEYB_2_DEST
	LD	A, TX_CMD_KB_2
	LD	(KEYB_CMD_BYTE), A
	LD	IX, KEYB_2_DEST
	LD	IY, KEYB_2_OLD
	CALL	KEYB_PROCESS
	MCOPY	KEYB_2_DEST, KEYB_2_OLD, 0x08
	RET

	; --- START PROC L0D0F ---
KEYB_READ_3:
	IN	A,(PIO_B_DATA)
	BIT	5,A
	RET	Z
	LD	A,KEYB_CMD_READ_DATA
	OUT	(KEYB_3_CMD),A
	LD	HL,KEYB_3_DEST
	LD	B,08h
	LD	C,KEYB_3_DATA
	INIR
	LD	A,KEYB_CMD_END
	OUT	(KEYB_3_CMD),A
	; TX_A	"3"
	; TX_A_KB	KEYB_3_DEST
	LD	A, TX_CMD_KB_3
	LD	(KEYB_CMD_BYTE), A
	LD	IX, KEYB_3_DEST
	LD	IY, KEYB_3_OLD
	CALL	KEYB_PROCESS
	MCOPY	KEYB_3_DEST, KEYB_3_OLD, 0x08
	RET

	; --- START PROC L7E26 ---
KEYB_WRITE_1:
	; CALL	L358D		; 7-segment encoding?
	LD	A,KEYB_CMD_WRITE_DISPLAY
	OUT	(KEYB_1_CMD),A
	LD	HL,KEYB_1_DISP
	LD	C,KEYB_1_DATA
	LD	B,08h
	OTIR
	RET

	; --- START PROC L7E37 ---
KEYB_WRITE_2:
	LD	A,KEYB_CMD_WRITE_DISPLAY
	OUT	(KEYB_2_CMD),A
	LD	HL,KEYB_2_DISP
	LD	C,KEYB_2_DATA
	LD	B,08h
	OTIR
	; Would copy lamp bitmap to shared memory
	; LD	HL,960Ah
	; LD	DE,0FA00h
	; LD	BC,001Eh
	; LDIR
	RET

	; Assumed "PWM" lamp brightness settings
	; Set bit 5 to 0000 0000
	; --- START PROC L3E07 ---
KEYB_2_PWM_5_0:
	LD	IY,KEYB_2_DISP
	RES	5,(IY+00h)	; 0
	RES	5,(IY+01h)	; 0
	RES	5,(IY+02h)	; 0
	RES	5,(IY+03h)	; 0
	RES	5,(IY+04h)	; 0
	RES	5,(IY+05h)	; 0
	RES	5,(IY+06h)	; 0
	RES	5,(IY+07h)	; 0
	RET

	; Set bit 5 to 1100 0001
	; --- START PROC L3E2C ---
KEYB_2_PWM_5_1:
	LD	IY,KEYB_2_DISP
	SET	5,(IY+00h)	; 1
	SET	5,(IY+01h)	; 1
	RES	5,(IY+02h)	; 0
	RES	5,(IY+03h)	; 0
	RES	5,(IY+04h)	; 0
	RES	5,(IY+05h)	; 0
	RES	5,(IY+06h)	; 0
	SET	5,(IY+07h)	; 1
	RET

	; Set bit 4 to 0000 0000
	; --- START PROC L3E51 ---
KEYB_2_PWM_4_0:
	LD	IY,KEYB_2_DISP
	RES	4,(IY+00h)	; 0
	RES	4,(IY+01h)	; 0
	RES	4,(IY+02h)	; 0
	RES	4,(IY+03h)	; 0
	RES	4,(IY+04h)	; 0
	RES	4,(IY+05h)	; 0
	RES	4,(IY+06h)	; 0
	RES	4,(IY+07h)	; 0
	RET

	; Set bit 4 to 1100 0001
	; --- START PROC L3E76 ---
KEYB_2_PWM_4_1:
	LD	IY,KEYB_2_DISP
	SET	4,(IY+00h)	; 1
	SET	4,(IY+01h)	; 1
	RES	4,(IY+02h)	; 0
	RES	4,(IY+03h)	; 0
	RES	4,(IY+04h)	; 0
	RES	4,(IY+05h)	; 0
	RES	4,(IY+06h)	; 0
	SET	4,(IY+07h)	; 1
	RET

	; Set bit 3 to 0000 0000
	; --- START PROC L3E9B ---
KEYB_2_PWM_3_0:
	LD	IX,KEYB_2_DISP
	RES	3,(IX+00h)	; 0
	RES	3,(IX+01h)	; 0
	RES	3,(IX+02h)	; 0
	RES	3,(IX+03h)	; 0
	RES	3,(IX+04h)	; 0
	RES	3,(IX+05h)	; 0
	RES	3,(IX+06h)	; 0
	RES	3,(IX+07h)	; 0
	RET

	; Set bit 3 to 1100 0001
	; --- START PROC L3EC0 ---
KEYB_2_PWM_3_1:
	LD	IX,KEYB_2_DISP
	SET	3,(IX+00h)	; 1
	SET	3,(IX+01h)	; 1
	RES	3,(IX+02h)	; 0
	RES	3,(IX+03h)	; 0
	RES	3,(IX+04h)	; 0
	RES	3,(IX+05h)	; 0
	RES	3,(IX+06h)	; 0
	SET	3,(IX+07h)	; 1
	RET

	; Set bit 2 to 0000 0000
	; --- START PROC L3EE5 ---
KEYB_2_PWM_2_0:
	LD	IX,KEYB_2_DISP
	RES	2,(IX+00h)	; 0
	RES	2,(IX+01h)	; 0
	RES	2,(IX+02h)	; 0
	RES	2,(IX+03h)	; 0
	RES	2,(IX+04h)	; 0
	RES	2,(IX+05h)	; 0
	RES	2,(IX+06h)	; 0
	RES	2,(IX+07h)	; 0
	RET

	; Set bit 2 to 1100 0001
	; --- START PROC L3F0A ---
KEYB_2_PWM_2_1:
	LD	IX,KEYB_2_DISP
	SET	2,(IX+00h)	; 1
	SET	2,(IX+01h)	; 1
	RES	2,(IX+02h)	; 0
	RES	2,(IX+03h)	; 0
	RES	2,(IX+04h)	; 0
	RES	2,(IX+05h)	; 0
	RES	2,(IX+06h)	; 0
	SET	2,(IX+07h)	; 1
	RET

	; Set bit 1 to 0000 0000
	; --- START PROC L3F2F ---
KEYB_2_PWM_1_0:
	LD	IX,KEYB_2_DISP
	RES	1,(IX+00h)
	RES	1,(IX+01h)
	RES	1,(IX+02h)
	RES	1,(IX+03h)
	RES	1,(IX+04h)
	RES	1,(IX+05h)
	RES	1,(IX+06h)
	RES	1,(IX+07h)
	RET

	; Set bit 1 to 1100 0001
	; --- START PROC L3F54 ---
KEYB_2_PWM_1_1:
	LD	IX,KEYB_2_DISP
	SET	1,(IX+00h)
	SET	1,(IX+01h)
	RES	1,(IX+02h)
	RES	1,(IX+03h)
	RES	1,(IX+04h)
	RES	1,(IX+05h)
	RES	1,(IX+06h)
	SET	1,(IX+07h)
	RET

	; --- START PROC L039A ---
KEYB_DELAY: PROC
	LD	HL,1000h
loop:	DEC	HL
	LD	A,H
	OR	L
	JR	NZ,loop
	RET
ENDP

	; gfedcba encoded 7 segment numbers and letters
KEYB_NUMBERS:
	DB	0x40	; 0
	DB	0x79	; 1
	DB	0x24	; 2
	DB	0x30	; 3
	DB	0x19	; 4
	DB	0x12	; 5
	DB	0x02	; 6
	DB	0x78	; 7
	DB	0x80	; 8
	DB	0x10	; 9
	DB	0x08	; A
	DB	0x03	; B
	DB	0x46	; C
	DB	0x21	; d
	DB	0x06	; E
	DB	0x0E	; F

KEYB_1_INIT_DATA:
	DB	04h	; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h	; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h	; 1110 0000 End mode set
	DB	0DFh	; 1101 1111 Fill display ram with 0xFF
KEYB_2_INIT_DATA:
	DB	04h	; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h	; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h	; 1110 0000 End mode set
	DB	0DFh	; 1101 1111 Fill display ram with 0xFF
KEYB_3_INIT_DATA:
	DB	04h	; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h	; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h	; 1110 0000 End mode set
	DB	0DFh	; 1101 1111 Fill display ram with 0xFF


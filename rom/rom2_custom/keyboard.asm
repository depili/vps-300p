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
	LD	A,KEYB_CMD_WRITE_DISPLAY        ; Fill keyboard 2 with 0x00
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
	; LD	C,01h
	; LD	IX,94F9h	; address or value?
	; CALL	L0D2E           ; Process the keypress
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
	; LD	C,02h
	; LD	IX,9501h	; address or value?
	; CALL	L0D2E
	RET

	; --- START PROC L0D0F ---
KEYB_READ_3:
	IN	A,(PIO_B_DATA)
	BIT	5,A
	LD	A,KEYB_CMD_READ_DATA
	OUT	(KEYB_3_CMD),A
	LD	HL,KEYB_3_DEST
	LD	B,08h
	LD	C,04h
	INIR
	LD	A,KEYB_CMD_END
	OUT	(KEYB_3_CMD),A
	; LD	C,03h
	; LD	IX,9509h	; address or value?
	; CALL	L0D2E
	RET

	; --- START PROC L7E26 ---
KEYB_WRITE_1:
	; CALL	L358D           ; 7-segment encoding?
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
	RES	5,(IY+00h)      ; 0
	RES	5,(IY+01h)      ; 0
	RES	5,(IY+02h)      ; 0
	RES	5,(IY+03h)      ; 0
	RES	5,(IY+04h)      ; 0
	RES	5,(IY+05h)      ; 0
	RES	5,(IY+06h)      ; 0
	RES	5,(IY+07h)      ; 0
	RET

	; Set bit 5 to 1100 0001
	; --- START PROC L3E2C ---
KEYB_2_PWM_5_1:
	LD	IY,KEYB_2_DISP
	SET	5,(IY+00h)      ; 1
	SET	5,(IY+01h)      ; 1
	RES	5,(IY+02h)      ; 0
	RES	5,(IY+03h)      ; 0
	RES	5,(IY+04h)      ; 0
	RES	5,(IY+05h)      ; 0
	RES	5,(IY+06h)      ; 0
	SET	5,(IY+07h)      ; 1
	RET

	; Set bit 4 to 0000 0000
	; --- START PROC L3E51 ---
KEYB_2_PWM_4_0:
	LD	IY,KEYB_2_DISP
	RES	4,(IY+00h)      ; 0
	RES	4,(IY+01h)      ; 0
	RES	4,(IY+02h)      ; 0
	RES	4,(IY+03h)      ; 0
	RES	4,(IY+04h)      ; 0
	RES	4,(IY+05h)      ; 0
	RES	4,(IY+06h)      ; 0
	RES	4,(IY+07h)      ; 0
	RET

	; Set bit 4 to 1100 0001
	; --- START PROC L3E76 ---
KEYB_2_PWM_4_1:
	LD	IY,KEYB_2_DISP
	SET	4,(IY+00h)      ; 1
	SET	4,(IY+01h)      ; 1
	RES	4,(IY+02h)      ; 0
	RES	4,(IY+03h)      ; 0
	RES	4,(IY+04h)      ; 0
	RES	4,(IY+05h)      ; 0
	RES	4,(IY+06h)      ; 0
	SET	4,(IY+07h)      ; 1
	RET

	; Set bit 3 to 0000 0000
	; --- START PROC L3E9B ---
KEYB_2_PWM_3_0:
	LD	IX,KEYB_2_DISP
	RES	3,(IX+00h)      ; 0
	RES	3,(IX+01h)      ; 0
	RES	3,(IX+02h)      ; 0
	RES	3,(IX+03h)      ; 0
	RES	3,(IX+04h)      ; 0
	RES	3,(IX+05h)      ; 0
	RES	3,(IX+06h)      ; 0
	RES	3,(IX+07h)      ; 0
	RET

	; Set bit 3 to 1100 0001
	; --- START PROC L3EC0 ---
KEYB_2_PWM_3_1:
	LD	IX,KEYB_2_DISP
	SET	3,(IX+00h)      ; 1
	SET	3,(IX+01h)      ; 1
	RES	3,(IX+02h)      ; 0
	RES	3,(IX+03h)      ; 0
	RES	3,(IX+04h)      ; 0
	RES	3,(IX+05h)      ; 0
	RES	3,(IX+06h)      ; 0
	SET	3,(IX+07h)      ; 1
	RET

	; Set bit 2 to 0000 0000
	; --- START PROC L3EE5 ---
KEYB_2_PWM_2_0:
	LD	IX,KEYB_2_DISP
	RES	2,(IX+00h)      ; 0
	RES	2,(IX+01h)      ; 0
	RES	2,(IX+02h)      ; 0
	RES	2,(IX+03h)      ; 0
	RES	2,(IX+04h)      ; 0
	RES	2,(IX+05h)      ; 0
	RES	2,(IX+06h)      ; 0
	RES	2,(IX+07h)      ; 0
	RET

	; Set bit 2 to 1100 0001
	; --- START PROC L3F0A ---
KEYB_2_PWM_2_1:
	LD	IX,KEYB_2_DISP
	SET	2,(IX+00h)      ; 1
	SET	2,(IX+01h)      ; 1
	RES	2,(IX+02h)      ; 0
	RES	2,(IX+03h)      ; 0
	RES	2,(IX+04h)      ; 0
	RES	2,(IX+05h)      ; 0
	RES	2,(IX+06h)      ; 0
	SET	2,(IX+07h)      ; 1
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

KEYB_1_INIT_DATA:
	DB	04h     ; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h    ; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h    ; 1110 0000 End mode set
	DB	0DFh    ; 1101 1111 Fill display ram with 0xFF
KEYB_2_INIT_DATA:
	DB	04h     ; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h    ; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h    ; 1110 0000 End mode set
	DB	0DFh    ; 1101 1111 Fill display ram with 0xFF
KEYB_3_INIT_DATA:
	DB	04h     ; 0000 0100 Mode set: 8x8bit character display, left entry, encoded scan sensor matrix
	DB	0A0h    ; 1010 0000 Display write inhibit: none
	DB	22h	; 0010 0010 Program clock: prescaler 00010
	DB	0E0h    ; 1110 0000 End mode set
	DB	0DFh    ; 1101 1111 Fill display ram with 0xFF

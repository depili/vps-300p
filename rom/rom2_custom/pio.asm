	; --- START PROC L01A2 ---
INIT_PIO:
	LD	HL,PIO_A_INIT_DATA
	LD	C,1Dh
	LD	B,03h
	OTIR
	LD	HL,PIO_B_INIT_DATA
	LD	C,1Fh
	LD	B,03h
	OTIR
	RET

	; --- START PROC L03B2 ---
PIO_A_SET_BIT_0:
	IN	A,(PIO_A_DATA)
	OR	01h
	OUT	(PIO_A_DATA),A
	RET

PIO_A_INIT_DATA:
	DB	0CFh	; Mode 3
	DB	80h	; Pin 7 input, others outputs
	DB	07h	; No interrupts
PIO_B_INIT_DATA:
	DB	0CFh	; Mode 3
	DB	0FFh	; All pins inputs
	DB	07h	; No interrupts

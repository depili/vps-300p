	; --- START PROC L03B9 ---
	; Does one read from the ADC, discards read data
INIT_ADC: PROC
	LD	A,ADC_CMD_TRIGGER
	OUT	(ADC_IO),A
loop:	IN	A,(PIO_A_DATA)
	BIT	7,A
	JR	NZ,loop
	IN	A,(ADC_IO)
	IN	A,(ADC_IO)
	RET
ENDP

	; --- START PROC L43FC ---
	; The order is wonky, does it matter?
ADC_READ_ALL:
	ADC_R	0x00, ADC_0_DEST
	ADC_R	0x04, ADC_4_DEST
	ADC_R	0x02, ADC_2_DEST
	ADC_R	0x06, ADC_6_DEST
	ADC_R	0x01, ADC_1_DEST
	ADC_R	0x05, ADC_5_DEST
	ADC_R	0x03, ADC_3_DEST
	; CALL	ADC_MUX_RESET
	CALL	ADC_PROCESS
	RET

	; Check for changes, send messages if needed
	; For now, only check and transmit the top byte
ADC_PROCESS:
	ADC_CHECK	0
	ADC_CHECK	1
	ADC_CHECK	2
	ADC_CHECK	3
	ADC_CHECK	4
	ADC_CHECK	5
	ADC_CHECK	6
	RET

ADC_COMPARE: PROC
;	LD	A, (IX+1)	; Old value, most significant 8 bits
;	LD	H, (IY+1)	; New value, most significant 8 bits
;	LD	E, (IY)		; New value, least significat 2 bits
;	LD	(IX+1), H	; Store the new value as "old"
;	LD	(IX), E
;	SRL	A		; Divide both old and new bytes by 2
;	SRL	H
;	CP	H		; Compare old and new top bytes, thus ignoring lowest bit
;	JR	NZ, send	; Send if changed enough
;	RET
send:	LD	H, (IY+1)
	LD	L, (IY)
	SRL	H
	RR	L
	SRL	L
	CALL	ADC_SEND
	RET
ENDP

	; Send the ADC value via the TX buffer
	; B = value to send (high byte)
	; C = ADC channel
ADC_SEND:
	LD	A, TX_CMD_ADC_0
	OR	C
	CALL	SIO_A_TX_BLOCKING
	LD	A, H
	CALL	SIO_A_TX_BLOCKING
	LD	A, L
	CALL	SIO_A_TX_BLOCKING
	RET

	; --- START PROC L444D ---
	; Not entirely sure
ADC_MUX_RESET: PROC
	; LD	A,60h
	; LD	B,0DFh
	; LD	C,07h
	; AND	B
	; OR	C
	LD	A, 0x47		; Folded the above constant, should select mux channel 7
	OUT	(ADC_MUX),A
	LD	B,0Ah
wait:	DJNZ	wait
	LD	A,60h
	OUT	(ADC_MUX),A
	RET
ENDP

	; --- START PROC L4460 ---
	; B  = 0xDF
	; C  = channel
	; IX = store location
ADC_READ: PROC
	LD	DE,0000h
	LD	A,10h
	LD	(ADC_READ_COUNTER),A
loop:	LD	A,60h
	AND	B
	OR	C
	OUT	(ADC_MUX),A		; Select mux channel
	LD	A,ADC_CMD_TRIGGER
	OUT	(ADC_IO),A		; Trigger ADC read
wait:	IN	A,(PIO_A_DATA)
	BIT	7,A
	JP	NZ,wait
	IN	A,(ADC_IO)
	LD	(IX+01h),A
	; LD	H,A
	IN	A,(ADC_IO)
	LD	(IX+00h),A
	; LD	L,A
	; SRL	H			; HL >> 6
	; RR	L
	; SRL	H
	; RR	L
	; SRL	H
	; RR	L
	; SRL	H
	; RR	L
	; SRL	H
	; RR	L
	; SRL	H
	; RR	L
	; ADD	HL,DE
	; EX	DE,HL
	; LD	A,(ADC_READ_COUNTER)
	; DEC	A
	; LD	(ADC_READ_COUNTER),A
	; JP	NZ,loop			; executed 16 times
	; EX	DE,HL
	; ADD	HL,HL			; HL * 4
	; ADD	HL,HL
	; LD	A,H
	; LD	(IX+01h),A
	; LD	A,L
	; AND	0C0h			; Mask leaves top 2 bytes
	; LD	(IX+00h),A
	RET
ENDP


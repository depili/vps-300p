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

ADC_R(ch,dest) MACRO
        LD      B, 0xDF
        LD      C, ch
        LD      IX, dest
        CALL    ADC_READ
ENDM


	; --- START PROC L43FC ---
        ; The order is wonky, does it matter?
ADC_READ_ALL:
        ADC_R   0x00, ADC_0_DEST
        ADC_R   0x04, ADC_4_DEST
        ADC_R   0x02, ADC_2_DEST
        ADC_R   0x06, ADC_6_DEST
        ADC_R   0x01, ADC_1_DEST
        ADC_R   0x05, ADC_5_DEST
        ADC_R   0x03, ADC_3_DEST
	CALL	ADC_MUX_RESET
	RET

	; --- START PROC L444D ---
        ; Not entirely sure
ADC_MUX_RESET: PROC
	; LD	A,60h
	; LD	B,0DFh
	; LD	C,07h
	; AND	B
	; OR	C
        LD      A, 0x47         ; Folded the above constant, should select mux channel 7
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
	OUT	(ADC_MUX),A             ; Select mux channel
	LD	A,ADC_CMD_TRIGGER
	OUT	(ADC_IO),A              ; Trigger ADC read
wait:	IN	A,(PIO_A_DATA)
	BIT	7,A
	JP	NZ,wait
	IN	A,(ADC_IO)
	LD	(IX+01h),A
	LD	H,A
	IN	A,(ADC_IO)
	LD	(IX+00h),A
	LD	L,A
	SRL	H                       ; HL >> 6
	RR	L
	SRL	H
	RR	L
	SRL	H
	RR	L
	SRL	H
	RR	L
	SRL	H
	RR	L
	SRL	H
	RR	L
	ADD	HL,DE
	EX	DE,HL
	LD	A,(ADC_READ_COUNTER)
	DEC	A
	LD	(ADC_READ_COUNTER),A
	JP	NZ,loop                 ; executed 16 times
	EX	DE,HL
	ADD	HL,HL                   ; HL * 4
	ADD	HL,HL
	LD	A,H
	LD	(IX+01h),A
	LD	A,L
	AND	0C0h                    ; Mask leaves top 2 bytes
	LD	(IX+00h),A
	RET
ENDP


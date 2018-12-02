ORG 0x0

START:
        DI
        IM	2
        LD	A,00h
        LD	I,A
        LD	SP,0F7FFh	; Stack
        CALL	INIT_F4_F0_F1
        ; CALL	L00E6           ; Memory check and init
        JR	INIT
; Second init here at 0x0012

        ; L001C
INIT:
        CALL	INIT_EI_RETI
        CALL	INIT_IO
        ; CALL	L0245           ; DI, re-init CTC 2, EI, wtf?
        ; EI
        JP	MAIN_LOOP


        SEEK 0x70
        ORG 0x70
        ; Interrupts
INTERRUPT_VECTORS:
        DW      INT_SIO_B_TX_EMPTY
        DW      INT_SIO_B_STATUS_CHANGE
        DW      INT_SIO_B_RX_AVAILABLE
        DW      INT_SIO_B_ERROR
        DW      INT_SIO_A_TX_EMPTY
        DW      INT_SIO_A_STATUS_CHANGE
        DW      INT_SIO_A_RX_AVAILABLE
        DW      INT_SIO_A_ERROR
CTC_VECTORS:
        DW      INT_CTC_CH1
        DW      INT_CTC_CH2
        DW      INT_CTC_CH3
        DW      INT_CTC_CH4
PIO_VECTORS:
        DW      INT_PIO_A
        DW      INT_PIO_B

        SEEK 0x90
        ORG 0x90


	; L008C
MAIN_LOOP: PROC
	; CALL	L0743           ; Send ping request if PIO B bit 3 is not set
	; CALL	L03C8           ; Memory setup
	; CALL	L06BA           ; Send initial data if ping-pong
	; CALL	L0812           ; Send 0x86 0x63 0x01 0x00
	; CALL	L081F           ; PIO B bit 3 check
loop:	; LD	A,(9461h)
	; AND	A
	; JR	Z,loop_end
	; LD	A,00h
	; LD	(9461h),A
	CALL	KEYB_READ_1
	CALL	KEYB_READ_2
	CALL	KEYB_READ_3
	; CALL	L3CDD           ; Weird compare that doesn't lead to anywhere?
	CALL	ADC_READ_ALL    ; ADC reads
	; CALL	L44B7           ; Process ADC readings
	; CALL	L476A           ; Process ADC readings
	; CALL	L47BE           ; Reads flag from shared memory, goes to the Brain(TM)
	; CALL	L3D13           ; Brain stuff?
	; CALL	KEYB_WRITE_1    ; Does also 7-segment encoding
	; CALL	KEYB_WRITE_2    ; Would also update the lamp bitmap
	; CALL	L7205           ; Brain stuff
	; CALL	L7CAA           ; PIO B bit 3 check and stuff
loop_end:
	; CALL	L7E50           ; Sets PIO bit 1 conditionally, beeper?
	; CALL	L7820           ; Big compare, serial command process?
	JR	loop
ENDP

	; --- START PROC L0186 ---
INIT_IO:
L0186:	CALL	INIT_CTC
	CALL	INIT_SIO
	CALL	INIT_PIO
	CALL	INIT_KEYB
	CALL	PIO_A_SET_BIT_0
	CALL	INIT_ADC
	; CALL	L042D           ; Zero shared memory
	; CALL	L0793           ; Checks PIO B bit 3, does stuff if it is high, dip switch?
	; CALL	L07E2           ; Checks PIO B bit 3, does stuff if it is high, dip switch?
	RET

        ; L0176
INIT_EI_RETI:
        LD	A,0FBh
        LD	(8000h),A
        LD	A,0EDh
        LD	(8001h),A
        LD	A,4Dh
        LD	(8002h),A
        RET


INIT_F4_F0_F1:
        LD	A,00h
        OUT	(0F4h),A
        LD	A,5Bh
        OUT	(0F0h),A
        LD	A,0B1h
        OUT	(0F1h),A
        RET

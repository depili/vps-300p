        ; Waits for TX buffer to be empty, then sends A
        ; Verfied working
SIO_A_TX_BLOCKING: PROC
        LD      D, A
        CALL    SIO_A_TX_CHECK
        LD      A, D
        OUT     (SIO_A_DATA), A ; Send the character
        RET
ENDP

        ; Wait until SIO A TX buffer is empty
        ; Verified working
SIO_A_TX_CHECK: PROC
        LD      A, 0x01         ; Select register 1
        OUT     (SIO_A_CMD), A
        IN      A, (SIO_A_CMD)  ; Read register 1
        BIT     0, A            ; Bit 0 is "All sent"
        JP      Z, SIO_A_TX_CHECK
        RET
ENDP

        ; --- START PROC L01C6 ---
INIT_SIO:
	LD	HL,SIO_A_INIT_DATA
	LD	C,SIO_A_CMD
	LD	B,0Dh
	OTIR
	LD	HL,SIO_B_INIT_DATA
	LD	C,1Bh
	LD	B,0Fh
	OTIR
	; CALL	L0B62           ; Memory setup
	; CALL	L0B82           ; Memory setup
	; CALL	L0B90           ; Memory setup
	; CALL	L0B72           ; Memory setup
	; LD	A,0FFh
	; LD	(9361h),A
	; LD	HL,0FFFFh	; address or value?
	; LD	(895Ch),HL
	; LD	HL,0000h	; address or value?
	; LD	(9465h),HL
	; LD	A,01h
	; LD	(9467h),A
	RET

SIO_A_INIT_DATA:
	DB	00h     ; Write register 0
	DB	00h
	DB	18h     ; CMD error reset
	DB	10h     ; CMD reset status interrupts
	DB	04h     ; Write register 4
	DB	45h     ; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
	DB	03h     ; Write register 3
	DB	0C1h    ; RX enable, 8 bits per character
	DB	05h     ; Write register 5
	DB	0EAh    ; TX crc disable, RTS low, TX enable, DTR low, external sync
	DB	10h     ; CMD reset status interrupts
	DB	01h     ; Write register 1
	DB	10h     ; receive interrupt on all characters (was 12h)
SIO_B_INIT_DATA:
	DB	00h     ; Write register 0
	DB	00h
	DB	18h     ; CMD error reset
	DB	10h     ; CMD reset status interrupts
	DB	04h     ; Write register 4
	DB	0C5h    ; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
	DB	03h     ; Write register 3
	DB	0C1h    ; RX enable, 8 bits per character
	DB	05h     ; Write register 5
	DB	0EAh    ; TX crc disable, RTS low, TX enable, DTR low, external sync
	DB	02h     ; Write register 2
	DB	INTERRUPT_VECTORS % 0x0100
	DB	10h     ; CMD reset status interrupts
	DB	01h     ; Write register 1
	DB	14h     ; Receive interrupt on all characters, status affects vector (Was 16h)
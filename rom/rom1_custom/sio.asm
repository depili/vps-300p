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

RX_INIT:
        LD      A, 0x00         ; Error in message, zero the counter and type
        LD      (RX_COUNTER), A
        LD      (RX_TYPE), A
        LD      HL, RX_POINTER+2
        LD      (RX_POINTER), HL
        RET

        ; Check the RX buffer
CHECK_RX: PROC
        LD      A, (RX_COUNTER)
        CP      RX_MAX_BYTES + 0x01
        JR      NC, err  ; Too many bytes
        AND     A
        JR      NZ, check_msg
        RET
err:    CALL    RX_INIT
        LD      A, "E"
        CALL    SIO_A_TX_BLOCKING
        RET
check_msg:
        LD      A, (RX_TYPE)
        CP      0x80            ; Lamp message
        JR      Z, lamp
        JR      err             ; Invalid message
lamp:                           ; Process a lamp message
        LD      A, "L"
        CALL    SIO_A_TX_BLOCKING
        LD      A, (RX_COUNTER)
        CP      0x03
        JR      NZ, return      ; not enough bytes for full message
        LD      IX, RX_TYPE     ; first byte of message
        LD      A, (IX+0x01)    ; Byte offset of the lamp data
        OUT     (SIO_A_DATA), A
        CP      LAMP_BYTES
        JP      NC, err         ; Too big offset
        LD      BC, 0x0000
        LD      C, A
        LD      HL, LAMP_DEST
        ADD     HL, BC
        LD      A, (IX+0x02)
        LD      (HL), A
        CALL    LAMP_UPDATE     ; Update the lamps
        LD      A, "O"
        CALL    SIO_A_TX_BLOCKING
        CALL    RX_INIT
return: RET
ENDP

        ; --- L015A ---
INIT_SIO:
        LD      HL,SIO_A_INIT_DATA
        LD      C,19h
        LD      B,0Dh
        OTIR
        LD      HL,SIO_B_INIT_DATA
        LD      C,1Bh
        LD      B,0Fh
        OTIR
        ; LD      HL,0000h
        ; LD      (9411h),HL
        ; LD      A,01h
        ; LD      (9413h),A
        RET

SIO_A_INIT_DATA:
        DB      00h     ; Write register 0
        DB      00h
        DB      18h     ; CMD error reset
        DB      10h     ; CMD reset status interrupts
        DB      04h     ; Write register 4
        DB      45h     ; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
        DB      03h     ; Write register 3
        DB      0C1h    ; RX enable, 8 bits per character
        DB      05h     ; Write register 5
        DB      0EAh    ; TX crc disable, RTS low, TX enable, DTR low, external sync
        DB      10h     ; CMD reset status interrupts
        DB      01h     ; Write register 1
        DB      10h     ; Transmit interrupt enable, receive interrupt on all characters
SIO_B_INIT_DATA:
        DB      00h     ; Write register 0
        DB      00h
        DB      18h     ; CMD error reset
        DB      10h     ; CMD reset status interrupts
        DB      04h     ; Write register 4
        DB      45h     ; 1,5 stop bits, even parity, 8bit programmed sync, data rate=16x clock rate
        DB      03h     ; Write register 3
        DB      0C1h    ; RX enable, 8 bits per character
        DB      05h     ; Write register 5
        DB      0EAh    ; TX crc disable, RTS low, TX enable, DTR low, external sync
        DB      02h     ; Write register 2
        DB      INTERRUPT_VECTORS % 0x0100
        DB      10h     ; CMD reset status interrupts
        DB      01h     ; Write register 1
        DB      14h     ; Transmit interrupt enable, receive interrupt on all characters, status affects vector

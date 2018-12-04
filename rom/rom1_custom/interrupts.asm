        ; -- INTERRUPT HANDLERS --
INT_SIO_B_TX_EMPTY:
        DI
        PUSH    AF
        LD      A, "z"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP      0x8000

INT_SIO_B_STATUS_CHANGE:
        DI
        PUSH    AF
        LD      A, "x"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP      0x8000

INT_SIO_B_RX_AVAILABLE:
        DI
        PUSH    AF
        IN      A, (SIO_B_DATA)
        LD      A, "c"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP 0x8000

INT_SIO_B_ERROR:
        DI
        PUSH    AF
        LD      A, SIO_CMD_ERROR_RST
        OUT     (SIO_B_CMD), A
        IN      A, (SIO_B_DATA)
        LD      A, "v"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP      0x8000

INT_SIO_A_TX_EMPTY:
        DI
        PUSH    AF
        LD      A, "T"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP      0x8000

INT_SIO_A_STATUS_CHANGE:
        DI
        PUSH    AF
        LD      A, "S"
        OUT     (SIO_A_DATA), A
        POP     AF
        JP      0x8000

INT_SIO_A_RX_AVAILABLE: PROC
        DI
        PUSH    AF
        PUSH    BC
        PUSH    DE
        PUSH    HL
        IN      A, (SIO_A_DATA)         ; Read the character
        OUT     (SIO_A_DATA), A
        LD      D, A                    ; Tuck the read character into D
        LD      A, (RX_COUNTER)
        AND     A
        JR      Z, first_byte           ; Seek for the start byte
        CP      RX_MAX_BYTES
        JR      Z, return               ; Buffer full
        LD      A, D                    ; Restore read character
store:  LD      HL, (RX_POINTER)
        LD      (HL), A
        INC     HL
        LD      (RX_POINTER), HL
        LD      HL, RX_COUNTER
        INC     (HL)
        LD      A, (HL)
        OUT     (SIO_A_DATA), A         ; Send the RX count
        JR      return
first_byte:
        LD      A, D                    ; Restore read character
        BIT     7, A
        JR      NZ, store               ; Valid first byte
return: POP     HL
        POP     DE
        POP     BC
        POP     AF
        JP      0x8000
ENDP

INT_SIO_A_ERROR:
        DI
        PUSH    AF
        LD      A, SIO_CMD_ERROR_RST
        OUT     (SIO_A_CMD), A
        IN      A, (SIO_A_DATA)
        LD      A, "!"
        OUT     (SIO_A_DATA), A

        POP     AF
        JP      0x8000

INT_CTC_CH1:
        DI
        JP      0x8000

INT_CTC_CH2:
        DI
        JP      0x8000

INT_CTC_CH3:
        DI
        JP      0x8000

INT_CTC_CH4:
        DI
        JP      0x8000

INT_PIO_A:
        DI
        JP      0x8000

INT_PIO_B:
        DI
        JP      0x8000

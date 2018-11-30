        ; -- INTERRUPT HANDLERS --
INT_SIO_B_TX_EMPTY:
        DI
        JP      0x8000

INT_SIO_B_STATUS_CHANGE:
        DI
        JP      0x8000

INT_SIO_B_RX_AVAILABLE:
        DI
        PUSH    AF
        IN      A, (SIO_B_DATA)
        POP     AF
        JP 0x8000

INT_SIO_B_ERROR:
        DI
        PUSH    AF
        LD      A, SIO_CMD_ERROR_RST
        OUT     (SIO_B_CMD), A
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
        PUSH    AF
        LD      A, (RX_COUNTER)
        AND     A
        JR      Z, first_byte           ; Seek for the start byte
        POP     AF
store:  LD      HL, (RX_POINTER)
        LD      (HL), A
        INC     HL
        LD      (RX_POINTER), HL
        LD      HL, RX_COUNTER
        INC     (HL)
        LD      A, (HL)
        OUT     (SIO_A_DATA), A         ; Send the RX count
return: POP     HL
        POP     DE
        POP     BC
        POP     AF
        JP      0x8000
first_byte:
        POP     AF
        BIT     7, A
        JR      NZ, store               ; Valid first byte
        JR      return
ENDP

INT_SIO_A_ERROR:
        DI
        PUSH    AF
        LD      A, "E"
        OUT     (SIO_A_DATA), A
        LD      A, SIO_CMD_ERROR_RST
        OUT     (SIO_A_CMD), A
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

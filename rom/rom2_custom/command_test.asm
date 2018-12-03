ORG 0x0000
RX_COUNTER              EQU     0x9000
RX_POINTER              EQU     0x9002
RX_TYPE                 EQU     0x9004  ; first byte of the message
RX_MAX_BYTES            EQU     0x9100 - RX_TYPE - 1
; LCD
LCD_FLAG                EQU     0xE800
LCD_SRC                 EQU     0xE801  ; Source for LCD data in shared memory
LCD_DEST                EQU     0xA000  ; Local memory copy destination for lcd data
LCD_BYTES               EQU     0xA0    ; 40 bytes per line, 4 lines = 160 = 0xA0 bytes
LCD_POINTER             EQU     0x8100  ; Pointer for writing bytes to the LCD_DEST buffer
; Lamps
LAMP_SRC                EQU     0xEA00  ; Source for lamp data in shared memory
LAMP_DEST               EQU     0x9500  ; Local memory destination for lamp data
LAMP_BYTES              EQU     0x1B    ; 27 bytes for lamp data

START:
        ZFILL   RX_COUNTER, 0x0100
        CALL    RX_INIT
        LD      A, 3 + 2 + 10 + 2
        LD      (RX_COUNTER), A
        LD      HL, LAMP_MSG
        LD      DE, RX_TYPE
        LD      BC, 3 + 2 + 10 +3
        LDIR
        CALL    CHECK_RX
        HALT

LAMP_MSG:
        DB      0x80
        DB      0x01
        DB      0xFF
LCD_MSG:
        DB      0x81
        DB      "A"
CRAP:
        DB      "ABBACDACDC"
LAMP_MSG2:
        DB      0x80
        DB      0x10
        DB      0x55

        ; Check the RX buffer
CHECK_RX: PROC
        LD      IX, RX_TYPE     ; Start of the buffer
        LD      IY, 0x0000      ; Total bytes consumed
process:
        LD      A, (RX_COUNTER)
        LD      C, 0x01         ; Default consume 1 byte
        AND     A
        JR      NZ, check_msg
        RET                     ; No bytes in buffer
err:    CALL    RX_INIT
        LD      A, "E"
        CALL    OUTCHAR
        RET
check_msg:
        LD      A, (IX)
        CP      0x80            ; Lamp message
        JR      Z, lamp
        CP      0x81            ; LCD message
        JP      Z, lcd
        JR      inc_bytes           ; Invalid message
lamp:
        LD      A, "L"
        CALL    OUTCHAR
        LD      A, (RX_COUNTER)
        CP      0x03
        JP      C, shift        ; not enough bytes for full message
        LD      A, (IX+0x01)    ; Byte offset of the lamp data
        PUSH    AF
        CALL    HEXOUT
        POP     AF
        CP      LAMP_BYTES
        JP      NC, err         ; Too big offset
        LD      BC, 0x0000
        LD      C, A
        LD      HL, LAMP_DEST
        ADD     HL, BC
        LD      A, (IX+0x02)
        LD      (HL), A
        LD      A, "O"
        CALL    OUTCHAR
        ; CALL    RX_INIT
        LD      C, 0x03
        JP      inc_bytes
lcd:
        LD      A, "D"
        CALL    OUTCHAR
        LD      A, (RX_COUNTER)
        CP      0x02
        JP      C, shift
        LD      A, (IX+0x01)
        CALL    OUTCHAR
        ; CALL    LCD_WRITE
        ; CALL    RX_INIT
        LD      C, 0x02
        JP      inc_bytes
inc_bytes:
        LD      A, "i"
        CALL    OUTCHAR
        ; LD      HL, RX_TYPE
        ; LD      DE, RX_TYPE
        LD      B, 0x00
        ADD     IX, BC
        ADD     IY, BC
        ; ADD     HL, BC

        LD      A, (RX_COUNTER)         ; Shift out C bytes from the buffer
        SUB     C
        LD      (RX_COUNTER), A
        ; LDIR

        CALL    HEXOUT
        LD      A, (RX_COUNTER)
        AND     A
        JP      NZ, process             ; Still bytes to check
shift:
        LD      A, "s"
        CALL    OUTCHAR
        LD      HL, RX_TYPE
        PUSH    IY
        POP     BC
        ADD     HL, BC
        LD      DE, RX_TYPE
        LD      A, (RX_COUNTER)
        LD      C, A
        LD      B, 0x00
        LDIR
        LD      HL, RX_TYPE     ; Set the write pointer
        LD      B, 0x00
        LD      C, A
        ADD     HL, BC
        LD      (RX_POINTER), HL
        RET
ENDP


RX_INIT:
        LD      A, 0x00         ; Error in message, zero the counter and type
        LD      (RX_COUNTER), A
        LD      (RX_TYPE), A
        LD      HL, RX_POINTER+2
        LD      (RX_POINTER), HL
        RET


; Output A as hex
HEXOUT:
	CALL	DIG2
	LD	A,10
	CALL	OUTCHAR
	LD	A,13
	CALL	OUTCHAR
	RET

DIG2:
	PUSH	AF		; Nibble-Paar ausdrucken
	RRCA			; unterstes Nibble retten
	RRCA			; oberes Nibble rechtsbuendig
	RRCA			; positionieren
	RRCA
	AND	00001111b
	ADD	A,90H		; binaer in ASCII (hex)
	DAA
	ADC	A,40H
	DAA
	CALL	OUTCHAR		; Zeichen ausgeben
	POP	AF		; jetzt unteres Nibble verarbeiten
	AND	00001111b	; Nibble maskieren
	ADD	A,90H		; binaer in ASCII (hex)
	DAA
	ADC	A,40H
	DAA
	CALL	OUTCHAR
	RET

OUTCHAR:			; Zeichen auf Console ausgeben
	OUT	(1),A
	RET

; Macro to zero fill memory
ZFILL(ptr,len) MACRO
        LD  HL,ptr
        LD  DE,ptr+1
        LD  BC,len-1
        LD  (HL),0
        LDIR
ENDM

; Memory copy
MCOPY(src,dest,len) MACRO
        LD HL,src
        LD DE,dest
        LD BC, len
        LDIR
ENDM

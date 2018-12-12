; Check given ADC channel for changes
ADC_CHECK(channel) MACRO
	LD	IX, ADC_0_OLD + (channel*2)
	LD	IY, ADC_0_DEST + (channel*2)
	LD	A, (IX)
	LD	B, (IY)
	LD	(IX), B				; Store the new value as "old"
	LD	C, channel
	CP	B
	CALL	NZ, ADC_SEND
ENDM

; Read an ADC channel
ADC_R(ch,dest) MACRO
	LD	B, 0xDF
	LD	C, ch
	LD	IX, dest
	CALL	ADC_READ
ENDM

; Transmit contents of a keyboard buffer via SIO A
TX_A_KB(keyboard_ptr) MACRO
	LD	IX, keyboard_ptr
	LD	A, (IX)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x01)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x02)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x03)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x04)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x05)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x06)
	CALL	SIO_A_TX_BLOCKING
	LD	A, (IX+0x07)
	CALL	SIO_A_TX_BLOCKING
ENDM

; Macro to zero fill memory
ZFILL(ptr,len) MACRO
	LD  HL, ptr
	LD  DE, ptr+1
	LD  BC, len-1
	LD  (HL), 0
	LDIR
ENDM

; Memory copy
MCOPY(src,dest,len) MACRO
	LD HL, src
	LD DE, dest
	LD BC, len
	LDIR
ENDM

; Send (blocking) data with SIO port A
TX_A(data) MACRO
	LD	A, data
	CALL	SIO_A_TX_BLOCKING
ENDM

CP_HL_C(value, jump) MACRO
	LD	A, H
	CP	value / 0x0100
	JP	C, jump
	LD	A, L
	CP	value % 0x0100
	JP	C, jump
ENDM

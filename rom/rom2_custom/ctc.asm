	; --- START PROC L0217 ---
INIT_CTC:
L0217:	LD	HL,CTC_1_DATA
	LD	C,10h
	LD	B,03h
	OTIR
	LD	HL,CTC_2_DATA
	LD	C,11h
	LD	B,02h
	OTIR
	LD	HL,CTC_3_DATA
	LD	C,12h
	LD	B,02h
	OTIR
	LD	HL,CTC_4_DATA
	LD	C,13h
	LD	B,02h
	OTIR
	RET

CTC_1_DATA:
	DB	CTC_VECTORS % 0x0100	 ; Interrupt vector location
	DB	05Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset WAS 0DDh
	DB	0Ch
CTC_2_DATA:
	DB	5Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	01h
CTC_3_DATA:
	DB	5Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	01h
CTC_4_DATA:
	DB	5Dh	; No interrupt, Counter, prescaler 16, rising edge, clk starts, time constant follows, no reset
	DB	01h
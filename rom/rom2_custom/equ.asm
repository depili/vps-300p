; SIO commands
SIO_CMD_NULL			EQU	0x00
SIO_CMD_TX_ABORT		EQU	0x08
SIO_CMD_RST_STATUS		EQU	0x10
SIO_CMD_CH_RESET		EQU	0x18
SIO_CMD_EI_RX			EQU	0x20
SIO_CMD_TX_IR_RST		EQU	0x28
SIO_CMD_ERROR_RST		EQU	0x30
SIO_CMD_RETI			EQU	0x38

; Serial receive
RX_COUNTER			EQU	0x9000
RX_POINTER			EQU	0x9002  ; Write pointer
RX_TYPE				EQU	0x9004  ; first byte of the message
RX_MAX_BYTES			EQU	0xFF - 0x06

; Serial transmit, 0xFF byte ring buffer
TX_COUNTER			EQU	0x9100	; Bytes in the TX buffer
TX_PTR_WRITE			EQU	0x9102	; Write pointer
TX_PTR_READ			EQU	0x9104	; Read pointer
TX_BUFFER_START			EQU	0x9200	; Start of the TX buffer, needs to be 0xXX00 as we use 8 bit compare

; LCD
LCD_FLAG			EQU	0xF800
LCD_SRC				EQU	0xA000  ; Source for LCD data
LCD_DEST			EQU	0xF801  ; Shared memory copy destination for lcd data
LCD_BYTES			EQU	0xA0    ; 40 bytes per line, 4 lines = 160 = 0xA0 bytes
LCD_POINTER			EQU	0x8100  ; Pointer for writing bytes to the LCD_SRC buffer
LCD_WRITE_DEST			EQU	LCD_DEST

; Lamps
LAMP_SRC			EQU	0xEA00  ; Source for lamp data in shared memory
LAMP_DEST			EQU	0x9500  ; Local memory destination for lamp data
LAMP_BYTES			EQU	0x1B    ; 27 bytes for lamp data

; Keyboard commands
KEYB_CMD_END			EQU	0xE0
KEYB_CMD_WRITE_DISPLAY  	EQU	0x90    ; Auto increment from address 0
KEYB_CMD_READ_DATA		EQU	0x50    ; Auto increment from address 0
; Keyboard memory
KEYB_1_DEST			EQU	0x9100
KEYB_1_DISP			EQU	0x9110
KEYB_2_DEST			EQU	0x9200
KEYB_2_DISP			EQU	0x9210
KEYB_3_DEST			EQU	0x9300
KEYB_3_DISP			EQU	0x9310

; ADC
ADC_CMD_TRIGGER			EQU	0xFF
ADC_0_DEST			EQU	0xA000
ADC_1_DEST			EQU	0xA002
ADC_2_DEST			EQU	0xA004
ADC_3_DEST			EQU	0xA006
ADC_4_DEST			EQU	0xA008
ADC_5_DEST			EQU	0xA00A
ADC_6_DEST			EQU	0xA00C
ADC_READ_COUNTER		EQU	0xA00E

; IO PORTS
KEYB_1_DATA			EQU	0x00
KEYB_1_CMD			EQU	0x01
KEYB_2_DATA			EQU	0x02
KEYB_2_CMD			EQU	0x03
KEYB_3_DATA			EQU	0x04
KEYB_3_CMD			EQU	0x05
ADC_MUX				EQU	0x08
ADC_IO				EQU	0x09
CTC_CH1				EQU	0x10
CTC_CH2				EQU	0x11
CTC_CH3				EQU	0x12
CTC_CH4				EQU	0x13
SIO_A_DATA			EQU	0x18
SIO_A_CMD			EQU	0x19
SIO_B_DATA			EQU	0x1A
SIO_B_CMD			EQU	0x1B
PIO_A_DATA			EQU	0x1C
PIO_A_CMD			EQU	0x1D
PIO_B_DATA			EQU	0x1E
PIO_B_CMD			EQU	0x1F


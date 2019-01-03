; SIO commands
SIO_CMD_NULL			EQU	0x00
SIO_CMD_TX_ABORT		EQU	0x08
SIO_CMD_RST_STATUS		EQU	0x10
SIO_CMD_CH_RESET		EQU	0x18
SIO_CMD_EI_RX			EQU	0x20
SIO_CMD_TX_IR_RST		EQU	0x28
SIO_CMD_ERROR_RST		EQU	0x30
SIO_CMD_RETI			EQU	0x38

; Serial protocol command bytes
TX_CMD_KB_1			EQU	0x91	; Keyboard events
TX_CMD_KB_2			EQU	0x92
TX_CMD_KB_3			EQU	0x93
TX_CMD_ADC_0			EQU	0xA0	; ADC changes
TX_CMD_ADC_1			EQU	0xA1
TX_CMD_ADC_2			EQU	0xA2
TX_CMD_ADC_3			EQU	0xA3
TX_CMD_ADC_4			EQU	0xA4
TX_CMD_ADC_5			EQU	0xA5
TX_CMD_ADC_6			EQU	0xA6

RX_CMD_LAMP			EQU	0x80	; Lamps via shared mem
RX_CMD_LCD			EQU	0x81	; LCD via shared meme
RX_CMD_LCD_HOME			EQU	0x82	; Reset LCD write pointer
RX_CMD_FRAME			EQU	0x83	; Frame rate 7 segment displays
RX_CMD_PATTERN			EQU	0x84	; Pattern number 7 segment displays
RX_CMD_PWM			EQU	0x85	; PWM via keyboard 2


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
LCD_SRC				EQU	0xC000  ; Source for LCD data
LCD_DEST			EQU	0xF801  ; Shared memory copy destination for lcd data
LCD_BYTES			EQU	0xA0    ; 40 bytes per line, 4 lines = 160 = 0xA0 bytes
LCD_POINTER			EQU	0x8100  ; Pointer for writing bytes to the LCD_SRC buffer
LCD_WRITE_DEST			EQU	LCD_DEST

; Lamps
LAMP_SRC			EQU	0xFA00  ; Source for lamp data in shared memory
LAMP_DEST			EQU	0x9500  ; Local memory destination for lamp data
LAMP_BYTES			EQU	0x1B    ; 27 bytes for lamp data

; Keyboard commands
KEYB_CMD_END			EQU	0xE0
KEYB_CMD_WRITE_DISPLAY  	EQU	0x90    ; Auto increment from address 0
KEYB_CMD_READ_DATA		EQU	0x50    ; Auto increment from address 0
; Keyboard memory
KEYB_CMD_BYTE			EQU	0xB000	; Command byte to send with encoded key event
KEYB_1_DEST			EQU	0xB100	; New data read from the keyboard
KEYB_1_OLD			EQU	0xB108	; Old data to compare
KEYB_1_DISP			EQU	0xB110	; Display data buffer
PATTERN_0			EQU	0xB110	; Pattern number, least significant digit
PATTERN_1			EQU	0xB111
PATTERN_2			EQU	0xB112	; Pattern number, most significant digit
FRAMERATE_0			EQU	0xB115	; Frame rate, least significant digit
FRAMERATE_1			EQU	0xB114
FRAMERATE_3			EQU	0xB113	; Frame rate, most significant digit
KEYB_2_DEST			EQU	0xB200
KEYB_2_OLD			EQU	0xB208
KEYB_2_DISP			EQU	0xB210
KEYB_3_DEST			EQU	0xB300
KEYB_3_OLD			EQU	0xB308
KEYB_3_DISP			EQU	0xB310

; ADC
ADC_CMD_TRIGGER			EQU	0xFF
ADC_0_DEST			EQU	0xA000	; Hue
ADC_0_OLD			EQU	0xA010	; Old value, we only send changes
ADC_1_DEST			EQU	0xA002	; Joystick Y (up-down)
ADC_1_OLD			EQU	0xA012
ADC_2_DEST			EQU	0xA004	; Clip
ADC_2_OLD			EQU	0xA014
ADC_3_DEST			EQU	0xA006	; T-bar
ADC_3_OLD			EQU	0XA016
ADC_4_DEST			EQU	0xA008	; Gain
ADC_4_OLD			EQU	0xA018
ADC_5_DEST			EQU	0xA00A	; Joystick Z (rotate)
ADC_5_OLD			EQU	0xA01A
ADC_6_DEST			EQU	0xA00C	; Joystick X (left-right)
ADC_6_OLD			EQU	0xA01C
ADC_READ_COUNTER		EQU	0xA00E

; LCD numbers and letters
SEG_OFF				EQU	0xFF	; All off
SEG_0				EQU	0x40	; 0
SEG_1				EQU	0x79	; 1
SEG_2				EQU	0x24	; 2
SEG_3				EQU	0x30	; 3
SEG_4				EQU	0x19	; 4
SEG_5				EQU	0x12	; 5
SEG_6				EQU	0x02	; 6
SEG_7				EQU	0x78	; 7
SEG_8				EQU	0x80	; 8
SEG_9				EQU	0x10	; 9
SEG_A				EQU	0x08	; A
SEG_B				EQU	0x03	; B
SEG_C				EQU	0x46	; C
SEG_D				EQU	0x21	; d
SEG_E				EQU	0x06	; E
SEG_F				EQU	0x0E	; F

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


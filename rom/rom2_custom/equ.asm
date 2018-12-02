; SIO commands
SIO_CMD_NULL            EQU     0x00
SIO_CMD_TX_ABORT        EQU     0x08
SIO_CMD_RST_STATUS      EQU     0x10
SIO_CMD_CH_RESET        EQU     0x18
SIO_CMD_EI_RX           EQU     0x20
SIO_CMD_TX_IR_RST       EQU     0x28
SIO_CMD_ERROR_RST       EQU     0x30
SIO_CMD_RETI            EQU     0x38

; Keyboard commands
KEYB_CMD_END            EQU     0xE0
KEYB_CMD_WRITE_DISPLAY  EQU     0x90    ; Auto increment from address 0
KEYB_CMD_READ_DATA      EQU     0x50    ; Auto increment from address 0
; Keyboard memory
KEYB_1_DEST             EQU     0x9100
KEYB_1_DISP             EQU     0x9110
KEYB_2_DEST             EQU     0x9200
KEYB_2_DISP             EQU     0x9210
KEYB_3_DEST             EQU     0x9300
KEYB_3_DISP             EQU     0x9310

; ADC
ADC_CMD_TRIGGER         EQU     0xFF
ADC_0_DEST              EQU     0xA000
ADC_1_DEST              EQU     0xA002
ADC_2_DEST              EQU     0xA004
ADC_3_DEST              EQU     0xA006
ADC_4_DEST              EQU     0xA008
ADC_5_DEST              EQU     0xA00A
ADC_6_DEST              EQU     0xA00C
ADC_READ_COUNTER        EQU     0xA00E

; IO PORTS
KEYB_1_DATA             EQU     0x00
KEYB_1_CMD              EQU     0x01
KEYB_2_DATA             EQU     0x02
KEYB_2_CMD              EQU     0x03
KEYB_3_DATA             EQU     0x04
KEYB_3_CMD              EQU     0x05
ADC_MUX                 EQU     0x08
ADC_IO                  EQU     0x09
CTC_CH1                 EQU     0x10
CTC_CH2                 EQU     0x11
CTC_CH3                 EQU     0x12
CTC_CH4                 EQU     0x13
SIO_A_DATA              EQU     0x18
SIO_A_CMD               EQU     0x19
SIO_B_DATA              EQU     0x1A
SIO_B_CMD               EQU     0x1B
PIO_A_DATA              EQU     0x1C
PIO_A_CMD               EQU     0x1D
PIO_B_DATA              EQU     0x1E
PIO_B_CMD               EQU     0x1F


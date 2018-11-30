; Macro to zero fill memory
ZFILL(ptr,len) MACRO
        LD  HL,ptr
        LD  DE,ptr+1
        LD  BC,len-1
        LD  (HL),0
        LDIR
ENDM

MCOPY(src,dest,len) MACRO
        LD HL,src
        LD DE,dest
        LD BC, len
        LDIR
ENDM

TX_A(data) MACRO
        LD      A, data
        CALL    SIO_A_TX_BLOCKING
ENDM

#!/bin/sh
FILES="equ.asm macros.asm main.asm ctc.asm pio.asm sio.asm interrupts.asm keyboard.asm adc.asm"
# Combine the files
cat $FILES > rom2_combined.asm
mono ../yaza.exe rom2_combined.asm && srec_cat rom2_combined.bin -binary -o rom2_combined.hex -intel

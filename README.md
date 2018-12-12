# vps-300p

Reverse engineering a controller unit for the VPS-300P video mixer. It started with the serial protocol between the mixer and the control unit and later moved to dissassembling the ROMs and writing new ones...


## Current status

We have dumped, dissassembled and puzzled over the ROMs. There are two Z80 compatible processors in the control unit.

### Verified IO
- [x] Rom1 SIO A out and in, interrupt ok
- [x] Rom2 SIO A out and in, interrupt ok
- [x] LCD display, rom1 cpu, PIO ports
- [x] Button indicator leds, rom1 IO 00-03
- [x] Keyboards, rom2 IO 00-06
- [x] ADC, channels mapped, rom2 IO 08 and 09
- [x] 7 segment displays, rom2 keyboard1 display, weird mapping
- [ ] Indicator lamp PWM on some buttons, rom2 keyboard2 display?
- [x] Beeper, rom2 PIO A bit 0, active low
- [x] Serial control for LCD (rudimentary)
- [x] Serial control for button indicator lamps
- [ ] Serial communication for keyboard events (send whole sensor matrix or encode?)
- [ ] Serial communication for ADC changes
- [ ] Serial communication for 7 segment control
- [ ] Shared memory communication, allow control with only one serial line
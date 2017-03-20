package decode

import (
	"fmt"
)

func Decode(msg []byte) string {
	decode := ""
	switch msg[0] {
	case 0x84:
		// Analog stuff?
		decode = decode_analog(msg)
	case 0x86:
		// Buttons, toggles and lamps?
		decode = decode_button_lamp(msg)
	case 0xE0:
		decode = "Ping?"
		if msg[1] == 00 && msg[2] == 00 {
			switch msg[3] {
			case 0x00:
				decode = "Ping1"
			case 0xf0:
				decode = "Ping2"
			default:
				decode = "UNKNOWN Ping?"
			}
		} else if msg[1] == 0 && msg[2] == 0x02 && msg[3] == 0xF0 {
			decode = "Ping3"
		} else {
			decode = "UNKNOWN Ping?"
		}
	}
	return decode
}

func toggle(name string, msg []byte) string {
	switch msg[2] {
	case 0x00:
		return fmt.Sprintf("%s toggle off", name)
	case 0x01:
		return fmt.Sprintf("%s toggle on", name)
	default:
		return fmt.Sprintf("%s toggle UNKNOWN STATE", name)
	}
}

// Button messages
func button(button string) string {
	return fmt.Sprintf("Button %s", button)
}

// Lamp messages
func lamp(lamp string) string {
	return fmt.Sprintf("Lamp %s", lamp)
}

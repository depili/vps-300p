package decode

import (
	"fmt"
)

// Messages with 0x86 as the first byte
func decode_button_lamp(msg []byte) string {
	decode := "UNKNOWN button or lamp?"

	switch msg[1] {
	case 0x01, 0x02:
		var row string
		if msg[1] == 0x01 {
			row = "PGM"
		} else {
			row = "PST"
		}
		switch msg[2] {
		case 0x00, 0x80:
			decode = button(fmt.Sprintf("%s Black", row))
		case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
			fallthrough
		case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
			decode = button(fmt.Sprintf("%s %d", row, msg[2]&0x0F))
		case 0x09, 0x89:
			decode = button(fmt.Sprintf("%s BKGD", row))
		default:
			decode = button(fmt.Sprintf("%s UNKNOWN", row))
		}
	case 0x03:
		// Pattern number
		decode = analog("Pattern number", msg)
	case 0x04:
		switch msg[2] {
		case 0x00:
			decode = "Mix toggle on"
		case 0x01:
			decode = "Wipe toggle on"
		}
	case 0x05:
		decode = toggle("Pattern REV", msg)
	case 0x08:
		switch msg[2] {
		case 0x00:
			decode = "Wipe border off"
		case 0x01:
			decode = "Wipe border matte"
		case 0x02:
			decode = "Wipe border design"
		}
	case 0x09:
		// top row of 6 buttons next to t-bar
		// key1 button sends nothing in current config
		switch msg[2] {
		case 0x00:
			decode = "BKGD, Key 1, Key 2 toggle off"
		case 0x01:
			decode = "BKGD toggle on"
		case 0x02:
			decode = "Key1 toggle on"
		}
	case 0x1A:
		// Key source
		decode = "Key source UNKNOWN"
		switch msg[2] {
		case 0x00:
			decode = "Key source 1"
		case 0x01:
			decode = "Key source 2"
		case 0x02:
			decode = "Key source LUM"
		case 0x03:
			// Speculation
			decode = "Key source chroma key?"
		}
	case 0x1B:
		// Keyer toggles?
		decode = "Keyer stuff?"
		switch msg[2] {
		case 0x00:
			decode = "Keyer video1 toggle"
		case 0x01:
			decode = "Keyer video2 toggle"
		case 0x04:
			decode = "Keyer matt toggle"
		}
	case 0x1C:
		decode = toggle("Keyer key invert", msg)
	case 0x1D:
		decode = toggle("Keyer mask", msg)
	case 0x1E:
		decode = toggle("Keyer mask invert", msg)
	case 0x1F:
		// Keyer edge toggle
		switch msg[2] {
		case 0x00:
			decode = "Keyer edge & shadow toggle off"
		case 0x01:
			decode = "Keyer edge toggle on"
		case 0x02:
			decode = "Keyer shadow toggle on"
		}
	case 0x20:
		switch msg[2] {
		case 0x00:
			decode = "Key2 edge width 1H"
		case 0x01:
			decode = "Key2 edge width 2H"
		case 0x02:
			decode = "Key3 edge width 4H"
		}
	case 0x28:
		if msg[2] == 0 && msg[3] == 0 {
			decode = "PGM button ACK"
		} else if msg[2] == 0x01 && msg[3] == 0 {
			decode = "PST button ACK"
		} else {
			decode = "UNKNOWN ACK?"
		}
	case 0x2C:
		decode = toggle("NEG lum", msg)
	case 0x2D:
		decode = toggle("NEG chroma", msg)
	case 0x32:
		mode := "UNKNOWN"
		switch msg[2] {
		case 0x00:
			mode = "hard"
		case 0x01:
			mode = "soft"
		case 0x02:
			mode = "stardust"
		case 0x03:
			mode = "block 1"
		case 0x04:
			mode = "block 2"
		case 0x05:
			mode = "block 3"
		}
		decode = fmt.Sprintf("Wipe2 edge %s", mode)
	case 0x33:
		if msg[2] == 0x01 {
			decode = button("CUT")
		}
	case 0x36:
		decode = analog("BKGD mode", msg)
	case 0x37:
		decode = analog("PGM frame rate", msg)
	case 0x38:
		decode = analog("Black frame rate", msg)
	case 0x39:
		decode = analog("DSK frame rate", msg)
	case 0x5A:
		switch msg[2] {
		case 0x00:
			decode = "Wipe spin fixed 0.5 spins"
		case 0x01:
			decode = "Wipe spin fixed 1 spins"
		case 0x02:
			decode = "Wipe spin fixed 1.5 spins"
		case 0x03:
			decode = "Wipe spin fixed 2 spins"
		}
	case 0x5B:
		mode := "UNKNOWN"
		switch msg[2] {
		case 0x00:
			mode = "off"
		case 0x01:
			mode = "fixed"
		case 0x02:
			mode = "free"
		}
		decode = fmt.Sprintf("Wipe spin mode %s", mode)
	case 0x5C:
		switch msg[2] {
		case 0x00:
			decode = "Wipe2 spin dir CCW"
		case 0x01:
			decode = "Wipe2 spin dir CW"
		}
	case 0x64, 0x65:
		// Upper row = 0x64, lower = 0x65
		var row string
		if msg[1] == 0x64 {
			row = "PGM"
		} else {
			row = "PST"
		}
		switch msg[2] {
		case 0x00, 0x60, 0x80:
			decode = lamp(fmt.Sprintf("%s Black", row))
		case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
			fallthrough
		case 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68:
			fallthrough
		case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
			decode = lamp(fmt.Sprintf("%s %d", row, msg[2]&0x0F))
		case 0x09, 0x69, 0x89:
			decode = lamp(fmt.Sprintf("%s BKGD", row))
		case 0x40:
			decode = lamp("PST bright Black")
		case 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48:
			fallthrough
		case 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8:
			decode = lamp(fmt.Sprintf("PST bright %d", msg[2]&0x0F))
		case 0x49:
			decode = lamp("PST bright BKGD")
		default:
			decode = fmt.Sprintf("UNKNOWN %s Lamp?", row)
		}
	case 0x67:
		decode = toggle("Key2 led", msg)
	case 0x68:
		decode = "Transition lamps: "
		if msg[2]&0x01 != 0 {
			decode = fmt.Sprintf("%s <low black trans>", decode)
		}
		if msg[2]&0x02 != 0 {
			decode = fmt.Sprintf("%s <bright black trans>", decode)
		}
		if msg[2]&0x04 != 0 {
			decode = fmt.Sprintf("%s <low DSK trans>", decode)
		}
		if msg[2]&0x08 != 0 {
			decode = fmt.Sprintf("%s <bright DSK trans>", decode)
		}
		if msg[2]&0x10 != 0 {
			decode = fmt.Sprintf("%s <low PGM trans>", decode)
		}
		if msg[2]&0x20 != 0 {
			decode = fmt.Sprintf("%s <bright PGM trans>", decode)
		}
	case 0x69:
		switch msg[2] {
		case 0x00:
			decode = lamp("PGM and PST 3D effect off")
		case 0x01:
			decode = lamp("PGM 3D effect")
		case 0x02:
			decode = lamp("PST 3D effect")
		}
	case 0x6A:
		// T-bar leds?
		decode = "T-bar led?"
		switch msg[2] {
		case 0x00:
			decode = "T-bar upwards arrow"
		case 0x01:
			decode = "T-bar downwards arrow"
		}
	}
	return decode
}

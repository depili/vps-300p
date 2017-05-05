package decode

import (
	"fmt"
)

// Messages with 0x84 as first byte are "analog" values
func decode_analog_84(msg []byte) string {
	name := fmt.Sprintf("UNKNOWN %02X", msg[1])
	switch msg[1] {
	case 0x01:
		name = "Hue"
	case 0x03:
		name = "Gain"
	case 0x08:
		name = "Clip"
	case 0x0E:
		name = "Key2 shadow pos X"
	case 0x0F:
		name = "Key2 shadow pos Y"
	case 0x11:
		name = "Key2 transp"
	case 0x1A:
		name = "Mask2 top"
	case 0x1B:
		name = "Mask2 bottom"
	case 0x1C:
		name = "Mask2 right"
	case 0x1D:
		name = "Mask2 left"
	case 0x21:
		name = "Key2 matt hue"
	case 0x22:
		name = "Key2 matt sat"
	case 0x23:
		name = "Key2 matt lum"
	case 0x27:
		name = "Key2 edge hue"
	case 0x28:
		name = "Key2 edge sat"
	case 0x29:
		name = "Key2 edge lum"
	case 0x2A:
		name = "BKGD hue"
	case 0x2B:
		name = "BKGD sat"
	case 0x2C:
		name = "BKGD lum"
	case 0x2D:
		name = "Wipe pos X (signed value!)"
	case 0x2E:
		name = "Wipe pos Y (signed value!)"
	case 0x31:
		name = "Wipe edge soft/width"
	case 0x32:
		name = "Wipe bord width"
	case 0x37:
		name = "Wipe border hue"
	case 0x38:
		name = "Wipe border sat"
	case 0x39:
		name = "Wipe border lum"
	case 0x3A:
		name = "Wipe free speed"
	case 0x41:
		name = "Paint Y"
	case 0x42:
		name = "Paint C"
	case 0x43:
		name = "Mosaic"
	case 0x4D:
		name = "T-bar"
	}
	return analog(name, msg)
}

// Decode a analog value message
func analog(control string, msg []byte) string {
	value := int(msg[2]) + ((int(msg[3]) & 0x7F) * 256)
	if msg[3]&0x80 > 0 {
		value = -32767 + value
	}
	return fmt.Sprintf("Analog %s value: %04d", control, value)
}

package main

import (
	"fmt"
	"github.com/fatih/color"
	ui "github.com/gizak/termui"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"os"
	"strings"
	"time"
)

var Options struct {
	Serial1 string `long:"serial1" description:"T1" default:"/dev/ttyUSB0"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	gTbar := ui.NewGauge()
	gTbar.Percent = 40
	gTbar.Width = 102
	gTbar.Height = 3
	gTbar.BorderLabel = "T-bar"
	gTbar.BarColor = ui.ColorRed
	gTbar.BorderFg = ui.ColorWhite
	gTbar.BorderLabelFg = ui.ColorCyan

	gJoyX := ui.NewGauge()
	gJoyX.Percent = 40
	gJoyX.Y = 3
	gJoyX.Width = 102
	gJoyX.Height = 3
	gJoyX.BorderLabel = "Joystick X"
	gJoyX.BarColor = ui.ColorRed
	gJoyX.BorderFg = ui.ColorWhite
	gJoyX.BorderLabelFg = ui.ColorCyan

	gJoyY := ui.NewGauge()
	gJoyY.Percent = 40
	gJoyY.Y = 6
	gJoyY.Width = 102
	gJoyY.Height = 3
	gJoyY.BorderLabel = "Joystick Y"
	gJoyY.BarColor = ui.ColorRed
	gJoyY.BorderFg = ui.ColorWhite
	gJoyY.BorderLabelFg = ui.ColorCyan

	gJoyZ := ui.NewGauge()
	gJoyZ.Percent = 40
	gJoyZ.Y = 9
	gJoyZ.Width = 102
	gJoyZ.Height = 3
	gJoyZ.BorderLabel = "Joystick Z"
	gJoyZ.BarColor = ui.ColorRed
	gJoyZ.BorderFg = ui.ColorWhite
	gJoyZ.BorderLabelFg = ui.ColorCyan

	gClip := ui.NewGauge()
	gClip.Percent = 40
	gClip.Y = 12
	gClip.Width = 102
	gClip.Height = 3
	gClip.BorderLabel = "Clip"
	gClip.BarColor = ui.ColorRed
	gClip.BorderFg = ui.ColorWhite
	gClip.BorderLabelFg = ui.ColorCyan

	gGain := ui.NewGauge()
	gGain.Percent = 40
	gGain.Y = 15
	gGain.Width = 102
	gGain.Height = 3
	gGain.BorderLabel = "Gain"
	gGain.BarColor = ui.ColorRed
	gGain.BorderFg = ui.ColorWhite
	gGain.BorderLabelFg = ui.ColorCyan

	gHue := ui.NewGauge()
	gHue.Percent = 40
	gHue.Y = 18
	gHue.Width = 102
	gHue.Height = 3
	gHue.BorderLabel = "Hue"
	gHue.BarColor = ui.ColorRed
	gHue.BorderFg = ui.ColorWhite
	gHue.BorderLabelFg = ui.ColorCyan

	pKB := ui.NewParagraph("Keyboard")
	pKB.Height = 22
	pKB.Width = 102
	pKB.Y = 21
	pKB.BorderLabel = "Keyboard"
	pKB.BorderFg = ui.ColorYellow
	kbText := make([]string, 20)

	s1Color := color.New(color.FgYellow)
	/*
		s2Color := color.New(color.FgRed)
		s3Color := color.New(color.FgBlue)
		s4Color := color.New(color.FgGreen)
	*/
	s1Chan := make(chan string)
	adcChan := make(chan []int)

	go serialListen(Options.Serial1, s1Chan, adcChan)

	filename := fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02T15:04:05"))
	s1Color.Printf("Writing log to file: %s\n", filename)
	f, err := os.Create(filename)
	check(err)

	uiEvents := ui.PollEvents()

	for {
		select {
		case s := <-s1Chan:
			msg := fmt.Sprintf("T1 => %s\n", s)
			_, err := f.WriteString(msg)
			check(err)
			// s1Color.Printf("%s", msg)
			for i := 0; i < 19; i++ {
				kbText[i] = kbText[i+1]
			}
			kbText[19] = s
			pKB.Text = strings.Join(kbText, "\n")
			ui.Render(pKB)
		case s := <-adcChan:
			gJoyX.Percent = s[0] * 100 / 1023
			gJoyX.Label = fmt.Sprintf("%d", s[0])
			gJoyY.Percent = s[1] * 100 / 1023
			gJoyY.Label = fmt.Sprintf("%d", s[1])
			gJoyZ.Percent = s[2] * 100 / 1023
			gJoyZ.Label = fmt.Sprintf("%d", s[2])
			gClip.Percent = s[3] * 100 / 1023
			gClip.Label = fmt.Sprintf("%d", s[3])
			gGain.Percent = s[4] * 100 / 1023
			gGain.Label = fmt.Sprintf("%d", s[4])
			gHue.Percent = s[5] * 100 / 1023
			gHue.Label = fmt.Sprintf("%d", s[5])
			gTbar.Percent = s[6] * 100 / 1023
			gTbar.Label = fmt.Sprintf("%d", s[6])
			ui.Render(gTbar, gJoyX, gJoyY, gJoyZ, gClip, gGain, gHue)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	}

}

func serialListen(port string, c chan string, adcChan chan []int) {
	serialConfig := serial.Config{
		Name:        port,
		Baud:        57600,
		Size:        8,
		StopBits:    1,
		ReadTimeout: time.Millisecond * 10,
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 256)
	msg := make([]byte, 4)
	msg_field := 0

	hue := 0
	gain := 0
	clip := 0
	tbar := 0
	joyx := 0
	joyy := 0
	joyz := 0

	serial.Flush()

	for {
		n, err := serial.Read(buf)
		if err == nil {
			for i := 0; i < n; i++ {
				if (buf[i] & 0xA0) == 0x80 {
					// High bit is set, start of message
					msg[0] = buf[i]
					msg_field = 1
				} else {
					msg[msg_field] = buf[i]
					msg_field = (msg_field + 1)
				}

				if (msg[0]&0xA0) == 0xA0 && (msg_field == 3) {
					// Analog message
					value := (int(msg[1]) * 8) + int(msg[2]>>4)
					decode := ""
					switch msg[0] {
					case 0xA0:
						decode = "Hue"
						hue = value
					case 0xA1:
						decode = "Joy Y"
						joyy = value
					case 0xA2:
						decode = "Clip"
						clip = value
					case 0xA3:
						decode = "T-bar"
						tbar = value
					case 0xA4:
						decode = "Gain"
						gain = value
					case 0xA5:
						decode = "Joy Z"
						joyz = value
					case 0xA6:
						decode = "Joy X"
						joyx = value
					}
					_ = decode
					/*
						c <- fmt.Sprintf("%02X %02X %02X %s %d", msg[0], msg[1], msg[2], decode, value)
						c <- fmt.Sprintf("Joystick X: %04d\tY: %04d\tZ: %04d\t", joyx, joyy, joyz)
						c <- fmt.Sprintf("Clip: %04d\tGain: %04d\tHue: %04d\tT-bar: %04d", clip, gain, hue, tbar)
					*/
					adc_msg := make([]int, 7)
					adc_msg[0] = joyx
					adc_msg[1] = joyy
					adc_msg[2] = joyz
					adc_msg[3] = clip
					adc_msg[4] = gain
					adc_msg[5] = hue
					adc_msg[6] = tbar
					adcChan <- adc_msg
					msg_field = 0
				} else if (msg[0]&0x90) == 0x90 && (msg_field == 2) {
					// Keyboard
					kb := ""
					keyname := ""
					up := "up"
					switch msg[0] {
					case 0x91:
						kb = "Keyboard 1"
						keyname = kb1_key(msg[1] & 0x3F)
					case 0x92:
						kb = "Keyboard 2"
						keyname = kb2_key(msg[1] & 0x3F)
					case 0x93:
						kb = "Keyboard 3"
						keyname = kb3_key(msg[1] & 0x3F)
					}
					if (msg[1] & 0x40) == 0x40 {
						up = "down"
					}
					c <- fmt.Sprintf("%02X %02X %s key %02X %s %s", msg[0], msg[1], kb, msg[1]&0x3F, keyname, up)
					msg_field = 0
				} else if msg_field == 4 {
					msg_field = 0
				}
			}
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func kb1_key(keycode byte) string {
	keymap := map[byte]string{
		0x39: "F_plusminus",
		0x31: "F_play",
		0x09: "F_f0",
		0x01: "F_f1",
		0x38: "F_f2",
		0x30: "F_f3",
		0x28: "F_f4",
		0x20: "F_f5",
		0x18: "F_f6",
		0x10: "F_f7",
		0x08: "F_f8",
		0x00: "F_f9",
		0x33: "P_rev",
		0x0C: "P_learn",
		0x3B: "P_num_lock",
		0x04: "P_shift",
		0x02: "P_patt_1",
		0x0A: "P_patt_2",
		0x12: "P_patt_3",
		0x1A: "P_patt_4",
		0x22: "P_patt_5",
		0x2A: "P_patt_6",
		0x32: "P_patt_7",
		0x3A: "P_patt_8",
		0x03: "P_patt_9",
		0x0B: "P_patt_10",
		0x21: "P_set_left",
		0x19: "P_set_middle",
		0x11: "P_set_right",
		0x2D: "K_source_1",
		0x25: "K_source_2",
		0x1D: "K_source_lum",
		0x36: "K_source_chroma_key",
		0x05: "K_video_1",
		0x0D: "K_video_2",
		0x3C: "K_3d_effect",
		0x34: "K_chroma_key",
		0x2C: "K_matt",
		0x24: "K_edge",
		0x1C: "K_shadow",
		0x14: "K_outline",
		0x26: "K_key_invert",
		0x1E: "K_mask",
		0x16: "K_mask_invert",
		0x06: "K_priority",
		0x3D: "K_key_1",
		0x35: "K_key_2",
		0x0E: "K_key_ch",
		0x3E: "K_clip",
		0x27: "J_def_x_y_z",
		0x2F: "J_def_x",
		0x37: "J_def_y",
		0x3F: "J_def_z",
		0x1F: "T_select",
		0x17: "T_left",
		0x0F: "T_middle",
		0x07: "T_right",
	}
	return keymap[keycode]
}

func kb2_key(keycode byte) string {
	keymap := map[byte]string{
		0x07: "PGM_black",
		0x05: "PGM_1",
		0x0D: "PGM_2",
		0x15: "PGM_3",
		0x1D: "PGM_4",
		0x25: "PGM_5",
		0x2D: "PGM_6",
		0x35: "PGM_7",
		0x3D: "PGM_8",
		0x0F: "PGM_bkgd",
		0x2C: "PGM_3d_effect",
		0x17: "PST_black",
		0x06: "PST_1",
		0x0E: "PST_2",
		0x16: "PST_3",
		0x1E: "PST_4",
		0x26: "PST_5",
		0x2E: "PST_6",
		0x36: "PST_7",
		0x3E: "PST_8",
		0x1F: "PST_bkgd",
		0x34: "PST_3d_effect",
		0x3B: "BKGD",
		0x0C: "KEY1",
		0x04: "KEY2",
		0x24: "MIX",
		0x1C: "WIPE",
		0x14: "CUT",
		0x2B: "BLACK_TRANS",
		0x23: "DSK_TRANS",
		0x33: "PGM_TRANS",
	}
	return keymap[keycode]
}

func kb3_key(keycode byte) string {
	keymap := map[byte]string{
		0x2F: "S_up",
		0x37: "S_down",
		0x0E: "S_store",
		0x06: "S_recall",
		0x3D: "S_play",
		0x35: "S_seq",
		0x2D: "S_pause",
		0x25: "S_linear",
		0x1D: "S_over",
		0x15: "S_insert",
		0x0D: "S_delete",
		0x05: "S_cancel",
		0x27: "S_0",
		0x1F: "S_1",
		0x17: "S_2",
		0x0F: "S_3",
		0x07: "S_4",
		0x3E: "S_5",
		0x36: "S_6",
		0x2E: "S_7",
		0x26: "S_8",
		0x1E: "S_9",
		0x00: "M_pos",
		0x08: "M_rot",
		0x10: "M_warp_1",
		0x18: "M_warp_2",
		0x20: "M_trail",
		0x28: "M_drop",
		0x30: "M_hilite",
		0x38: "M_dve_edge",
		0x01: "M_sub_eff_1",
		0x09: "M_sub_eff_2",
		0x11: "M_wipe_1",
		0x19: "M_wipe_2",
		0x21: "M_clear",
		0x29: "M_key_2",
		0x31: "M_bkgd",
		0x39: "M_shift",
	}
	return keymap[keycode]
}

package main

import (
	"fmt"
	"github.com/fatih/color"
	ui "github.com/gizak/termui"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"os"
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
	gTbar.Width = 80
	gTbar.Height = 3
	gTbar.BorderLabel = "T-bar"
	gTbar.BarColor = ui.ColorRed
	gTbar.BorderFg = ui.ColorWhite
	gTbar.BorderLabelFg = ui.ColorCyan

	gJoyX := ui.NewGauge()
	gJoyX.Percent = 40
	gJoyX.Y = 3
	gJoyX.Width = 80
	gJoyX.Height = 3
	gJoyX.BorderLabel = "Joystick X"
	gJoyX.BarColor = ui.ColorRed
	gJoyX.BorderFg = ui.ColorWhite
	gJoyX.BorderLabelFg = ui.ColorCyan

	gJoyY := ui.NewGauge()
	gJoyY.Percent = 40
	gJoyY.Y = 6
	gJoyY.Width = 80
	gJoyY.Height = 3
	gJoyY.BorderLabel = "Joystick Y"
	gJoyY.BarColor = ui.ColorRed
	gJoyY.BorderFg = ui.ColorWhite
	gJoyY.BorderLabelFg = ui.ColorCyan

	gJoyZ := ui.NewGauge()
	gJoyZ.Percent = 40
	gJoyZ.Y = 9
	gJoyZ.Width = 80
	gJoyZ.Height = 3
	gJoyZ.BorderLabel = "Joystick Z"
	gJoyZ.BarColor = ui.ColorRed
	gJoyZ.BorderFg = ui.ColorWhite
	gJoyZ.BorderLabelFg = ui.ColorCyan

	gClip := ui.NewGauge()
	gClip.Percent = 40
	gClip.Y = 12
	gClip.Width = 80
	gClip.Height = 3
	gClip.BorderLabel = "Clip"
	gClip.BarColor = ui.ColorRed
	gClip.BorderFg = ui.ColorWhite
	gClip.BorderLabelFg = ui.ColorCyan

	gGain := ui.NewGauge()
	gGain.Percent = 40
	gGain.Y = 15
	gGain.Width = 80
	gGain.Height = 3
	gGain.BorderLabel = "Gain"
	gGain.BarColor = ui.ColorRed
	gGain.BorderFg = ui.ColorWhite
	gGain.BorderLabelFg = ui.ColorCyan

	gHue := ui.NewGauge()
	gHue.Percent = 40
	gHue.Y = 18
	gHue.Width = 80
	gHue.Height = 3
	gHue.BorderLabel = "Hue"
	gHue.BarColor = ui.ColorRed
	gHue.BorderFg = ui.ColorWhite
	gHue.BorderLabelFg = ui.ColorCyan

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
			s1Color.Printf("%s", msg)
		case s := <-adcChan:
			gJoyX.Percent = s[0] / 10
			gJoyY.Percent = s[1] / 10
			gJoyZ.Percent = s[2] / 10
			gClip.Percent = s[3] / 10
			gGain.Percent = s[4] / 10
			gHue.Percent = s[5] / 10
			gTbar.Percent = s[6] / 10
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
		Baud:        38400,
		Parity:      'O',
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
				if (buf[i] & 0x80) == 0x80 {
					// High bit is set, start of message
					msg[0] = buf[i]
					msg_field = 1
				} else {
					msg[msg_field] = buf[i]
					msg_field = (msg_field + 1)
				}

				if (msg[0]&0xA0) == 0xA0 && (msg_field == 3) {
					// Analog message
					value := (int(msg[1]) * 8) + int(msg[2]>>5)
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
					up := "up"
					switch msg[0] {
					case 0x91:
						kb = "Keyboard 1"
					case 0x92:
						kb = "Keyboard 2"
					case 0x93:
						kb = "Keyboard 3"
					}
					if (msg[1] & 0x40) == 0x40 {
						up = "down"
					}
					c <- fmt.Sprintf("%02X %02X %s %02X %s", msg[0], msg[1], kb, msg[1]&0x1F, up)
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

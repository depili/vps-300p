package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"os"
	"time"
)

var Options struct {
	Serial1 string `long:"serial1" description:"T1" default:"/dev/ttyUSB0"`
	Serial2 string `long:"serial2" description:"T2" default:"/dev/ttyUSB0"`
	Serial3 string `long:"serial3" description:"R1" default:"/dev/ttyUSB0"`
	Serial4 string `long:"serial4" description:"R2" default:"/dev/ttyUSB0"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	s1Color := color.New(color.FgYellow)
	s2Color := color.New(color.FgRed)
	s3Color := color.New(color.FgBlue)
	s4Color := color.New(color.FgGreen)

	s1Chan := make(chan string)
	s2Chan := make(chan string)
	s3Chan := make(chan string)
	s4Chan := make(chan string)

	go serialListen(Options.Serial1, s1Chan)
	go serialListen(Options.Serial2, s2Chan)
	go serialListen(Options.Serial3, s3Chan)
	go serialListen(Options.Serial4, s4Chan)

	filename := fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02T15:04:05"))
	fmt.Printf("Writing log to file: %s\n", filename)
	f, err := os.Create(filename)
	check(err)

	for {
		select {
		case s := <-s1Chan:
			msg := fmt.Sprintf("T1 => %s\n", s)
			_, err := f.WriteString(msg)
			check(err)
			s1Color.Printf("%s", msg)
		case s := <-s2Chan:
			msg := fmt.Sprintf("T2 => %s\n", s)
			_, err := f.WriteString(msg)
			check(err)
			s2Color.Printf("%s", msg)
		case s := <-s3Chan:
			msg := fmt.Sprintf("R1 <= %s\n", s)
			_, err := f.WriteString(msg)
			check(err)
			s3Color.Printf("%s", msg)
		case s := <-s4Chan:
			msg := fmt.Sprintf("R2 <= %s\n", s)
			_, err := f.WriteString(msg)
			check(err)
			s4Color.Printf("%s", msg)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func serialListen(port string, c chan string) {
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

	serial.Flush()

	for {
		n, err := serial.Read(buf)
		if err == nil {
			for i := 0; i < n; i++ {
				msg[msg_field] = buf[i]
				msg_field = (msg_field + 1) % 4
				if msg_field == 0 {
					decode := ""
					switch msg[0] {
					case 0x84:
						// Analog stuff?
						switch msg[1] {
						case 0x01:
							decode = analog("Hue", msg)
						case 0x3:
							decode = analog("Gain", msg)
						case 0x08:
							decode = analog("Clip", msg)
						case 0x4D:
							decode = analog("T-bar", msg)
						default:
							decode = analog(fmt.Sprintf("UNKNOWN %02x", msg[1]), msg)
						}
					case 0x86:
						// Buttons and lamps?
						switch msg[1] {
						case 0x01, 0x02:
							var row string
							if msg[1] == 0x01 {
								row = "PGM"
							} else {
								row = "PST"
							}
							switch msg[2] {
							case 0x80:
								decode = button(fmt.Sprintf("%s Black", row))
							case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
								decode = button(fmt.Sprintf("%s %d", row, msg[2]-0x80))
							case 0x89:
								decode = button(fmt.Sprintf("%s BKGD", row))
							default:
								decode = button(fmt.Sprintf("%s UNKNOWN", row))
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
						case 0x28:
							if msg[2] == 0 && msg[3] == 0 {
								decode = "PGM button ACK"
							} else if msg[2] == 0x01 && msg[3] == 0 {
								decode = "PST button ACK"
							} else {
								decode = "UNKNOWN ACK?"
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
							case 0x80:
								decode = lamp(fmt.Sprintf("%s Black", row))
							case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
								decode = lamp(fmt.Sprintf("%s %d", row, msg[2]-0x80))
							case 0x89:
								decode = lamp(fmt.Sprintf("%s BKGD", row))
							case 0xC0:
								decode = lamp(fmt.Sprintf("%s Bright Black", row))
							case 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8:
								decode = lamp(fmt.Sprintf("%s Bright %d", row, msg[2]-0xC0))
							case 0xC9:
								decode = lamp(fmt.Sprintf("%s BKGD", row))
							default:
								decode = fmt.Sprintf("UNKNOWN %s Lamp?", row)
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
						default:
							decode = "UNKNOWN button or lamp?"
						}
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
					c <- fmt.Sprintf("%02X %02X %02X %02X %s", msg[0], msg[1], msg[2], msg[3], decode)
				}
			}
		}
	}
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

// Decode a analog value message
func analog(control string, msg []byte) string {
	value := int(msg[3])*256 + int(msg[2])
	return fmt.Sprintf("Analog %s value: %04d", control, value)
}

// Button messages
func button(button string) string {
	return fmt.Sprintf("Button %s", button)
}

// Lamp messages
func lamp(lamp string) string {
	return fmt.Sprintf("Lamp %s", lamp)
}

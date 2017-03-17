package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
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

	s1Color := color.FgYellow
	s2Color := color.FgRed
	s3Color := color.FgBlue
	s4Color := color.FgGreen

	go serialListen(Options.Serial1, "T1", s1Color)
	go serialListen(Options.Serial2, "T2", s2Color)
	go serialListen(Options.Serial3, "R1", s3Color)
	go serialListen(Options.Serial4, "R2", s4Color)

	for {
		// wait
		time.Sleep(1 * time.Second)
	}
}

func serialListen(port, name string, c color.Attribute) {
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

	for {
		n, err := serial.Read(buf)
		if err == nil {
			color.Set(c)
			for i := 0; i < n; i++ {
				msg[msg_field] = buf[i]
				msg_field = (msg_field + 1) % 4
				if msg_field == 0 {
					switch msg[0] {
					case 0x84:
						// Analog stuff?
						switch msg[1] {
						case 0x01:
							analog(name, "Hue", msg)
						case 0x3:
							analog(name, "Gain", msg)
						case 0x08:
							analog(name, "Clip", msg)
						case 0x4D:
							// T-bar
							analog(name, "T-bar", msg)
						default:
							analog(name, fmt.Sprintf("UNKNOWN %02x", msg[1]), msg)
						}
					case 0x86:
						// Crosspoint stuff
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
								button(name, fmt.Sprintf("%s Black", row))
							case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
								button(name, fmt.Sprintf("%s %d", row, msg[2]-0x80))
							case 0x89:
								button(name, fmt.Sprintf("%s BKGD", row))
							}
						case 0x28:
							if msg[2] == 0 && msg[3] == 0 {
								fmt.Printf("%s: PGM button ACK\n", name)
							} else if msg[2] == 0x01 && msg[3] == 0 {
								fmt.Printf("%s: PST button ACK\n", name)
							} else {
								unknown_msg(name, msg)
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
								lamp(name, fmt.Sprintf("%s Black", row), msg)
							case 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88:
								lamp(name, fmt.Sprintf("%s %d", row, msg[2]-0x80), msg)
							case 0x89:
								lamp(name, fmt.Sprintf("%s BKGD", row), msg)
							default:
								fmt.Printf("%s: Crosspoint lamp: %02x %02x %02x\n", name, msg[1], msg[2], msg[3])
							}
						default:
							fmt.Printf("%s: Crosspoint button: %02x %02x %02x\n", name, msg[1], msg[2], msg[3])
						}
					case 0xE0:
						if msg[1] == 00 && msg[2] == 00 {
							switch msg[3] {
							case 0x00:
								fmt.Printf("%s: Ping1?\n", name)
							case 0xf0:
								fmt.Printf("%s: Ping2?\n", name)
							default:
								unknown_msg(name, msg)
							}
						} else if msg[1] == 0 && msg[2] == 0x02 && msg[3] == 0xF0 {
							fmt.Printf("%s: Ping3?\n", name)
						} else {
							unknown_msg(name, msg)
						}
					default:
						unknown_msg(name, msg)
					}
				}
			}
			color.Unset()
		}
	}
}

func unknown_msg(name string, msg []byte) {
	fmt.Printf("%s: %02x %02x %02x %02x\n", name, msg[0], msg[1], msg[2], msg[3])
}

func analog(name, control string, msg []byte) {
	value := int(msg[3])*256 + int(msg[2])
	fmt.Printf("%s: Analog %s value: %04d\n", name, control, value)
}

func button(name, button string) {
	fmt.Printf("%s: Button %s\n", name, button)
}

func lamp(name, lamp string, msg []byte) {
	fmt.Printf("%s: Lamp %s %02x %02x\n", name, lamp, msg[2], msg[3])
}

func decode_command(b byte) {
	switch b {
	case 0x41, 0x42, 0x43, 0x44:
		fmt.Printf("READ crosspoint")
	case 0xc1, 0xc2, 0xc3, 0xc4:
		fmt.Printf("WRITE crosspoint")
	case 0x45:
		fmt.Printf("READ analog control")
	case 0xc5:
		fmt.Printf("WRITE analog control")
	case 0x46, 0x47:
		fmt.Printf("READ panel/lamp control")
	case 0xc6, 0xc7:
		fmt.Printf("WRITE panel/lamp control")
	case 0x48:
		fmt.Printf("READ wipe pattern")
	case 0xc8:
		fmt.Printf("WRITE wipe pattern")
	case 0x4a:
		fmt.Printf("READ transition mode")
	case 0xca:
		fmt.Printf("WRITE transition mode")
	case 0x4c, 0x4d, 0x7d:
		fmt.Printf("READ transition rate")
	case 0xcc, 0xcd, 0xfd:
		fmt.Printf("WRITE transition rate")
	case 0xda:
		fmt.Printf("WRITE effect memory store")
	case 0xdb:
		fmt.Printf("WRITE effect memory recall")
	case 0x6d:
		fmt.Printf("READ field mode")
	case 0xed:
		fmt.Printf("WRITE field mode")
	case 0xee:
		fmt.Printf("WRITE status update")
	case 0xf2:
		fmt.Printf("WRITE all stop")
	case 0x78:
		fmt.Printf("READ lamp status map")
	case 0xf8:
		fmt.Printf("WRITE lamp status map")
	case 0xfb:
		fmt.Printf("WRITE panel key select")
	case 0x7e:
		fmt.Printf("READ effect memory transfer")
	case 0xfe:
		fmt.Printf("WRITE effect memory transfer")
	default:
		fmt.Printf("UNKNOWN: %02x", b)
	}
}

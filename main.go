package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"time"
)

var Options struct {
	Serial1 string `long:"serial1" description:"Serial1" default:"/dev/ttyUSB0"`
	Serial2 string `long:"serial2" description:"Serial2" default:"/dev/ttyUSB0"`
	Serial3 string `long:"serial3" description:"Serial2" default:"/dev/ttyUSB0"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	s1Color := color.FgYellow
	s2Color := color.FgRed
	s3Color := color.FgBlue

	go serialListen(Options.Serial1, "S1", s1Color)
	go serialListen(Options.Serial2, "S2", s2Color)
	go serialListen(Options.Serial3, "S3", s3Color)

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
		// ReadTimeout:    options.SerialTimeout,
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 256)
	for {
		n, err := serial.Read(buf)
		if err == nil {
			color.Set(c)
			sb := fmt.Sprintf("%s: ", name)
			for i := 0; i < n; i++ {
				sb = sb + fmt.Sprintf("%02x ", buf[i])
			}
			fmt.Println(sb)
			color.Unset()
		}
	}
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

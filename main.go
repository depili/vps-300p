package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"time"
)

var Options struct {
	SerialName string `long:"serial-name" description:"Serial device for arduino" default:"/dev/ttyUSB0"`
	SerialBaud int    `long:"serial-baud" value-name:"BAUD" default:"57600"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	serialConfig := serial.Config{
		Name:        Options.SerialName,
		Baud:        Options.SerialBaud,
		ReadTimeout: time.Millisecond * 10,
		// ReadTimeout:    options.SerialTimeout,
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 128)
	for {
		n, err := serial.Read(buf)
		if err != nil {
			continue
		}
		for i := 0; i < n; i++ {
			fmt.Printf("%02x ", buf[i])
		}
		fmt.Printf("\n")
	}
}

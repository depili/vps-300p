package main

import (
	"fmt"
	"github.com/depili/vps-300p/decode"
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
					decode := decode.Decode(msg)
					c <- fmt.Sprintf("%02X %02X %02X %02X %s", msg[0], msg[1], msg[2], msg[3], decode)
				}
			}
		}
	}
}

package main

import (
	"fmt"
	"github.com/depili/vps-300p/decode"
	"github.com/depili/vps-300p/static"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"os"
	"time"
)

var Options struct {
	Serial1 string `long:"serial1" description:"T1 & R1" default:"/dev/ttyUSB0"`
	Serial2 string `long:"serial2" description:"T2 & R2" default:"/dev/ttyUSB0"`
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	T1Color := color.New(color.FgYellow)
	T2Color := color.New(color.FgRed)
	R1Color := color.New(color.FgBlue)
	R2Color := color.New(color.FgGreen)

	R1Chan := make(chan []byte, 20)
	R2Chan := make(chan []byte, 20)
	T1Chan := make(chan []byte, 20)
	T2Chan := make(chan []byte, 20)

	go serialWorker(Options.Serial1, R1Chan, T1Chan)
	go serialWorker(Options.Serial2, R2Chan, T2Chan)

	filename := fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02T15:04:05"))
	fmt.Printf("Writing log to file: %s\n", filename)
	f, err := os.Create(filename)
	check(err)

	ticker := time.NewTicker(time.Second)
	lamp := byte(0)
	initialized := false

	for {
		select {
		case <-ticker.C:
			if initialized {
				msg := []byte{0x86, 0x64, lamp, 0x00}
				write_msg("T1 <=", msg, f, T1Color)
				T1Chan <- msg
				msg = []byte{0x86, 0x6A, lamp % 2, 0x00}
				write_msg("T1 <=", msg, f, T1Color)
				T1Chan <- msg
				lamp = (lamp + 1) % 10
			}
		case s := <-R1Chan:
			write_msg("R1 =>", s, f, R1Color)
			if p, buf := ping_1(s); p == true {
				T1Chan <- buf
				write_msg("T1 <=", buf, f, T1Color)
			}
			if decode.MsgEq(s, []byte{0x86, 0x63, 0x01, 0x00}) {
				// Initialization request

				for _, buf := range static.T1_init_data() {
					T1Chan <- buf
					write_msg("T1 <=", buf, f, T1Color)
				}
				for _, buf := range static.T2_init_data() {
					T2Chan <- buf
					write_msg("T2 <=", buf, f, T2Color)
				}
				msg := []byte{0xE1, 0x00, 0xFF, 0x7F}
				T1Chan <- msg
				write_msg("T1 <=", msg, f, T1Color)
				initialized = true
			}
		case s := <-R2Chan:
			write_msg("R2 =>", s, f, R2Color)
			if p, buf := ping_2(s); p == true {
				T2Chan <- buf
				write_msg("T2 <=", buf, f, T2Color)
			}
		}
	}
}

func write_msg(line string, buf []byte, f *os.File, c *color.Color) {
	msg := fmt.Sprintf("%s %02X %02X %02X %02X %s\n", line, buf[0], buf[1], buf[2],
		buf[3], decode.Decode(buf))
	_, err := f.WriteString(msg)
	check(err)
	c.Printf("%s", msg)
}

func ping_1(msg []byte) (bool, []byte) {
	if msg[0] == 0xE0 && msg[1] == 0 && msg[2] == 0 && msg[3] == 0 {
		return true, []byte{0xE0, 0x00, 0x02, 0xF0}
	} else {
		return false, nil
	}
}

func ping_2(msg []byte) (bool, []byte) {
	if msg[0] == 0xE0 && msg[1] == 0 && msg[2] == 0 && msg[3] == 0 {
		return true, []byte{0xE0, 0x00, 0x00, 0xF0}
	} else {
		return false, nil
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func serialWorker(port string, r chan []byte, w chan []byte) {
	serialConfig := serial.Config{
		Name:        port,
		Baud:        38400,
		Parity:      'O',
		Size:        8,
		StopBits:    1,
		ReadTimeout: time.Millisecond,
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 256)
	msg := make([]byte, 4)
	msg_field := 0

	serial.Flush()

	go serialWriter(serial, w)

	for {
		n, err := serial.Read(buf)
		if err == nil {
			for i := 0; i < n; i++ {
				msg[msg_field] = buf[i]
				msg_field = (msg_field + 1) % 4
				if msg_field == 0 {
					r <- msg
				}
			}
		}
	}
}

func serialWriter(serial *serial.Port, w chan []byte) {
	for buf := range w {
		_, err := serial.Write(buf)
		check(err)
	}
}
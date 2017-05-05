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
}

var parser = flags.NewParser(&Options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	T1Color := color.New(color.FgYellow)
	R1Color := color.New(color.FgBlue)

	R1Chan := make(chan []byte)
	T1Chan := make(chan []byte)

	go serialWorker(Options.Serial1, R1Chan, T1Chan)

	filename := fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02T15:04:05"))
	fmt.Printf("Writing log to file: %s\n", filename)
	f, err := os.Create(filename)
	check(err)

	ticker := time.NewTicker(1000 * time.Millisecond)
	lamp := byte(0)
	initialized := false

	for {
		select {
		case <-ticker.C:
			if initialized {
				flashLeds(lamp, T1Chan, T1Color, f)
				lamp++
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
				msg := []byte{0xE1, 0x00, 0xFF, 0x7F}
				T1Chan <- msg
				write_msg("T1 <=", msg, f, T1Color)
				initialized = true
			}
		}
	}
}

func flashLeds(lamp byte, T1Chan chan []byte, T1Color *color.Color, f *os.File) {
	// PGM inputs
	msg := []byte{0x86, 0x64, lamp % 10, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	// T-bar up/down arrow
	msg = []byte{0x86, 0x6A, lamp % 2, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	// 3D effect toggle
	msg = []byte{0x86, 0x69, lamp % 3, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	// Key1 source
	msg = []byte{0x86, 0x0C, lamp % 4, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	// Key2 source
	msg = []byte{0x86, 0x1A, lamp % 4, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	// Key1 insert
	msg = []byte{0x86, 0x11, lamp % 4, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	msg = []byte{0x86, 0x0D, lamp % 5, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg
	msg = []byte{0x86, 0x27, lamp % 2, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg

	msg = []byte{0x86, 0x66, lamp % 2, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg

	msg = []byte{0x86, 0x67, lamp % 2, 0x00}
	write_msg("T1 <=", msg, f, T1Color)
	T1Chan <- msg

	if lamp&0x01 > 0 {
		msg = []byte{0x86, 0x0E, 0x01, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
	} else {
		msg = []byte{0x86, 0x0E, 0x00, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
	}
	if lamp&0x02 > 0 {
		msg = []byte{0x86, 0x0F, 0x01, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
	} else {
		msg = []byte{0x86, 0x0F, 0x00, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
	}
	if lamp&0x04 > 0 {
		msg = []byte{0x86, 0x10, 0x01, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
	} else {
		msg = []byte{0x86, 0x10, 0x00, 0x00}
		write_msg("T1 <=", msg, f, T1Color)
		T1Chan <- msg
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
		return true, []byte{0xE0, 0xFF, 0xFF, 0xF0}
	} else {
		return false, nil
	}
}

func ping_2(msg []byte) (bool, []byte) {
	if msg[0] == 0xE0 && msg[1] == 0 && msg[2] == 0 && msg[3] == 0 {
		return true, []byte{0xE0, 0xFF, 0xFF, 0xF0}
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
		Baud:        115200,
		Parity:      'N',
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

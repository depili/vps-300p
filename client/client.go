package client

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

type client struct {
	serial1                            *serial.Port
	serial2                            *serial.Port
	R1Chan, R2Chan, t1Chan, t2Chan     chan []byte
	r1Color, r2Color, t1Color, t2Color *color.Color
}

func Init(serial1, serial2 string) *client {
	var client = client{
		t1Color: color.New(color.FgYellow),
		t2Color: color.New(color.FgRed),
		r1Color: color.New(color.FgBlue),
		r2Color: color.New(color.FgGreen),
		R1Chan:  make(chan []byte),
		R2Chan:  make(chan []byte),
		t1Chan:  make(chan []byte),
		t2Chan:  make(chan []byte),
	}

	go serialWorker(serial1, 1)
	go serialWorker(serial2, 2)
}

func (client *client) serialWorker(port string, int line) {
	serialConfig := serial.Config{
		Name:        port,
		Baud:        38400,
		Parity:      'O',
		Size:        8,
		StopBits:    1,
		ReadTimeout: time.Millisecond,
	}

	r := client.r1chan
	w := client.t1chan
	if line == 2 {
		r = client.r2chan
		w = client.t2chan
	}

	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	if line == 2 {
		client.serial2 = serial
	} else {
		client.serial1 = serial
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

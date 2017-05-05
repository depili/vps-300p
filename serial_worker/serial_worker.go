package serial_worker

import (
	"fmt"
	"github.com/depili/vps-300p/decode"
	"github.com/fatih/color"
	"github.com/tarm/serial"
)

func Init(serialConfig serial.Config, r chan []byte, w chan []byte) {
	serial, err := serial.OpenPort(&serialConfig)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 256)
	msg := make([]byte, 4)
	msg_field := 0

	serial.Flush()

	go serialWriter(serial, w)

	RxColor := color.New(color.FgBlue)

	for {
		n, err := serial.Read(buf)
		if err == nil {
			for i := 0; i < n; i++ {
				msg[msg_field] = buf[i]
				msg_field = (msg_field + 1) % 4
				if msg_field == 0 {
					r <- msg
					write_msg("<=", msg, RxColor)
					msg = make([]byte, 4)
				}
			}
		}
	}
}

func serialWriter(serial *serial.Port, w chan []byte) {
	TxColor := color.New(color.FgYellow)

	for buf := range w {
		_, err := serial.Write(buf)
		write_msg("=>", buf, TxColor)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write_msg(line string, buf []byte, c *color.Color) {
	msg := fmt.Sprintf("%s %02X %02X %02X %02X %s\n", line, buf[0], buf[1], buf[2],
		buf[3], decode.Decode(buf))
	c.Printf("%s", msg)
}

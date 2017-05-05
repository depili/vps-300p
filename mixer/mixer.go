package mixer

import (
	"fmt"
	"github.com/depili/vps-300p/decode"
	"github.com/depili/vps-300p/serial_worker"
	"github.com/depili/vps-300p/static"
	"github.com/tarm/serial"
)

// Transition types
const (
	TransMix  = iota
	TransWipe = iota
)

const (
	SourceBlack = iota
	Source1     = iota
	Source2     = iota
	Source3     = iota
	Source4     = iota
	Source5     = iota
	Source6     = iota
	Source7     = iota
	Source8     = iota
	SourceBkgd  = iota
)

type mixer struct {
	PGM        byte // Active PGM selection
	PST        byte // Active PST selection
	rxChan     chan []byte
	TxChan     chan []byte
	TransType  int
	transDir   bool
	TransValue int
}

func Init(serialConfig serial.Config) *mixer {
	var mixer = mixer{
		PGM:        0,
		PST:        0,
		TransType:  TransWipe,
		transDir:   true,
		TransValue: 0,
		rxChan:     make(chan []byte),
		TxChan:     make(chan []byte),
	}

	go serial_worker.Init(serialConfig, mixer.rxChan, mixer.TxChan)
	go mixer.stateKeeper()

	return &mixer
}

func (mixer *mixer) stateKeeper() {
	for {
		select {
		case msg := <-mixer.rxChan:
			switch msg[0] {
			case 0x84:
				// Analog readings
				switch msg[1] {
				case 0x4D:
					// T-bar
					value := analog(msg)
					if mixer.transDir {
						if value == 1023 {
							// Completed transition
							mixer.transDir = false
							mixer.TransValue = 1023 - value
							mixer.cut()
							mixer.TxChan <- []byte{0x86, 0x6A, 0x01, 0x00} // Upwards arrow
						} else {
							mixer.TransValue = value
						}
					} else {
						if value == 0 {
							// Completed transition
							mixer.transDir = true
							mixer.TransValue = value
							mixer.cut()
							mixer.TxChan <- []byte{0x86, 0x6A, 0x00, 0x00} // Downwards arrow
						} else {
							mixer.TransValue = 1023 - value
						}
					}
				}
			case 0x86:
				switch msg[1] {
				case 0x01:
					mixer.PGM = msg[2] & 0x0F
				case 0x02:
					mixer.PST = msg[2] & 0x0F
				case 0x04:
					if msg[2] == 0x00 {
						mixer.TransType = TransMix
					} else {
						mixer.TransType = TransWipe
					}
				case 0x28:
					// Update lamps on the controller
					mixer.send_sources()
				case 0x33:
					// Cut button
					mixer.cut()
				default:
					fmt.Printf("**** UNKNOWN MESSAGE ****\n")
				}
			}
			if decode.MsgEq(msg, []byte{0x86, 0x63, 0x01, 0x00}) {
				// Initialization request

				for _, buf := range static.T1_init_data() {
					mixer.TxChan <- buf
				}
				msg := []byte{0xE1, 0x00, 0xFF, 0x7F}
				mixer.TxChan <- msg
			}
		}
	}
}

func (mixer *mixer) Print() {
	fmt.Printf("Mixer state:\n")
	fmt.Printf("\tPGM: %d\n", mixer.PGM)
	fmt.Printf("\tPST: %d\n", mixer.PST)
	trans_type := "Wipe"
	if mixer.TransType == TransMix {
		trans_type = "Mix"
	}
	fmt.Printf("\tTransition type: %s\n", trans_type)
}

func (mixer *mixer) cut() {
	p := mixer.PGM
	mixer.PGM = mixer.PST
	mixer.PST = p
	mixer.send_sources()
}

func (mixer *mixer) send_sources() {
	mixer.TxChan <- []byte{0x86, 0x64, byte(mixer.PGM), 0x00}
	mixer.TxChan <- []byte{0x86, 0x65, byte(mixer.PST), 0x00}
}

func analog(msg []byte) int {
	value := int(msg[2]) + ((int(msg[3]) & 0x7F) * 256)
	if msg[3]&0x80 > 0 {
		value = -32767 + value
	}
	return value
}

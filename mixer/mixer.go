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
	PGM       byte // Active PGM selection
	PST       byte // Active PST selection
	rxChan    chan []byte
	TxChan    chan []byte
	TransType int
}

func Init(serialConfig serial.Config) *mixer {
	var mixer = mixer{
		PGM:    0,
		PST:    0,
		rxChan: make(chan []byte),
		TxChan: make(chan []byte),
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
					mixer.TxChan <- []byte{0x86, 0x64, byte(mixer.PGM), 0x00}
					mixer.TxChan <- []byte{0x86, 0x65, byte(mixer.PST), 0x00}
				case 0x33:
					// Cut button
					p := mixer.PGM
					mixer.PGM = mixer.PST
					mixer.PST = p

					mixer.TxChan <- []byte{0x86, 0x64, byte(mixer.PGM), 0x00}
					mixer.TxChan <- []byte{0x86, 0x65, byte(mixer.PST), 0x00}
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

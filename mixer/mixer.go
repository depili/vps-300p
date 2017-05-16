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

const (
	LayerBKGD = 0x01
	LayerKey1 = 0x04
	LayerKey2 = 0x02
)

type mixer struct {
	rxChan chan []byte
	TxChan chan []byte
	State  *MixerState
}

func Init(serialConfig serial.Config) *mixer {
	var mixer = mixer{
		rxChan: make(chan []byte),
		TxChan: make(chan []byte),
	}

	// TODO: Initialize the controller for the same values!
	mixer.State = &MixerState{
		PGM:    0,
		PST:    0,
		Type:   TransWipe,
		Dir:    true,
		Value:  0,
		Layers: 0,
	}

	keyState := &KeyerState{
		MattSat:      0,
		MattLum:      0,
		MattHue:      0,
		Edge:         0,
		EdgeSat:      0,
		EdgeLum:      0,
		EdgeHue:      0,
		Transparency: 0,
		ShadowX:      0,
		ShadowY:      0,
	}

	mixer.State.Key1State = keyState.copy()
	mixer.State.Key2State = keyState.copy()

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
				case 0x1E:
					// Key1 matte hue
					mixer.State.Key1State.MattHue = analog(msg)
				case 0x1F:
					// Key1 matte saturation
					mixer.State.Key1State.MattSat = analog(msg)
				case 0x20:
					// Key1 matte luminance
					mixer.State.Key1State.MattLum = analog(msg)
				case 0x21:
					// Key2 matte hue
					mixer.State.Key2State.MattHue = analog(msg)
				case 0x22:
					// Key2 matte saturation
					mixer.State.Key2State.MattSat = analog(msg)
				case 0x23:
					// Key2 matte luminance
					mixer.State.Key2State.MattLum = analog(msg)
				case 0x4D:
					// T-bar
					value := analog(msg)
					if mixer.State.Dir {
						if value == 1023 {
							// Completed transition
							s := mixer.State.transComplete()
							mixer.TxChan <- []byte{0x86, 0x6A, 0x01, 0x00} // Upwards arrow
							mixer.State = &s
						} else {
							mixer.State.Value = value
						}
					} else {
						if value == 0 {
							s := mixer.State.transComplete()
							mixer.TxChan <- []byte{0x86, 0x6A, 0x00, 0x00} // Downwards arrow
							mixer.State = &s
						} else {
							mixer.State.Value = 1023 - value
						}
					}
					mixer.send_sources()
				}
			case 0x86:
				switch msg[1] {
				case 0x01:
					mixer.State.PGM = msg[2] & 0x0F
				case 0x02:
					mixer.State.PST = msg[2] & 0x0F
				case 0x03:
					// Pattern number
					mixer.State.Pattern = analog(msg)
				case 0x04:
					// Transition type mix/wipe
					if msg[2] == 0x00 {
						mixer.State.Type = TransMix
					} else {
						mixer.State.Type = TransWipe
					}
				case 0x05:
					// Pattern reverse toggle
					if msg[2] == 0x00 {
						mixer.State.PatternRev = false
					} else {
						mixer.State.PatternRev = true
					}
				case 0x09:
					// Select transition layers
					mixer.State.Layers = msg[2]
				case 0x28:
					// Update lamps on the controller
					mixer.send_sources()
				case 0x33:
					// Cut button
					mixer.cut()
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
	fmt.Printf("\tPGM: %d\n", mixer.State.PGM)
	fmt.Printf("\tPST: %d\n", mixer.State.PST)
	trans_type := "Wipe"
	if mixer.State.Type == TransMix {
		trans_type = "Mix"
	}
	fmt.Printf("\tTransition type: %s\n", trans_type)
	fmt.Printf("\tTransition value: %d\n", mixer.State.Value)
}

func (mixer *mixer) cut() {
	s := mixer.State.cut()
	mixer.State = &s
	mixer.send_sources()
}

func (mixer *mixer) send_sources() {
	state := mixer.State
	pst := byte(state.PST)
	if state.Value != 0 && state.Layers&LayerBKGD != 0 {
		pst |= 0xC0
	}

	mixer.TxChan <- []byte{0x86, 0x64, byte(state.PGM), 0x00}
	mixer.TxChan <- []byte{0x86, 0x65, pst, 0x00}

	if state.Key1 {
		mixer.TxChan <- []byte{0x86, 0x66, 0x01, 0x00}
	} else {
		mixer.TxChan <- []byte{0x86, 0x66, 0x00, 0x00}
	}

	if state.Key2 {
		mixer.TxChan <- []byte{0x86, 0x67, 0x01, 0x00}
	} else {
		mixer.TxChan <- []byte{0x86, 0x67, 0x00, 0x00}
	}
}

func analog(msg []byte) int {
	value := int(msg[2]) + ((int(msg[3]) & 0x7F) * 256)
	if msg[3]&0x80 > 0 {
		value = -32767 + value
	}
	return value
}

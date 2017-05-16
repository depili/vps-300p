package main

import (
	"fmt"
	// "github.com/depili/go-rgb-led-matrix/bdf"
	"github.com/depili/vps-300p/mixer"
	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
	"image"
	"os"
	"os/signal"
	"time"
)

const matrix_height = 64
const matrix_width = 128

var Options struct {
	Serial   string `short:"s" long:"serial" description:"Arduino bridge serial" default:"/dev/ttyUSB0"`
	Matrix   string `short:"m" long:"matrix" description:"Matrix to connect to" default:"127.0.0.1:9999"`
	TickRate uint   `short:"t" long:"tick-rate" description:"Delay between frame updates (ms)" default:"10"`
}

var parser = flags.NewParser(&Options, flags.Default)
var stateChan chan *mixer.MixerState

func main() {
	if _, err := parser.Parse(); err != nil {
		panic(err)
	}

	serialConfig := serial.Config{
		Name:        Options.Serial,
		Baud:        38400,
		Parity:      'N',
		Size:        8,
		StopBits:    1,
		ReadTimeout: time.Millisecond,
	}

	stateChan = make(chan *mixer.MixerState)

	mixer := mixer.Init(serialConfig)
	fmt.Printf("Mixer connection open.\n")
	mixer.Print()

	// Trap SIGINT aka Ctrl-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	pgmChan := make(chan *image.NRGBA, 1)
	pstChan := make(chan *image.NRGBA, 1)
	key1Chan := make(chan [][]bool, 1)

	// Channels for requestin updates on various effects
	var effChans [11]chan bool
	for i := range effChans {
		effChans[i] = make(chan bool)
	}

	// Start effect workers
	go runBlack(effChans[0], pgmChan, pstChan)
	go runPlasma(effChans[1], pgmChan, pstChan)
	go runFlame(effChans[2], pgmChan, pstChan)
	go runBlack(effChans[3], pgmChan, pstChan)
	go runBlack(effChans[4], pgmChan, pstChan)
	go runBlack(effChans[5], pgmChan, pstChan)
	go runBlack(effChans[6], pgmChan, pstChan)
	go runBlack(effChans[7], pgmChan, pstChan)
	go runBlack(effChans[8], pgmChan, pstChan)
	go runBlack(effChans[9], pgmChan, pstChan)

	go runText(effChans[10], key1Chan, key1Chan)

	go imageWorker(pgmChan, pstChan, key1Chan)

	ticker := time.NewTicker(time.Millisecond * time.Duration(Options.TickRate))

	// Send the PGM framerate
	fps := 1000.0 / float64(Options.TickRate)
	mixer.TxChan <- []byte{0x86, 0x37, byte(fps * 2), 0x00}

	for {
		select {
		case <-ticker.C:
			state := mixer.State
			effChans[state.PGM] <- true
			effChans[state.PST] <- false
			effChans[10] <- true
			stateChan <- state
		case <-sigChan:
			// SIGINT received, shutdown gracefully
			os.Exit(1)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

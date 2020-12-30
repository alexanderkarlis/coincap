package main

import (
	"log"
	"os"

	ui "github.com/gizak/termui/v3"

	"github.com/alexanderkarlis/coincap/streamer"
)

func main() {
	done := make(chan struct{})
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	coinStreamer := streamer.NewStreamer(os.Args[1:])
	go coinStreamer.GetCoinData(done)
	<-done
}

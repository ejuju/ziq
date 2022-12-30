package main

import (
	"time"

	"github.com/ejuju/ziq/pkg/audio"
	"github.com/ejuju/ziq/pkg/wave"
)

func main() {
	// Create a sine wave that oscillates at 440 hertz.
	osc := wave.OscillateSine(wave.Const(440.0))

	// Play wave with ffplay.
	config := audio.FFPlayPlayerConfig{Wave: osc, Duration: time.Second}
	player, err := audio.NewFFPlayPlayer(config)
	if err != nil {
		panic(err)
	}
	err = player.Play()
	if err != nil {
		panic(err)
	}
}

package main

import (
	"time"

	"github.com/ejuju/ziq/pkg/audio"
	"github.com/ejuju/ziq/pkg/wave"
)

func main() {
	totalDuration := 10 * time.Second

	kick := wave.MustImportWav("audio_files/kick2.wav")
	kick = wave.Amplitude(kick, wave.Const(0.5))
	kick = wave.Loop(kick, time.Second/2)

	freq := wave.Lerp(440, 880, totalDuration/3)
	sine1 := wave.Amplitude(wave.OscillateSine(freq), wave.Const(0.5))

	mix := wave.Combine(sine1, kick)
	mix = wave.Loop(mix, totalDuration/2)

	// Play wave with ffplay.
	config := audio.FFPlayPlayerConfig{Wave: mix, Duration: totalDuration}
	player, err := audio.NewFFPlayPlayer(config)
	if err != nil {
		panic(err)
	}
	err = player.Play()
	if err != nil {
		panic(err)
	}
}

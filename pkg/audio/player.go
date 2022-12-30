package audio

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ejuju/ziq/pkg/wave"
)

type FFPlayPlayerConfig struct {
	Wave       wave.Wave
	SampleRate int
	Duration   time.Duration
}

// FFplayPlayer uses ffplay to play the provided frames.
// It produces a .pcm file under the hood to encode the output sound wave.
// This file can be saved by setting the SaveFile field to true.
type FFPlayPlayer struct {
	config FFPlayPlayerConfig
}

func NewFFPlayPlayer(config FFPlayPlayerConfig) (*FFPlayPlayer, error) {
	_, err := exec.LookPath("ffplay")
	if err != nil {
		return nil, fmt.Errorf("ffplay executable lookup: %w", err)
	}
	if config.Wave == nil {
		return nil, errors.New("no wave was provided")
	}
	if config.Duration <= 0 {
		return nil, fmt.Errorf("invalid duration: %s", config.Duration)
	}
	if config.SampleRate <= 0 {
		config.SampleRate = 44100
	}

	return &FFPlayPlayer{config: config}, nil
}

func (p FFPlayPlayer) Play() error {
	// get output frames
	frames := Frames(p.config.Wave, p.config.SampleRate, 0, p.config.Duration)

	// Create tmp file
	f, err := os.CreateTemp(os.TempDir(), "audio_*.pcm")
	if err != nil {
		panic(err)
	}

	// Encode PCM output to file
	err = WritePCM(f, frames)
	if err != nil {
		return fmt.Errorf("encode PCM pulses: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// Read output file with ffplay (by launching ffplay from the CLI)
	cmdstr := strings.Split(newFFPlayCommand(p.config.SampleRate, f.Name()), " ")
	_, err = exec.Command(cmdstr[0], cmdstr[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("play PCM file using ffplay: %w", err)
	}
	return nil
}

// newFFPlayCommand returns the command string used to play a PCM file with ffplay.
func newFFPlayCommand(sampleRate int, filepath string) string {
	return "ffplay" + " " +
		"-f f64le" + " " +
		"-ar " + strconv.Itoa(sampleRate) + " " +
		"-autoexit" + " " +
		"-showmode 1" + " " +
		filepath
}

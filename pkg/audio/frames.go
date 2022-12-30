package audio

import (
	"time"

	"github.com/ejuju/ziq/pkg/wave"
)

func Frames(src wave.Wave, framesPerSec int, start, end time.Duration) []float64 {
	frames := []float64{}
	step := float64(time.Second) / float64(framesPerSec) // step == time per frame
	for i := float64(start); i < float64(start+end); i += step {
		val := src(time.Duration(i))
		frames = append(frames, val)
	}
	return frames
}

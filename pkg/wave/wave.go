package wave

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/go-audio/wav"
)

// Wave represents a signal used to produce sound.
//
// It can be an oscillation. Or can also be a value provided to another signal.
//
// For example: One wave could provide the value for a frequency that evolves
// over time (frequency modulation wave).
// You can then produce a sine oscillation (= another wave) using the previously
// created frequency modulation wave.
type Wave func(time.Duration) float64

// A utility wave that always produces the same value.
func Const(value float64) Wave { return func(d time.Duration) float64 { return value } }

// Sinusoidal oscillation wave.
// Can be used to produce notes.
func OscillateSine(frequency Wave) Wave {
	return func(x time.Duration) float64 {
		return math.Sin(math.Pi * 2 * x.Seconds() * frequency(x))
	}
}

// Controls the amplitude of the source wave using another wave.
func Amplitude(src, multiplyBy Wave) Wave {
	return func(x time.Duration) float64 { return src(x) * multiplyBy(x) }
}

// Repeats the source wave from the beginning when period duration has elapsed.
func Loop(src Wave, period time.Duration) Wave {
	return func(x time.Duration) float64 { return src(x % period) }
}

// Returns the wave value until the duration has elapsed.
// Then it returns the provided value.
func Limit(before, after Wave, d time.Duration) Wave {
	return func(x time.Duration) float64 {
		if x >= d {
			return after(x)
		}
		return before(x)
	}
}

// Combines several waves together (additive synthesis)
func Combine(waves ...Wave) Wave {
	return func(x time.Duration) float64 {
		sum := 0.0
		for _, w := range waves {
			sum += w(x)
		}
		return sum / float64(len(waves))
	}
}

// Linear interpollation wave
func Lerp(start, end float64, d time.Duration) Wave {
	return func(x time.Duration) float64 {
		diffY := end - start
		elapsedFromTotal := float64(x) / float64(d)
		return diffY*elapsedFromTotal + start
	}
}

// Shift the source wave by a certain duration.
func Shift(src Wave, by time.Duration) Wave {
	return func(x time.Duration) float64 { return src(x + by) }
}

// Speed returns a wave with a modified time speed.
// Accelerates (when by > 1) or slows down (when by < 1) the time.
func Speed(src Wave, by float64) Wave {
	return func(x time.Duration) float64 { return src(time.Duration(float64(x) * by)) }
}

func pcmFramesToWave(sampleRate int, frames []float64) Wave {
	timePerFrame := time.Second / time.Duration(sampleRate)
	audioDuration := time.Duration(len(frames)) * timePerFrame

	return func(x time.Duration) float64 {
		if x > audioDuration {
			return 0
		}
		return frames[x/timePerFrame]
	}
}

// Creates a wave from a PCM audio file.
func ImportPCM(filepath string, sampleRate int) (Wave, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open file: %s: %w", filepath, err)
	}
	defer f.Close()

	rawfile, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read file: %s: %w", filepath, err)
	}

	frames := []float64{}
	for i := 0; i < len(rawfile); i += 8 {
		frames = append(frames, math.Float64frombits(binary.LittleEndian.Uint64(rawfile[i:i+8])))
	}

	return pcmFramesToWave(sampleRate, frames), nil
}

func MustImportPCM(filepath string, sampleRate int) Wave {
	out, err := ImportPCM(filepath, sampleRate)
	if err != nil {
		panic(err)
	}
	return out
}

func ImportWav(filepath string) (Wave, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("open file: %s: %w", filepath, err)
	}
	defer f.Close()

	pcmBuffer, err := wav.NewDecoder(f).FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("decode wav to pcm: %w", err)
	}
	numChannels := pcmBuffer.PCMFormat().NumChannels
	if numChannels != 1 {
		return nil, fmt.Errorf("num channels should be 1: %d", numChannels)
	}
	frames := []float64{}
	for _, srcframe := range pcmBuffer.AsFloatBuffer().Data {
		frames = append(frames, srcframe/18_000.0)
	}
	return pcmFramesToWave(pcmBuffer.Format.SampleRate, frames), nil
}

func MustImportWav(filepath string) Wave {
	out, err := ImportWav(filepath)
	if err != nil {
		panic(err)
	}
	return out
}

// type ChainedWave struct {
// 	waves   []Wave
// 	pattern []time.Duration
// }

// func NewChainedWave(waves []Wave, pattern ...time.Duration) (*ChainedWave, error) {
// 	if len(waves) < 2 {
// 		return nil, errors.New("missing waves")
// 	}
// 	if len(waves) == 0 {
// 		return nil, errors.New("missing pattern")
// 	}
// 	return &ChainedWave{waves: waves, pattern: pattern}, nil
// }

// func (cw *ChainedWave) At(x time.Duration) float64 {
// 	maxDuration := time.Duration(0)
// 	x = time.Duration(math.Mod(float64(x), float64(maxDuration)))

// 	countDuration := time.Duration(0.0)
// 	for i, segment := range cw.pattern {
// 		if !(x >= countDuration && x < countDuration+segment) {
// 			countDuration += segment
// 			continue
// 		}
// 		index := i
// 		return cw.waves[index].At(x - countDuration)
// 	}
// 	panic("wave not found")
// }

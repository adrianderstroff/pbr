package main

import (
	"math/rand"
	"time"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
)

// Noise holding the generator with a random seed.
type Noise struct {
	generator *rand.Rand
}

// MakeNoise initialized the seed with the current time.
func MakeNoise() Noise {
	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)

	return Noise{
		generator: generator,
	}
}

// MakeNoiseTexture a noise texture of the specified image dimensions
func (noise *Noise) MakeNoiseTexture(width, height int) (texture.Texture, error) {
	// create random data
	var data = make([]uint8, width*height*4)

	// random number generator
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (x + y*width) * 4
			data[idx+0] = uint8(noise.generator.Intn(256))
			data[idx+1] = uint8(noise.generator.Intn(256))
			data[idx+2] = uint8(noise.generator.Intn(256))
			data[idx+3] = uint8(noise.generator.Intn(256))
		}
	}

	tex, err := texture.MakeFromData(data, width, height, gl.RGBA, gl.RGBA)
	if err != nil {
		return texture.Texture{}, err
	}

	return tex, nil
}

// MakeNoiseSlice returns an array of random numbers of the defined length.
func (noise *Noise) MakeNoiseSlice(len int) []float32 {
	data := make([]float32, len)

	// random number generator
	for i := 0; i < len; i++ {
		data[i] = noise.generator.Float32()
	}

	return data
}

// MakeConstantSlice returns an array filled with the value of the defined length.
func (noise *Noise) MakeConstantSlice(len int, val float32) []float32 {
	data := make([]float32, len)

	// random number generator
	for i := 0; i < len; i++ {
		data[i] = val
	}

	return data
}

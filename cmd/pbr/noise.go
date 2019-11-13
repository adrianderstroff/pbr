package main

import (
	"math/rand"
	"time"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
)

/**
 * Creates a noise texture of the specified image dimensions
 */
func MakeNoiseTexture(width, height int) (texture.Texture, error) {
	// create random data
	var data = make([]uint8, width*height*4)

	// random number generator
	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (x + y*width) * 4
			data[idx+0] = uint8(generator.Intn(256))
			data[idx+1] = uint8(generator.Intn(256))
			data[idx+2] = uint8(generator.Intn(256))
			data[idx+3] = uint8(generator.Intn(256))
		}
	}

	tex, err := texture.MakeFromData(data, width, height, gl.RGBA, gl.RGBA)
	if err != nil {
		return texture.Texture{}, err
	}

	return tex, nil
}

// MakeNoiseSlice returns an array of random numbers of the defined length.
func MakeNoiseSlice(len int) []float32 {
	data := make([]float32, len)

	// random number generator
	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)
	for i := 0; i < len; i++ {
		data[i] = generator.Float32()
	}

	return data
}

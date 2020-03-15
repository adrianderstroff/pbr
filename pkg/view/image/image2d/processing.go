package image2d

import (
	"github.com/adrianderstroff/pbr/pkg/cgm"
)

// ConvertToPowerOfTwo subsamples an image to be quadratic and be a power of two.
func (img *Image2D) ConvertToPowerOfTwo() {
	width := img.GetWidth()
	height := img.GetHeight()
	channels := img.GetChannels()
	bytedepth := img.GetByteDepth()

	// determine appropriate power of two dimensions and take the smaller one
	nwidth := closestPowerOfTwoSmallerThanDimension(width)
	nheight := closestPowerOfTwoSmallerThanDimension(height)
	dim := cgm.Mini(nwidth, nheight)

	// determine the skip for sampling the image
	skipX := float64(width) / float64(dim)
	skipY := float64(height) / float64(dim)

	var data []uint8
	for y := 0; y < dim; y++ {
		ny := int(float64(y) * skipY)
		for x := 0; x < dim; x++ {
			nx := int(float64(x) * skipX)
			// get the start of the pixel
			idx := img.getIdx(nx, ny)

			// iterate over all color channels
			for c := 0; c < channels; c++ {
				// iterate over byte depth
				for d := 0; d < bytedepth; d++ {
					// get the color and depth offset
					off := c*bytedepth + d

					data = append(data, img.data[idx+off])
				}
			}
		}
	}

	// overwrite data
	img.data = data
	img.width = dim
	img.height = dim
}

// IsPowerOfTwo returns true if width and height are both powers or two
func (img *Image2D) IsPowerOfTwo() bool {
	return img.width == closestPowerOfTwoSmallerThanDimension(img.width) &&
		img.height == closestPowerOfTwoSmallerThanDimension(img.height)
}

// IsQuadratic returns true if width equals height.
func (img *Image2D) IsQuadratic() bool {
	return img.width == img.height
}

func closestPowerOfTwoSmallerThanDimension(dim int) int {
	powerOfTwo := 1

	for (powerOfTwo * 2) <= dim {
		powerOfTwo *= 2
	}

	return powerOfTwo
}

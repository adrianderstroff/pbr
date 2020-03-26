package image2d

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/adrianderstroff/pbr/pkg/cgm"
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"
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

// getIdx turns the x and y indices into a 1D index with respect to the channels
// and byte depth.
func (img *Image2D) getIdx(x, y int) int {
	return (x + y*img.width) * img.channels * img.bytedepth
}

// getIdx turns the x and y indices into a 1D index.
func (img *Image2D) getOIdx(x, y int) int {
	return (x + y*img.width) * img.channels
}

// extractData grabs the image data from the image.Image.
// it returns the image data, channels, bytedepth and an error.
func extractData(img image.Image, rect image.Rectangle, fname string,
	channels int) ([]uint8, int, int, error) {

	// exctract data values
	var (
		data      []uint8
		bytedepth int
	)

	if fname == "hdr" {
		bytedepth = 4
		channels = 3

		colormodel := img.ColorModel()
		switch colormodel {
		case hdrcolor.RGBModel:
			rgb := hdr.NewRGB(rect)
			draw.Draw(rgb, rect, img, image.Pt(0, 0), draw.Src)
			data = float32SliceToUint8Slice(rgb.Pix)
			break
		case hdrcolor.XYZModel:
			rgb := hdr.NewXYZ(rect)
			draw.Draw(rgb, rect, img, image.Pt(0, 0), draw.Src)
			data = float32SliceToUint8Slice(rgb.Pix)
			break
		default:
			return data, channels, bytedepth, errors.New("hdr color model is not supported")
		}
	} else {
		bytedepth = 1

		// determine number of channels if not already provided
		if channels == -1 {
			colormodel := img.ColorModel()
			channels = 4
			if colormodel == color.AlphaModel ||
				colormodel == color.Alpha16Model ||
				colormodel == color.GrayModel ||
				colormodel == color.Gray16Model {
				channels = 1
			}
		}

		switch channels {
		case 1:
			gray := image.NewGray(rect)
			draw.Draw(gray, rect, img, image.Pt(0, 0), draw.Src)
			data = gray.Pix
		case 4:
			rgba := image.NewRGBA(rect)
			draw.Draw(rgba, rect, img, image.Pt(0, 0), draw.Src)
			data = rgba.Pix
		}
	}

	return data, channels, bytedepth, nil
}

// checkDimensions checks if width, height and number of channels is in an
// appropriate range.
func checkDimensions(width, height, channels int) error {
	if width < 1 || height < 1 {
		return errors.New("width and height must be bigger than 0")
	}

	if channels < 1 || channels > 4 {
		return errors.New("number of channels must be between 1 and 4")
	}

	return nil
}

// getColorModel returns the name of the respective color model.
func getColorModel(model color.Model) string {
	colorname := "undefined"
	switch model {
	case color.RGBAModel:
		colorname = "RGBA"
		break
	case color.RGBA64Model:
		colorname = "RGBA64"
		break
	case color.NRGBAModel:
		colorname = "NRGBA"
		break
	case color.NRGBA64Model:
		colorname = "NRGBA64"
		break
	case color.AlphaModel:
		colorname = "Alpha"
		break
	case color.Alpha16Model:
		colorname = "Alpha16"
		break
	case color.GrayModel:
		colorname = "Gray"
		break
	case color.Gray16Model:
		colorname = "Gray16"
		break
	}
	return colorname
}

// getChannelsName returns the name of the channel.
func getChannelsName(channels int) string {
	c := "Unknown Channel Number"
	switch channels {
	case 1:
		c = "RED"
		break
	case 2:
		c = "RG"
		break
	case 3:
		c = "RGB"
		break
	case 4:
		c = "RGBA"
		break
	}
	return c
}

// getPixelTypeFromByteDepth returns the appropriate pixel type for the given
// bytedepth. So far online a bytedepth of 1 and 4 is supported.
func getPixelTypeFromByteDepth(bytedepth int) (int, error) {
	switch bytedepth {
	case 1:
		return gl.UNSIGNED_BYTE, nil
	case 4:
		return gl.FLOAT, nil
	}

	return 0, errors.New("bytedepth not supported")
}

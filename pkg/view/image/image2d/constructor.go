package image2d

import (
	"image"
	"os"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
)

// Make constructs a white image of the specified width and height and number of channels.
func Make(width, height, channels int) (Image2D, error) {
	// early return if invalid dimensions had been specified
	err := checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// create image data
	var data []uint8
	length := width * height
	for i := 0; i < length; i++ {
		for c := 0; c < channels; c++ {
			data = append(data, 255)
		}
	}

	return Image2D{
		pixelType: uint32(gl.UNSIGNED_BYTE),
		width:     width,
		height:    height,
		channels:  channels,
		bytedepth: 1,
		data:      data,
	}, nil
}

// MakeFromData constructs an image of the specified width and height and the specified data.
func MakeFromData(width, height, channels int, data []uint8) (Image2D, error) {
	// data is stored as rgba value even if data is one channel only
	bytedepth := len(data) / (width * height * channels)

	// early return if invalid dimensions had been specified
	err := checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// get the right pixeltype
	pixeltype, err := getPixelTypeFromByteDepth(bytedepth)
	if err != nil {
		return Image2D{}, err
	}

	return Image2D{
		pixelType: uint32(pixeltype),
		width:     width,
		height:    height,
		channels:  channels,
		bytedepth: bytedepth,
		data:      data,
	}, nil
}

// MakeFromPathFixedChannels constructs the image data from the specified path.
// the second parameter fixes the number of channels of the output image to the specified number.
// If there is no image at the specified path an error is returned instead.
func MakeFromPathFixedChannels(path string, channels int) (Image2D, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return Image2D{}, err
	}
	defer file.Close()

	// decode image
	img, fname, err := image.Decode(file)
	if err != nil {
		return Image2D{}, err
	}

	// get image dimensions
	rect := img.Bounds()
	size := rect.Size()
	width := size.X
	height := size.Y

	// early return if invalid dimensions had been specified
	err = checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// exctract data values
	data, channels, bytedepth, err := extractData(img, rect, fname, channels)
	if err != nil {
		return Image2D{}, err
	}

	// determine pixel type
	pixeltype, err := getPixelTypeFromByteDepth(bytedepth)
	if err != nil {
		return Image2D{}, err
	}

	return Image2D{
		pixelType: uint32(pixeltype),
		width:     width,
		height:    height,
		channels:  channels,
		bytedepth: bytedepth,
		data:      data,
	}, nil
}

// MakeFromPath constructs the image data from the specified path.
// If there is no image at the specified path an error is returned instead.
func MakeFromPath(path string) (Image2D, error) {
	// load image file
	file, err := os.Open(path)
	if err != nil {
		return Image2D{}, err
	}
	defer file.Close()

	// decode image
	img, fname, err := image.Decode(file)
	if err != nil {
		return Image2D{}, err
	}

	// get image dimensions
	rect := img.Bounds()
	size := rect.Size()
	width := size.X
	height := size.Y

	// determine number of channels
	data, channels, bytedepth, err := extractData(img, rect, fname, -1)

	// early return if invalid dimensions had been specified
	err = checkDimensions(width, height, channels)
	if err != nil {
		return Image2D{}, err
	}

	// determine pixel type
	pixeltype, err := getPixelTypeFromByteDepth(bytedepth)
	if err != nil {
		return Image2D{}, err
	}

	return Image2D{
		pixelType: uint32(pixeltype),
		width:     width,
		height:    height,
		channels:  channels,
		bytedepth: bytedepth,
		data:      data,
	}, nil
}

// MakeFromFrameBuffer grabs the data from the frame buffer and copies it into
// an image.
func MakeFromFrameBuffer(width, height, channels int) (Image2D, error) {
	data := make([]uint8, width*height*channels)
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	img := Image2D{
		pixelType: uint32(gl.UNSIGNED_BYTE),
		width:     width,
		height:    height,
		channels:  channels,
		bytedepth: 1,
		data:      data,
	}
	img.FlipY()
	return img, nil
}

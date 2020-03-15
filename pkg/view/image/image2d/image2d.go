// Package image2d provides classes for creating and storing 2d images.
package image2d

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // used to decode jpeg
	"image/png"
	"os"
	"unsafe"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"

	// import for side effects
	_ "github.com/mdouchement/hdr/codec/rgbe"
)

// Image2D stores the dimensions, data format and it's pixel data.
// It can be used to manipulate single pixels and is used to
// upload it's data to a texture.
type Image2D struct {
	pixelType uint32
	width     int
	height    int
	channels  int
	bytedepth int
	data      []uint8
}

// SaveToPath saves the image at the specified path in the png format.
// The specified image path has to have the fileextension .png.
// An error is thrown if the path is not valid or any of the specified
// directories don't exist.
func (img *Image2D) SaveToPath(path string) error {
	// create a file at the specified path
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// write data back into the golang image format
	rect := image.Rect(0, 0, img.width, img.height)
	switch img.channels {
	case 1:
		// fill image data
		out := image.NewGray(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 2:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = 0
				out.Pix[idxdst+3] = 255
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 3:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = img.data[idxsrc+2]
				out.Pix[idxdst+3] = 255
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	case 4:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
				out.Pix[idx+1] = img.data[idx+1]
				out.Pix[idx+2] = img.data[idx+2]
				out.Pix[idx+3] = img.data[idx+3]
			}
		}

		// write image into file
		if err := png.Encode(file, out); err != nil {
			return err
		}
	}

	return nil
}

// ToImage converts an Image2D into a Go image.Image.
func (img *Image2D) ToImage() (image.Image, error) {
	// write data back into the golang image format
	rect := image.Rect(0, 0, img.width, img.height)
	switch img.channels {
	case 1:
		// fill image data
		out := image.NewGray(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
			}
		}

		return out, nil
	case 2:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = 0
				out.Pix[idxdst+3] = 255
			}
		}

		return out, nil
	case 3:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				idxdst := (x + y*img.width) * 4
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+1]
				out.Pix[idxdst+2] = img.data[idxsrc+2]
				out.Pix[idxdst+3] = 255
			}
		}

		return out, nil
	case 4:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				out.Pix[idx] = img.data[idx]
				out.Pix[idx+1] = img.data[idx+1]
				out.Pix[idx+2] = img.data[idx+2]
				out.Pix[idx+3] = img.data[idx+3]
			}
		}
		return out, nil
	}

	emptyRect := image.Rectangle{image.Point{0, 0}, image.Point{0, 0}}
	return image.NewRGBA(emptyRect), errors.New("Unsupported number of channels")
}

// FlipX changes the order of the columns by swapping the first column of a row
// with the last column of the same row, the second column of this row with the
// second last column of this row etc.
func (img *Image2D) FlipX() {
	var tempdata []uint8
	for row := 0; row < img.height; row++ {
		for col := img.width - 1; col >= 0; col-- {
			idx := img.getIdx(col, row)
			for c := 0; c < img.channels; c++ {
				for bd := 0; bd < img.bytedepth; bd++ {
					off := c*img.bytedepth + bd
					tempdata = append(tempdata, img.data[idx+off])
				}
			}
		}
	}
	img.data = tempdata
}

// FlipY changes the order of the rows by swapping the first row with the
// last row, the second row with the second last row etc.
func (img *Image2D) FlipY() {
	var tempdata []uint8
	for row := img.height - 1; row >= 0; row-- {
		for col := 0; col < img.width; col++ {
			idx := img.getIdx(col, row)
			for c := 0; c < img.channels; c++ {
				for bd := 0; bd < img.bytedepth; bd++ {
					off := c*img.bytedepth + bd
					tempdata = append(tempdata, img.data[idx+off])
				}
			}
		}
	}
	img.data = tempdata
}

// GetWidth returns the width of the image.
func (img *Image2D) GetWidth() int {
	return img.width
}

// GetHeight returns the height of the image.
func (img *Image2D) GetHeight() int {
	return img.height
}

// GetChannels return the number of the channels of the image.
func (img *Image2D) GetChannels() int {
	return img.channels
}

// GetByteDepth returns the number of bytes a channel consists of.
func (img *Image2D) GetByteDepth() int {
	return img.bytedepth
}

// GetPixelType gets the data type of the pixel data.
func (img *Image2D) GetPixelType() uint32 {
	return img.pixelType
}

// GetDataPointer returns an pointer to the beginning of the image data.
func (img *Image2D) GetDataPointer() unsafe.Pointer {
	return gl.Ptr(img.data)
}

// GetData returns a copy of the image's data
func (img *Image2D) GetData() []uint8 {
	cpy := make([]uint8, len(img.data))
	copy(cpy, img.data)
	return cpy
}

// GetR returns the red value of the pixel at (x,y).
func (img *Image2D) GetR(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx]
}

// GetG returns the green value of the pixel at (x,y).
func (img *Image2D) GetG(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+1]
}

// GetB returns the blue value of the pixel at (x,y).
func (img *Image2D) GetB(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+2]
}

// GetA returns the alpha value of the pixel at (x,y).
func (img *Image2D) GetA(x, y int) uint8 {
	idx := img.getIdx(x, y)
	return img.data[idx+3]
}

// GetRGB returns the RGB values of the pixel at (x,y).
func (img *Image2D) GetRGB(x, y int) (uint8, uint8, uint8) {
	idx := img.getIdx(x, y)
	return img.data[idx],
		img.data[idx+1],
		img.data[idx+2]
}

// GetRGBA returns the RGBA value of the pixel at (x,y).
func (img *Image2D) GetRGBA(x, y int) (uint8, uint8, uint8, uint8) {
	idx := img.getIdx(x, y)
	return img.data[idx],
		img.data[idx+1],
		img.data[idx+2],
		img.data[idx+3]
}

// SetR sets the red value of the pixel at (x,y).
func (img *Image2D) SetR(x, y int, r uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
}

// SetG sets the green value of the pixel at (x,y).
func (img *Image2D) SetG(x, y int, g uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+1] = g
}

// SetB sets the blue value of the pixel at (x,y).
func (img *Image2D) SetB(x, y int, b uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+2] = b
}

// SetA sets the alpha value of the pixel at (x,y).
func (img *Image2D) SetA(x, y int, a uint8) {
	idx := img.getIdx(x, y)
	img.data[idx+3] = a
}

// SetRGB sets the RGB values of the pixel at (x,y).
func (img *Image2D) SetRGB(x, y int, r, g, b uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
	img.data[idx+1] = g
	img.data[idx+2] = b
}

// SetRGBA sets the RGBA values of the pixel at (x,y).
func (img *Image2D) SetRGBA(x, y int, r, g, b, a uint8) {
	idx := img.getIdx(x, y)
	img.data[idx] = r
	img.data[idx+1] = g
	img.data[idx+2] = b
	img.data[idx+3] = a
}

func (img Image2D) String() string {
	c := getChannelsName(img.channels)
	d := img.bytedepth * 8
	return fmt.Sprintf("Image2D (%v,%v) %v %vbit", img.width, img.height, c, d)
}

// getIdx turns the x and y indices into a 1D index.
func (img *Image2D) getIdx(x, y int) int {
	return (x + y*img.width) * img.channels * img.bytedepth
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

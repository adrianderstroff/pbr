// Package image2d provides classes for creating and storing 2d images.
package image2d

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // used to decode jpeg
	"image/png"
	"os"
	"path/filepath"
	"unsafe"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"

	// import for side effects

	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/codec/rgbe"
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

	// grab file extension
	extension := filepath.Ext(path)
	ishdr := extension == ".hdr"

	// extract image.Image from the Image2D
	// write data back into the golang image format
	rect := image.Rect(0, 0, img.width, img.height)

	if ishdr {
		out := hdr.NewRGB(rect)

		// make sure that byte depth is 4 bytes
		if img.bytedepth != 4 {
			return errors.New("hdr image has to have a byte depth of 4")
		}

		// make sure the image has 3 channels
		if img.channels != 3 {
			return errors.New("hdr image has to have 3 color channels")
		}

		// fill hdr image data
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				oidx := img.getOIdx(x, y)
				idx := img.getIdx(x, y)

				for c := 0; c < img.channels; c++ {
					// turn bytes to float32
					cidx := idx + c*img.bytedepth
					bytes := make([]byte, img.bytedepth)
					for d := 0; d < img.bytedepth; d++ {
						bytes[d] = img.data[cidx+d]
					}
					val := bytesToFloat32(bytes)

					out.Pix[oidx+c] = val
				}
			}
		}

		return rgbe.Encode(file, out)
	}

	switch img.channels {
	case 1:
		// fill image data
		out := image.NewGray(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				oidx := img.getOIdx(x, y)
				out.Pix[oidx] = img.data[idx]
			}
		}
		return png.Encode(file, out)
	case 2:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				offsrc := img.bytedepth
				idxdst := img.getOIdx(x, y)
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+offsrc]
				out.Pix[idxdst+2] = 0
				out.Pix[idxdst+3] = 255
			}
		}
		return png.Encode(file, out)
	case 3:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idxsrc := img.getIdx(x, y)
				offsrc := img.bytedepth
				idxdst := img.getOIdx(x, y)
				out.Pix[idxdst] = img.data[idxsrc]
				out.Pix[idxdst+1] = img.data[idxsrc+offsrc]
				out.Pix[idxdst+2] = img.data[idxsrc+2*offsrc]
				out.Pix[idxdst+3] = 255
			}
		}
		return png.Encode(file, out)
	case 4:
		// fill image data
		out := image.NewRGBA(rect)
		for y := 0; y < img.height; y++ {
			for x := 0; x < img.width; x++ {
				idx := img.getIdx(x, y)
				off := img.bytedepth
				oidx := img.getOIdx(x, y)
				out.Pix[oidx] = img.data[idx]
				out.Pix[oidx+1] = img.data[idx+off]
				out.Pix[oidx+2] = img.data[idx+2*off]
				out.Pix[oidx+3] = img.data[idx+3*off]
			}
		}
		return png.Encode(file, out)
	}

	return errors.New("invalid number of channels")
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

func (img Image2D) String() string {
	c := getChannelsName(img.channels)
	d := img.bytedepth * 8
	return fmt.Sprintf("Image2D (%v,%v) %v %vbit", img.width, img.height, c, d)
}

// Package texture provides classes for creating and storing images and textures.
package texture

import (
	"unsafe"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/image/image2d"
)

// MakeEmpty creates a Texture with no image data.
func MakeEmpty() Texture {
	return Texture{0, gl.TEXTURE_2D, 0}
}

// Make creates a texture the given width and height.
// Internalformat, format and pixelType specifed the layout of the data.
// Internalformat is the format of the texture on the GPU.
// Format is the format of the pixeldata that provided to this function.
// Pixeltype specifies the data type of a single component of the pixeldata.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func Make(width, height int, internalformat int32, format, pixelType uint32, data unsafe.Pointer, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t)

	// specify a texture image
	gl.TexImage2D(gl.TEXTURE_2D, 0, internalformat, int32(width), int32(height), 0, format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColor creates a color texture of the specified size.
func MakeColor(width, height int) Texture {
	return Make(width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepth creates a depth texture of the specfied size.
func MakeDepth(width, height int) Texture {
	tex := Make(width, height, gl.DEPTH_COMPONENT, gl.DEPTH_COMPONENT, gl.UNSIGNED_BYTE, nil,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
	return tex
}

// MakeFromPathFixedChannels creates a texture with the image data specifed in path.
// The number is enforced no matter how many channels the image in the specified
// file actually has.
func MakeFromPathFixedChannels(path string, channels int, internalformat int32, format uint32) (Texture, error) {
	image, err := image2d.MakeFromPathFixedChannels(path, channels)
	if err != nil {
		return Texture{}, err
	}

	image.FlipY()

	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// MakeFromPath creates a texture with the image data specifed in path.
func MakeFromPath(path string, internalformat int32, format uint32) (Texture, error) {
	image, err := image2d.MakeFromPath(path)
	if err != nil {
		return Texture{}, err
	}

	image.FlipY()

	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// MakeFromImage grabs the dimensions and information from the image
func MakeFromImage(image *image2d.Image2D, internalformat int32, format uint32) Texture {
	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
}

// MakeFromData creates a texture
func MakeFromData(data []uint8, width, height int, internalformat int32, format uint32) (Texture, error) {
	image, err := image2d.MakeFromData(width, height, data)
	if err != nil {
		return Texture{}, err
	}

	return Make(image.GetWidth(), image.GetHeight(), internalformat, format,
		image.GetPixelType(), image.GetDataPointer(), gl.NEAREST, gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

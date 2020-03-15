// Package texture provides classes for creating and storing images and textures.
package texture

import (
	"unsafe"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/image/image2d"
	"github.com/adrianderstroff/pbr/pkg/view/image/image3d"
)

// Make3D constructs a 3D texture of the width and height of each image per
// slice and depth describing the number of slices.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded. The data layout
// is slices first then rows and lastly columns.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image. r specified the
// behaviour between the slices.
func Make3D(width, height, depth, internalformat int32, format, pixelType uint32,
	data unsafe.Pointer, min, mag, s, t, r int32) Texture {

	texture := Texture{0, gl.TEXTURE_3D, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_T, t)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_R, r)

	// specify a texture image
	gl.TexImage3D(gl.TEXTURE_3D, 0, internalformat, width, height, depth, 0,
		format, pixelType, data)

	// unbind texture
	texture.Unbind()

	return texture
}

// Make3DFromPath creates a 3D texture with the data of the images specifed by
// the provided paths.
func Make3DFromPath(paths []string, internalformat int32, format uint32) (Texture, error) {
	// load images from the specified paths and accumulate the loaded data
	images := []image2d.Image2D{}
	data := []uint8{}
	for _, path := range paths {
		image, err := image2d.MakeFromPath(path)
		if err != nil {
			return Texture{}, err
		}

		image.FlipY()

		data = append(data, image.GetData()...)
		images = append(images, image)
	}

	image := images[0]
	layers := int32(len(paths))
	return Make3D(int32(image.GetWidth()), int32(image.GetHeight()), layers,
		internalformat, format, image.GetPixelType(), gl.Ptr(data), gl.NEAREST,
		gl.NEAREST, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Make3DFromImage creates a 3D texture with the data of the 3D image.
func Make3DFromImage(image3d *image3d.Image3D, internalformat int32, format uint32) (Texture, error) {
	// load images from the specified paths and accumulate the loaded data
	data := image3d.GetData()

	return Make3D(int32(image3d.GetWidth()), int32(image3d.GetHeight()),
		int32(image3d.GetSlices()), internalformat, format,
		image3d.GetPixelType(), gl.Ptr(data), gl.NEAREST, gl.NEAREST,
		gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Make3DFromData creates a 3D texture with the data of the 3D image.
func Make3DFromData(data []uint8, width, height, slices int, internalformat int32,
	format uint32) (Texture, error) {

	return Make3D(int32(width), int32(height), int32(slices), internalformat,
		format, gl.UNSIGNED_BYTE, gl.Ptr(data), gl.NEAREST, gl.NEAREST,
		gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE), nil
}

// Package texture provides classes for creating and storing images and textures.
package texture

import (
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/image/image2d"
)

// MakeEmptyCubeMap creates an empty cube map of the specified dimension,
// internal format on the GPU, format of the image data and the pixel type of
// the data. The image dimension has to be a power of two.
func MakeEmptyCubeMap(dim int, internalformat int32, format, pixelType uint32) (Texture, error) {

	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// generate textures for all sides
	for i := 0; i < 6; i++ {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)

		gl.TexImage2D(target, 0, internalformat, int32(dim), int32(dim), 0,
			format, pixelType, nil)
	}

	// format texture
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// unset active texture
	tex.Unbind()

	return tex, nil
}

// MakeCubeMap creates a cube map with the images specfied from the path.
// For usage with skyboxes where textures are on the inside of the cube, set the
// inside parameter to true to flip all textures horizontally, otherwise set
// this parameter to false. The internal format describes the format of the
// texture on the gpu. If the images are not quadratic and a power of two,
// they will automatically be subsampled to fit the required size.
func MakeCubeMap(right, left, top, bottom, front, back string, inside bool,
	internalformat int32) (Texture, error) {

	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// load images
	imagePaths := []string{right, left, top, bottom, front, back}
	for i, path := range imagePaths {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)

		// loads an image from the specified path
		img, err := image2d.MakeFromPath(path)
		if err != nil {
			return Texture{}, err
		}

		if !img.IsPowerOfTwo() || !img.IsQuadratic() {
			//return Texture{}, errors.New("image is not power of two or quadratic")
			img.ConvertToPowerOfTwo()
		}

		// if inside (e.g. for skyboxes) flip images horizontally
		if inside {
			img.FlipX()
		}

		format := determineFormat(img.GetChannels())

		gl.TexImage2D(target, 0, internalformat, int32(img.GetWidth()),
			int32(img.GetHeight()), 0, uint32(format), img.GetPixelType(),
			img.GetDataPointer())
	}

	// format texture
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// unset active texture
	tex.Unbind()

	return tex, nil
}

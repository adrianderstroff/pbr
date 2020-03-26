// Package texture provides classes for creating and storing images and textures.
package texture

import (
	"fmt"

	"github.com/adrianderstroff/pbr/pkg/view/image/image2d"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
)

// Texture holds no to several images.
type Texture struct {
	handle uint32
	target uint32
	texPos uint32 // e.g. gl.TEXTURE0
}

// GetHandle returns the OpenGL of this texture.
func (tex *Texture) GetHandle() uint32 {
	return tex.handle
}

// Delete destroys the Texture.
func (tex *Texture) Delete() {
	gl.DeleteTextures(1, &tex.handle)
}

// GenMipmap generates mipmap levels.
// Chooses the two mipmaps that most closely match the size of the pixel being
// textured and uses the GL_LINEAR criterion to produce a texture value.
func (tex *Texture) GenMipmap() {
	tex.Bind(0)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

// SetMinMagFilter sets the filter to determine which behaviour is used for
// level of detail functions.
func (tex *Texture) SetMinMagFilter(min, mag int32) {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, mag)
	tex.Unbind()
}

// SetWrap1D sets the behavior at the 1D texure borders
func (tex *Texture) SetWrap1D(s int32) {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_WRAP_S, s)
	tex.Unbind()
}

// SetWrap2D sets the behavior at the 2D texure borders
func (tex *Texture) SetWrap2D(s, t int32) {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t)
	tex.Unbind()
}

// SetWrap3D sets the behavior at the 3D texure borders
func (tex *Texture) SetWrap3D(s, t, r int32) {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_T, t)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_R, r)
	tex.Unbind()
}

// Bind makes the texure available at the specified position.
func (tex *Texture) Bind(index uint32) {
	tex.texPos = gl.TEXTURE0 + index
	gl.ActiveTexture(tex.texPos)
	gl.BindTexture(tex.target, tex.handle)
}

// Unbind makes the texture unavailable for reading.
func (tex *Texture) Unbind() {
	tex.texPos = 0
	gl.BindTexture(tex.target, 0)
}

// DownloadImage2D texture data from the GPU into an Image2D.
func (tex *Texture) DownloadImage2D(format, pixeltype uint32) (image2d.Image2D, error) {
	// bind texture for using the following functions
	tex.Bind(0)
	defer tex.Unbind()

	// grab texture dimensions
	var (
		width  int32
		height int32
	)
	gl.GetTexLevelParameteriv(tex.target, 0, gl.TEXTURE_WIDTH, &width)
	gl.GetTexLevelParameteriv(tex.target, 0, gl.TEXTURE_HEIGHT, &height)

	// grab sizes from format and pixel type
	bytesize := byteSizeFromPixelType(pixeltype)
	channels := channelsFromFormat(format)

	fmt.Printf("Texture Format (%v,%v) %v channels %vbit\n", width, height,
		channels, bytesize*8)

	// initialize data
	data := make([]uint8, width*height*int32(channels*bytesize))

	// download data into buffer
	gl.GetTexImage(tex.target, 0, format, pixeltype, gl.Ptr(data))

	img, err := image2d.MakeFromData(int(width), int(height), channels, data)
	if err != nil {
		return image2d.Image2D{}, err
	}

	return img, nil
}

// DownloadCubeMapImages extracts texture data from the GPU into 6 Image2D for
// each side of the cube map.
func (tex *Texture) DownloadCubeMapImages(format, pixeltype uint32) ([]image2d.Image2D, error) {
	// bind texture for using the following functions
	tex.Bind(0)
	defer tex.Unbind()

	// grab texture dimensions
	var (
		width  int32
		height int32
	)
	gl.GetTexLevelParameteriv(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.TEXTURE_WIDTH, &width)
	gl.GetTexLevelParameteriv(gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, gl.TEXTURE_HEIGHT, &height)

	// grab sizes from format and pixel type
	bytesize := byteSizeFromPixelType(pixeltype)
	channels := channelsFromFormat(format)

	// download all sides of the cubemap
	cubeMapImages := make([]image2d.Image2D, 6)
	for i := 0; i < 6; i++ {
		// download data of a cubemap side into buffer
		var target uint32 = gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)
		data := make([]uint8, width*height*int32(channels*bytesize))
		gl.GetTexImage(target, 0, format, pixeltype, gl.Ptr(data))

		// create image2d from data
		img, err := image2d.MakeFromData(int(width), int(height), channels, data)
		if err != nil {
			return []image2d.Image2D{}, err
		}

		// add to slice
		cubeMapImages[i] = img
	}

	return cubeMapImages, nil
}

func channelsFromFormat(format uint32) int {
	var channels int = -1
	switch format {
	case gl.RED:
		channels = 1
		break
	case gl.RG:
		channels = 2
		break
	case gl.RGB:
		channels = 3
		break
	case gl.RGBA:
		channels = 4
		break
	}
	return channels
}

func byteSizeFromPixelType(pixeltype uint32) int {
	var bytesize int = 3
	switch pixeltype {
	case gl.BYTE, gl.UNSIGNED_BYTE:
		bytesize = 1
		break
	case gl.SHORT, gl.UNSIGNED_SHORT:
		bytesize = 2
		break
	case gl.INT, gl.UNSIGNED_INT, gl.FLOAT:
		bytesize = 4
		break
	}
	return bytesize
}

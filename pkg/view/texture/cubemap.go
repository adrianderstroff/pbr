// Package texture provides classes for creating and storing images and textures.
package texture

import (
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/image/image2d"
)

// MakeCubeMap creates a cube map with the images specfied from the path.
// For usage with skyboxes where textures are on the inside of the cube, set the
// inside parameter to true to flip all textures horizontally, otherwise set this
// parameter to false.
func MakeCubeMap(right, left, top, bottom, front, back string, inside bool) (Texture, error) {
	tex := Texture{0, gl.TEXTURE_CUBE_MAP, 0}

	// generate cube map texture
	gl.GenTextures(1, &tex.handle)
	tex.Bind(0)

	// load images
	imagePaths := []string{right, left, top, bottom, front, back}
	for i, path := range imagePaths {
		target := gl.TEXTURE_CUBE_MAP_POSITIVE_X + uint32(i)
		image, err := image2d.MakeFromPath(path)
		if err != nil {
			return Texture{}, err
		}
		// if inside (e.g. for skyboxes) flip images horizontally
		if inside {
			image.FlipX()
		}
		gl.TexImage2D(target, 0, gl.RGBA, int32(image.GetWidth()), int32(image.GetHeight()),
			0, gl.RGBA, image.GetPixelType(), image.GetDataPointer())
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

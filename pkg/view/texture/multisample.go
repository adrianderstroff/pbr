// Package texture provides classes for creating and storing images and textures.
package texture

import gl "github.com/adrianderstroff/pbr/pkg/core/gl"

// MakeMultisample creates a multisample texture of the given width and height and the number of samples that should be used.
// Internalformat, format and pixelType specifed the layout of the data.
// Data is pointing to the data that is going to be uploaded.
// Min and mag specify the behaviour when down and upscaling the texture.
// S and t specify the behaviour at the borders of the image.
func MakeMultisample(width, height, samples int, format uint32, min, mag, s, t int32) Texture {
	texture := Texture{0, gl.TEXTURE_2D_MULTISAMPLE, 0}

	// generate and bind texture
	gl.GenTextures(1, &texture.handle)
	texture.Bind(0)

	// set texture properties
	/* gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MIN_FILTER, min)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_MAG_FILTER, mag)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_S, s)
	gl.TexParameteri(gl.TEXTURE_2D_MULTISAMPLE, gl.TEXTURE_WRAP_T, t) */

	// specify a texture image
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, int32(samples), format, int32(width), int32(height), false)

	// unbind texture
	texture.Unbind()

	return texture
}

// MakeColorMultisample creates a multisample color texture of the given width and height and the number of samples that should be used.
func MakeColorMultisample(width, height, samples int) Texture {
	return MakeMultisample(width, height, samples, gl.RGBA,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

// MakeDepthMultisample creates a multisample depth texture of the given width and height and the number of samples that should be used.
func MakeDepthMultisample(width, height, samples int) Texture {
	return MakeMultisample(width, height, samples, gl.DEPTH_COMPONENT,
		gl.LINEAR, gl.LINEAR, gl.CLAMP_TO_BORDER, gl.CLAMP_TO_BORDER)
}

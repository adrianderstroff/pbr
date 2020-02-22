// Package texture provides classes for creating and storing images and textures.
package texture

import (
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
// Chooses the two mipmaps that most closely match the size of the pixel being textured and uses the GL_LINEAR criterion to produce a texture value.
func (tex *Texture) GenMipmap() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

// GenMipmapNearest generates mipmap levels.
// Chooses the mipmap that most closely matches the size of the pixel being textured and uses the GL_LINEAR criterion to produce a texture value.
func (tex *Texture) GenMipmapNearest() {
	tex.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
	gl.GenerateMipmap(tex.target)
	tex.Unbind()
}

// SetMinMagFilter sets the filter to determine which behaviour is used for level of detail functions.
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

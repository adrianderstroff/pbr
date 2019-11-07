// Package fbo is a buffer that can hold different kind of textures.
// It can be bound and rendered to instead of the windows frame buffer.
package fbo

import (
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	tex "github.com/adrianderstroff/pbr/pkg/view/texture"
)

// FBO can hold multiple color textues and up to one depth texture.
// It can be bound as alternative render target contrary to the default frame buffer.
type FBO struct {
	handle        uint32
	isBound       bool
	colorTextures map[uint32]*tex.Texture
	depthTexture  *tex.Texture
	textureType   uint32
}

// GetColorTexture returns the color texture at the position of the specified index.
func (fbo *FBO) GetColorTexture(index uint32) *tex.Texture {
	return fbo.colorTextures[index]
}

// GetDepthTexture returns the attached depth texture.
func (fbo *FBO) GetDepthTexture() *tex.Texture {
	return fbo.depthTexture
}

// MakeEmpty FBO creates an empty FBO that works as the default frame buffer.
func MakeEmpty() FBO {
	fbo := FBO{
		handle:        0,
		isBound:       false,
		colorTextures: map[uint32]*tex.Texture{},
		depthTexture:  nil,
		textureType:   gl.TEXTURE_2D,
	}
	gl.GenFramebuffers(1, &fbo.handle)
	return fbo
}

// Make creates an FBO with one color and depth texture of the specified width and height.
func Make(width, height int) FBO {
	fbo := FBO{
		handle:        0,
		isBound:       false,
		colorTextures: map[uint32]*tex.Texture{},
		depthTexture:  nil,
		textureType:   gl.TEXTURE_2D,
	}
	gl.GenFramebuffers(1, &fbo.handle)
	color := tex.MakeColor(width, height)
	depth := tex.MakeDepth(width, height)
	fbo.AttachColorTexture(&color, 0)
	fbo.AttachDepthTexture(&depth)
	return fbo
}

// MakeEmptyMultisample make an empty multisampled frame buffer.
func MakeEmptyMultisample() FBO {
	fbo := FBO{
		handle:        0,
		isBound:       false,
		colorTextures: map[uint32]*tex.Texture{},
		depthTexture:  nil,
		textureType:   gl.TEXTURE_2D_MULTISAMPLE,
	}
	gl.GenFramebuffers(1, &fbo.handle)
	return fbo
}

// MakeMultisample make a multisampled frame buffer with a color and a depth texture of the specified width, height and samples attached.
func MakeMultisample(width, height, samples int) FBO {
	fbo := FBO{
		handle:        0,
		isBound:       false,
		colorTextures: map[uint32]*tex.Texture{},
		depthTexture:  nil,
		textureType:   gl.TEXTURE_2D_MULTISAMPLE,
	}
	gl.GenFramebuffers(1, &fbo.handle)
	color := tex.MakeColorMultisample(width, height, samples)
	depth := tex.MakeDepthMultisample(width, height, samples)
	fbo.AttachColorTexture(&color, 0)
	fbo.AttachDepthTexture(&depth)
	return fbo
}

// Delete destroys this FBO and all associated color and depth textures.
func (fbo *FBO) Delete() {
	// delete textures
	if fbo.colorTextures != nil {
		for _, colTex := range fbo.colorTextures {
			if colTex != nil {
				colTex.Delete()
			}
		}
	}
	if fbo.depthTexture != nil {
		fbo.depthTexture.Delete()
	}

	// unbind fbo
	if fbo.isBound {
		fbo.Unbind()
	}

	// delete buffer
	gl.DeleteFramebuffers(1, &fbo.handle)
}

// Clear clears the color and depth buffer if this FBO has been bound before.
func (fbo *FBO) Clear() {
	if fbo.isBound {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
}

// Bind sets this FBO as current render target.
func (fbo *FBO) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo.handle)
	fbo.isBound = true
}

// Unbind sets the default frame buffer as current render target.
func (fbo *FBO) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	fbo.isBound = false
}

// AttachColorTexture adds a color texture at the position specified by index.
func (fbo *FBO) AttachColorTexture(texture *tex.Texture, index uint32) {
	fbo.Bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0+index, fbo.textureType, texture.GetHandle(), 0)
	drawBuffers := []uint32{gl.COLOR_ATTACHMENT0 + index}
	gl.DrawBuffers(1, &drawBuffers[0])
	fbo.Unbind()

	// delete texture if it would be overwritten
	if previoustexture, ok := fbo.colorTextures[index]; ok {
		previoustexture.Delete()
	}

	// add handle
	fbo.colorTextures[index] = texture
}

// AttachDepthTexture adds a depth texture to the FBO.
func (fbo *FBO) AttachDepthTexture(texture *tex.Texture) {
	fbo.Bind()
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, fbo.textureType, texture.GetHandle(), 0)
	fbo.Unbind()

	// delete previous texture
	if fbo.depthTexture != nil {
		fbo.depthTexture.Delete()
	}

	// add handle
	fbo.depthTexture = texture
}

// Checks if the framebuffer is complete
func (fbo *FBO) IsComplete() bool {
	fbo.Bind()
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	fbo.Unbind()
	return status == gl.FRAMEBUFFER_COMPLETE
}

// CopyToScreen copies all color and depth textures to the default frame buffer.
func (fbo *FBO) CopyToScreen(index uint32, x, y, width, height int32) {
	fbo.CopyToScreenRegion(index, x, y, width, height, x, y, width, height)
}

// CopyToScreenRegion copies all color and depth textures within a region specified by the position (x1,y1) and width w1 and height h1
// to the default frame buffer in the region (x2,y2) and the width w2 and height h2.
func (fbo *FBO) CopyToScreenRegion(index uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.DrawBuffer(gl.BACK)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyToFBO copies all color and depth textures to another FBO.
func (fbo *FBO) CopyToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyToFBORegion(other, x, y, width, height, x, y, width, height)
}

// CopyToFBORegion copies all color and depth textures within a region specified by the position (x1,y1) and width w1 and height h1
// to another FBO in the region (x2,y2) and the width w2 and height h2.
func (fbo *FBO) CopyToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyColorToFBO copies all color textures to another FBO.
func (fbo *FBO) CopyColorToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyColorToFBORegion(other, x, y, width, height, x, y, width, height)
}

// CopyColorToFBORegion copies all color textures within a region specified by the position (x1,y1) and width w1 and height h1
// to another FBO in the region (x2,y2) and the width w2 and height h2.
func (fbo *FBO) CopyColorToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyColorToFBOSmooth copies all color textures to another FBO using linear interpolation.
func (fbo *FBO) CopyColorToFBOSmooth(other *FBO, x, y, width, height int32) {
	fbo.CopyColorToFBORegionSmooth(other, x, y, width, height, x, y, width, height)
}

// CopyColorToFBORegionSmooth copies all color textures within a region specified by the position (x1,y1) and width w1 and height h1
// to another FBO in the region (x2,y2) and the width w2 and height h2 using linear interpolation.
func (fbo *FBO) CopyColorToFBORegionSmooth(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.LINEAR,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyAttachmentColorToFBO copies a color texture specified by index1 to the color texture of another FBO at index2.
func (fbo *FBO) CopyAttachmentColorToFBO(other *FBO, index1, index2 uint32, x, y, width, height int32) {
	fbo.CopyColorAttachmentToFBORegion(other, index1, index2, x, y, width, height, x, y, width, height)
}

// CopyColorAttachmentToFBORegion copies a texture specfied by index2 within a region specified by the position (x1,y1) and width w1 and height h1
// to the color texture of another FBO at index2 in the region (x2,y2) and the width w2 and height h2 using linear interpolation.
func (fbo *FBO) CopyColorAttachmentToFBORegion(other *FBO, index1, index2 uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index1)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0 + index2)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyAttachmentColorToFBOSmooth copies a color texture specified by index1 to the color texture of another FBO at index2 using linear interpolation.
func (fbo *FBO) CopyAttachmentColorToFBOSmooth(other *FBO, index1, index2 uint32, x, y, width, height int32) {
	fbo.CopyAttachmentColorToFBORegionSmooth(other, index1, index2, x, y, width, height, x, y, width, height)
}

// CopyAttachmentColorToFBORegionSmooth copies a texture specfied by index2 within a region specified by the position (x1,y1) and width w1 and height h1
// to the color texture of another FBO at index2 in the region (x2,y2) and the width w2 and height h2 using linear interpolation using linear interpolation.
func (fbo *FBO) CopyAttachmentColorToFBORegionSmooth(other *FBO, index1, index2 uint32, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.ReadBuffer(gl.COLOR_ATTACHMENT0 + index1)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0 + index2)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.COLOR_BUFFER_BIT,
		gl.LINEAR,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// CopyDepthToFBO copies the depth texture to another FBO.
func (fbo *FBO) CopyDepthToFBO(other *FBO, x, y, width, height int32) {
	fbo.CopyDepthToFBORegion(other, x, y, width, height, x, y, width, height)
}

// CopyDepthToFBORegion copies the depth texture within a region specified by the position (x1,y1) and width w1 and height h1
// to another FBO in the region (x2,y2) and the width w2 and height h2.
func (fbo *FBO) CopyDepthToFBORegion(other *FBO, x1, y1, w1, h1, x2, y2, w2, h2 int32) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo.handle)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, other.handle)
	gl.BlitFramebuffer(
		x1, y1, x1+w1, y1+h1,
		x2, y2, x2+w2, y2+h2,
		gl.DEPTH_BUFFER_BIT,
		gl.NEAREST,
	)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

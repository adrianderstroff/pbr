package main

import (
	"errors"

	"github.com/adrianderstroff/pbr/pkg/buffer/fbo"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// IblPass encapsulates all relevant data for rendering a mesh using physically based rendering.
type IblPass struct {
	texturedshader shader.Shader
	cubemap        texture.Texture
	// dimensions
	width  int
	height int
	// uniform variables
	wireframe bool
	// pbr textures
	albedotexture    texture.Texture
	normaltexture    texture.Texture
	metallictexture  texture.Texture
	roughnesstexture texture.Texture
	aotexture        texture.Texture
	// time
	time float32
	// deferred rendering
	gbuffer  fbo.FBO
	imageidx int32
}

// MakeIblPass creates a pbr pass
func MakeIblPass(width, height int, shaderpath, texturepath string, cubemap *texture.Texture) IblPass {
	// create shaders
	sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	texturedshader, err := shader.Make(shaderpath+"/pbr/test/main.vert", shaderpath+"/pbr/test/ibl.frag")
	if err != nil {
		panic(err)
	}
	texturedshader.AddRenderable(sphere)

	// load pbr material
	albedotexture, err := texture.MakeFromPathFixedChannels(texturepath+"/albedo.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	albedotexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	albedotexture.GenMipmap()
	normaltexture, err := texture.MakeFromPathFixedChannels(texturepath+"/normal.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	normaltexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	normaltexture.GenMipmap()
	metallictexture, err := texture.MakeFromPathFixedChannels(texturepath+"/metallic.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	metallictexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	metallictexture.GenMipmap()
	roughnesstexture, err := texture.MakeFromPathFixedChannels(texturepath+"/roughness.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	roughnesstexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	roughnesstexture.GenMipmap()
	aotexture, err := texture.MakeFromPathFixedChannels(texturepath+"/ao.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	aotexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	aotexture.GenMipmap()

	// update textured shader
	texturedshader.Use()
	texturedshader.UpdateFloat32("uGlobalRoughness", 0.1)
	texturedshader.Release()

	// setup g-buffer
	gbuffer := fbo.MakeEmpty()
	depthtex := texture.MakeDepth(width, height)
	colortex := texture.MakeColor(width, height)
	albedotex := texture.MakeColor(width, height)
	normaltex := texture.MakeColor(width, height)
	metallictex := texture.MakeColor(width, height)
	roughnesstex := texture.MakeColor(width, height)
	aotex := texture.MakeColor(width, height)
	gbuffer.AttachDepthTexture(&depthtex)
	gbuffer.AttachColorTexture(&colortex, 0)
	gbuffer.AttachColorTexture(&albedotex, 1)
	gbuffer.AttachColorTexture(&normaltex, 2)
	gbuffer.AttachColorTexture(&metallictex, 3)
	gbuffer.AttachColorTexture(&roughnesstex, 4)
	gbuffer.AttachColorTexture(&aotex, 5)
	if !gbuffer.IsComplete() {
		panic(errors.New("gbuffer incomplete"))
	}

	return IblPass{
		texturedshader: texturedshader,
		cubemap:        *cubemap,
		// dimensions
		width:  width,
		height: height,
		// debug
		wireframe: false,
		// pbr textures
		albedotexture:    albedotexture,
		normaltexture:    normaltexture,
		metallictexture:  metallictexture,
		roughnesstexture: roughnesstexture,
		aotexture:        aotexture,
		// random
		time: 0,
		// deferred rendering
		gbuffer:  gbuffer,
		imageidx: 0,
	}
}

// SetState updates the state of the pass
func (rmp *IblPass) SetState(state State) {
	rmp.texturedshader.Use()
	rmp.texturedshader.UpdateInt32("uSamples", state.samples)
	rmp.texturedshader.UpdateFloat32("uGlobalRoughness", state.globalroughness)
	rmp.texturedshader.Release()

	rmp.imageidx = state.imageidx

	rmp.wireframe = state.wireframe
}

// Render does the pbr pass
func (rmp *IblPass) Render(camera camera.Camera) {
	rmp.gbuffer.Bind()
	rmp.gbuffer.Clear()
	rmp.RenderObj(camera)
	rmp.gbuffer.Unbind()

	// copy the selected texture from the gbuffer to the framebuffer
	w, h := int32(rmp.width), int32(rmp.height)
	rmp.gbuffer.CopyColorToScreen(uint32(rmp.imageidx), 0, 0, w, h)
	rmp.gbuffer.CopyDepthToScreen(0, 0, w, h)
}

// RenderObj actually renders the object
func (rmp *IblPass) RenderObj(camera camera.Camera) {
	if rmp.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	rmp.cubemap.Bind(0)
	rmp.albedotexture.Bind(1)
	rmp.normaltexture.Bind(2)
	rmp.metallictexture.Bind(3)
	rmp.roughnesstexture.Bind(4)
	rmp.aotexture.Bind(5)

	rmp.texturedshader.Use()
	rmp.texturedshader.UpdateMat4("V", camera.GetView())
	rmp.texturedshader.UpdateMat4("P", camera.GetPerspective())
	rmp.texturedshader.UpdateMat4("M", mgl32.Ident4())
	rmp.texturedshader.UpdateVec3("uCameraPos", camera.GetPos())
	rmp.texturedshader.Render()
	rmp.texturedshader.Release()

	rmp.cubemap.Unbind()
	rmp.albedotexture.Unbind()
	rmp.normaltexture.Unbind()
	rmp.metallictexture.Unbind()
	rmp.roughnesstexture.Unbind()
	rmp.aotexture.Unbind()

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

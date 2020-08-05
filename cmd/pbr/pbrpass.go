package main

import (
	"errors"

	"github.com/adrianderstroff/pbr/pkg/buffer/fbo"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/io/obj"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// PbrPass encapsulates all relevant data for rendering a mesh using physically based rendering.
type PbrPass struct {
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
	gbuffer fbo.FBO
}

// MakePbrPass creates a pbr pass
func MakePbrPass(width, height int, meshpath, shaderpath, texturepath string, cubemap *texture.Texture) PbrPass {
	// create shaders
	//sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	gun, err := obj.Load(meshpath+"gun.obj", false, false)
	if err != nil {
		panic(err)
	}
	texturedshader, err := shader.Make(shaderpath+"/pbr/ibl/main.vert", shaderpath+"/pbr/ibl/main.frag")
	if err != nil {
		panic(err)
	}
	texturedshader.AddRenderable(gun)

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

	return PbrPass{
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
		gbuffer: gbuffer,
	}
}

// SetState updates the state of the pass
func (rmp *PbrPass) SetState(state *State) {
	rmp.texturedshader.Use()
	rmp.texturedshader.UpdateInt32("uSamples", state.samples)
	rmp.texturedshader.UpdateFloat32("uGlobalRoughness", state.roughness)
	rmp.texturedshader.Release()

	rmp.wireframe = state.wireframe
}

// Render does the pbr pass
func (rmp *PbrPass) Render(camera camera.Camera) {
	rmp.gbuffer.Bind()
	rmp.gbuffer.Clear()
	rmp.RenderObj(camera)
	rmp.gbuffer.Unbind()

	// copy gbuffer textures to the screen
	w, h := int32(rmp.width), int32(rmp.height)
	rmp.gbuffer.CopyColorToScreen(0, 0, 0, w, h)
	rmp.gbuffer.CopyDepthToScreen(0, 0, w, h)
	rmp.gbuffer.CopyColorToScreenRegion(1, 0, 0, w, h, 0, 0, w/5, h/5)
	rmp.gbuffer.CopyDepthToScreenRegion(0, 0, w, h, 0, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(2, 0, 0, w, h, w*1/5, 0, w/5, h/5)
	rmp.gbuffer.CopyDepthToScreenRegion(0, 0, w, h, w*1/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(3, 0, 0, w, h, w*2/5, 0, w/5, h/5)
	rmp.gbuffer.CopyDepthToScreenRegion(0, 0, w, h, w*2/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(4, 0, 0, w, h, w*3/5, 0, w/5, h/5)
	rmp.gbuffer.CopyDepthToScreenRegion(0, 0, w, h, w*3/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(5, 0, 0, w, h, w*4/5, 0, w/5, h/5)
	rmp.gbuffer.CopyDepthToScreenRegion(0, 0, w, h, w*4/5, 0, w/5, h/5)
}

// RenderObj actually renders the object
func (rmp *PbrPass) RenderObj(camera camera.Camera) {
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

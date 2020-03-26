package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// PbrPass encapsulates all relevant data for rendering a mesh using physically based rendering.
type PbrPass struct {
	texturedshader shader.Shader
	cubemap        texture.Texture
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
	// random
	noisetexture texture.Texture
}

// MakePbrPass creates a pbr pass
func MakePbrPass(width, height int, shaderpath, texturepath string, cubemap *texture.Texture) PbrPass {
	// create shaders
	sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	texturedshader, err := shader.Make(shaderpath+"/pbr/ibl/main.vert", shaderpath+"/pbr/ibl/main.frag")
	if err != nil {
		panic(err)
	}
	texturedshader.AddRenderable(sphere)

	// load pbr material
	albedotexture, err := texture.MakeFromPathFixedChannels(texturepath+"/albedo.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	albedotexture.GenMipmap()
	normaltexture, err := texture.MakeFromPathFixedChannels(texturepath+"/normal.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	normaltexture.GenMipmap()
	metallictexture, err := texture.MakeFromPathFixedChannels(texturepath+"/metallic.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	metallictexture.GenMipmap()
	roughnesstexture, err := texture.MakeFromPathFixedChannels(texturepath+"/roughness.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	roughnesstexture.GenMipmap()
	aotexture, err := texture.MakeFromPathFixedChannels(texturepath+"/ao.png", 4, gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	//aotexture.GenMipmap()

	// random value array
	noise := MakeNoise()
	r := noise.MakeNoiseSlice(100)
	r1 := noise.MakeNoiseSlice(100)
	r2 := noise.MakeNoiseSlice(100)

	// update textured shader
	texturedshader.Use()
	texturedshader.UpdateFloat32Slice("uRandR", r)
	texturedshader.UpdateFloat32Slice("uRandX", r1)
	texturedshader.UpdateFloat32Slice("uRandY", r2)
	texturedshader.UpdateFloat32("uGlobalRoughness", 0.1)
	texturedshader.Release()

	// random texture
	noisetexture, err := noise.MakeNoiseTexture(2048, 2048)
	if err != nil {
		panic(err)
	}
	noisetexture.SetMinMagFilter(gl.LINEAR, gl.LINEAR)

	return PbrPass{
		texturedshader: texturedshader,
		cubemap:        *cubemap,
		wireframe:      false,
		// pbr textures
		albedotexture:    albedotexture,
		normaltexture:    normaltexture,
		metallictexture:  metallictexture,
		roughnesstexture: roughnesstexture,
		aotexture:        aotexture,
		// random
		time:         0,
		noisetexture: noisetexture,
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
	rmp.noisetexture.Bind(6)

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
	rmp.noisetexture.Unbind()

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

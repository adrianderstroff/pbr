package main

import (
	"math"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// PbrPass encapsulates all relevant data for rendering a mesh using physically based rendering.
type PbrPass struct {
	raymarchshader shader.Shader
	cubemap        texture.Texture
	// uniform variables
	samples         int32
	globalroughness float32
	wireframe       bool
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
	raymarchshader, err := shader.Make(shaderpath+"/pbr/variant4/main.vert", shaderpath+"/pbr/variant4/main.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(sphere)

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
	aotexture.GenMipmap()

	// random value array
	noise := MakeNoise()
	r := noise.MakeNoiseSlice(100)
	r1 := noise.MakeNoiseSlice(100)
	r2 := noise.MakeNoiseSlice(100)
	raymarchshader.Use()
	raymarchshader.UpdateFloat32Slice("uRandR", r)
	raymarchshader.UpdateFloat32Slice("uRandX", r1)
	raymarchshader.UpdateFloat32Slice("uRandY", r2)
	raymarchshader.UpdateFloat32("uGlobalRoughness", 0.1)
	raymarchshader.Release()

	// random texture
	noisetexture, err := noise.MakeNoiseTexture(2048, 2048)
	if err != nil {
		panic(err)
	}
	noisetexture.SetMinMagFilter(gl.LINEAR, gl.LINEAR)

	return PbrPass{
		raymarchshader: raymarchshader,
		cubemap:        *cubemap,
		// uniform variables
		samples:         10,
		globalroughness: 0.1,
		wireframe:       false,
		// pbr textures
		albedotexture:    albedotexture,
		normaltexture:    normaltexture,
		metallictexture:  metallictexture,
		roughnesstexture: roughnesstexture,
		aotexture:        aotexture,
		time:             0,
		noisetexture:     noisetexture,
	}
}

// Render does the pbr pass
func (rmp *PbrPass) Render(camera camera.Camera) {
	if rmp.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	/* r1 := MakeConstantSlice(100, rmp.time)
	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateFloat32Slice("uRandX", r1)
	rmp.raymarchshader.UpdateFloat32Slice("uRandY", r1)
	rmp.raymarchshader.Release()
	rmp.time = cgm.Mod32(rmp.time+0.0001, 1.0) */

	rmp.cubemap.Bind(0)
	rmp.albedotexture.Bind(1)
	rmp.normaltexture.Bind(2)
	rmp.metallictexture.Bind(3)
	rmp.roughnesstexture.Bind(4)
	rmp.aotexture.Bind(5)
	rmp.noisetexture.Bind(6)

	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateMat4("V", camera.GetView())
	rmp.raymarchshader.UpdateMat4("P", camera.GetPerspective())
	rmp.raymarchshader.UpdateMat4("M", mgl32.Ident4())
	rmp.raymarchshader.UpdateVec3("uCameraPos", camera.GetPos())
	rmp.raymarchshader.Render()
	rmp.raymarchshader.Release()

	rmp.cubemap.Unbind()
	rmp.albedotexture.Unbind()
	rmp.normaltexture.Unbind()
	rmp.metallictexture.Unbind()
	rmp.roughnesstexture.Unbind()
	rmp.aotexture.Unbind()
	rmp.noisetexture.Unbind()

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (rmp *PbrPass) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (rmp *PbrPass) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (rmp *PbrPass) OnMouseScroll(x, y float64) bool {
	return false
}

// OnResize is a callback handler that is called every time the window is resized.
func (rmp *PbrPass) OnResize(width, height int) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (rmp *PbrPass) OnKeyPress(key, action, mods int) bool {
	if action == int(glfw.Release) {
		return false
	}

	// update global density
	if key == int(glfw.KeyQ) {
		rmp.samples--
		rmp.samples = int32(math.Max(1, float64(rmp.samples)))
	} else if key == int(glfw.KeyW) {
		rmp.samples++
		rmp.samples = int32(math.Min(100, float64(rmp.samples)))
	} else if key == int(glfw.KeyE) {
		rmp.globalroughness += 0.001
		rmp.globalroughness = float32(math.Min(1.0, float64(rmp.globalroughness)))
	} else if key == int(glfw.KeyR) {
		rmp.globalroughness -= 0.001
		rmp.globalroughness = float32(math.Max(0.0, float64(rmp.globalroughness)))
	} else if key == int(glfw.KeyT) {
		rmp.wireframe = !rmp.wireframe
	}

	// update uniforms
	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateInt32("uSamples", rmp.samples)
	rmp.raymarchshader.UpdateFloat32("uGlobalRoughness", rmp.globalroughness)
	rmp.raymarchshader.Release()

	return false
}

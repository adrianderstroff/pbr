package main

import (
	"math"

	"github.com/adrianderstroff/pbr/pkg/cgm"
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
	texturedshader shader.Shader
	simpleshader   shader.Shader
	cubemap        texture.Texture
	// uniform variables
	samples         int32
	globalroughness float32
	wireframe       bool
	usesimple       bool
	// pbr textures
	albedotexture    texture.Texture
	normaltexture    texture.Texture
	metallictexture  texture.Texture
	roughnesstexture texture.Texture
	aotexture        texture.Texture
	// simple parameters
	metallic       float32
	roughness      float32
	lightintensity float32
	// time
	time float32
	// random
	noisetexture texture.Texture
}

// MakePbrPass creates a pbr pass
func MakePbrPass(width, height int, shaderpath, texturepath string, cubemap *texture.Texture) PbrPass {
	// create shaders
	sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	texturedshader, err := shader.Make(shaderpath+"/pbr/variant4/main.vert", shaderpath+"/pbr/variant4/main.frag")
	if err != nil {
		panic(err)
	}
	simpleshader, err := shader.Make(shaderpath+"/pbr/simple/main.vert", shaderpath+"/pbr/simple/main.frag")
	if err != nil {
		panic(err)
	}
	texturedshader.AddRenderable(sphere)
	simpleshader.AddRenderable(sphere)

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

	// update textured shader
	texturedshader.Use()
	texturedshader.UpdateFloat32Slice("uRandR", r)
	texturedshader.UpdateFloat32Slice("uRandX", r1)
	texturedshader.UpdateFloat32Slice("uRandY", r2)
	texturedshader.UpdateFloat32("uGlobalRoughness", 0.1)
	texturedshader.Release()

	// update simple shader
	var metallic float32 = 0.0
	var roughness float32 = 0.5
	var lightintensity float32 = 10
	simpleshader.Use()
	simpleshader.UpdateFloat32Slice("uRandR", r)
	simpleshader.UpdateFloat32Slice("uRandX", r1)
	simpleshader.UpdateFloat32Slice("uRandY", r2)
	simpleshader.UpdateVec3("uAlbedo", mgl32.Vec3{1.00, 0.71, 0.29})
	simpleshader.UpdateFloat32("uMetallic", metallic)
	simpleshader.UpdateFloat32("uRoughness", roughness)
	simpleshader.UpdateVec3("uLightColor", mgl32.Vec3{lightintensity, lightintensity, lightintensity})
	simpleshader.Release()

	// random texture
	noisetexture, err := noise.MakeNoiseTexture(2048, 2048)
	if err != nil {
		panic(err)
	}
	noisetexture.SetMinMagFilter(gl.LINEAR, gl.LINEAR)

	return PbrPass{
		texturedshader: texturedshader,
		simpleshader:   simpleshader,
		cubemap:        *cubemap,
		// uniform variables
		samples:         10,
		globalroughness: 0.1,
		wireframe:       false,
		usesimple:       false,
		// pbr textures
		albedotexture:    albedotexture,
		normaltexture:    normaltexture,
		metallictexture:  metallictexture,
		roughnesstexture: roughnesstexture,
		aotexture:        aotexture,
		// simple parameters
		metallic:       metallic,
		roughness:      roughness,
		lightintensity: lightintensity,
		// random
		time:         0,
		noisetexture: noisetexture,
	}
}

// Render does the pbr pass
func (rmp *PbrPass) Render(camera camera.Camera) {
	if rmp.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	if rmp.usesimple {
		rmp.RenderSimple(camera)
	} else {
		rmp.RenderTextured(camera)
	}

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

// RenderTextured performs PBR with texture parameters.
func (rmp *PbrPass) RenderTextured(camera camera.Camera) {
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
}

// RenderSimple performs PBR with simple parameters.
func (rmp *PbrPass) RenderSimple(camera camera.Camera) {
	rmp.cubemap.Bind(0)
	rmp.noisetexture.Bind(1)

	rmp.simpleshader.Use()
	rmp.simpleshader.UpdateMat4("V", camera.GetView())
	rmp.simpleshader.UpdateMat4("P", camera.GetPerspective())
	rmp.simpleshader.UpdateMat4("M", mgl32.Ident4())
	rmp.simpleshader.UpdateVec3("uCameraPos", camera.GetPos())
	rmp.simpleshader.Render()
	rmp.simpleshader.Release()

	rmp.cubemap.Unbind()
	rmp.noisetexture.Unbind()
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
	} else if key == int(glfw.KeyS) {
		rmp.usesimple = !rmp.usesimple
	} else if key == int(glfw.KeyK) {
		rmp.roughness = cgm.Clamp(rmp.roughness+0.01, 0, 1)
	} else if key == int(glfw.KeyL) {
		rmp.roughness = cgm.Clamp(rmp.roughness-0.01, 0, 1)
	} else if key == int(glfw.KeyM) {
		rmp.metallic = 1 - rmp.metallic
	} else if key == int(glfw.KeyI) {
		rmp.lightintensity = rmp.lightintensity + 0.1
	} else if key == int(glfw.KeyO) {
		rmp.lightintensity = cgm.Max32(rmp.lightintensity-0.1, 0.0)
	}

	// update uniforms
	rmp.texturedshader.Use()
	rmp.texturedshader.UpdateInt32("uSamples", rmp.samples)
	rmp.texturedshader.UpdateFloat32("uGlobalRoughness", rmp.globalroughness)
	rmp.texturedshader.Release()

	rmp.simpleshader.Use()
	rmp.simpleshader.UpdateFloat32("uMetallic", rmp.metallic)
	rmp.simpleshader.UpdateFloat32("uRoughness", rmp.roughness)
	rmp.simpleshader.UpdateVec3("uLightColor", mgl32.Vec3{rmp.lightintensity,
		rmp.lightintensity, rmp.lightintensity})
	rmp.simpleshader.Release()

	return false
}

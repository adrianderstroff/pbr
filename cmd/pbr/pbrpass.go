package main

import (
	"github.com/adrianderstroff/pbr/pkg/cgm"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/box"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type RaymarchingPass struct {
	raymarchshader shader.Shader
	cubemap        texture.Texture
	// uniform variables
	globaldensity  float32
	globalcoverage float32
	// pbr textures
	albedotexture    texture.Texture
	normaltexture    texture.Texture
	metallictexture  texture.Texture
	roughnesstexture texture.Texture
	aotexture        texture.Texture
	noisetexture     texture.Texture
}

func MakePbrPass(width, height int, shaderpath, texturepath string, cubemap *texture.Texture) RaymarchingPass {
	// create shaders
	box := box.Make(1, 1, 1, false, gl.TRIANGLES)
	raymarchshader, err := shader.Make(shaderpath+"/pbr/main.vert", shaderpath+"/pbr/main.frag")
	if err != nil {
		panic(err)
	}
	raymarchshader.AddRenderable(box)

	// load pbr material
	albedotexture, err := texture.MakeFromPath(texturepath+"/albedo.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	normaltexture, err := texture.MakeFromPath(texturepath+"/normal.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	metallictexture, err := texture.MakeFromPath(texturepath+"/metallic.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	roughnesstexture, err := texture.MakeFromPath(texturepath+"/roughness.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	aotexture, err := texture.MakeFromPath(texturepath+"/ao.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	noisetexture, err := MakeNoiseTexture(width, height)
	if err != nil {
		panic(err)
	}

	return RaymarchingPass{
		raymarchshader: raymarchshader,
		cubemap:        *cubemap,
		// uniform variables
		globaldensity:  0.5,
		globalcoverage: 0.5,
		// pbr textures
		albedotexture:    albedotexture,
		normaltexture:    normaltexture,
		metallictexture:  metallictexture,
		roughnesstexture: roughnesstexture,
		aotexture:        aotexture,
		noisetexture:     noisetexture,
	}
}

func (rmp *RaymarchingPass) Render(camera camera.Camera) {
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
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (rmp *RaymarchingPass) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (rmp *RaymarchingPass) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (rmp *RaymarchingPass) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (rmp *RaymarchingPass) OnKeyPress(key, action, mods int) bool {
	// update global density
	if key == int(glfw.KeyQ) {
		rmp.globaldensity -= 0.01
	} else if key == int(glfw.KeyE) {
		rmp.globaldensity += 0.01
	}
	rmp.globaldensity = cgm.Clamp(rmp.globaldensity, 0, 1)

	// update global coverage
	if key == int(glfw.KeyZ) {
		rmp.globalcoverage -= 0.01
	} else if key == int(glfw.KeyC) {
		rmp.globalcoverage += 0.01
	}
	rmp.globalcoverage = cgm.Clamp(rmp.globalcoverage, 0, 1)

	// update uniforms
	rmp.raymarchshader.Use()
	rmp.raymarchshader.UpdateFloat32("uGlobalDensity", rmp.globaldensity)
	rmp.raymarchshader.UpdateFloat32("uGlobalCoverage", rmp.globalcoverage)
	rmp.raymarchshader.Release()

	return false
}

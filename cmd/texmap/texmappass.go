package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// TexMapPass encapsulates all relevant data for rendering a textured sphere.
type TexMapPass struct {
	texturingshader shader.Shader
	// uniform variables
	wireframe bool
	// pbr textures
	albedotexture texture.Texture
}

// MakeTexMapPass creates a rendering pass
func MakeTexMapPass(width, height int, shaderpath, texturepath string) TexMapPass {
	// create shaders
	sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	texturingshader, err := shader.Make(shaderpath+"/texture/main.vert", shaderpath+"/texture/main.frag")
	if err != nil {
		panic(err)
	}
	texturingshader.AddRenderable(sphere)

	// load texture
	albedotexture, err := texture.MakeFromPath(texturepath+"/albedo.png", gl.RGBA, gl.RGBA)
	if err != nil {
		panic(err)
	}
	albedotexture.GenMipmap()

	return TexMapPass{
		texturingshader: texturingshader,
		wireframe:       false,
		albedotexture:   albedotexture,
	}
}

// Render does the pbr pass
func (tmp *TexMapPass) Render(camera camera.Camera) {
	if tmp.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}

	tmp.albedotexture.Bind(0)

	tmp.texturingshader.Use()
	tmp.texturingshader.UpdateMat4("V", camera.GetView())
	tmp.texturingshader.UpdateMat4("P", camera.GetPerspective())
	tmp.texturingshader.UpdateMat4("M", mgl32.Ident4())
	tmp.texturingshader.UpdateVec3("uCameraPos", camera.GetPos())
	tmp.texturingshader.Render()
	tmp.texturingshader.Release()

	tmp.albedotexture.Unbind()

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (tmp *TexMapPass) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (tmp *TexMapPass) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (tmp *TexMapPass) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (tmp *TexMapPass) OnKeyPress(key, action, mods int) bool {
	if action == int(glfw.Release) {
		return false
	}

	// update global density
	if key == int(glfw.KeyW) {
		tmp.wireframe = !tmp.wireframe
	}

	return false
}

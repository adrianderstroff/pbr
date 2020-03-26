package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/go-gl/mathgl/mgl32"
)

// NormalPass encapsulates all relevant data for rendering the normals of a mesh.
type NormalPass struct {
	// shader
	normalshader shader.Shader
	M            mgl32.Mat4
}

// MakeNormalPass creates a normal pass
func MakeNormalPass(width, height int, shaderpath string) NormalPass {
	// create shaders
	sphere := sphere.Make(5, 5, 0.1, gl.TRIANGLES)
	normalshader, err := shader.Make(shaderpath+"/normal/main.vert", shaderpath+"/normal/main.frag")
	if err != nil {
		panic(err)
	}
	normalshader.AddRenderable(sphere)

	// create render pass
	return NormalPass{
		normalshader: normalshader,
		M:            mgl32.Translate3D(0, 1, 0),
	}
}

// Render does the pbr pass
func (np *NormalPass) Render(camera camera.Camera) {
	np.normalshader.Use()
	np.normalshader.UpdateMat4("V", camera.GetView())
	np.normalshader.UpdateMat4("P", camera.GetPerspective())
	np.normalshader.UpdateMat4("M", np.M)
	np.normalshader.Render()
	np.normalshader.Release()
}

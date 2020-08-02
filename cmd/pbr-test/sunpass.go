package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/go-gl/mathgl/mgl32"
)

// SunPass encapsulates all relevant data for rendering the pos position.
type SunPass struct {
	// shader
	sunshader shader.Shader
	sunpos    mgl32.Vec3
}

// MakeSunPass creates a pbr pass
func MakeSunPass(width, height int, shaderpath string) SunPass {
	// create shaders
	sphere := sphere.Make(5, 5, 0.5, gl.TRIANGLES)
	sunshader, err := shader.Make(shaderpath+"/flat/main.vert", shaderpath+"/flat/main.frag")
	if err != nil {
		panic(err)
	}
	sunshader.AddRenderable(sphere)

	// create render pass
	return SunPass{
		sunshader: sunshader,
	}
}

// SetState updates the color of the material
func (sp *SunPass) SetState(state State) {
	sp.sunpos = state.lightpos
}

// Render does the pbr pass
func (sp *SunPass) Render(camera camera.Camera) {
	sp.sunshader.Use()
	sp.sunshader.UpdateMat4("V", camera.GetView())
	sp.sunshader.UpdateMat4("P", camera.GetPerspective())
	sp.sunshader.UpdateMat4("M", mgl32.Translate3D(sp.sunpos.X(), sp.sunpos.Y(), sp.sunpos.Z()))
	sp.sunshader.UpdateVec3("uCameraPos", camera.GetPos())
	sp.sunshader.Render()
	sp.sunshader.Release()
}

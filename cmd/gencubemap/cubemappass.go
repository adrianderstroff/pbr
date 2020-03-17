package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/cube"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// CubemapPass contains all resources used for this shader pass
type CubemapPass struct {
	cubemapshader shader.Shader
	cubemap       *texture.Texture
}

// MakeCubemapPass creates the cubemap pass with the specified paths
func MakeCubemapPass(shaderpath string, cubemap *texture.Texture) CubemapPass {
	// create shaders
	box := cube.Make(50, 50, 50, true, gl.TRIANGLES)
	cubemapshader, err := shader.Make(shaderpath+"/cubemap/main.vert", shaderpath+"/cubemap/main.frag")
	if err != nil {
		panic(err)
	}
	cubemapshader.AddRenderable(box)

	return CubemapPass{
		cubemapshader: cubemapshader,
		cubemap:       cubemap,
	}
}

// Render executes the draw command
func (cmp *CubemapPass) Render(camera camera.Camera) {
	cmp.cubemapshader.Use()
	cmp.cubemap.Bind(0)
	cmp.cubemapshader.UpdateMat4("M", mgl32.Ident4())
	cmp.cubemapshader.UpdateMat4("V", camera.GetView())
	cmp.cubemapshader.UpdateMat4("P", camera.GetPerspective())
	cmp.cubemapshader.Render()
	cmp.cubemap.Unbind()
	cmp.cubemapshader.Release()
}

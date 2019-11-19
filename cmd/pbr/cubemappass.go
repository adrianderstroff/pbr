package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/box"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

type CubemapPass struct {
	cubemapshader shader.Shader
	cubemap       texture.Texture
}

func MakeCubemapPass(shaderpath, cubemappath string) CubemapPass {
	// create shaders
	box := box.Make(50, 50, 50, true, gl.TRIANGLES)
	cubemapshader, err := shader.Make(shaderpath+"/cubemap/main.vert", shaderpath+"/cubemap/main.frag")
	if err != nil {
		panic(err)
	}
	cubemapshader.AddRenderable(box)

	// create cubemap
	var (
		right  = cubemappath + "right.jpg"
		left   = cubemappath + "left.jpg"
		top    = cubemappath + "top.jpg"
		bottom = cubemappath + "bottom.jpg"
		front  = cubemappath + "front.jpg"
		back   = cubemappath + "back.jpg"
	)
	cubemap, err := texture.MakeCubeMap(right, left, top, bottom, front, back, false)
	if err != nil {
		panic(err)
	}

	return CubemapPass{
		cubemapshader: cubemapshader,
		cubemap:       cubemap,
	}
}

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

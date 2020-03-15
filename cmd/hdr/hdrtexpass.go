package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/cube"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// HDRTexPass encapsulates all relevant data for rendering a hdr texture.
type HDRTexPass struct {
	texturingshader shader.Shader
	albedotexture   texture.Texture
}

// MakeHDRTexPass creates a rendering pass
func MakeHDRTexPass(width, height int, shaderpath, texturepath string) HDRTexPass {
	// create shaders
	cube := cube.Make(1, 1, 1, false, gl.TRIANGLES)
	texturingshader, err := shader.Make(shaderpath+"/texture/main.vert", shaderpath+"/texture/main.frag")
	if err != nil {
		panic(err)
	}
	texturingshader.AddRenderable(cube)

	// load texture
	//albedotexture, err := texture.MakeHDRFromPath(texturepath+"/leadenhall_market_1k.hdr", gl.RGB, gl.RGB)
	albedotexture, err := texture.MakeFromPath(texturepath+"/501-free-hdri-skies-com.hdr", gl.RGB, gl.RGB)
	if err != nil {
		panic(err)
	}
	albedotexture.SetWrap2D(gl.REPEAT, gl.REPEAT)
	albedotexture.GenMipmap()

	return HDRTexPass{
		texturingshader: texturingshader,
		albedotexture:   albedotexture,
	}
}

// Render does the pbr pass
func (tmp *HDRTexPass) Render(camera camera.Camera) {
	tmp.albedotexture.Bind(0)

	tmp.texturingshader.Use()
	tmp.texturingshader.UpdateMat4("V", camera.GetView())
	tmp.texturingshader.UpdateMat4("P", camera.GetPerspective())
	tmp.texturingshader.UpdateMat4("M", mgl32.Ident4())
	tmp.texturingshader.UpdateVec3("uCameraPos", camera.GetPos())
	tmp.texturingshader.Render()
	tmp.texturingshader.Release()

	tmp.albedotexture.Unbind()
}

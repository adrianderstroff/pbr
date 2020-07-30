package main

import (
	"github.com/adrianderstroff/pbr/pkg/buffer/fbo"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/cube"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// GenPass encapsulates all relevant data for turning an equirectangular
// texture into a set of cube map textures.
type GenPass struct {
	shader       shader.Shader
	equirecttex  texture.Texture
	emptyCubemap texture.Texture
	framebuffer  fbo.FBO
	resolution   int
}

// MakeGenPass creates a rendering pass
func MakeGenPass(shaderpath, texturepath string, resolution int) GenPass {
	// create shaders
	cube := cube.Make(1, 1, 1, true, gl.TRIANGLES)
	texturingshader, err := shader.Make(shaderpath+"/gencubemap/main.vert", shaderpath+"/gencubemap/main.frag")
	if err != nil {
		panic(err)
	}
	texturingshader.AddRenderable(cube)

	// load hdr texture
	equirecttex, err := texture.MakeFromPath(texturepath, gl.RGB, gl.RGB)
	if err != nil {
		panic(err)
	}

	// create framebuffer
	depthTexture := texture.MakeDepth(resolution, resolution)
	framebuffer := fbo.MakeEmpty()
	framebuffer.AttachDepthTexture(&depthTexture)
	framebuffer.Bind()

	// load cubemap texture
	emptyCubemap, err := texture.MakeEmptyCubeMap(resolution, gl.RGB32F, gl.RGB, gl.FLOAT)
	if err != nil {
		panic(err)
	}

	return GenPass{
		shader:       texturingshader,
		equirecttex:  equirecttex,
		emptyCubemap: emptyCubemap,
		framebuffer:  framebuffer,
		resolution:   resolution,
	}
}

// GetCubeMap returns the cube map
func (tmp *GenPass) GetCubeMap() *texture.Texture {
	return &tmp.emptyCubemap
}

// Render does the pbr pass
func (tmp *GenPass) Render() {
	P := mgl32.Perspective(mgl32.DegToRad(90.0), 1.0, 0.1, 10.0)
	Vs := make([]mgl32.Mat4, 6)
	Vs[0] = mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0)
	Vs[1] = mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0)
	Vs[2] = mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1)
	Vs[3] = mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1)
	Vs[4] = mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0)
	Vs[5] = mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0)

	tmp.shader.Use()
	tmp.equirecttex.Bind(0)

	gl.Viewport(0, 0, int32(tmp.resolution), int32(tmp.resolution))
	tmp.framebuffer.Bind()

	// render to all sides of the cube map
	for i := 0; i < 6; i++ {
		// update uniforms
		tmp.shader.UpdateMat4("V", Vs[i])
		tmp.shader.UpdateMat4("P", P)

		// select cubemap side
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			tmp.emptyCubemap.GetHandle(), 0)

		tmp.framebuffer.Clear()

		tmp.shader.Render()
	}

	tmp.framebuffer.Unbind()
	tmp.equirecttex.Unbind()
	tmp.shader.Release()

	// restore viewport
	gl.Viewport(0, 0, 800, 600)
}

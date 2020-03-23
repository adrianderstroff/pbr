package main

import (
	"errors"

	"github.com/adrianderstroff/pbr/pkg/buffer/fbo"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/io/obj"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/go-gl/mathgl/mgl32"
)

// PbrPass encapsulates all relevant data for rendering a mesh using physically based rendering.
type PbrPass struct {
	// shader
	pbrshader shader.Shader
	// uniform variables
	samples         int32
	globalroughness float32
	wireframe       bool
	usesimple       bool
	// simple parameters
	albedo    mgl32.Vec3
	metallic  float32
	roughness float32
	// light
	lightpos       mgl32.Vec3
	lightintensity mgl32.Vec3
	// dimensions
	width  int
	height int
	// deferred rendering
	gbuffer fbo.FBO
}

// MakePbrPass creates a pbr pass
func MakePbrPass(width, height int, shaderpath, texturepath, objpath string) PbrPass {
	// create shaders
	//sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	bunny, err := obj.Load(objpath+"dragon.obj", true, true)
	if err != nil {
		panic(err)
	}

	pbrshader, err := shader.Make(shaderpath+"/pbr/simple/main.vert", shaderpath+"/pbr/simple/main.frag")
	if err != nil {
		panic(err)
	}
	pbrshader.AddRenderable(bunny)

	// setup g-buffer
	gbuffer := fbo.MakeEmpty()
	depthtex := texture.MakeDepth(width, height)
	colortex := texture.MakeColor(width, height)
	diffusetex := texture.MakeColor(width, height)
	speculartex := texture.MakeColor(width, height)
	dtex := texture.MakeColor(width, height)
	gtex := texture.MakeColor(width, height)
	ftex := texture.MakeColor(width, height)
	gbuffer.AttachDepthTexture(&depthtex)
	gbuffer.AttachColorTexture(&colortex, 0)
	gbuffer.AttachColorTexture(&diffusetex, 1)
	gbuffer.AttachColorTexture(&speculartex, 2)
	gbuffer.AttachColorTexture(&dtex, 3)
	gbuffer.AttachColorTexture(&gtex, 4)
	gbuffer.AttachColorTexture(&ftex, 5)
	if !gbuffer.IsComplete() {
		panic(errors.New("gbuffer incomplete"))
	}

	// create render pass
	return PbrPass{
		// shaders
		pbrshader: pbrshader,
		// uniform variables
		samples:         10,
		globalroughness: 0.1,
		wireframe:       false,
		usesimple:       false,
		// dimensions
		width:  width,
		height: height,
		// deferred rendering
		gbuffer: gbuffer,
	}
}

// SetState updates the color of the material
func (rmp *PbrPass) SetState(state State) {
	albedo := state.albedo
	roughness := state.roughness
	metalness := state.metalness
	lightpos := state.lightpos
	lightintensity := state.lightintensity

	color3 := mgl32.Vec3{albedo.X(), albedo.Y(), albedo.Z()}
	rmp.pbrshader.Use()
	rmp.pbrshader.UpdateVec3("uAlbedo", color3)
	rmp.pbrshader.UpdateFloat32("uRoughness", roughness)
	rmp.pbrshader.UpdateFloat32("uMetallic", metalness)
	rmp.pbrshader.UpdateVec3("uLightPos", lightpos)
	rmp.pbrshader.UpdateVec3("uLightColor", lightintensity)
	rmp.pbrshader.Release()
}

// Render does the pbr pass
func (rmp *PbrPass) Render(camera camera.Camera) {
	rmp.gbuffer.Bind()
	rmp.gbuffer.Clear()

	rmp.pbrshader.Use()
	rmp.pbrshader.UpdateMat4("V", camera.GetView())
	rmp.pbrshader.UpdateMat4("P", camera.GetPerspective())
	rmp.pbrshader.UpdateMat4("M", mgl32.Ident4())
	rmp.pbrshader.UpdateVec3("uCameraPos", camera.GetPos())
	rmp.pbrshader.Render()
	rmp.pbrshader.Release()

	rmp.gbuffer.Unbind()

	w, h := int32(rmp.width), int32(rmp.height)
	rmp.gbuffer.CopyColorToScreen(0, 0, 0, w, h)
	rmp.gbuffer.CopyDepthToScreen(0, 0, w, h)
	rmp.gbuffer.CopyColorToScreenRegion(1, 0, 0, w, h, 0, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(2, 0, 0, w, h, w*1/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(3, 0, 0, w, h, w*2/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(4, 0, 0, w, h, w*3/5, 0, w/5, h/5)
	rmp.gbuffer.CopyColorToScreenRegion(5, 0, 0, w, h, w*4/5, 0, w/5, h/5)
}

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
	pbrshader    shader.Shader
	normalshader shader.Shader
	// uniform variables
	samples         int32
	globalroughness float32
	wireframe       bool
	rendernormal    bool
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
	// rendered image index
	imageidx int
}

// MakePbrPass creates a pbr pass
func MakePbrPass(width, height int, shaderpath, texturepath, objpath string) PbrPass {
	// create shaders
	//sphere := sphere.Make(20, 20, 1, gl.TRIANGLES)
	sphere, err := obj.Load(objpath+"dragon.obj", true, true)
	if err != nil {
		panic(err)
	}

	// set up pbr shader
	pbrshader, err := shader.Make(shaderpath+"/pbr/test/main.vert", shaderpath+"/pbr/test/main.frag")
	if err != nil {
		panic(err)
	}
	pbrshader.AddRenderable(sphere)

	// set up normal shader
	normalshader, err := shader.Make(shaderpath+"/normal/main.vert", shaderpath+"/normal/main.frag")
	if err != nil {
		panic(err)
	}
	normalshader.AddRenderable(sphere)

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
		pbrshader:    pbrshader,
		normalshader: normalshader,
		// uniform variables
		samples:         10,
		globalroughness: 0.1,
		wireframe:       false,
		rendernormal:    false,
		// dimensions
		width:  width,
		height: height,
		// deferred rendering
		gbuffer: gbuffer,
		// rendered image index
		imageidx: 0,
	}
}

// SetState updates the color of the material
func (pbr *PbrPass) SetState(state State) {
	albedo := state.albedo
	roughness := state.roughness
	metalness := state.metalness
	lightpos := state.lightpos
	lightintensity := state.lightintensity

	pbr.imageidx = int(state.imageidx)
	pbr.rendernormal = state.normal

	color3 := mgl32.Vec3{albedo.X(), albedo.Y(), albedo.Z()}
	pbr.pbrshader.Use()
	pbr.pbrshader.UpdateVec3("uAlbedo", color3)
	pbr.pbrshader.UpdateFloat32("uRoughness", roughness)
	pbr.pbrshader.UpdateFloat32("uMetallic", metalness)
	pbr.pbrshader.UpdateVec3("uLightPos", lightpos)
	pbr.pbrshader.UpdateVec3("uLightColor", lightintensity)
	pbr.pbrshader.Release()
}

// Render does the pbr pass
func (pbr *PbrPass) Render(camera camera.Camera) {
	if pbr.rendernormal {
		pbr.normalshader.Use()
		pbr.normalshader.UpdateMat4("P", camera.GetPerspective())
		pbr.normalshader.UpdateMat4("V", camera.GetView())
		pbr.normalshader.UpdateMat4("M", mgl32.Ident4())
		pbr.normalshader.Render()
		pbr.normalshader.Release()
	} else {
		pbr.gbuffer.Bind()
		pbr.gbuffer.Clear()

		// invoke pbr test shader
		pbr.pbrshader.Use()
		pbr.pbrshader.UpdateMat4("P", camera.GetPerspective())
		pbr.pbrshader.UpdateMat4("V", camera.GetView())
		pbr.pbrshader.UpdateMat4("M", mgl32.Ident4())
		pbr.pbrshader.UpdateVec3("uCameraPos", camera.GetPos())
		pbr.pbrshader.Render()
		pbr.pbrshader.Release()

		pbr.gbuffer.Unbind()

		// copy the selected texture from the gbuffer to the framebuffer
		w, h := int32(pbr.width), int32(pbr.height)
		pbr.gbuffer.CopyColorToScreen(uint32(pbr.imageidx), 0, 0, w, h)
		pbr.gbuffer.CopyDepthToScreen(0, 0, w, h)
	}
}

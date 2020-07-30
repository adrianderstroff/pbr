package main

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/shader"
	"github.com/adrianderstroff/pbr/pkg/scene/camera"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/cylinder"
	"github.com/adrianderstroff/pbr/pkg/view/mesh/sphere"
	"github.com/go-gl/mathgl/mgl32"
)

// RenderPass for rendering the cosine distribution
type RenderPass struct {
	flatshader shader.Shader
}

// MakeRenderPass creates a render pass for rendering the cosine distribution.
func MakeRenderPass(shaderpath string, raycount int) RenderPass {
	// create shader
	flatshader, err := shader.Make(shaderpath+"/normal/main.vert", shaderpath+"/normal/main.frag")
	if err != nil {
		panic(err)
	}

	// create geometry
	sphere := sphere.Make(20, 25, 1, gl.TRIANGLES)
	flatshader.AddRenderable(sphere)

	// make cube aabb
	aabb := AABB{mgl32.Vec3{-25, -25, -25}, mgl32.Vec3{25, 25, 25}}

	// create random numbers
	noise := MakeNoise()
	rand1 := noise.MakeNoiseSlice(100)
	rand2 := noise.MakeNoiseSlice(100)

	// intersect sphere
	pos := mgl32.Vec3{0, 0, 2}
	dir := pos.Mul(-1).Normalize()
	ray := Ray{pos, dir}
	rsphere := Sphere{mgl32.Vec3{0, 0, 0}, 1}
	hitinfo, didhit := intersectSphere(&ray, &rsphere, 0, 1000)

	// send cosine distributed rays and intersect them with the cube map
	if didhit {
		for i := 0; i < raycount; i++ {
			// next random numbers
			r1 := rand1[i]
			r2 := rand2[i]

			//fmt.Println(fmt.Sprintf("%f, %f", r1, r2))

			// reflect ray
			refl := cosineDistribution(&hitinfo, r1, r2, 1)
			reflray := Ray{hitinfo.p, refl}

			end := hitinfo.p.Add(refl.Mul(1))

			// calculate intersection with bounding box
			pbox := intersectAABB(&reflray, &aabb)

			cylinder1 := cylinder.Make(pos, hitinfo.p, 0.12, gl.TRIANGLES)
			cylinder2 := cylinder.Make(hitinfo.p, end, 0.10, gl.TRIANGLES)
			cylinder3 := cylinder.Make(hitinfo.p, pbox, 0.05, gl.TRIANGLES)
			flatshader.AddRenderable(cylinder1)
			flatshader.AddRenderable(cylinder2)
			flatshader.AddRenderable(cylinder3)
		}
	}

	return RenderPass{
		flatshader: flatshader,
	}
}

// Render does the actual rendering
func (rp *RenderPass) Render(camera camera.Camera) {
	rp.flatshader.Use()
	rp.flatshader.UpdateMat4("V", camera.GetView())
	rp.flatshader.UpdateMat4("P", camera.GetPerspective())
	rp.flatshader.UpdateMat4("M", mgl32.Ident4())
	rp.flatshader.UpdateVec3("uCameraPos", camera.GetPos())
	rp.flatshader.Render()
	rp.flatshader.Release()
}

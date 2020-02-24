package skybox

import (
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	mesh "github.com/adrianderstroff/pbr/pkg/view/mesh"
	tex "github.com/adrianderstroff/pbr/pkg/view/texture"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/geometry"
)

// Make constructs a skybox made from a quad with the cube map textures
// specified by the provided paths as well as the rendering mode.
func Make(sidelength float32, right, left, top, bottom, front, back string, mode uint32) (mesh.Mesh, error) {
	// make geometry
	geometry := makeCubeGeometry(sidelength, sidelength, sidelength)
	// make texture
	cubemap, err := tex.MakeCubeMap(right, left, top, bottom, front, back, true)
	if err != nil {
		return mesh.Mesh{}, err
	}
	textures := []tex.Texture{cubemap}
	// make mesh
	mesh := mesh.Make(geometry, textures, mode)
	// add actions
	prerender := func() {
		gl.DepthMask(false)
	}
	postrender := func() {
		gl.DepthMask(true)
	}
	mesh.SetPreRenderAction(prerender)
	mesh.SetPostRenderAction(postrender)
	return mesh, nil
}

// MakeFromDirectory constructs a skybox made from a quad with the specified side length
// in all 3 dimensions as well as the  the cube map textures specified by the provided
// directory and fileending as well as the rendering mode.
// The specified directory has to have all images in the same file format and the names
// of the files have to be named right, left, top, bottom, front and back respectively.
func MakeFromDirectory(sidelength float32, dir, fileending string, mode uint32) (mesh.Mesh, error) {
	right := dir + "right." + fileending
	left := dir + "left." + fileending
	top := dir + "top." + fileending
	bottom := dir + "bottom." + fileending
	front := dir + "front." + fileending
	back := dir + "back." + fileending
	return Make(sidelength, right, left, top, bottom, front, back, mode)
}

// makeCubeGeometry creates a cube with the specified width, height and depth.
// If the normals should be inside the cube the inside parameter should be true.
func makeCubeGeometry(width, height, depth float32) mesh.Geometry {
	// half side lengths
	halfWidth := width / 2.0
	halfHeight := height / 2.0
	halfDepth := depth / 2.0

	// vertex positions
	v1 := []float32{-halfWidth, halfHeight, halfDepth}
	v2 := []float32{-halfWidth, -halfHeight, halfDepth}
	v3 := []float32{halfWidth, halfHeight, halfDepth}
	v4 := []float32{halfWidth, -halfHeight, halfDepth}
	v5 := []float32{-halfWidth, halfHeight, -halfDepth}
	v6 := []float32{-halfWidth, -halfHeight, -halfDepth}
	v7 := []float32{halfWidth, halfHeight, -halfDepth}
	v8 := []float32{halfWidth, -halfHeight, -halfDepth}

	positions := geometry.Combine(
		// right
		v7, v8, v3,
		v3, v8, v4,
		// left
		v1, v2, v5,
		v5, v2, v6,
		// top
		v7, v3, v5,
		v5, v3, v1,
		// bottom
		v4, v8, v2,
		v2, v8, v6,
		// front
		v3, v4, v1,
		v1, v4, v2,
		//back
		v5, v6, v7,
		v7, v6, v8,
	)

	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := geometry.Repeat(geometry.Combine(t1, t2, t3, t3, t2, t4), 6)

	// normals
	right := []float32{-1.0, 0.0, 0.0}
	left := []float32{1.0, 0.0, 0.0}
	top := []float32{0.0, -1.0, 0.0}
	bottom := []float32{0.0, 1.0, 0.0}
	front := []float32{0.0, 0.0, -1.0}
	back := []float32{0.0, 0.0, 1.0}
	normals := mesh.Combine(
		mesh.Repeat(right, 6),
		mesh.Repeat(left, 6),
		mesh.Repeat(top, 6),
		mesh.Repeat(bottom, 6),
		mesh.Repeat(front, 6),
		mesh.Repeat(back, 6),
	)

	// setup data
	data := [][]float32{
		positions,
		uvs,
		normals,
	}

	// setup layout
	layout := []mesh.VertexAttribute{
		mesh.MakeVertexAttribute("pos", gl.FLOAT, 3, gl.STATIC_DRAW),
		mesh.MakeVertexAttribute("uv", gl.FLOAT, 2, gl.STATIC_DRAW),
		mesh.MakeVertexAttribute("normal", gl.FLOAT, 3, gl.STATIC_DRAW),
	}

	return mesh.MakeGeometry(layout, data)
}

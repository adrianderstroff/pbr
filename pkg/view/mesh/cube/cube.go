// Package box is used for creating a simple box mesh.
package cube

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	mesh "github.com/adrianderstroff/pbr/pkg/view/mesh"
)

// Make constructs a box with the specified dimensions. If inside is true
// then the triangles are specified in an order in which the normals will
// point inwards.
func Make(width, height, depth float32, inside bool, mode uint32) mesh.Mesh {
	geometry := makeCubeGeometry(width, height, depth, inside)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

// makeCubeGeometry creates a cube with the specified width, height and depth.
// If the normals should be inside the cube the inside parameter should be true.
func makeCubeGeometry(width, height, depth float32, inside bool) mesh.Geometry {
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
	positions := mesh.Combine(
		// right
		v3, v4, v7,
		v7, v4, v8,
		// left
		v5, v6, v1,
		v1, v6, v2,
		// top
		v5, v1, v7,
		v7, v1, v3,
		// bottom
		v2, v6, v4,
		v4, v6, v8,
		// front
		v1, v2, v3,
		v3, v2, v4,
		// back
		v7, v8, v5,
		v5, v8, v6,
	)
	if inside {
		positions = mesh.Combine(
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
	}

	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := mesh.Repeat(mesh.Combine(t1, t2, t3, t3, t2, t4), 6)

	// normals
	right := []float32{1.0, 0.0, 0.0}
	left := []float32{-1.0, 0.0, 0.0}
	top := []float32{0.0, 1.0, 0.0}
	bottom := []float32{0.0, -1.0, 0.0}
	front := []float32{0.0, 0.0, 1.0}
	back := []float32{0.0, 0.0, -1.0}
	// swap normals if inside is true
	if inside {
		right, left = left, right
		top, bottom = bottom, top
		front, back = back, front
	}
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

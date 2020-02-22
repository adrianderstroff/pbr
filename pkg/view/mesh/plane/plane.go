// Package plane is used for creating a simple plane mesh.
package plane

import (
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	mesh "github.com/adrianderstroff/pbr/pkg/view/mesh"
)

// Make constructs a plane with the specified dimensions. The plane is on the
// x-y axis and the normal points up the y-axis.
func Make(width, height float32, mode uint32) mesh.Mesh {
	geometry := makePlaneGeometry(width, height)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

// Make creates a Quad with the specified width and height on the x-y plane
// with the normal pointing up the y-axis
func makePlaneGeometry(width, height float32) mesh.Geometry {
	// half side lengths
	halfWidth := width / 2.0
	halfHeight := height / 2.0

	// vertex positions
	v1 := []float32{-halfWidth, 0, halfHeight}
	v2 := []float32{-halfWidth, 0, -halfHeight}
	v3 := []float32{halfWidth, 0, halfHeight}
	v4 := []float32{halfWidth, 0, -halfHeight}
	positions := mesh.Combine(v1, v2, v3, v3, v2, v4)

	// tex coordinates
	t1 := []float32{0.0, 1.0}
	t2 := []float32{0.0, 0.0}
	t3 := []float32{1.0, 1.0}
	t4 := []float32{1.0, 0.0}
	uvs := mesh.Combine(t1, t2, t3, t3, t2, t4)

	// normals
	up := []float32{0.0, 1.0, 0.0}
	normals := mesh.Repeat(up, 6)

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

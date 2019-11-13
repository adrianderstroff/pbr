// Package box is used for creating a simple box mesh.
package sphere

import (
	"github.com/adrianderstroff/pbr/pkg/view/geometry/gsphere"
	mesh "github.com/adrianderstroff/pbr/pkg/view/mesh"
)

// Make constructs a sphere of the specified horizontal and vertical
// resolution. The resolution should be bigger or equal to 1. Also the
// radius of the sphere has to be specified, it should be bigger than 0.
// The mode can be gl.Triangles, gl.TriangleStrip etc.
func Make(hres, vres int, radius float32, mode uint32) mesh.Mesh {
	geometry := gsphere.Make(hres, vres, radius)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

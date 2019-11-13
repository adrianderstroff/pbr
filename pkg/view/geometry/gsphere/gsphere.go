// Package gsphere creates sphere geometry
package gsphere

import (
	"math"

	"github.com/adrianderstroff/pbr/pkg/cgm"
	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	geometry "github.com/adrianderstroff/pbr/pkg/view/geometry"
	"github.com/go-gl/mathgl/mgl32"
)

// Make creates a Sphere with the specified horizontal and vertical resolution and a radius.
// The resolutions have to be 1 or greater
func Make(hres, vres int, radius float32) geometry.Geometry {
	// enforce boundary conditions
	hres = int(math.Max(float64(hres), 1))
	vres = int(math.Max(float64(vres), 1))

	// half side lengths
	w := 2*hres + 1
	h := 2*vres + 1

	var positions []float32
	var uvs []float32
	var normals []float32

	// all other rings
	var rings = make([][]mgl32.Vec3, h)
	var tempuvs = make([][]mgl32.Vec2, h)
	for y := 0; y < h; y++ {
		rings = append(rings, make([]mgl32.Vec3, w))
		tempuvs = append(tempuvs, make([]mgl32.Vec2, w))

		for x := 0; x <= w; x++ {
			// uv coordinates
			u := float32(x) / float32(w)
			v := 1.0 - float32(y)/float32(h-1)

			// spherical coordinates
			theta := 2 * math.Pi * u
			phi := math.Pi * v

			// spherical to cartesian
			px := radius * cgm.Cos32(theta) * cgm.Sin32(phi)
			py := radius * cgm.Cos32(phi)
			pz := radius * (-cgm.Sin32(theta)) * cgm.Sin32(phi)

			// uv coordinates
			//n := mgl32.Vec3{px, py, pz}.Normalize()
			//s := 0.5 + float32(math.Atan2(float64(n.X()), float64(n.Z())))/(2*math.Pi)
			//t := 0.5 - float32(math.Asin(float64(n.Y())))/math.Pi

			// add to arrays
			rings[y] = append(rings[y], mgl32.Vec3{px, py, pz})
			tempuvs[y] = append(tempuvs[y], mgl32.Vec2{u, v})
		}
	}

	// create the vertex data
	for y := 1; y < h; y++ {
		for x := 1; x <= w; x++ {
			// ^ 1------3
			// | |    / |
			// | | /    |
			// y 2------4
			//   x------>
			p1 := rings[y][x-1]
			p2 := rings[y-1][x-1]
			p3 := rings[y][x]
			p4 := rings[y-1][x]

			n1 := p1.Normalize()
			n2 := p2.Normalize()
			n3 := p3.Normalize()
			n4 := p4.Normalize()

			uv1 := tempuvs[y][x-1]
			uv2 := tempuvs[y-1][x-1]
			uv3 := tempuvs[y][x]
			uv4 := tempuvs[y-1][x]

			if y == h-1 {
				// add positions
				positions = append(positions, p3.X(), p3.Y(), p3.Z())
				positions = append(positions, p2.X(), p2.Y(), p2.Z())
				positions = append(positions, p4.X(), p4.Y(), p4.Z())

				// add uvs
				uvs = append(uvs, uv3.X(), uv3.Y())
				uvs = append(uvs, uv2.X(), uv2.Y())
				uvs = append(uvs, uv4.X(), uv4.Y())

				// add normals
				normals = append(normals, n3.X(), n3.Y(), n3.Z())
				normals = append(normals, n2.X(), n2.Y(), n2.Z())
				normals = append(normals, n4.X(), n4.Y(), n4.Z())
			} else if y == 1 {
				// add positions
				positions = append(positions, p1.X(), p1.Y(), p1.Z())
				positions = append(positions, p2.X(), p2.Y(), p2.Z())
				positions = append(positions, p3.X(), p3.Y(), p3.Z())

				// add uvs
				uvs = append(uvs, uv1.X(), uv1.Y())
				uvs = append(uvs, uv2.X(), uv2.Y())
				uvs = append(uvs, uv3.X(), uv3.Y())

				// add normals
				normals = append(normals, n1.X(), n1.Y(), n1.Z())
				normals = append(normals, n2.X(), n2.Y(), n2.Z())
				normals = append(normals, n3.X(), n3.Y(), n3.Z())
			} else {
				// add positions
				positions = append(positions, p1.X(), p1.Y(), p1.Z())
				positions = append(positions, p2.X(), p2.Y(), p2.Z())
				positions = append(positions, p3.X(), p3.Y(), p3.Z())
				positions = append(positions, p3.X(), p3.Y(), p3.Z())
				positions = append(positions, p2.X(), p2.Y(), p2.Z())
				positions = append(positions, p4.X(), p4.Y(), p4.Z())

				// add uvs
				uvs = append(uvs, uv1.X(), uv1.Y())
				uvs = append(uvs, uv2.X(), uv2.Y())
				uvs = append(uvs, uv3.X(), uv3.Y())
				uvs = append(uvs, uv3.X(), uv3.Y())
				uvs = append(uvs, uv2.X(), uv2.Y())
				uvs = append(uvs, uv4.X(), uv4.Y())

				// add normals
				normals = append(normals, n1.X(), n1.Y(), n1.Z())
				normals = append(normals, n2.X(), n2.Y(), n2.Z())
				normals = append(normals, n3.X(), n3.Y(), n3.Z())
				normals = append(normals, n3.X(), n3.Y(), n3.Z())
				normals = append(normals, n2.X(), n2.Y(), n2.Z())
				normals = append(normals, n4.X(), n4.Y(), n4.Z())
			}
		}
	}

	// setup data
	data := [][]float32{
		positions,
		uvs,
		normals,
	}

	// setup layout
	layout := []geometry.VertexAttribute{
		geometry.MakeVertexAttribute("pos", gl.FLOAT, 3, gl.STATIC_DRAW),
		geometry.MakeVertexAttribute("uv", gl.FLOAT, 2, gl.STATIC_DRAW),
		geometry.MakeVertexAttribute("normal", gl.FLOAT, 3, gl.STATIC_DRAW),
	}

	return geometry.Make(layout, data)
}

// Package cylinder is used for creating a simple cylinder mesh.
package cylinder

import (
	"math"

	"github.com/adrianderstroff/pbr/pkg/cgm"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	mesh "github.com/adrianderstroff/pbr/pkg/view/mesh"
	"github.com/go-gl/mathgl/mgl32"
)

// CoordinateSystem consists of 3 orthogonal normalized vectors that form a
// basis in R3.
type CoordinateSystem struct {
	b mgl32.Vec3
	t mgl32.Vec3
	n mgl32.Vec3
}

// Make constructs a cylinder specified by the start and end points p1 and p2
// as well as the radius of the cylinder.
// The mode can be gl.Triangles, gl.TriangleStrip etc.
func Make(p1, p2 mgl32.Vec3, radius float32, mode uint32) mesh.Mesh {
	geometry := makeCylinderGeometry(p1, p2, radius)
	mesh := mesh.Make(geometry, nil, mode)
	return mesh
}

// MakeCoordinateSystem calculates a orthographic coordinate system of the form
// b, t, n with n being the normalized direction, b the bitangent and t the
// tangent.
func MakeCoordinateSystem(dir mgl32.Vec3) CoordinateSystem {
	n := dir
	t := mgl32.Vec3{0, 1, 0}
	if n.Dot(mgl32.Vec3{0, 1, 0}) == 1 {
		t = mgl32.Vec3{1, 0, 0}
	}
	b := t.Cross(n)
	t = n.Cross(b)

	return CoordinateSystem{b, t, n}
}

// calcPos takes a position p and adds an offset to it. the offset is calculated
// from the index i and an delta angle step in degrees.
func calcPos(p mgl32.Vec3, radius float32, i int, step float32, coords CoordinateSystem) mgl32.Vec3 {
	a := (float32(i) * step) * (math.Pi / 180)
	x := cgm.Cos32(a)
	y := cgm.Sin32(a)

	vx := coords.t.Mul(x)
	vy := coords.b.Mul(y)
	offset := vx.Add(vy).Mul(radius)

	return p.Add(offset)
}

// Make creates a cylinder with the specified start and end position p1 and p2
// as well as the radius of the cylinder.
func makeCylinderGeometry(p1, p2 mgl32.Vec3, radius float32) mesh.Geometry {
	var positions []float32
	var uvs []float32
	var normals []float32

	// create all points
	var temppos []mgl32.Vec3

	// calculate coordinate system
	dir := p2.Sub(p1).Normalize()
	coords := MakeCoordinateSystem(dir)

	// resolution of the cylinder
	res := 20
	astep := 360 / float32(res)

	// add first point of the first circle
	temppos = append(temppos, calcPos(p1, radius, 0, astep, coords))

	// create first circle
	for i := 1; i <= res; i++ {
		i1, i2 := i-1, i%res

		pc1 := temppos[i1]
		pc2 := calcPos(p1, radius, i2, astep, coords)

		// add temp positions
		temppos = append(temppos, pc2)

		// add vertex attributes
		positions = append(positions, p1.X(), p1.Y(), p1.Z())
		positions = append(positions, pc1.X(), pc1.Y(), pc1.Z())
		positions = append(positions, pc2.X(), pc2.Y(), pc2.Z())
		uvs = append(uvs, 0, 0)
		uvs = append(uvs, float32(i2)/float32(res), 1)
		uvs = append(uvs, float32(i1)/float32(res), 1)
		normals = append(normals, -coords.n.X(), -coords.n.Y(), -coords.n.Z())
		normals = append(normals, -coords.n.X(), -coords.n.Y(), -coords.n.Z())
		normals = append(normals, -coords.n.X(), -coords.n.Y(), -coords.n.Z())
	}

	// add first point of the second circle
	temppos = append(temppos, calcPos(p2, radius, 0, astep, coords))

	// create second circle
	for i := 1; i <= res; i++ {
		i1, i2 := i-1, i%res

		pc2 := temppos[i1+(res+1)]
		pc1 := calcPos(p2, radius, i2, astep, coords)

		// add temp positions
		temppos = append(temppos, pc1)

		// add vertex attributes
		positions = append(positions, p2.X(), p2.Y(), p2.Z())
		positions = append(positions, pc1.X(), pc1.Y(), pc1.Z())
		positions = append(positions, pc2.X(), pc2.Y(), pc2.Z())
		uvs = append(uvs, 0, 0)
		uvs = append(uvs, float32(i2)/float32(res), 1)
		uvs = append(uvs, float32(i1)/float32(res), 1)
		normals = append(normals, coords.n.X(), coords.n.Y(), coords.n.Z())
		normals = append(normals, coords.n.X(), coords.n.Y(), coords.n.Z())
		normals = append(normals, coords.n.X(), coords.n.Y(), coords.n.Z())
	}

	// create tube
	// 1---3  ring 2
	// | / |
	// 2---4  ring 1
	for i := 1; i <= res; i++ {
		i2, i1 := i-1, i%res

		pc1 := temppos[(i1)+(res+1)]
		pc2 := temppos[i1]
		pc3 := temppos[i2+(res+1)]
		pc4 := temppos[i2]

		n1 := pc1.Sub(p2).Normalize()
		n2 := pc3.Sub(p2).Normalize()

		// triangle 1
		positions = append(positions, pc1.X(), pc1.Y(), pc1.Z())
		positions = append(positions, pc2.X(), pc2.Y(), pc2.Z())
		positions = append(positions, pc3.X(), pc3.Y(), pc3.Z())
		uvs = append(uvs, float32(i1)/float32(res), 1)
		uvs = append(uvs, float32(i1)/float32(res), 0)
		uvs = append(uvs, float32(i2)/float32(res), 1)
		normals = append(normals, n1.X(), n1.Y(), n1.Z())
		normals = append(normals, n1.X(), n1.Y(), n1.Z())
		normals = append(normals, n2.X(), n2.Y(), n2.Z())

		// triangle 2
		positions = append(positions, pc3.X(), pc3.Y(), pc3.Z())
		positions = append(positions, pc2.X(), pc2.Y(), pc2.Z())
		positions = append(positions, pc4.X(), pc4.Y(), pc4.Z())
		uvs = append(uvs, float32(i2)/float32(res), 1)
		uvs = append(uvs, float32(i1)/float32(res), 0)
		uvs = append(uvs, float32(i2)/float32(res), 0)
		normals = append(normals, n2.X(), n2.Y(), n2.Z())
		normals = append(normals, n1.X(), n1.Y(), n1.Z())
		normals = append(normals, n2.X(), n2.Y(), n2.Z())
	}

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

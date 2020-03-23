package obj

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/view/mesh"
	"github.com/go-gl/mathgl/mgl32"
)

// Face respresents a polygon that holds its vertices, normals and uv
// coordinates.
type Face struct {
	Positions []int
	UVs       []int
	Normals   []int
}

// Load a mesh from an .obj file.
func Load(filepath string, invert, smooth bool) (mesh.Mesh, error) {
	// setup temp variables
	faces := []Face{}
	tpositions := []float32{}
	tnormals := []float32{}
	tuvs := []float32{}

	// extract all vertex attributes and faces from the file
	err := extract(filepath, &faces, &tpositions, &tnormals, &tuvs)
	if err != nil {
		return mesh.Mesh{}, err
	}

	// break down faces consisting of polygons with more than 3 vertices into
	// a set of triangles
	sanitizeFaces(&faces)

	// generate object from faces and vertex attributes
	positions, uvs, normals := generateObject(faces, tpositions, tnormals, tuvs, smooth)

	fmt.Printf("p: %v\n", len(positions))
	fmt.Printf("t: %v\n", len(uvs))
	fmt.Printf("n: %v\n", len(normals))

	// calc center of gravity
	positions = center(positions)

	// invert normals if requested
	if invert {
		flipNormals(&normals)
	}

	// setup data
	mesh := createMesh(positions, uvs, normals)
	return mesh, nil
}

func extract(filepath string, faces *[]Face, positions, normals, uvs *[]float32) error {

	// opening the file
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return err
	}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// tokens are separated by whitespace
		tokens := strings.Split(scanner.Text(), " ")

		// depending on the type of the vertex attribute parse token to float
		// and add it to the according slice. the face attribute f is a bit
		// more involved as it describes the indices of the vertices, normals
		// and uv coordinates for the specific face. a face can be a triangle
		// or any higher polygon. only convex polygons are supported tho.
		switch tokens[0] {
		case "v":
			for _, token := range tokens[1:] {
				parseAdd(token, positions)
			}
		case "vn":
			for _, token := range tokens[1:] {
				parseAdd(token, normals)
			}
		case "vt":
			for _, token := range tokens[1:] {
				parseAdd(token, uvs)
			}
		case "f":
			face := Face{}
			// iterate over vertices
			for _, token := range tokens[1:] {

				// vertex attribute indices
				vertexattridxs := []int{-1, -1, -1}

				// grab the indices for the respective vertex attribute index.
				// each face token consists of a tripplet p/t/n of indices
				// describing the indices for position,texture (uv) coordinate
				// and normal. the first index must be defined. the other two
				// indices are optional.
				ftokens := strings.Split(token, "/")
				for i, ftoken := range ftokens {
					idx, err := strconv.Atoi(ftoken)
					if err == nil {
						vertexattridxs[i] = idx
					}
				}

				// the vertex position has to be specified. if not skip this
				// face vertex.
				if vertexattridxs[0] == -1 {
					continue
				}
				// if vertex uv index is not specified then set it to be the
				// same as the index of the vertex position.
				if vertexattridxs[1] == -1 {
					vertexattridxs[1] = vertexattridxs[0]
				}
				// if vertex normal index is not specified then set it to be the
				// same as the index of the vertex position.
				if vertexattridxs[2] == -1 {
					vertexattridxs[2] = vertexattridxs[0]
				}

				// added vertex infos
				face.Positions = append(face.Positions, vertexattridxs[0])
				face.UVs = append(face.UVs, vertexattridxs[1])
				face.Normals = append(face.Normals, vertexattridxs[2])
			}
			*faces = append(*faces, face)
		}
	}

	// was reading from file was successful ?
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func sanitizeFaces(faces *[]Face) {
	var newfaces []Face

	// if a polygon consists of more than three vertices then break it down into
	// multiple triangles. this approach assumes that the polygon is convex.
	for _, face := range *faces {
		if len(face.Positions) == 3 {
			newfaces = append(newfaces, face)
		} else {
			for fidx := 2; fidx < len(face.Positions); fidx++ {
				newface := Face{
					Positions: []int{face.Positions[0], face.Positions[fidx-1], face.Positions[fidx]},
					UVs:       []int{face.UVs[0], face.UVs[fidx-1], face.UVs[fidx]},
					Normals:   []int{face.Normals[0], face.Normals[fidx-1], face.Normals[fidx]},
				}
				newfaces = append(newfaces, newface)
			}
		}
	}

	// overwrite faces
	*faces = newfaces
}

func generateObject(faces []Face, tpositions, tnormals, tuvs []float32,
	smooth bool) ([]float32, []float32, []float32) {

	// setup obj properties
	positions := []float32{}
	uvs := []float32{}
	normals := []float32{}

	// create neighborhood info if smooth is enabled and tnormals are empty
	neighborhood := map[int][]int{}
	if smooth && len(tnormals) == 0 {
		for fidx, face := range faces {
			// save neighbor faces for current vertex
			for _, nidx := range face.Normals {
				if _, ok := neighborhood[nidx]; !ok {
					neighborhood[nidx] = []int{fidx}
				} else {
					neighbors := neighborhood[nidx]
					neighbors = append(neighbors, fidx)
				}
			}
		}
	}

	// build structure from faces
	for _, face := range faces {
		// extract vertex positions
		extractVertexAttribute(&positions, tpositions, face.Positions, 3)

		// extract or create normals
		if len(tnormals) > 0 {
			extractVertexAttribute(&normals, tnormals, face.Normals, 3)
		} else {
			if smooth {
				// normals of each face vertex
				n1 := smoothNormal(faces, tpositions, neighborhood, face.Normals[0])
				n2 := smoothNormal(faces, tpositions, neighborhood, face.Normals[1])
				n3 := smoothNormal(faces, tpositions, neighborhood, face.Normals[2])
				normals = append(normals, n1.X(), n1.Y(), n1.Z())
				normals = append(normals, n2.X(), n2.Y(), n2.Z())
				normals = append(normals, n3.X(), n3.Y(), n3.Z())
			} else {
				n := cross(&face, tpositions)
				normals = append(normals, n.X(), n.Y(), n.Z())
				normals = append(normals, n.X(), n.Y(), n.Z())
				normals = append(normals, n.X(), n.Y(), n.Z())
			}
		}

		// extract or create texture (uv) coordinates
		if len(tuvs) > 0 {
			extractVertexAttribute(&uvs, tuvs, face.UVs, 2)
		} else {
			uvs = append(uvs, 0, 0)
			uvs = append(uvs, 1, 0)
			uvs = append(uvs, 1, 1)
		}
	}

	return positions, uvs, normals
}

func parseAdd(token string, slice *[]float32) {
	v, err := strconv.ParseFloat(token, 32)
	if err == nil {
		*slice = append(*slice, float32(v))
	}
}

func extractVertexAttribute(outSlice *[]float32, inSlice []float32, indices []int, offset int) {
	// get indices
	idx1 := indices[0] - 1
	idx2 := indices[1] - 1
	idx3 := indices[2] - 1
	idxs := []int{idx1, idx2, idx3}

	// copy elements from in to out slice
	for i := 0; i < 3; i++ {
		idx := idxs[i]
		for o := 0; o < offset; o++ {
			*outSlice = append(*outSlice, inSlice[idx*offset+o])
		}
	}
}

func cross(face *Face, positions []float32) mgl32.Vec3 {
	// grab positions
	idx1 := face.Positions[0] - 1
	idx2 := face.Positions[1] - 1
	idx3 := face.Positions[2] - 1
	p1 := componentsToVec3(positions, idx1)
	p2 := componentsToVec3(positions, idx2)
	p3 := componentsToVec3(positions, idx3)

	// calc directions
	v1 := p1.Sub(p2)
	v2 := p3.Sub(p2)

	// calc cross product
	n := v1.Cross(v2)

	// return normal
	return n.Normalize()
}

func componentsToVec3(positions []float32, idx int) mgl32.Vec3 {
	v1 := positions[idx*3+0]
	v2 := positions[idx*3+1]
	v3 := positions[idx*3+2]
	return mgl32.Vec3{v1, v2, v3}
}

func smoothNormal(faces []Face, positions []float32, neighborhood map[int][]int, nidx int) mgl32.Vec3 {
	n := mgl32.Vec3{0}

	for _, fidx := range neighborhood[nidx] {
		nface := faces[fidx]
		tnormal := cross(&nface, positions)
		n = n.Add(tnormal)
	}

	return n.Normalize()
}

func center(vertices []float32) []float32 {
	var (
		x float32 = 0.0
		y float32 = 0.0
		z float32 = 0.0

		minX float64 = math.Inf(1)
		maxX float64 = math.Inf(-1)
		minY float64 = math.Inf(1)
		maxY float64 = math.Inf(-1)
		minZ float64 = math.Inf(1)
		maxZ float64 = math.Inf(-1)
	)
	vertexCount := len(vertices) / 3
	for i := 0; i < vertexCount; i++ {
		posX := vertices[i*3+0]
		posY := vertices[i*3+1]
		posZ := vertices[i*3+2]
		x += posX
		y += posY
		z += posZ
		minX = math.Min(minX, float64(posX))
		minY = math.Min(minY, float64(posY))
		minZ = math.Min(minZ, float64(posZ))
		maxX = math.Max(maxX, float64(posX))
		maxY = math.Max(maxY, float64(posY))
		maxZ = math.Max(maxZ, float64(posZ))
	}
	x /= float32(vertexCount)
	y /= float32(vertexCount)
	z /= float32(vertexCount)

	// center vertices
	for i := 0; i < vertexCount; i++ {
		diffX := (maxX - minX) / 2
		diffY := (maxY - minY) / 2
		diffZ := (maxZ - minZ) / 2
		diff := float32(math.Max(diffX, math.Max(diffY, diffZ)))
		vertices[i*3+0] = (vertices[i*3+0] - x) / diff
		vertices[i*3+1] = (vertices[i*3+1] - y) / diff
		vertices[i*3+2] = (vertices[i*3+2] - z) / diff
	}

	return vertices
}

func flipNormals(normals *[]float32) {
	for i := 0; i < len(*normals); i++ {
		(*normals)[i] = -(*normals)[i]
	}
}

func createMesh(positions, uvs, normals []float32) mesh.Mesh {
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

	// create the mesh
	geometry := mesh.MakeGeometry(layout, data)
	return mesh.Make(geometry, nil, gl.TRIANGLES)
}

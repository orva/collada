package collada

import (
	"errors"
)

type Vertex struct {
	X float32
	Y float32
	Z float32
}

type Mesh struct {
	// TODO Implement tristrips and trifans
	// TODO Normals
	// Contains vertices, but not all of them. You need to use VertexPrimitives
	// to build actual triangles.
	Vertices []Vertex
	// Contains indexes to vertices
	VertexPrimitives []int
	Id               string
	Name             string
}

// Create single triangulated Mesh from given MeshData
func NewMesh(m *MeshData, id, name string) (*Mesh, error) {
	var err error

	vertices, err := m.vertices()
	if err != nil {
		return nil, err
	}

	primitives, err := m.Triangles.primitives()
	if err != nil {
		return nil, err
	}

	mesh := &Mesh{
		Vertices:         vertices,
		VertexPrimitives: primitives,
		Id:               id,
		Name:             name,
	}
	return mesh, nil
}

func (m *MeshData) triangles() ([]Vertex, error) {
	if m.Triangles == nil {
		return nil, errors.New("Mesh is not triangulated")
	}

	var err error

	vertices, err := m.vertices()
	if err != nil {
		return nil, err
	}

	primitives, err := m.Triangles.primitives()
	if err != nil {
		return nil, err
	}

	vertOffset := m.Triangles.semantic("VERTEX")
	retval := make([]Vertex, 0)
	baseIndex := 0
	stride := len(m.Triangles.Inputs)
	for count := 0; count < m.Triangles.Count; count++ {
		// There should always be valid amount of indices, crash and burn if
		// there isn't

		// Fetch right vertex with index from primitives array.
		vertex := vertices[primitives[baseIndex+vertOffset]]
		retval = append(retval, vertex)
		baseIndex += stride
	}

	return retval, nil
}

// Apply SourceData.Accessor to SourceData.FloatArray to get vertices.
func (m *MeshData) vertices() ([]Vertex, error) {
	var err error

	acc, err := m.vertexAccessor()
	if err != nil {
		return nil, err
	}

	vFloats, err := m.vertexFloats()
	if err != nil {
		return nil, err
	}

	xInd := acc.paramIndex("X")
	yInd := acc.paramIndex("Y")
	zInd := acc.paramIndex("Z")

	retval := make([]Vertex, 0)
	baseIndex := acc.Offset
	count := 0
	for {
		// Lets just crash with invalid index if we have too few vertex points.
		vertex := Vertex{
			X: vFloats[baseIndex+xInd],
			Y: vFloats[baseIndex+yInd],
			Z: vFloats[baseIndex+zInd],
		}
		retval = append(retval, vertex)

		baseIndex += acc.Stride
		count += 1
		if count >= acc.Count {
			break
		}
	}

	return retval, nil
}

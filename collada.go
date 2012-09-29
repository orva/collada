package collada

import (
	"errors"
)

type Vertex struct {
	X float64
	Y float64
	Z float64
}

type Mesh struct {
	// Contains vertices which you can use to draw triangles. At this point no
	// tristrips or trifans are implemented.
	// TODO Implement tristrips and trifans
	// TODO Normals
	Vertices []Vertex
	Id       string
	Name     string
}

func newMesh(m *MeshData, id, name string) (*Mesh, error) {
	if m.Triangles == nil {
		return nil, errors.New("Mesh is not triangulated")
	}

	return nil, errors.New("Not implemented")
}

func (m *MeshData) triangles() ([]Vertex, error) {
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

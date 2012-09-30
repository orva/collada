package collada

import (
	. "launchpad.net/gocheck"
)

type ColladaSuite struct {
	data *ColladaData
	mesh *Mesh
}

const (
	NAME = "mesh name"
	ID   = "mesh id"
)

var _ = Suite(&ColladaSuite{})

func (s *ColladaSuite) SetUpTest(c *C) {
	s.data, _ = ParseColladaData("test-data/cube_triangulate.dae")
	s.mesh, _ = NewMesh(s.data.Geometries[0].Mesh, ID, NAME)
}

func (s *ColladaSuite) TestNewMeshAttributes(c *C) {
	c.Check(s.mesh.Id, Equals, ID)
	c.Check(s.mesh.Name, Equals, NAME)
}

func (s *ColladaSuite) TestNewMeshVertices(c *C) {
	expectedMesh, _ := s.data.Geometries[0].Mesh.vertices()
	c.Check(s.mesh.Vertices, DeepEquals, expectedMesh)
}

func (s *ColladaSuite) TestNewMeshPrimitives(c *C) {
	expected, _ := s.data.Geometries[0].Mesh.Triangles.primitives()
	c.Check(s.mesh.VertexPrimitives, DeepEquals, expected)
}

func (s *ColladaSuite) TestGettingVertices(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.vertices()
	c.Check(len(actual), Equals, 8)
}

func (s *ColladaSuite) TestGettingTriangles(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.triangles()
	c.Check(len(actual), Equals, 12)
}

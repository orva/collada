package collada

import (
	. "launchpad.net/gocheck"
)

type ColladaSuite struct {
	data *ColladaData
}

var _ = Suite(&ColladaSuite{})

func (s *ColladaSuite) SetUpTest(c *C) {
	s.data, _ = ParseColladaData("test-data/cube_triangulate.dae")
}

func (s *ColladaSuite) TestNewMesh(c *C) {
	mesh, _ := NewMesh(s.data.Geometries[0].Mesh, "id", "name")
	expected, _ := s.data.Geometries[0].Mesh.triangles()
	c.Check(mesh.Vertices, DeepEquals, expected)
	c.Check(mesh.Id, Equals, "id")
	c.Check(mesh.Name, Equals, "name")
}

func (s *ColladaSuite) TestGettingVertices(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.vertices()
	c.Check(len(actual), Equals, 8)
}

func (s *ColladaSuite) TestGettingTriangles(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.triangles()
	c.Check(len(actual), Equals, 12)
}

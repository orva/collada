package collada

import (
	. "launchpad.net/gocheck"
)

type ColladaSuite struct {
	data *ColladaData
}

var _ = Suite(&ColladaSuite{})

func (s *ColladaSuite) SetUpTest(c *C) {
	s.data, _ = parseColladaData("test-data/cube_triangulate.dae")
}

func (s *ColladaSuite) TestNewMesh(c *C) {
	_, err := newMesh(s.data.Geometries[0].Mesh, "", "")
	c.Check(err.Error(), Equals, "Not implemented")
}

func (s *ColladaSuite) TestGettingVertices(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.vertices()
	c.Check(len(actual), Equals, 8)
}

func (s *ColladaSuite) TestGettingTriangles(c *C) {
	actual, _ := s.data.Geometries[0].Mesh.triangles()
	c.Check(len(actual), Equals, 12)
}

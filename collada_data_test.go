package collada

import (
	. "launchpad.net/gocheck"
)

type ColladaDataSuite struct {
	data *ColladaData
}

var _ = Suite(&ColladaDataSuite{})

func (s *ColladaDataSuite) SetUpTest(c *C) {
	s.data, _ = ParseColladaData("test-data/cube_triangulate.dae")
}

func (s *ColladaDataSuite) TestXMLName(c *C) {
	c.Check(s.data.XMLName.Local, Equals, "COLLADA")
	c.Check(s.data.XMLName.Space, Equals,
		"http://www.collada.org/2005/11/COLLADASchema")
}

func (s *ColladaDataSuite) TestThatFileNeedsToBeColladaFile(c *C) {
	_, err := ParseColladaData("test-data/invalid_top_level.dae")
	c.Assert(err, NotNil)
	tmp := &NotValidColladaFileError{}
	c.Assert(err.Error(), Equals, tmp.Error())
}

func (s *ColladaDataSuite) TestSimpleGeometries(c *C) {
	c.Check(len(s.data.Geometries), Equals, 1)
}

func (s *ColladaDataSuite) TestGeometryAttributes(c *C) {
	c.Check(s.data.Geometries[0].Id, Equals, "box-lib")
	c.Check(s.data.Geometries[0].Name, Equals, "box")
}

func (s *ColladaDataSuite) TestGeometryFields(c *C) {
	c.Check(s.data.Geometries[0].Mesh, NotNil)
}

func (s *ColladaDataSuite) TestTriangleData(c *C) {
	triangles := s.data.Geometries[0].Mesh.Triangles
	c.Check(triangles.Count, Equals, 12)
	c.Check(triangles.Inputs[0].Offset, Equals, 0)
	c.Check(triangles.Inputs[0].Semantic, Equals, "VERTEX")
	c.Check(triangles.Inputs[0].Source, Equals, "#box-lib-vertices")
	c.Check(triangles.Inputs[1].Offset, Equals, 1)
	c.Check(triangles.Inputs[1].Semantic, Equals, "NORMAL")
	c.Check(triangles.Inputs[1].Source, Equals, "#box-lib-normals")
	c.Check(triangles.P, Equals, "0 0 2 1 3 2 0 0 3 2 1 3 0 4 1 5 5 6 0 4 5 6 4 7 6 8 7 9 3 10 6 8 3 10 2 11 0 12 4 13 6 14 0 12 6 14 2 15 3 16 7 17 5 18 3 16 5 18 1 19 5 20 7 21 6 22 5 20 6 22 4 23")
}

func (s *ColladaDataSuite) TestTriangleDataPrimitives(c *C) {
	primitives, _ := s.data.Geometries[0].Mesh.Triangles.primitives()
	expected :=
		[]int{0, 0, 2, 1, 3, 2, 0, 0, 3, 2, 1, 3, 0, 4, 1, 5, 5, 6, 0, 4, 5, 6, 4, 7, 6, 8, 7, 9, 3, 10, 6, 8, 3, 10, 2, 11, 0, 12, 4, 13, 6, 14, 0, 12, 6, 14, 2, 15, 3, 16, 7, 17, 5, 18, 3, 16, 5, 18, 1, 19, 5, 20, 7, 21, 6, 22, 5, 20, 6, 22, 4, 23}
	c.Check(primitives, DeepEquals, expected)
}

func (s *ColladaDataSuite) TestVertexData(c *C) {
	vertices := s.data.Geometries[0].Mesh.Vertices
	c.Check(vertices.Input.Semantic, Equals, "POSITION")
	c.Check(vertices.Input.Source, Equals, "#box-lib-positions")
	c.Check(vertices.Input.Offset, Equals, 0)
}

func (s *ColladaDataSuite) TestMeshSourceAttributes(c *C) {
	s0 := s.data.Geometries[0].Mesh.Sources[0]
	s1 := s.data.Geometries[0].Mesh.Sources[1]

	c.Check(s0.Id, Equals, "box-lib-positions")
	c.Check(s0.Name, Equals, "position")
	c.Check(s1.Id, Equals, "box-lib-normals")
	c.Check(s1.Name, Equals, "normal")
}

func (s *ColladaDataSuite) TestMeshSourceAccessor(c *C) {
	a0 := s.data.Geometries[0].Mesh.Sources[0].Accessor
	a1 := s.data.Geometries[0].Mesh.Sources[1].Accessor

	c.Check(a0.Count, Equals, 8)
	c.Check(a0.Offset, Equals, 0)
	c.Check(a0.Source, Equals, "#box-lib-positions-array")
	c.Check(a0.Stride, Equals, 3)
	c.Assert(len(a0.Params), Equals, 3)

	c.Check(a1.Count, Equals, 24)
	c.Check(a1.Offset, Equals, 0)
	c.Check(a1.Source, Equals, "#box-lib-normals-array")
	c.Check(a1.Stride, Equals, 3)
	c.Check(len(a1.Params), Equals, 3)
}

func (s *ColladaDataSuite) TestMeshSourceAccessorParams(c *C) {
	p0 := s.data.Geometries[0].Mesh.Sources[0].Accessor.Params
	p1 := s.data.Geometries[0].Mesh.Sources[1].Accessor.Params

	c.Assert(len(p0), Equals, 3)
	c.Assert(len(p1), Equals, 3)

	c.Check(p0[0].Name, Equals, "X")
	c.Check(p0[1].Name, Equals, "Y")
	c.Check(p0[2].Name, Equals, "Z")
	c.Check(p0[0].Type, Equals, "float")
	c.Check(p0[1].Type, Equals, "float")
	c.Check(p0[2].Type, Equals, "float")

	c.Check(p1[0].Name, Equals, "X")
	c.Check(p1[1].Name, Equals, "Y")
	c.Check(p1[2].Name, Equals, "Z")
	c.Check(p1[0].Type, Equals, "float")
	c.Check(p1[1].Type, Equals, "float")
	c.Check(p1[2].Type, Equals, "float")
}

func (s *ColladaDataSuite) TestMeshSourceFloatArray(c *C) {
	arr := s.data.Geometries[0].Mesh.Sources[0].FloatArr
	c.Check(arr, NotNil)
	c.Check(arr.Id, Equals, "box-lib-positions-array")
	c.Check(arr.Count, Equals, 24)
	c.Check(arr.Data, Equals, "-50 50 50 50 50 50 -50 -50 50 50 -50 50 -50 50 -50 50 50 -50 -50 -50 -50 50 -50 -50")
}

func (s *ColladaDataSuite) TestSourceFloatExtract(c *C) {
	srcData := s.data.Geometries[0].Mesh.Sources[0]
	expected := []float32{-50, 50, 50, 50, 50, 50, -50, -50, 50, 50, -50, 50, -50, 50, -50, 50, 50, -50, -50, -50, -50, 50, -50, -50}
	actual, _ := srcData.extractFloats()

	c.Check(len(actual), Equals, 24)
	c.Check(actual, DeepEquals, expected)
}

func (s *ColladaDataSuite) TestSourceFloatExtractErrors(c *C) {
	srcData := s.data.Geometries[0].Mesh.Sources[0]
	srcData.FloatArr = nil
	_, err := srcData.extractFloats()

	c.Check(err, NotNil)
}

func (s *ColladaDataSuite) TestMeshVertexFloats(c *C) {
	meshData := s.data.Geometries[0].Mesh
	expected := []float32{-50, 50, 50, 50, 50, 50, -50, -50, 50, 50, -50, 50, -50, 50, -50, 50, 50, -50, -50, -50, -50, 50, -50, -50}
	actual, _ := meshData.vertexFloats()

	c.Check(actual, DeepEquals, expected)
}

func (s *ColladaDataSuite) TestMeshVertexFloatsErrors(c *C) {
	meshData := s.data.Geometries[0].Mesh
	input := meshData.Vertices.Input
	var err error
	var expected InvalidColladaId

	input.Source = ""
	_, err = meshData.vertexFloats()
	expected = InvalidColladaId{input.Source}
	c.Check(err.Error(), Equals, expected.Error())

	input.Source = "bogus"
	_, err = meshData.vertexFloats()
	expected = InvalidColladaId{input.Source}
	c.Check(err.Error(), Equals, expected.Error())
}

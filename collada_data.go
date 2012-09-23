package collada

import (
	"encoding/xml"
	"io/ioutil"
)

// ColladaData contains raw data parsed from collada xml. No work has yet been
// done to make this data anyhow usable for graphics pipeline.
//
// All types ending with *Data are considered unusable for graphics pipeline
// consumption.
//
// Specification for collada can be found in 'doc' directory.
type ColladaData struct {
	XMLName    xml.Name
	Geometries []GeometryData `xml:"library_geometries>geometry"`
}

// Check if parsed file contains right toplevel name.
func (c *ColladaData) isCollada() bool {
	return c.XMLName.Local == "COLLADA"
}

// See p.82 in collada spec.
type GeometryData struct {
	Id   string   `xml:"id,attr"`
	Name string   `xml:"name,attr"`
	Mesh MeshData `xml:"mesh"`
}

// See p.129 in collada spec.
type MeshData struct {
	Vertices  VertexData   `xml:"vertices"`
	Triangles TriangleData `xml:"triangles"`
	Sources   []SourceData `xml:"source"`
}

// See p.196 in collada spec.
type VertexData struct {
	Input InputData `xml:"input"`
}

// See p.188 in collada spec.
type TriangleData struct {
	Count  int         `xml:"count,attr"`
	Inputs []InputData `xml:"input"`
	P      string      `xml:"p"`
}

// See p.87 in collada spec.
type InputData struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   int    `xml:"offset,attr"`
}

// See p.177 in collada spec.
type SourceData struct {
	Id       string         `xml:"id,attr"`
	Name     string         `xml:"name,attr"`
	FloatArr FloatArrayData `xml:"float_array"`
	// In xml Accessor is inside 'technique_common' element, but here we drop
	// that indirection because in source section there can be only 'accessor'
	// element.
	Accessor AccessorData `xml:"technique_common>accessor"`
}

// See p.188 and p.77 in collada spec.
type FloatArrayData struct {
	Id    string `xml:"id,attr"`
	Count int    `xml:"count,attr"`
	Data  string `xml:",innerxml"`
}

// See p.45 in collada spec.
type AccessorData struct {
	Count  uint        `xml:"count,attr"`
	Offset uint        `xml:"offset,attr"`
	Source string      `xml:"source,attr"`
	Stride uint        `xml:"stride,attr"`
	Params []ParamData `xml:"param"`
}

// See p.144 in collada spec.
type ParamData struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

func parseColladaData(path string) (*ColladaData, error) {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data ColladaData
	if err = xml.Unmarshal(xmlFile, &data); err != nil {
		return nil, err
	} else if !data.isCollada() {
		return nil, &NotValidColladaFileError{}
	}

	return &data, err
}

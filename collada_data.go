package collada

import (
	"encoding/xml"
	"io/ioutil"
	"errors"
	"strings"
	"strconv"
	"fmt"
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
func (c *ColladaData) isCollada() bool { return c.XMLName.Local == "COLLADA" }

// See p.82 in collada spec.
type GeometryData struct {
	Id   string    `xml:"id,attr"`
	Name string    `xml:"name,attr"`
	Mesh *MeshData `xml:"mesh"`
}

// See p.129 in collada spec.
type MeshData struct {
	Vertices  *VertexData   `xml:"vertices"`
	Triangles *TriangleData `xml:"triangles"`
	Sources   []SourceData  `xml:"source"`
}

func (m *MeshData) vertexFloats() ([]float64, error) {
	sourceId := m.Vertices.Input.Source
	src, err := m.source(sourceId)
	if err != nil {
		return nil, err
	}

	return src.extractFloats()
}

func (m *MeshData) vertexAccessor() (*AccessorData, error) {
	sourceId := m.Vertices.Input.Source
	src, err := m.source(sourceId)
	if err != nil {
		return nil, err
	}

	return src.Accessor, nil
}

func (m *MeshData) source(id string) (*SourceData, error) {
	if len(id) == 0 {
		return nil, &InvalidColladaId{id}
	}

	cleanId := id[1:] // Strip leading hash
	for _, src := range m.Sources {
		if src.Id == cleanId {
			return &src, nil
		}
	}

	return nil, &InvalidColladaId{id}
}

// See p.196 in collada spec.
type VertexData struct {
	Input *InputData `xml:"input"`
}

// See p.188 in collada spec.
type TriangleData struct {
	Count  int         `xml:"count,attr"`
	Inputs []InputData `xml:"input"`
	P      string      `xml:"p"`
}

func (t *TriangleData) primitives() ([]int, error) {
	intStrs := strings.Split(t.P, " ")
	retval := make([]int, 0)
	for _, raw := range intStrs {
		primitive, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		retval = append(retval, primitive)
	}
	return retval, nil
}

func (t *TriangleData) semantic(semantic string) int {
	for index, input := range t.Inputs {
		if input.Semantic == semantic {
			return index
		}
	}

	return -1
}

// See p.87 in collada spec.
type InputData struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   int    `xml:"offset,attr"`
}

// See p.177 in collada spec.
type SourceData struct {
	Id       string          `xml:"id,attr"`
	Name     string          `xml:"name,attr"`
	FloatArr *FloatArrayData `xml:"float_array"`
	// In xml Accessor is inside 'technique_common' element, but here we drop
	// that indirection because in source section there can be only 'accessor'
	// element.
	Accessor *AccessorData `xml:"technique_common>accessor"`
}

func (s *SourceData) extractFloats() ([]float64, error) {
	if s.FloatArr == nil {
		return nil, errors.New("SourceData '" + s.Id + "' doesn't contain FloatArray")
	}

	retval := make([]float64, 0)
	floatStrs := strings.Split(s.FloatArr.Data, " ")
	for _, fstr := range floatStrs {
		fl, err := strconv.ParseFloat(fstr, 32)
		if err != nil {
			return nil, err
		}
		retval = append(retval, fl)
	}

	if len(retval) != s.FloatArr.Count {
		errStr := fmt.Sprint("Parsed only %d elements from #%s. %d was required",
			len(retval), s.FloatArr.Id, s.FloatArr.Count)
		return nil, errors.New(errStr)
	}

	return retval, nil
}

// See p.188 and p.77 in collada spec.
type FloatArrayData struct {
	Id    string `xml:"id,attr"`
	Count int    `xml:"count,attr"`
	Data  string `xml:",innerxml"`
}

// See p.45 in collada spec.
type AccessorData struct {
	Count  int         `xml:"count,attr"`
	Offset int         `xml:"offset,attr"`
	Source string      `xml:"source,attr"`
	Stride int         `xml:"stride,attr"`
	Params []ParamData `xml:"param"`
}

func (a *AccessorData) paramIndex(name string) int {
	for index, param := range a.Params {
		if param.Name == name {
			return index
		}
	}

	return -1
}

// See p.144 in collada spec.
type ParamData struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

func ParseColladaData(path string) (*ColladaData, error) {
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

package collada

type NotValidColladaFileError struct{}

func (e *NotValidColladaFileError) Error() string {
	return "File was not valid collada file, missed top level COLLADA element"
}

type InvalidColladaId struct {
	Id string
}

func (e *InvalidColladaId) Error() string {
	return "Couldn't find source element with id: " + e.Id
}

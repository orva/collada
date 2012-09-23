package collada

type NotValidColladaFileError struct{}

func (e *NotValidColladaFileError) Error() string {
	return "File was not valid collada file, missed top level COLLADA element"
}

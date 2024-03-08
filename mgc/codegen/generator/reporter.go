package generator

type Reporter interface {
	Generate(path, message string)
	Error(path, message string, err error)
}

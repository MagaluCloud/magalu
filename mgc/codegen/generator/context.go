package generator

import "github.com/spf13/afero"

type GeneratorContext struct {
	ModuleName string
	Reporter   Reporter
	FS         afero.Fs
}

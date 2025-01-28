package ui

import (
	"fmt"
	"io"
	"path/filepath"

	generator "github.com/MagaluCloud/magalu/mgc/codegen/generator"
)

type Reporter struct {
	BaseDir string
	Out     io.Writer
	Err     io.Writer
}

var _ generator.Reporter = (*Reporter)(nil)

func NewReporter(baseDir string, out io.Writer, err io.Writer) *Reporter {
	return &Reporter{
		BaseDir: baseDir,
		Out:     out,
		Err:     err,
	}
}

func (r *Reporter) Generate(p, message string) {
	if r.Out == nil {
		return
	}
	p, _ = filepath.Rel(r.BaseDir, p)
	fmt.Fprintf(r.Out, "Generate %s: %s\n", p, message)
}

func (r *Reporter) Error(p, message string, err error) {
	if r.Err == nil {
		return
	}
	p, _ = filepath.Rel(r.BaseDir, p)
	fmt.Fprintf(r.Err, "Error: generate %s: %s, error=%s\n", p, message, err.Error())
}

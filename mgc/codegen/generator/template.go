package generator

import (
	"bytes"
	"go/format"
	"path"
	"strings"
	"text/template"
)

func templteIndentFunc(indent int, prefix string, rest string) string {
	if indent < 1 {
		indent = 2
	}
	if prefix == "" {
		prefix = "\t"
	}

	return strings.Repeat(prefix, indent) + rest
}

var templateFuncs = template.FuncMap{
	"indent": templteIndentFunc,
}

func templateMust(name, contents string) (t *template.Template) {
	t, err := template.New(name).Funcs(templateFuncs).Parse(contents)
	if err != nil {
		panic(err.Error())
	}
	return t
}

func templateWrite[T any](ctx *GeneratorContext, name string, t *template.Template, data T) (err error) {
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		ctx.Reporter.Error(name, "failed to generate file", err)
		return
	}

	var source []byte
	if path.Ext(name) != ".go" {
		source = buf.Bytes()
	} else {
		source, err = format.Source(buf.Bytes())
		if err != nil {
			ctx.Reporter.Error(name, "failed to format file", err)
			_ = replaceFileIfNeeded(ctx, name, buf.Bytes()) // write it anyway, so we can check what's up
			return
		}
	}

	return replaceFileIfNeeded(ctx, name, source)
}

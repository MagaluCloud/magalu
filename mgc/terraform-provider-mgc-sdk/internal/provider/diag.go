package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

type Diagnostics diag.Diagnostics

func (d *Diagnostics) Contains(other diag.Diagnostic) bool {
	return (*diag.Diagnostics)(d).Contains(other)
}

func (d *Diagnostics) Append(toAppend ...diag.Diagnostic) {
	(*diag.Diagnostics)(d).Append(toAppend...)
}

func (d *Diagnostics) AddError(summary, detail string) {
	d.Append(NewErrorDiagnostic(summary, detail))
}

func (d *Diagnostics) AddLocalError(summary, detail string) {
	d.Append(NewLocalErrorDiagnostic(summary, detail))
}

func (d *Diagnostics) AddAttributeError(path path.Path, summary, detail string) {
	(*diag.Diagnostics)(d).AddAttributeError(path, summary, detail)
}

func (d *Diagnostics) AddLocalAttributeError(path path.Path, summary, detail string) {
	d.Append(NewLocalAttributeErrorDiagnostic(path, summary, detail))
}

func (d *Diagnostics) AddWarning(summary, detail string) {
	d.Append(NewWarningDiagnostic(summary, detail))
}

func (d *Diagnostics) HasError() bool {
	for _, cur := range *d {
		if cur.Severity() == diag.SeverityError {
			return true
		}

		if _, ok := cur.(LocalErrorDiagnostic); ok {
			return true
		}
	}

	return false
}

func (d *Diagnostics) Errors() Diagnostics {
	return Diagnostics((*diag.Diagnostics)(d).Errors())
}

func (d *Diagnostics) AppendErrorReturn(summary, detail string) Diagnostics {
	return d.AppendReturn(diag.NewErrorDiagnostic(summary, detail))
}

func (d *Diagnostics) AppendLocalErrorReturn(summary, detail string) Diagnostics {
	return d.AppendReturn(NewLocalErrorDiagnostic(summary, detail))
}

func (d *Diagnostics) AppendWarningReturn(summary, detail string) Diagnostics {
	return d.AppendReturn(diag.NewWarningDiagnostic(summary, detail))
}

func (d *Diagnostics) AppendReturn(toAppend ...diag.Diagnostic) Diagnostics {
	d.Append(toAppend...)
	return *d
}

func (d *Diagnostics) AppendCheckError(toAppend ...diag.Diagnostic) bool {
	d.Append(toAppend...)
	return d.HasError()
}

func (d Diagnostics) DemoteErrorsToWarnings() Diagnostics {
	demoted := make(Diagnostics, len(d))
	for i, d := range d {
		if d.Severity() == diag.SeverityError {
			demoted[i] = NewWarningDiagnostic(d.Summary(), d.Detail())
		} else {
			demoted[i] = d
		}
	}
	return demoted
}

func NewErrorDiagnostic(summary, detail string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(summary, detail)
}

func NewWarningDiagnostic(summary, detail string) diag.Diagnostic {
	return diag.NewWarningDiagnostic(summary, detail)
}

func NewErrorDiagnostics(summary, detail string) Diagnostics {
	return Diagnostics{NewErrorDiagnostic(summary, detail)}
}

func NewWarningDiagnostics(summary, detail string) Diagnostics {
	return Diagnostics{NewWarningDiagnostic(summary, detail)}
}

// Local Errors are sent to Terraform as Warnings (see 'Severity()'), but our custom
// 'Diagnostics' type returns 'true' in 'HasErrors()' if there's any present. This helps
// us avoid the server state being out of sync with the local state
type LocalErrorDiagnostic struct {
	summary string
	detail  string
}

func (e LocalErrorDiagnostic) Severity() diag.Severity {
	return diag.SeverityWarning
}

func (e LocalErrorDiagnostic) Summary() string {
	return e.summary
}

func (e LocalErrorDiagnostic) Detail() string {
	return e.detail
}

func (e LocalErrorDiagnostic) Equal(other diag.Diagnostic) bool {
	led, ok := other.(LocalErrorDiagnostic)
	if !ok {
		return false
	}

	return led.Summary() == e.Summary() && led.Detail() == e.Detail()
}

var _ diag.Diagnostic = (*LocalErrorDiagnostic)(nil)

func NewLocalErrorDiagnostic(summary, detail string) diag.Diagnostic {
	return LocalErrorDiagnostic{summary, detail}
}

func NewLocalErrorDiagnostics(summary, detail string) Diagnostics {
	return Diagnostics{LocalErrorDiagnostic{summary, detail}}
}

type LocalAttributeErrorDiagnostic struct {
	LocalErrorDiagnostic
	path path.Path
}

func (e *LocalAttributeErrorDiagnostic) Path() path.Path {
	return e.path
}

var _ diag.DiagnosticWithPath = (*LocalAttributeErrorDiagnostic)(nil)

func NewLocalAttributeErrorDiagnostic(path path.Path, summary, detail string) diag.Diagnostic {
	return LocalAttributeErrorDiagnostic{LocalErrorDiagnostic{summary, detail}, path}
}

func NewLocalAttributeErrorDiagnostics(path path.Path, summary, detail string) Diagnostics {
	return Diagnostics{LocalAttributeErrorDiagnostic{LocalErrorDiagnostic{summary, detail}, path}}
}

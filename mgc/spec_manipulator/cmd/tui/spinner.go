package tui

import (
	"github.com/pterm/pterm"
)

type Spinner struct {
	spinner *pterm.SpinnerPrinter
	multi   pterm.MultiPrinter
}

func NewSpinner() *Spinner {
	multi := pterm.DefaultMultiPrinter
	spinner := pterm.DefaultSpinner.WithWriter(multi.NewWriter())
	return &Spinner{
		spinner: spinner,
		multi:   multi,
	}
}

func (s *Spinner) Start(message string) {
	s.spinner.Start(message)
	s.multi.Start()
}

func (s *Spinner) UpdateText(text string) {
	s.spinner.UpdateText(text)
}

func (s *Spinner) Success(text string) {
	s.spinner.Success(text)
	s.multi.Stop()
}

func (s *Spinner) Fail(err error) {
	s.spinner.Fail(err)
	s.multi.Stop()
}

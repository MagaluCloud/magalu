package progress_report

import (
	"context"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type unitsProgressReport struct {
	units uint64
	total uint64
	err   error
}

type UnitsReporter struct {
	name           string
	total          uint64
	reportProgress ReportProgress
	reportChan     chan unitsProgressReport
}

func NewUnitsReporter(
	ctx context.Context,
	name string,
	total uint64,
) *UnitsReporter {
	return &UnitsReporter{
		name:           name,
		total:          total,
		reportProgress: FromContext(ctx),
	}
}

func (r *UnitsReporter) Start() {
	r.End()
	r.reportChan = make(chan unitsProgressReport)
	go r.progressReportSubroutine()
}

// Report the amount of new units progressed, if any, an update to the total value, if any, and an error, if any.
// Nil-pointer safe
func (r *UnitsReporter) Report(units uint64, total uint64, err error) {
	if r == nil {
		return
	}
	r.reportChan <- unitsProgressReport{units: units, total: total, err: err}
}

func (r *UnitsReporter) End() {
	if r.reportChan == nil {
		return
	}

	select {
	case <-r.reportChan:
		break
	default:
		close(r.reportChan)
	}
}

func (r *UnitsReporter) progressReportSubroutine() {
	total := r.total
	progress := uint64(0)

	// Report we're starting progress
	r.reportProgress(r.name, progress, 1, UnitsNone, nil)

	var errors utils.MultiError
	for report := range r.reportChan {
		progress += report.units
		total += report.total

		if report.err != nil {
			errors = append(errors, report.err)
		}

		r.reportProgress(r.name, progress, total, UnitsNone, nil)
	}

	if len(errors) > 0 {
		r.reportProgress(r.name, progress, total, UnitsNone, errors)
		return
	}

	r.reportProgress(r.name, total, total, UnitsNone, ErrorProgressDone)
}

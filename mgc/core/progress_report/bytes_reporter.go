package progress_report

import (
	"context"
	"errors"
	"io"
)

const bytesProgressReporterKey contextKey = "magalu.cli/core/progressreport/bytesreporter"

type bytesProgressReport struct {
	bytes uint64
	size  uint64
	err   error
}

type BytesReporter struct {
	name           string
	size           uint64
	reportProgress ReportProgress
	reportChan     chan bytesProgressReport
}

func NewBytesReporter(
	ctx context.Context,
	name string,
	size uint64,
) *BytesReporter {
	return &BytesReporter{
		name:           name,
		size:           size,
		reportProgress: FromContext(ctx),
	}
}

func NewBytesReporterContext(
	ctx context.Context,
	reporter *BytesReporter,
) context.Context {
	return context.WithValue(ctx, bytesProgressReporterKey, reporter)
}

func BytesReporterFromContext(ctx context.Context) *BytesReporter {
	if reporter, ok := ctx.Value(bytesProgressReporterKey).(*BytesReporter); ok {
		return reporter
	}
	return nil
}

func (r *BytesReporter) Start() {
	r.End()
	r.reportChan = make(chan bytesProgressReport)
	go r.progressReportSubroutine()
}

// Report the amount of new bytes progressed, if any, and an error, if any.
// Nil-pointer safe
func (r *BytesReporter) Report(bytes uint64, err error) {
	if r == nil {
		return
	}
	r.reportChan <- bytesProgressReport{bytes: bytes, err: err}
}

func (r *BytesReporter) IncreaseTotal(size uint64, err error) {
	if r == nil {
		return
	}
	r.reportChan <- bytesProgressReport{size: size, err: err}
}

func (r *BytesReporter) End() {
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

func (r *BytesReporter) progressReportSubroutine() {
	bytesDone := uint64(0)
	total := r.size

	// Report we're starting progress
	r.reportProgress(r.name, bytesDone, r.size, UnitsBytes, nil)

	var err error

	for report := range r.reportChan {
		bytesDone += report.bytes
		total += report.size
		if report.err != nil && !errors.Is(report.err, io.EOF) {
			err = report.err
		}
		r.reportProgress(r.name, bytesDone, total, UnitsBytes, nil)
	}
	// Set DONE flag
	if err == nil {
		err = ErrorProgressDone
	}

	r.reportProgress(r.name, bytesDone, r.size, UnitsBytes, err)
}

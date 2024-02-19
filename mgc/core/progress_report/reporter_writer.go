package progress_report

import (
	"io"
)

type ReportWrite func(n uint64, err error)

type reporterWriter struct {
	parent         io.Writer
	reportProgress ReportWrite
}

func NewReporterWriter(parent io.Writer, reportProgress ReportWrite) *reporterWriter {
	return &reporterWriter{
		parent:         parent,
		reportProgress: reportProgress,
	}
}

func (rw *reporterWriter) Unwrap() io.Writer {
	return rw.parent
}

func (rw *reporterWriter) Write(p []byte) (n int, err error) {
	n, err = rw.parent.Write(p)
	rw.reportProgress(uint64(n), err)
	return
}

func (rw *reporterWriter) Close() error {
	if closer, ok := rw.parent.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

var _ io.WriteCloser = (*reporterWriter)(nil)

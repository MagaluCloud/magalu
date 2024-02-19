package progress_report

import (
	"io"
)

type ReportRead func(n uint64, err error)

type reporterReader struct {
	parent         io.Reader
	reportProgress ReportRead
}

// Wraps an io.Reader in another Reader which reports the amount of bytes read anytime
// the parent is read
func NewReporterReader(parent io.Reader, reportProgress ReportRead) *reporterReader {
	return &reporterReader{
		parent:         parent,
		reportProgress: reportProgress,
	}
}

func (pr *reporterReader) Unwrap() io.Reader {
	return pr.parent
}

// BEGIN io.ReadCloser implementation

func (pr *reporterReader) Read(p []byte) (n int, err error) {
	n, err = pr.parent.Read(p)
	pr.reportProgress(uint64(n), err)
	return
}

func (pr *reporterReader) Close() error {
	if closer, ok := pr.parent.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

var _ io.ReadCloser = (*reporterReader)(nil)

// END io.ReadCloser implementation

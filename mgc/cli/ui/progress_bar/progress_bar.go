package progress_bar

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core/progress_report"
	"github.com/jedib0t/go-pretty/v6/progress"
)

func parseUnits(units progress_report.Units) progress.Units {
	switch units {
	case progress_report.UnitsNone:
		return progress.UnitsDefault
	case progress_report.UnitsBytes:
		return progress.UnitsBytes
	default:
		return progress.UnitsDefault
	}
}

type ProgressBar struct {
	progress.Writer
	trackers sync.Map
}

var updateFrequency = progress.DefaultUpdateFrequency

func New() *ProgressBar {
	writer := progress.NewWriter()
	writer.SetAutoStop(true)
	writer.SetUpdateFrequency(updateFrequency)
	writer.SetMessageWidth(60)
	writer.SetTrackerPosition(progress.PositionRight)
	writer.SetTrackerLength(progress.DefaultLengthTracker)
	writer.SetSortBy(progress.SortByPercent)
	return &ProgressBar{
		Writer: writer,
	}
}

// TODO: sometimes the progress bar does not render the final update.
// Investigate why p.Stop() misses this last render
func (pb *ProgressBar) Finalize() {
	if pb.Length() > 0 {
		time.Sleep(updateFrequency)
	}
	pb.Stop()
}

func (pb *ProgressBar) Flush() {
	pb.Finalize()
	go pb.Render()
}

func (pb *ProgressBar) ReportProgress(msg string, done, total uint64, units progress_report.Units, reportErr error) {
	tracker, found := pb.trackers.LoadOrStore(msg,
		&progress.Tracker{
			Message: msg,
			Total:   int64(total),
			Units:   parseUnits(units),
		},
	)
	castTracker := tracker.(*progress.Tracker)
	if !found {
		pb.Writer.AppendTracker(castTracker)
	}

	castTracker.UpdateTotal(int64(total))
	castTracker.SetValue(int64(done))
	if reportErr != nil && !errors.Is(reportErr, io.EOF) {
		if errors.Is(reportErr, progress_report.ErrorProgressDone) {
			castTracker.MarkAsDone()
			return
		}
		// Accept error reports even after report is done
		if castTracker.IsDone() {
			castTracker.Reset()
			castTracker.Start()
		}
		castTracker.MarkAsErrored()
		castTracker.UpdateMessage(fmt.Sprintf("%s [%s]", msg, reportErr.Error()))
	}
}

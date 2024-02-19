package pipeline

import (
	"context"
	"fmt"
	"io"
)

type ChunkWriter interface {
	io.Writer
	io.WriterAt
}

type WriteableChunk struct {
	Writer      ChunkWriter
	StartOffset int64
	EndOffset   int64
}

func PrepareWriteChunks(
	ctx context.Context,
	w io.WriterAt,
	size int64,
	chunkSize int64,
) (outputChan <-chan WriteableChunk) {
	ch := make(chan WriteableChunk)
	outputChan = ch

	logger := FromContext(ctx).Named("PrepareWriteChunks").With(
		"size", size,
		"chunkSize", chunkSize,
		"outputChan", fmt.Sprintf("%#v", outputChan),
	)
	ctx = NewContext(ctx, logger)

	generator := func() {
		defer func() {
			logger.Info("closing output channel")
			close(ch)
		}()

		var i int64
		for i = 0; i < size; i += chunkSize {
			// Make sure that writing range is accurate to chunk's size and file size
			end := i + chunkSize - 1
			if end >= size {
				end = size - 1
			}

			select {
			case <-ctx.Done():
				logger.Debugw("context.Done()", "err", ctx.Err())
				return

			case ch <- WriteableChunk{io.NewOffsetWriter(w, i), i, end}:
				logger.Debugw("prepared section for write", "start", i, "end", end)
			}
		}
		logger.Debug("finished preparing sections for write")
	}

	logger.Info("start")
	go generator()
	return
}

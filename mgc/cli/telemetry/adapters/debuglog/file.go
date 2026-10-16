package debuglog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

const (
	FileName = "telemetry-debug.txt"

	maxDebugLogBytes = 1 << 20 // 1 MiB
	fileMode         = 0o600
	dirMode          = 0o755
	timestampLayout  = "2006-01-02T15:04:05.000Z07:00"
)

// New devolve um logger que acrescenta uma linha por chamada no arquivo de path.
// O arquivo é zerado antes de uma gravação que o faria passar de 1 MiB, e falhas
// de escrita são ignoradas porque o log nunca pode atrapalhar o envio
func New(path string, now func() time.Time, pid int) telemetry.DebugLogger {
	return func(msg string, keysAndValues ...any) {
		_ = appendLine(path, formatLine(now().UTC(), pid, msg, keysAndValues))
	}
}

// Discard é o logger do processo de envio fora do modo debug
func Discard(string, ...any) {}

func formatLine(at time.Time, pid int, msg string, keysAndValues []any) string {
	var line strings.Builder
	fmt.Fprintf(&line, "%s pid=%d %s", at.Format(timestampLayout), pid, msg)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 == len(keysAndValues) {
			fmt.Fprintf(&line, " %s", formatValue(keysAndValues[i]))
			break
		}
		fmt.Fprintf(&line, " %v=%s", keysAndValues[i], formatValue(keysAndValues[i+1]))
	}
	line.WriteByte('\n')
	return line.String()
}

func formatValue(value any) string {
	text := fmt.Sprint(value)
	if text == "" || strings.ContainsAny(text, " \t\n\"=") {
		return strconv.Quote(text)
	}
	return text
}

func appendLine(path, line string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirMode); err != nil {
		return err
	}

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if info, err := os.Stat(path); err == nil && info.Size()+int64(len(line)) > maxDebugLogBytes {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, fileMode)
	if err != nil {
		return err
	}
	_, err = file.WriteString(line)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}

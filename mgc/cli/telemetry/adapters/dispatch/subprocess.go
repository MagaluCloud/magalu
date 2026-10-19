package dispatch

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
)

const (
	payloadPattern  = "mgc-telemetry-*.json"
	payloadFileMode = 0o600
)

// Subprocess entrega o evento a um processo filho, que é o próprio binário da CLI
// chamado com SenderArg. A CLI não espera esse processo terminar
type Subprocess struct {
	executable func() (string, error)
	start      func(*exec.Cmd) error
	tempDir    func() string
	// debugLogPath devolve o caminho do log de debug do processo de envio, ou vazio
	// fora do modo debug. É uma função porque o modo debug só é conhecido depois
	// que as flags do comando foram lidas
	debugLogPath func() string
	debug        telemetry.DebugLogger
}

var _ telemetry.Dispatcher = (*Subprocess)(nil)

func NewSubprocess(debugLogPath func() string, debug telemetry.DebugLogger) *Subprocess {
	if debugLogPath == nil {
		debugLogPath = func() string { return "" }
	}
	if debug == nil {
		debug = func(string, ...any) {}
	}
	return &Subprocess{
		executable:   os.Executable,
		start:        (*exec.Cmd).Start,
		tempDir:      os.TempDir,
		debugLogPath: debugLogPath,
		debug:        debug,
	}
}

func (s *Subprocess) Dispatch(event telemetry.Event) error {
	executable, err := s.executable()
	if err != nil {
		return fmt.Errorf("locate the cli executable: %w", err)
	}

	debugLogPath := s.debugLogPath()
	data, err := Encode(event, debugLogPath)
	if err != nil {
		return err
	}

	path, err := writePayload(s.tempDir(), data)
	if err != nil {
		return err
	}

	// Stdin, Stdout e Stderr ficam nil, então o processo de envio nunca escreve no terminal
	cmd := exec.Command(executable, SenderArg, path)
	cmd.Dir = s.tempDir()

	// sysProcAttr vem de spawn_unix.go ou spawn_windows.go, escolhido pelas build constraints do Go
	// https://pkg.go.dev/cmd/go#hdr-Build_constraints
	// https://github.com/golang/go/blob/master/src/internal/syslist/syslist.go
	cmd.SysProcAttr = sysProcAttr()

	if err := s.start(cmd); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("start the telemetry sender: %w", err)
	}

	pid := 0
	if cmd.Process != nil {
		pid = cmd.Process.Pid
		_ = cmd.Process.Release()
	}
	s.debug("telemetry: dispatched", dispatchedFields(pid, len(data), debugLogPath)...)
	return nil
}

func dispatchedFields(pid, size int, debugLogPath string) []any {
	kv := []any{"pid", pid, "bytes", size}
	if debugLogPath != "" {
		kv = append(kv, "debug_log", debugLogPath)
	}
	return kv
}

func writePayload(dir string, data []byte) (string, error) {
	file, err := os.CreateTemp(dir, payloadPattern)
	if err != nil {
		return "", err
	}
	path := file.Name()

	err = file.Chmod(payloadFileMode)
	if err == nil {
		_, err = file.Write(data)
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return path, nil
}

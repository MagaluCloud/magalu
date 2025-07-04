package http

import (
	"encoding/json"
	"fmt"
	"os"
)

type LogSensitive string

const logSensitiveEnvVar = "MGC_LOG_SENSITIVE"

var shouldLogSensitiveStatus *bool

func shouldLogSensitive() bool {
	if shouldLogSensitiveStatus == nil {
		b := os.Getenv(logSensitiveEnvVar) == "1"
		shouldLogSensitiveStatus = &b
	}
	return *shouldLogSensitiveStatus
}

func (s LogSensitive) MarshalJSON() ([]byte, error) {
	var text string
	if shouldLogSensitive() {
		text = string(s)
	} else {
		text = fmt.Sprintf("[REDACTED %d CHARS]", len(s))
	}

	return json.Marshal(text)
}

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogHttpHeaders http.Header

func isHeaderSensitive(canonicalKey string) bool {
	switch canonicalKey {
	case "Authorization":
		return true
	case "X-Api-Key":
		return true
	default:
		return false
	}
}

func (h LogHttpHeaders) MarshalJSON() ([]byte, error) {
	logSensitive := shouldLogSensitive()
	b := bytes.Buffer{}
	b.WriteByte('{')

	for key, list := range h {
		valueListLength := len(list)

		if valueListLength == 0 {
			continue
		}

		if b.Len() > 1 {
			b.WriteByte(',')
		}

		s, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		b.Write(s)
		b.WriteByte(':')

		if !logSensitive && isHeaderSensitive(key) {
			if valueListLength == 1 {
				s = ([]byte)(fmt.Sprintf(`"[REDACTED %d CHARS]"`, len(list[0])))
			} else {
				s = ([]byte)(fmt.Sprintf(`"[REDACTED %d ENTRIES]"`, valueListLength))
			}
		} else {
			if valueListLength == 1 {
				s, err = json.Marshal(list[0])
			} else {
				s, err = json.Marshal(list)
			}
		}
		if err != nil {
			return nil, err
		}
		b.Write(s)
	}

	b.WriteByte('}')
	return b.Bytes(), nil
}

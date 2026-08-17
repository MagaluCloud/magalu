package common

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

var candidateEncodings = []*charmap.Charmap{
	charmap.Windows1252,
	charmap.Windows1250,
	charmap.Windows1251,
	charmap.Windows1253,
	charmap.Windows1254,
	charmap.Windows1257,
	charmap.ISO8859_1,
	charmap.CodePage850,
}

func FixFilenameEncoding(name string) (fixed string, converted bool, err error) {
	if utf8.ValidString(name) {
		return name, false, nil
	}

	for _, enc := range candidateEncodings {
		decoded, decErr := enc.NewDecoder().String(name)
		if decErr == nil && utf8.ValidString(decoded) && !hasUndecodableRune(decoded) {
			return decoded, true, nil
		}
	}

	pos, b := firstInvalidByte(name)
	return "", false, fmt.Errorf("filename contains invalid UTF-8 bytes (position %d, byte 0x%02X) — rename the file and try again", pos, b)
}

func hasUndecodableRune(s string) bool {
	for _, r := range s {
		if r == utf8.RuneError || unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func firstInvalidByte(s string) (int, byte) {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size <= 1 {
			return i, s[i]
		}
		i += size
	}
	return -1, 0
}

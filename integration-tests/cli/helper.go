package cli

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateNanoIDWithPrefix(prefix string) string {
	id, _ := gonanoid.New()
	return prefix + "-" + id
}

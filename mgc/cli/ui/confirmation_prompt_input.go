package ui

import (
	"github.com/erikgeiser/promptkit/textinput"
)

func RunPromptInput(message string) (string, error) {
	input := textinput.New(message)
	ready, err := input.RunPrompt()
	if err != nil {
		return "", err
	}
	return ready, nil
}

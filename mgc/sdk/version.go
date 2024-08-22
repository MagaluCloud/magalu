package sdk

import (
	"os/exec"
	"strings"
)

var version string = func() string {
	tag, err := getLatestGitTag()
	if err != nil {
		return "unknown"
	}
	return strings.Trim(tag, " \t\n\r")
}()

func getLatestGitTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

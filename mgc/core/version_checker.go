package core

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	githubReleases   = "https://github.com/MagaluCloud/mgccli/releases/latest"
	updateInterval   = 2 * time.Hour
	versionLastCheck = "version_last_check"
)

type configValue func(key string, value interface{}) error
type getCall func(url string) (resp *http.Response, err error)

type VersionChecker struct {
	getCall   getCall
	getConfig configValue
	setConfig configValue
}

// NewVersionChecker creates a new VersionChecker
// getHttp - the getCall client to use for requests
// getConfig - a function to get a config value
// setConfig - a function to set a config value
func NewVersionChecker(
	getHttp getCall, getConfig, setConfig configValue,
) *VersionChecker {
	return &VersionChecker{getCall: getHttp, getConfig: getConfig, setConfig: setConfig}
}

// CheckVersion checks if the current version is outdated
// and prints a message if it is.
// It also sets the last check time to the current time.
// currentVersion - the current version of the cli
// args - optional command line arguments
func (v *VersionChecker) CheckVersion(currentVersion string, args ...string) {
	logger().Debug("Checking for updates")
	lastCheckTime := v.getLastCheckTime()

	isVersionCmd := len(args) > 0 && slices.Contains(args, "--version")

	if !isVersionCmd && time.Since(lastCheckTime) < updateInterval {
		logger().Debug("Skipping update check")
		return
	}

	latestVersion, err := v.getLatestVersion()
	if err != nil {
		logger().Debugw("cannot getConfig latest version", "err", err)
		return
	}

	cv, err := semver.NewVersion(strings.SplitN(currentVersion, " ", 2)[0])
	if err != nil {
		fmt.Println("Invalid current version:", err)
		return
	}

	latestSemVersion, err := semver.NewVersion(latestVersion)
	if err != nil {
		logger().Debugw("cannot parse latest version", "err", err)
		return
	}

	if cv.LessThan(latestSemVersion) {
		v.setCurrentTime()
		fmt.Printf(
			"⚠️ You are using an outdated version of mgc cli. "+
				"Please update to the latest version: %s \n\n\n", latestVersion,
		)
		return
	}
	logger().Debug("No updates available")
	v.setCurrentTime()
}

func (v *VersionChecker) getLastCheckTime() time.Time {
	var lastCheck string
	_ = v.getConfig(versionLastCheck, &lastCheck)
	lastCheckTime, err := time.Parse(time.RFC3339, lastCheck)
	if err != nil {
		v.setCurrentTime()
	}
	return lastCheckTime
}

func (v *VersionChecker) getLatestVersion() (string, error) {
	response, err := v.getCall(githubReleases)

	if err != nil {
		return "", err
	}
	location := response.Header.Get("Location")
	if location == "" {
		return "", errors.New("no location header")
	}
	return path.Base(location), nil
}

func (v *VersionChecker) setCurrentTime() {
	err := v.setConfig(versionLastCheck, time.Now().Format(time.RFC3339))
	if err != nil {
		logger().Debugw("cannot set last check time", "err", err)
	}
}

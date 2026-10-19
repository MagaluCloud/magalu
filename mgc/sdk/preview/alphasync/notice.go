package alphasync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	stateFile      = "state.json"
	noticeInterval = 24 * time.Hour
	noticeTimeout  = 3 * time.Second
)

type state struct {
	LastCheck time.Time `json:"last_check"`
}

// CheckUpdates tells, at most once a day, which installed modules have a newer
// version or were taken away. It never changes what is installed, and says
// nothing when it cannot tell. A module pinned to a version is only told about
// when it, or its version, is taken away. The client of cfg must not retry: the
// retries would outlast the deadline of the check.
func CheckUpdates(ctx context.Context, cfg Config) {
	if err := checkUpdates(ctx, cfg, time.Now(), os.Stderr); err != nil {
		logger().Debugw("checking the alpha modules for updates failed", "error", err)
	}
}

func checkUpdates(ctx context.Context, cfg Config, now time.Time, out io.Writer) error {
	index, err := readIndex(cfg.Dir)
	if err != nil || len(index) == 0 || cfg.BaseURL == "" {
		return err
	}
	if now.Sub(lastCheck(cfg.Dir)) < noticeInterval {
		return nil
	}
	// Saved before checking, so that a source failing is not tried again on
	// every command.
	data, err := json.Marshal(state{LastCheck: now})
	if err != nil {
		return err
	}
	if err := writeFile(cfg.Dir, stateFile, data); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, noticeTimeout)
	defer cancel()
	manifests := manifestCache(ctx, cfg)

	var updates, needNewerCLI []string
	for _, module := range slices.Sorted(maps.Keys(index)) {
		entry := index[module]
		manifest, err := manifests(entry.Manifest)
		if err != nil {
			return err
		}
		newest, listed := manifest.Find(module, "")
		change := fmt.Sprintf("%s %s -> %s", module, entry.Version, newest.Version)
		switch {
		case !listed:
			updates = append(updates, module+" was taken away")
		case takenAway(manifest, entry):
			updates = append(updates, entry.label()+" was taken away")
		case entry.Pinned || newest.Version == entry.Version:
		case tooOld(cfg.CLIVersion, newest.MinCLIVersion) != "":
			needNewerCLI = append(needNewerCLI, change)
		default:
			updates = append(updates, change)
		}
	}
	if len(updates) > 0 {
		fmt.Fprintf(out, "Alpha updates: %s. Run 'mgc preview upgrade'.\n", strings.Join(updates, ", "))
	}
	if len(needNewerCLI) > 0 {
		fmt.Fprintf(out, "Alpha updates that need a newer mgc: %s. Update mgc first.\n", strings.Join(needNewerCLI, ", "))
	}
	return nil
}

func lastCheck(dir string) time.Time {
	var saved state
	if data, err := os.ReadFile(filepath.Join(dir, stateFile)); err == nil {
		_ = json.Unmarshal(data, &saved)
	}
	return saved.LastCheck
}

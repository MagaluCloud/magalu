package alphasync

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/core"
)

type syncParams struct {
	Invite string `json:"invite,omitempty" jsonschema:"description=The invite Magalu sent you. Defaults to the one of the last sync,example=3f2a9c1e-5b7d-4e8a-9c0f-1a2b3c4d5e6f" mgc:"positional"`
}

func NewCommand(cfg Config) core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "sync",
			Summary: "Download the alpha commands your invite gives access to",
			Description: `Downloads the alpha commands your invite gives access to and replaces the
installed ones, so a command your invite lost access to is removed`,
			Observations: "Pass your invite once: later syncs reuse it.",
		},
		func(ctx context.Context, params syncParams, _ struct{}) (*Result, error) {
			if cfg.BaseURL == "" {
				return nil, errNoSource
			}
			invite, err := resolveInvite(cfg.Dir, params.Invite)
			if err != nil {
				return nil, err
			}

			result, err := Sync(ctx, cfg, invite)
			if err != nil {
				return nil, err
			}
			for _, skipped := range result.Skipped {
				fmt.Fprintf(os.Stderr, "Warning: alpha module %q was not installed: %s\n", skipped.Module, skipped.Reason)
			}
			return result, nil
		},
	)
}

func resolveInvite(dir, given string) (string, error) {
	if given != "" {
		return given, nil
	}

	saved, err := SavedInvite(dir)
	if err != nil {
		return "", err
	}
	if saved == "" {
		return "", errors.New("Please include your invite in the command: mgc alpha sync [invite]")
	}
	return saved, nil
}

package object_lock

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type setObjectLockParams struct {
	Object          mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Specifies the object whose lock is being requested" mgc:"positional"`
	RetainUntilDate string           `json:"retain_until_date" jsonschema:"description=Timestamp in ISO 8601 format,example=2025-10-03T00:00:00"`
	Mode            string           `json:"mode,omitempty" jsonschema:"description=Lock mode,enum=COMPLIANCE,enum=GOVERNANCE,default=COMPLIANCE,required" mgc:"hidden"`
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "set number of either days or years to lock new objects for",
		},
		setObjectLocking,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set Object Locking for object %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setObjectLocking(ctx context.Context, params setObjectLockParams, cfg common.Config) (result core.Value, err error) {
	req, err := newSetObjectLockingRequest(ctx, params, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func parseISODate(dateStr string) (string, error) {
	formats := []string{
		time.RFC3339,                       // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05",              // "2006-01-02T15:04:05"
		"2006-01-02",                       // "2006-01-02"
		"2006-01-02 15:04:05",              // "2006-01-02 15:04:05"
		"2006-01-02T15:04:05.000000Z07:00", // "2006-01-02T15:04:05.000000Z07:00"
		"2006-01-02T15:04:05Z",             // "2006-01-02T15:04:05Z"
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date.UTC().Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("invalid date format: %s", dateStr)
}

func newSetObjectLockingRequest(ctx context.Context, p setObjectLockParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPath(cfg, common.NewBucketNameFromURI(p.Object), p.Object.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(url), nil)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := req.URL.Query()
	query.Add("retention", "")
	req.URL.RawQuery = query.Encode()

	getBody := func() (io.ReadCloser, error) {
		parsedTimeStr, err := parseISODate(p.RetainUntilDate)
		if err != nil {
			return nil, core.UsageError{Err: err}
		}
		parsedTime, err := time.Parse(time.RFC3339, parsedTimeStr)
		if err != nil {
			return nil, core.UsageError{Err: err}
		}

		bodyObj := common.DefaultObjectRetentionBody(parsedTime)
		if p.Mode == string(common.ObjectLockModeGovernance) {
			bodyObj.Mode = common.ObjectLockModeGovernance
		}

		body, err := xml.MarshalIndent(bodyObj, "", "  ")
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(body)
		return io.NopCloser(reader), nil
	}

	req.Body, err = getBody()
	if err != nil {
		return nil, err
	}
	req.GetBody = getBody

	return req, nil
}

package object_lock

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type objectLockMode string

const (
	objectLockModeGovernance = objectLockMode("GOVERNANCE")
	objectLockModeCompliance = objectLockMode("COMPLIANCE")
)

type objectLockRetention struct {
	XMLName         xml.Name `xml:"Retention"`
	Namespace       string   `xml:"xmlns,attr"`
	Mode            objectLockMode
	RetainUntilDate string `xml:",omitempty"`
}

func defaultObjectLockingBody(retainUntilDate time.Time) objectLockRetention {
	return objectLockRetention{
		Namespace:       "http://s3.amazonaws.com/doc/2006-03-01/",
		Mode:            objectLockModeCompliance,
		RetainUntilDate: retainUntilDate.UTC().Format("2006-01-02T15:04:05"),
	}

}

type setObjectLockParams struct {
	Object          common.ObjectError `json:"dst" jsonschema:"description=Specifies the object whose lock is being requested" mgc:"positional"`
	RetainUntilDate string             `json:"retain_until_date,omitempty" jsonschema:"description=Timestamp in ISO 8601 format,example=2025-10-03T00:00:00Z"`
	Days            int                `json:"days,omitempty" jsonschema:"description=Number of days to retain the object"`
	Years           int                `json:"years,omitempty" jsonschema:"description=Number of years to retain the object"`
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
	if params.RetainUntilDate == "" && params.Days == 0 && params.Years == 0 {
		return nil, fmt.Errorf("You must provide one of the following: `retain_until_date`, `days`, or `years`")
	}

	if (params.RetainUntilDate != "" && (params.Days > 0 || params.Years > 0)) ||
		(params.Days > 0 && params.Years > 0) {
		return nil, fmt.Errorf("You must provide only one of the following: `retain_until_date`, `days`, or `years`, not multiple")
	}

	req, err := newSetObjectLockingRequest(ctx, params, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func newSetObjectLockingRequest(ctx context.Context, p setObjectLockParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostWithPath(cfg, common.NewBucketNameFromURI(p.Object.Url), p.Object.Url.Path())
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(url), nil)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := req.URL.Query()
	query.Add("object-lock", "")
	req.URL.RawQuery = query.Encode()

	getBody := func() (io.ReadCloser, error) {
		var parsedTime time.Time

		if p.RetainUntilDate != "" {
			parsedTime, err = time.Parse("2006-01-02T15:04:05", p.RetainUntilDate)
			if err != nil {
				return nil, core.UsageError{Err: err}
			}
		} else {
			parsedTime = time.Now()
			if p.Days > 0 {
				parsedTime = parsedTime.AddDate(0, 0, p.Days)
			}
			if p.Years > 0 {
				parsedTime = parsedTime.AddDate(p.Years, 0, 0)
			}
		}
		bodyObj := defaultObjectLockingBody(parsedTime)
		body, err := xml.MarshalIndent(bodyObj, "", "  ")
		fmt.Println(string(body))
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

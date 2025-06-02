package cors

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type CORSConfiguration struct {
	XMLName   xml.Name   `xml:"CORSConfiguration"`
	XMLns    string     `xml:"xmlns,attr"`
    CORSRules []CORSRule `xml:"CORSRule"`
}

type CORSRule struct {
	AllowedOrigins []string `xml:"AllowedOrigin"`
	AllowedMethods []string `xml:"AllowedMethod"`
	AllowedHeaders []string `xml:"AllowedHeader,omitempty"`
	ExposeHeaders  []string `xml:"ExposeHeader,omitempty"`
	MaxAgeSeconds  int      `xml:"MaxAgeSeconds,omitempty"`
}

type setBucketCorsParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to set permissions for,example=my-bucket" mgc:"positional"`
	Cors   map[string]any    `json:"cors" jsonschema:"description=CORS file path to be uploaded,example=@./cors.json or ./cors.json" mgc:"positional"`
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set CORS document for the specified bucket. The CORS can be provided as a direct JSON string or a file path using @./cors.json.",
		},
		setCors,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set CORS for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setCors(ctx context.Context, params setBucketCorsParams, cfg common.Config) (result core.Value, err error) {
	req, err := newSetBucketCorsRequest(ctx, params, cfg)
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

func newSetBucketCorsRequest(ctx context.Context, p setBucketCorsParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("cors", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	getBody := func() (io.ReadCloser, error) {
		jsonBytes, err := json.Marshal(p.Cors)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal CORS input: %w", err)
		}

		var corsConfig CORSConfiguration
		corsConfig.XMLns = "http://s3.amazonaws.com/doc/2006-03-01/"

		err = json.Unmarshal(jsonBytes, &corsConfig)
		if err != nil {
			return nil, fmt.Errorf("invalid CORS JSON: %w", err)
		}

		xmlBytes, err := xml.Marshal(corsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to XML: %w", err)
		}

		//xmlWithHeader := append([]byte(xml.Header), xmlBytes...)

		reader := bytes.NewReader(xmlBytes)
		return io.NopCloser(reader), nil
	}

	req.Body, err = getBody()
	if err != nil {
		return nil, err
	}
	req.GetBody = getBody

	return req, nil
}

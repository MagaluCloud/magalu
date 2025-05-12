package label

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type Unit struct{}

type TaggingXML struct {
	XMLName xml.Name `xml:"Tagging"`
	XMLNS   string   `xml:"xmlns,attr"`
	Tags    []Tag    `xml:"TagSet>Tag"`
}

type setBucketLabelParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to set labels for,example=my-bucket" mgc:"positional"`
	Labels string            `json:"label" jsonschema:"description=Label values, comma separated, without whitespaces" mgc:"positional"`
}

type setBucketTaggingParams struct {
	Bucket common.BucketName
	TagSet TagSet
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set labels for the specified bucket",
		},
		setLabels,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set labels for bucket %q", result.Source().Parameters["bucket"])
	})

	return exec
})

func setLabels(ctx context.Context, params setBucketLabelParams, cfg common.Config) (_ core.Value, err error) {
	res, err := getTags(ctx, GetBucketLabelParams{Bucket: params.Bucket}, cfg)
	if err != nil {
		return
	}
	labels, hasMGCLabels := findMGCLabels(res.Tags)
	if !hasMGCLabels {
		res.Tags = append(res.Tags, Tag{Key: "MGC_LABELS", Value: params.Labels})
	} else {
		savedLabels := strings.Split(labels, ",")
		inputLabels := strings.Split(params.Labels, ",")
		labelSet := make(map[string]Unit, len(inputLabels)+len(savedLabels))
		for _, label := range inputLabels {
			labelSet[strings.Trim(label, " ")] = Unit{}
		}
		for _, label := range savedLabels {
			labelSet[strings.Trim(label, " ")] = Unit{}
		}
		var buffer bytes.Buffer
		for label := range labelSet {
			buffer.WriteString(label)
			buffer.WriteString(",")
		}
		tagLabel := strings.Trim(buffer.String(), ",")
		updateMGCLabels(&res.Tags, tagLabel)
	}

	taggingParams := setBucketTaggingParams{Bucket: params.Bucket, TagSet: res}
	req, err := newSetBucketTaggingRequest(ctx, taggingParams, cfg)
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

func findMGCLabels(tags []Tag) (labels string, hasMGCLabel bool) {
	for _, tag := range tags {
		if tag.Key == "MGC_LABELS" {
			if tag.Value != "" {
				return tag.Value, true
			}
			return labels, true
		}
	}
	return labels, false
}

func updateMGCLabels(tags *[]Tag, labels string) {
	for i := range *tags {
		if (*tags)[i].Key == "MGC_LABELS" {
			(*tags)[i].Value = labels
			return
		}
	}
}

func newSetBucketTaggingRequest(ctx context.Context, p setBucketTaggingParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("tagging", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	getBody := func() (io.ReadCloser, error) {
		tagsXML := TaggingXML{Tags: p.TagSet.Tags}
		body, err := xml.Marshal(tagsXML)
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

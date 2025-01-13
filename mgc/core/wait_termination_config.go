package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type WaitTerminationConfig struct {
	MaxRetries         int           `json:"maxRetries,omitempty"`
	Interval           time.Duration `json:"interval,omitempty"`
	JSONPathQuery      string        `json:"jsonPathQuery,omitempty"`
	TemplateQuery      string        `json:"templateQuery,omitempty"`
	ErrorJSONPathQuery string        `json:"errorJsonPathQuery,omitempty"`
	ErrorTemplateQuery string        `json:"errorTemplateQuery,omitempty"`
}

var defaultWaitTermination = WaitTerminationConfig{MaxRetries: 30, Interval: time.Second}

func (c *WaitTerminationConfig) Build(exec Executor, getDocument func(result ResultWithValue) any) (tExec TerminatorExecutor, err error) {
	maxRetries := c.MaxRetries
	if maxRetries <= 0 {
		maxRetries = defaultWaitTermination.MaxRetries
	}
	interval := c.Interval
	if interval <= 0 {
		interval = defaultWaitTermination.Interval
	}

	var expChecker func(value any) (bool, error)
	if c.JSONPathQuery != "" && c.TemplateQuery != "" {
		err = errors.New("cannot specify both jsonPathQuery and templateQuery")
	} else if c.JSONPathQuery != "" {
		expChecker, err = utils.CreateJsonPathChecker(c.JSONPathQuery)
	} else if c.TemplateQuery != "" {
		expChecker, err = utils.CreateTemplateChecker(c.TemplateQuery)
	} else {
		err = errors.New("need one of jsonPathQuery or templateQuery")
	}

	var errorChecker func(value any) (bool, error)
	if c.ErrorJSONPathQuery != "" && c.TemplateQuery != "" {
		err = errors.New("cannot specify both errorJsonPathQuery and errorTemplateQuery")
	} else if c.ErrorJSONPathQuery != "" {
		errorChecker, err = utils.CreateJsonPathChecker(c.ErrorJSONPathQuery)
	} else if c.ErrorTemplateQuery != "" {
		errorChecker, err = utils.CreateTemplateChecker(c.ErrorJSONPathQuery)
	}

	if err != nil {
		return nil, err
	}

	return NewTerminatorExecutorWithCheck(exec, maxRetries, interval, func(ctx context.Context, exec Executor, result ResultWithValue) (terminated bool, err error) {
		doc := getDocument(result)

		terminated, err = expChecker(doc)
		if terminated || err != nil {
			return
		}

		if errorChecker != nil {
			terminated, err = errorChecker(doc)
			if err != nil {
				return
			}

			if terminated {
				err = fmt.Errorf(
					"wait termination for %q exited due to error condition returning true with document %#v",
					exec.Name(),
					doc,
				)
				return
			}
		}

		return
	}), nil
}

func (c *WaitTerminationConfig) UnmarshalJSON(data []byte) (err error) {
	m := map[string]any{}
	err = json.Unmarshal(data, &m) // decoding interval to time.Duration would fail
	if err != nil {
		return
	}
	return utils.DecodeValue(m, c) // nicely decodes time.Duration
}

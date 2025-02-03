package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var ErrorResultHasNoValue = errors.New("result has no value")

type RetryUntilCheck func(ctx context.Context, value Value) (finished bool, err error)

type RetryUntilConfig struct {
	MaxRetries    int           `json:"maxRetries,omitempty"`
	Interval      time.Duration `json:"interval,omitempty"`
	JSONPathQuery string        `json:"jsonPathQuery,omitempty"`
	TemplateQuery string        `json:"templateQuery,omitempty"`
}

func (c *RetryUntilConfig) Build() (r *RetryUntil, err error) {
	if c == nil {
		return nil, nil
	}

	var check RetryUntilCheck
	if c.JSONPathQuery != "" && c.TemplateQuery != "" {
		err = errors.New("cannot specify both jsonPathQuery and templateQuery")
	} else if c.JSONPathQuery != "" {
		check, err = NewRetryUntilCheckFromJsonPath(c.JSONPathQuery)
	} else if c.TemplateQuery != "" {
		check, err = NewRetryUntilCheckFromTemplate(c.TemplateQuery)
	} else {
		err = errors.New("need one of jsonPathQuery or templateQuery")
	}

	if err != nil {
		return nil, err
	}

	return &RetryUntil{
		MaxRetries: c.MaxRetries,
		Interval:   c.Interval,
		Check:      check,
	}, nil
}

func (c *RetryUntilConfig) UnmarshalJSON(data []byte) (err error) {
	m := map[string]any{}
	err = json.Unmarshal(data, &m) // decoding interval to time.Duration would fail
	if err != nil {
		return
	}
	return utils.DecodeValue(m, c) // nicely decodes time.Duration
}

var _ json.Unmarshaler = (*RetryUntilConfig)(nil)

type RetryUntil struct {
	MaxRetries int
	Interval   time.Duration
	Check      RetryUntilCheck
}

type RetryUntilCb func() (result Result, err error)

func (r *RetryUntil) Run(ctx context.Context, cb RetryUntilCb) (result Result, err error) {
	if r == nil {
		return cb()
	}

	for i := 0; i < r.MaxRetries; i++ {
		result, err = cb()
		if err != nil {
			return result, err
		}
		resultWithValue, ok := ResultAs[ResultWithValue](result)
		if !ok {
			return result, ErrorResultHasNoValue
		}
		finished, err := r.Check(ctx, resultWithValue.Value())
		if err != nil {
			return result, err
		}
		if finished {
			return result, nil
		}

		timer := time.NewTimer(r.Interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}

	msg := fmt.Sprintf("exceeded maximum retries %d with interval %s", r.MaxRetries, r.Interval)
	return nil, FailedTerminationError{Result: result, Message: msg}
}

func NewRetryUntilCheckFromJsonPath(expression string) (check RetryUntilCheck, err error) {
	jpChecker, err := utils.CreateJsonPathChecker(expression)
	if err != nil {
		return nil, err
	}

	check = func(ctx context.Context, value Value) (finished bool, err error) {
		return jpChecker(value)
	}

	return
}

func NewRetryUntilCheckFromTemplate(expression string) (check RetryUntilCheck, err error) {
	tmplChecker, err := utils.CreateTemplateChecker(expression)
	if err != nil {
		return nil, err
	}

	check = func(ctx context.Context, value Value) (finished bool, err error) {
		return tmplChecker(value)
	}

	return
}

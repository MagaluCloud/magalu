package openapi

import (
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
)

type linkSpecResolver struct {
	httpResult         mgcHttpPkg.HttpResult
	findParameterValue func(location, name string) (core.Value, bool)
}

func (s *linkSpecResolver) resolve(value core.Value) (core.Value, bool, error) {
	switch specVal := value.(type) {
	case string:
		rtExp := linkRtExpression{specVal, s.httpResult, s.findParameterValue}
		return rtExp.resolve()
	default:
		// Treat as raw value
		return specVal, true, nil
	}
}

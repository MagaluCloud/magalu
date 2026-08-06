package openapi

import (
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"go.uber.org/zap"
)

// requestBodyCompatShim patches an assembled request body in place, adding
// fields that the operation schema does not declare.
//
// Shims exist so that a single CLI build keeps working against regions that
// already migrated to a new contract and regions that did not. They must only
// restate information the user already provided -- never introduce new
// semantics -- and every one of them is temporary: each entry below carries the
// condition for its own removal.
//
// See also updateBody() in operation.go, which rewrites bodies for IAM.
type requestBodyCompatShim func(logger *zap.SugaredLogger, body map[string]core.Value)

// Keyed by operationId. Every key here is checked against the embedded specs by
// TestRequestBodyCompatShimsExistInEmbeddedSpecs, so that a spec refresh that
// renames an operation fails the build instead of silently disabling the shim.
var requestBodyCompatShims = map[string]requestBodyCompatShim{
	"create_route_v1_vpcs__vpc_id__route_table_routes_post": addLegacyRoutePortID,
}

const legacyRouteTargetTypePort = "port_id"

// VPC route creation moved from "port_id" to "targets": {"id", "type"}, but not
// every region is migrated at once. Regions that are not migrated only
// understand "port_id" and discard "targets"; migrated regions only understand
// "targets" and discard "port_id". Sending both keeps one CLI build correct on
// both, and it is only possible when the target is a port -- any other target
// type has no legacy representation and must fail on old regions.
//
// The target value itself is never modified: whatever the user typed goes to
// the API as-is, and the API remains the authority on its accepted values.
//
// It stays inert until the spec declares "targets", so it can land before the
// contract change. Remove it once every region is migrated and "port_id" is
// gone from the API.
func addLegacyRoutePortID(logger *zap.SugaredLogger, body map[string]core.Value) {
	if _, ok := body["port_id"]; ok {
		return
	}

	target, ok := body["targets"].(map[string]any)
	if !ok {
		return
	}

	targetType, ok := target["type"].(string)
	if !ok || !strings.EqualFold(targetType, legacyRouteTargetTypePort) {
		return
	}

	id, ok := target["id"].(string)
	if !ok || id == "" {
		return
	}

	body["port_id"] = id
	logger.Debugw("legacy route compatibility: derived port_id from targets", "port_id", id)
}

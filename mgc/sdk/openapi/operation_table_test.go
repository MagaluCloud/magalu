package openapi

import (
	"fmt"
	"strings"
	"testing"

	"slices"
)

func operationTableToStrings(t *operationTable, prefix string) (lines []string) {
	for _, o := range t.childOperations {
		lines = append(lines, fmt.Sprintf("%s%s - %s %s", prefix, o.key, o.desc.method, o.desc.pathKey))
	}
	for _, child := range t.childTables {
		lines = append(lines, operationTableToStrings(child, prefix+child.name+" ")...)
	}
	return lines
}

func checkOperationTable(t *testing.T, operations []*operationDesc, expected []string) {
	slices.Sort(expected)

	table := newOperationTable("", operations)
	got := operationTableToStrings(table, "")
	slices.Sort(got)
	if !slices.Equal(expected, got) {
		e := strings.Join(expected, "\n")
		g := strings.Join(got, "\n")
		t.Errorf("diverging results:\nEXPECTED:\n%s\n\nGOT:\n%s\n", e, g)
	}
}

// BEGIN: Test block-storage.openapi.yaml resources

func Test_operationTree_block_storage_block_storage(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/volumes", method: "get"},
		{pathKey: "/v0/volumes", method: "post"},
		{pathKey: "/v0/volumes/{id}", method: "delete"},
		{pathKey: "/v0/volumes/{id}", method: "get"},
		{pathKey: "/v0/volumes/{id}", method: "patch"},
		{pathKey: "/v0/volumes/{id}/attach/{virtual_machine_id}", method: "post"},
		{pathKey: "/v0/volumes/{id}/detach/{virtual_machine_id}", method: "post"},
	}
	expected := []string{
		"list - get /v0/volumes",
		"create - post /v0/volumes",
		"delete - delete /v0/volumes/{id}",
		"get - get /v0/volumes/{id}",
		"update - patch /v0/volumes/{id}",
		"attach - post /v0/volumes/{id}/attach/{virtual_machine_id}",
		"detach - post /v0/volumes/{id}/detach/{virtual_machine_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_multi_version(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/volumes", method: "get"},
		{pathKey: "/v1/volumes", method: "post"},
		{pathKey: "/v0/volumes/{id}", method: "delete"},
		{pathKey: "/v0/volumes/{id}", method: "get"},
		{pathKey: "/v0/volumes/{id}", method: "patch"},
		{pathKey: "/v0/volumes/{id}/attach/{virtual_machine_id}", method: "post"},
		{pathKey: "/v0/volumes/{id}/detach/{virtual_machine_id}", method: "post"},
	}
	expected := []string{
		"list - get /v0/volumes",
		"create - post /v1/volumes",
		"delete - delete /v0/volumes/{id}",
		"get - get /v0/volumes/{id}",
		"update - patch /v0/volumes/{id}",
		"attach - post /v0/volumes/{id}/attach/{virtual_machine_id}",
		"detach - post /v0/volumes/{id}/detach/{virtual_machine_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_block_storage_snapshots(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/snapshots", method: "get"},
		{pathKey: "/v0/snapshots", method: "post"},
		{pathKey: "/v0/snapshots/{snapshot_id}", method: "delete"},
		{pathKey: "/v0/snapshots/{snapshot_id}", method: "post"},
	}
	expected := []string{
		"list - get /v0/snapshots",
		"create - post /v0/snapshots",
		"delete - delete /v0/snapshots/{snapshot_id}",
		"create-snapshot-id - post /v0/snapshots/{snapshot_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_block_storage_usage(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/usage", method: "get"},
	}
	expected := []string{
		"list - get /v0/usage",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_block_storage_volume_types(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/volume_types", method: "get"},
		{pathKey: "/v0/volume_types", method: "post"},
		{pathKey: "/v0/volume_types_all", method: "get"},
	}
	expected := []string{
		"list - get /v0/volume_types",
		"create - post /v0/volume_types",
		"list-all - get /v0/volume_types_all",
	}
	checkOperationTable(t, operations, expected)
}

// END: Test block-storage.openapi.yaml resources

// BEGIN: Test dbaas.openapi.yaml resources

func Test_operationTree_dbaas_backups(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v1/backups", method: "get"},
		{pathKey: "/v1/backups/{backup_id}", method: "delete"},
		{pathKey: "/v1/backups/{backup_id}", method: "get"},
	}
	expected := []string{
		"list - get /v1/backups",
		"delete - delete /v1/backups/{backup_id}",
		"get - get /v1/backups/{backup_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_dbaas_datastores(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v1/datastores", method: "get"},
		{pathKey: "/v1/datastores/{datastore_id}", method: "get"},
	}
	expected := []string{
		"list - get /v1/datastores",
		"get - get /v1/datastores/{datastore_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_dbaas_flavors(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v1/flavors", method: "get"},
		{pathKey: "/v1/flavors/{flavor_id}", method: "get"},
	}
	expected := []string{
		"list - get /v1/flavors",
		"get - get /v1/flavors/{flavor_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_dbaas_instances(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v1/instances", method: "get"},
		{pathKey: "/v1/instances", method: "post"},
		{pathKey: "/v1/instances/{id}", method: "delete"},
		{pathKey: "/v1/instances/{id}", method: "get"},
		{pathKey: "/v1/instances/{id}", method: "patch"},
		{pathKey: "/v1/instances/{id}/backups", method: "get"},
		{pathKey: "/v1/instances/{id}/backups", method: "post"},
		{pathKey: "/v1/instances/{id}/backups/{backup_id}", method: "delete"},
		{pathKey: "/v1/instances/{id}/backups/{backup_id}", method: "get"},
		{pathKey: "/v1/instances/{id}/resize", method: "post"},
		{pathKey: "/v1/instances/{id}/restores", method: "post"},
		{pathKey: "/v1/instances/{id}/start", method: "post"},
		{pathKey: "/v1/instances/{id}/stop", method: "post"},
	}
	expected := []string{
		"list - get /v1/instances",
		"create - post /v1/instances",
		"delete - delete /v1/instances/{id}",
		"get - get /v1/instances/{id}",
		"update - patch /v1/instances/{id}",
		"backups list - get /v1/instances/{id}/backups",
		"backups create - post /v1/instances/{id}/backups",
		"backups delete - delete /v1/instances/{id}/backups/{backup_id}",
		"backups get - get /v1/instances/{id}/backups/{backup_id}",
		"resize - post /v1/instances/{id}/resize",
		"restores - post /v1/instances/{id}/restores",
		"start - post /v1/instances/{id}/start",
		"stop - post /v1/instances/{id}/stop",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_dbaas_replicas(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v1/replicas", method: "get"},
		{pathKey: "/v1/replicas", method: "post"},
		{pathKey: "/v1/replicas/{replica_id}", method: "delete"},
		{pathKey: "/v1/replicas/{replica_id}", method: "get"},
		{pathKey: "/v1/replicas/{replica_id}/resize", method: "post"},
		{pathKey: "/v1/replicas/{replica_id}/start", method: "post"},
		{pathKey: "/v1/replicas/{replica_id}/stop", method: "post"},
	}
	expected := []string{
		"list - get /v1/replicas",
		"create - post /v1/replicas",
		"delete - delete /v1/replicas/{replica_id}",
		"get - get /v1/replicas/{replica_id}",
		"resize - post /v1/replicas/{replica_id}/resize",
		"start - post /v1/replicas/{replica_id}/start",
		"stop - post /v1/replicas/{replica_id}/stop",
	}
	checkOperationTable(t, operations, expected)
}

// END: Test dbaas.openapi.yaml resources

// BEGIN: Test mke.openapi.yaml resources

func Test_operationTree_mke_cluster(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/clusters", method: "get"},
		{pathKey: "/v0/clusters", method: "post"},
		{pathKey: "/v0/clusters/{cluster_id}", method: "delete"},
		{pathKey: "/v0/clusters/{cluster_id}", method: "get"},
		{pathKey: "/v0/clusters/{cluster_id}/kubeconfig", method: "get"},
	}
	expected := []string{
		"list - get /v0/clusters",
		"create - post /v0/clusters",
		"delete - delete /v0/clusters/{cluster_id}",
		"get - get /v0/clusters/{cluster_id}",
		"kubeconfig - get /v0/clusters/{cluster_id}/kubeconfig",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_mke_info(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/info/flavors", method: "get"},
		{pathKey: "/v0/info/versions", method: "get"},
	}
	expected := []string{
		"flavors - get /v0/info/flavors",
		"versions - get /v0/info/versions",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_mke_nodepool(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/clusters/{cluster_id}/node_pools", method: "get"},
		{pathKey: "/v0/clusters/{cluster_id}/node_pools", method: "post"},
		{pathKey: "/v0/clusters/{cluster_id}/node_pools/{node_pool_id}", method: "delete"},
		{pathKey: "/v0/clusters/{cluster_id}/node_pools/{node_pool_id}", method: "get"},
		{pathKey: "/v0/clusters/{cluster_id}/node_pools/{node_pool_id}", method: "patch"},
		{pathKey: "/v0/clusters/{cluster_id}/node_pools/{node_pool_id}/nodes", method: "get"},
	}
	expected := []string{
		"list - get /v0/clusters/{cluster_id}/node_pools",
		"create - post /v0/clusters/{cluster_id}/node_pools",
		"delete - delete /v0/clusters/{cluster_id}/node_pools/{node_pool_id}",
		"get - get /v0/clusters/{cluster_id}/node_pools/{node_pool_id}",
		"update - patch /v0/clusters/{cluster_id}/node_pools/{node_pool_id}",
		"nodes - get /v0/clusters/{cluster_id}/node_pools/{node_pool_id}/nodes",
	}
	checkOperationTable(t, operations, expected)
}

// END: Test mke.openapi.yaml resources

// BEGIN: Test network.openapi.yaml resources

func Test_operationTree_network_port(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/ports", method: "get"},
		{pathKey: "/v0/ports", method: "post"},
		{pathKey: "/v0/ports/all", method: "delete"},
		{pathKey: "/v0/ports/{port_id}", method: "delete"},
		{pathKey: "/v0/ports/{port_id}", method: "get"},
		{pathKey: "/v0/ports/{port_id}/attach/{security_group_id}", method: "post"},
		{pathKey: "/v0/ports/{port_id}/detach/{security_group_id}", method: "post"},
		{pathKey: "/v0/vpcs/{vpc_id}/ports", method: "get"},
		{pathKey: "/v0/vpcs/{vpc_id}/ports", method: "post"},
	}
	expected := []string{
		"ports list - get /v0/ports",
		"ports create - post /v0/ports",
		"ports delete-all - delete /v0/ports/all",
		"ports delete - delete /v0/ports/{port_id}",
		"ports get - get /v0/ports/{port_id}",
		"ports attach - post /v0/ports/{port_id}/attach/{security_group_id}",
		"ports detach - post /v0/ports/{port_id}/detach/{security_group_id}",
		"vpcs-ports list - get /v0/vpcs/{vpc_id}/ports",
		"vpcs-ports create - post /v0/vpcs/{vpc_id}/ports",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_public_ip(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/public_ips", method: "get"},
		{pathKey: "/v0/public_ips/{public_ip_id}", method: "delete"},
		{pathKey: "/v0/public_ips/{public_ip_id}", method: "get"},
		{pathKey: "/v0/public_ips/{public_ip_id}/attach/{port_id}", method: "post"},
		{pathKey: "/v0/public_ips/{public_ip_id}/detach/{port_id}", method: "post"},
		{pathKey: "/v0/vpcs/{vpc_id}/public_ips", method: "get"},
		{pathKey: "/v0/vpcs/{vpc_id}/public_ips", method: "post"},
	}
	expected := []string{
		"public-ips list - get /v0/public_ips",
		"public-ips delete - delete /v0/public_ips/{public_ip_id}",
		"public-ips get - get /v0/public_ips/{public_ip_id}",
		"public-ips attach - post /v0/public_ips/{public_ip_id}/attach/{port_id}",
		"public-ips detach - post /v0/public_ips/{public_ip_id}/detach/{port_id}",
		"vpcs-public-ips list - get /v0/vpcs/{vpc_id}/public_ips",
		"vpcs-public-ips create - post /v0/vpcs/{vpc_id}/public_ips",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_quotas(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/usage", method: "get"},
	}
	expected := []string{
		"list - get /v0/usage",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_router(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/routers", method: "get"},
		{pathKey: "/v0/routers/default", method: "post"},
	}
	expected := []string{
		"list - get /v0/routers",
		"create-default - post /v0/routers/default",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_rule(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/rules/{rule_id}", method: "delete"},
		{pathKey: "/v0/rules/{rule_id}", method: "get"},
		{pathKey: "/v0/security_groups/{security_group_id}/rules", method: "get"},
		{pathKey: "/v0/security_groups/{security_group_id}/rules", method: "post"},
	}
	expected := []string{
		"rules delete - delete /v0/rules/{rule_id}",
		"rules get - get /v0/rules/{rule_id}",
		"security-groups-rules list - get /v0/security_groups/{security_group_id}/rules",
		"security-groups-rules create - post /v0/security_groups/{security_group_id}/rules",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_security_group(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/security_groups", method: "get"},
		{pathKey: "/v0/security_groups", method: "post"},
		{pathKey: "/v0/security_groups/default", method: "post"},
		{pathKey: "/v0/security_groups/{security_group_id}", method: "delete"},
		{pathKey: "/v0/security_groups/{security_group_id}", method: "get"},
		{pathKey: "/v0/security_groups_all", method: "delete"},
		{pathKey: "/v0/vpcs/{vpc_id}/security_groups", method: "get"},
		{pathKey: "/v0/vpcs/{vpc_id}/security_groups", method: "post"},
	}
	expected := []string{
		"security-groups list - get /v0/security_groups",
		"security-groups create - post /v0/security_groups",
		"security-groups create-default - post /v0/security_groups/default",
		"security-groups delete - delete /v0/security_groups/{security_group_id}",
		"security-groups get - get /v0/security_groups/{security_group_id}",
		"security-groups delete-all - delete /v0/security_groups_all",
		"vpcs-security-groups list - get /v0/vpcs/{vpc_id}/security_groups",
		"vpcs-security-groups create - post /v0/vpcs/{vpc_id}/security_groups",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_subnets(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/subnets/{subnet_id}", method: "delete"},
		{pathKey: "/v0/subnets/{subnet_id}", method: "get"},
		{pathKey: "/v0/subnets/{subnet_id}", method: "patch"},
		{pathKey: "/v0/vpcs/{vpc_id}/subnets", method: "get"},
		{pathKey: "/v0/vpcs/{vpc_id}/subnets", method: "post"},
	}
	expected := []string{
		"subnets delete - delete /v0/subnets/{subnet_id}",
		"subnets get - get /v0/subnets/{subnet_id}",
		"subnets update - patch /v0/subnets/{subnet_id}",
		"vpcs-subnets list - get /v0/vpcs/{vpc_id}/subnets",
		"vpcs-subnets create - post /v0/vpcs/{vpc_id}/subnets",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_network_vpc(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/vpcs", method: "get"},
		{pathKey: "/v0/vpcs", method: "post"},
		{pathKey: "/v0/vpcs/all", method: "delete"},
		{pathKey: "/v0/vpcs/default", method: "post"},
		{pathKey: "/v0/vpcs/{vpc_id}", method: "delete"},
		{pathKey: "/v0/vpcs/{vpc_id}", method: "get"},
	}
	expected := []string{
		"list - get /v0/vpcs",
		"create - post /v0/vpcs",
		"delete-all - delete /v0/vpcs/all",
		"create-default - post /v0/vpcs/default",
		"delete - delete /v0/vpcs/{vpc_id}",
		"get - get /v0/vpcs/{vpc_id}",
	}
	checkOperationTable(t, operations, expected)
}

// END: Test network.openapi.yaml resources

// BEGIN: Test virtual-machine.openapi.yaml resources

func Test_operationTree_virtual_machine_images(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/images", method: "get"},
		{pathKey: "/v0/images", method: "post"},
		{pathKey: "/v0/images/{image_id}", method: "delete"},
	}
	expected := []string{
		"list - get /v0/images",
		"create - post /v0/images",
		"delete - delete /v0/images/{image_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_virtual_machine_instance_types(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/instance_types", method: "get"},
		{pathKey: "/v0/instance_types", method: "post"},
		{pathKey: "/v0/instance_types/{instance_type_id}", method: "delete"},
		{pathKey: "/v0/instance_types_all", method: "get"},
	}
	expected := []string{
		"list - get /v0/instance_types",
		"create - post /v0/instance_types",
		"delete - delete /v0/instance_types/{instance_type_id}",
		"list-all - get /v0/instance_types_all",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_virtual_machine_instances(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/instances", method: "get"},
		{pathKey: "/v0/instances", method: "post"},
		{pathKey: "/v0/instances/{id}", method: "delete"},
		{pathKey: "/v0/instances/{id}", method: "get"},
		{pathKey: "/v0/instances/{id}", method: "patch"},
		{pathKey: "/v0/instances/{id}/events", method: "get"},
		{pathKey: "/v0/instances/{id}/events/{event_id}", method: "get"},
	}
	expected := []string{
		"list - get /v0/instances",
		"create - post /v0/instances",
		"delete - delete /v0/instances/{id}",
		"get - get /v0/instances/{id}",
		"update - patch /v0/instances/{id}",
		"events list - get /v0/instances/{id}/events",
		"events get - get /v0/instances/{id}/events/{event_id}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_virtual_machine_keypairs(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/keypairs", method: "get"},
		{pathKey: "/v0/keypairs", method: "post"},
		{pathKey: "/v0/keypairs/{keypair_name}", method: "delete"},
		{pathKey: "/v0/keypairs/{keypair_name}", method: "post"},
	}
	expected := []string{
		"list - get /v0/keypairs",
		"create - post /v0/keypairs",
		"delete - delete /v0/keypairs/{keypair_name}",
		"create-keypair-name - post /v0/keypairs/{keypair_name}",
	}
	checkOperationTable(t, operations, expected)
}

func Test_operationTree_virtual_machine_usage(t *testing.T) {
	operations := []*operationDesc{
		{pathKey: "/v0/usage", method: "get"},
	}
	expected := []string{
		"list - get /v0/usage",
	}
	checkOperationTable(t, operations, expected)
}

// END: Test virtual-machine.openapi.yaml resources

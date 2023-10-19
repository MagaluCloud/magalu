package provider

type logKey string

func (c logKey) String() string {
	return string(c)
}

var (
	rpcField                = logKey("rpc-call")
	resourceNameField       = logKey("resource-name")
	actionResourceNameField = logKey("action-resource-name")
)

var (
	schemaGenSubsystem = logKey("schema-generation")
)

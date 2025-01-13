package generator

import (
	"errors"
	"os"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

const (
	dirMode         = 0755
	fileMode        = 0644
	createFileFlags = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	backupFmt       = "%s~"
)

var (
	errNotDir        = errors.New("not a directory")
	errMissingSchema = errors.New("missing schema")
)

var (
	nullSchema = mgcSchemaPkg.NewNullSchema()
	anySchema  = mgcSchemaPkg.NewAnySchema()
)

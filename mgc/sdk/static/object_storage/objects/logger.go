package objects

import mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"

var logger = mgcLoggerPkg.NewLazy[downloadAllObjectsParams]()

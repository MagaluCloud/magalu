package core

import mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"

var logger = mgcLoggerPkg.NewLazy[Descriptor]()

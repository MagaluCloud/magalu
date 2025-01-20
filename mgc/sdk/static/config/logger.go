package config

import mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"

var logger = mgcLoggerPkg.NewLazy[configSetParams]()

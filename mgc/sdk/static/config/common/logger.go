package common

import mgcLoggerPkg "magalu.cloud/core/logger"

type pkgSymbol int

var logger = mgcLoggerPkg.NewLazy[pkgSymbol]()

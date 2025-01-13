package common

import mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"

type pkgSymbol int

var logger = mgcLoggerPkg.NewLazy[pkgSymbol]()

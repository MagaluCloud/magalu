package objects

import (
	"go.uber.org/zap"
	mgcLoggerPkg "magalu.cloud/core/logger"
)

type pkgSymbol struct{}

var pkgLogger *zap.SugaredLogger

func logger() *zap.SugaredLogger {
	if pkgLogger == nil {
		pkgLogger = mgcLoggerPkg.New[pkgSymbol]()
	}
	return pkgLogger
}

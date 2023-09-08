package objects

import (
	"go.uber.org/zap"
	corelogger "magalu.cloud/core/logger"
)

type pkgSymbol struct{}

var pkgLogger *zap.SugaredLogger

func logger() *zap.SugaredLogger {
	if pkgLogger == nil {
		pkgLogger = corelogger.New[pkgSymbol]()
	}
	return pkgLogger
}

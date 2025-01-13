package logger

import (
	"reflect"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"go.uber.org/zap"
)

var rootLogger *zap.SugaredLogger

func NewLazy[T any]() func() *zap.SugaredLogger {
	return utils.NewLazyLoader(func() *zap.SugaredLogger {
		return New[T]()
	})
}

func New[T any]() *zap.SugaredLogger {
	return Root().Named(reflect.TypeOf(new(T)).Elem().PkgPath())
}

func Root() *zap.SugaredLogger {
	if rootLogger == nil {
		rootLogger = zap.Must(zap.NewProduction()).Sugar()
	}
	return rootLogger
}

func SetRoot(logger *zap.SugaredLogger) {
	rootLogger = logger
}

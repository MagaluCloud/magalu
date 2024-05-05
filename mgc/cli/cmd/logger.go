package cmd

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	mgcLoggerPkg "magalu.cloud/core/logger"
	mgcSdk "magalu.cloud/sdk"
	"moul.io/zapfilter"
)

var logger = mgcLoggerPkg.NewLazy[osArgParser]()

func newLogConfig() zap.Config {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)                 // it's widely used, zapfilter will default to "warn+:*"
	zapConfig.Encoding = "console"                                         // after all, it's a CLI
	zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()            // use the development encoder config
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // use colored level names
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // use ISO8601 format for timestamps
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder // use string duration encoder
	zapConfig.EncoderConfig.CallerKey = ""                                 // do not show file:line
	zapConfig.EncoderConfig.TimeKey = ""                                   // do not show timestamp by default
	return zapConfig
}

func initLogger(sdk *mgcSdk.Sdk, filterRules string) error {
	zapConfig := newLogConfig()

	if err := sdk.Config().Get(loggerConfigKey, &zapConfig); err != nil {
		return fmt.Errorf("unable to get logger configuration: %w", err)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return fmt.Errorf(
			"unable to build logger with current configuration: %w\nTo fix this, you'll need to alter the configuration file manually: %s",
			err,
			sdk.Config().FilePath(),
		)
	}

	filterOpt := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapfilter.NewFilteringCore(c, zapfilter.MustParseRules(filterRules))
	})

	logger = logger.WithOptions(filterOpt)
	mgcLoggerPkg.SetRoot(logger.Sugar())

	return nil
}

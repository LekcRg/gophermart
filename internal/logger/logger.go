package logger

import (
	"fmt"

	"github.com/LekcRg/gophermart/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(cfg config.Config) {
	var zapCfg zap.Config
	if cfg.IsDev {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	zl, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	Log = zl

	cfgString := fmt.Sprintf("%+v", cfg)
	Log.Info(cfgString)
}

package log

import (
	"github.com/pingcap/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(level string) *zap.Logger {
	cfg := &log.Config{
		Level: level,
	}
	gl, props, err := log.InitLogger(cfg, zap.AddStacktrace(zapcore.DPanicLevel))
	if err != nil {
		panic(err)
	}
	log.ReplaceGlobals(gl, props)
	return gl
}

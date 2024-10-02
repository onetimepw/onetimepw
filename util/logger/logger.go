package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func NewLogger(env string) (*zap.Logger, error) {
	var ws []zapcore.WriteSyncer

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	ws = append(ws, zapcore.AddSync(os.Stdout))

	level := zap.InfoLevel
	if env == "dev" || env == "local" {
		level = zap.DebugLevel
	}

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			zapcore.NewMultiWriteSyncer(ws...),
			level,
		),
	)

	zap.ReplaceGlobals(logger)

	return logger, nil
}

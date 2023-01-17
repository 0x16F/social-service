package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func HumanizedTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2001-02-03 12:04:05"))
}

func NewLogger() *Logger {
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(fmt.Sprintf("logs/%s.log", time.Now().Format("2001-02-03")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = HumanizedTimeEncoder

	encoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(file), config.Level),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.Level),
	)

	return &Logger{zap.New(core, zap.AddCaller()).Sugar()}
}

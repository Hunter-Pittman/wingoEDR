package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
Logging level reference
    DebugLevel: are usually present only on development environments.
    InfoLevel: default logging priority.
    WarnLevel: more important than InfoLevel, but still doesn't need individual human attention.
    ErrorLevel: these are high-priority and shouldn't be present in the application.
    DPanicLevel: these are particularly important errors and in the development environment logger will panic.
    PanicLevel: logs a message, then panics.
    FatalLevel: logs a message, then calls os.Exit(1).
*/

func InitLogger() {
	createLogDirectory()
	writerSync := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
}

func createLogDirectory() {
	path, _ := os.Getwd()

	if _, err := os.Stat(fmt.Sprintf("%s\\logs", path)); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(path+"\\logs\\wingoEDR_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05z0700"))
	})

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

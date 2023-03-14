package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// Add a call to the SIEM if one is available
type RemoteLogger struct {
	url string
}

type logMessage struct {
	Message string `json:"message"`
	Level   string `json:"level"`
}

func (rl *RemoteLogger) Write(p []byte) (n int, err error) {
	// Marshal the log message as JSON
	message := &logMessage{Message: string(p), Level: "error"}
	body, err := json.Marshal(message)
	if err != nil {
		return 0, err
	}

	// Send a POST request to the remote server with the log message as the body
	resp, err := http.Post(rl.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// Discard the response body
	_, _ = ioutil.ReadAll(resp.Body)
	return len(p), nil
}

func InitLogger() {
	createLogDirectory()
	writerSync := getLogWriter()
	encoder := getEncoder()

	//siemUrl := config.GetSiemUrl()

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		//zapcore.NewCore(encoder, zapcore.AddSync(&RemoteLogger{url: siemUrl}), zapcore.ErrorLevel),
	)
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

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

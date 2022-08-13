package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	logFileJSON *os.File
}

func NewLogger(defaultLogLevel zapcore.Level) (logger *Logger, err error) {
	logger = &Logger{}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	//logger.logFileJSON, err = os.OpenFile("./logs/log-json.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	return
	//}

	core := logger.getLoggerTee(config, defaultLogLevel)
	logger.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return
}

func (logger Logger) getLoggerTee(config zapcore.EncoderConfig, defaultLogLevel zapcore.Level) (core zapcore.Core) {
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	if logger.logFileJSON == nil {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		)
		return
	}

	//fileEncoderJSON := zapcore.NewJSONEncoder(config)
	//writer := zapcore.AddSync(logger.logFileJSON)
	core = zapcore.NewTee(
		//zapcore.NewCore(fileEncoderJSON, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return
}

func (logger Logger) Close() error {
	//return logger.logFileJSON.Close()
	return nil
}

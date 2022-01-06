package log

import (
	"go.uber.org/zap"
)

var logger Logger = zap.NewExample()

func GetLog() Logger {
	return logger
}

func L() Logger {
	return GetLog()
}

func SetLog(l Logger) {
	logger = l
}

func Field(key string, data string) zap.Field {
	return zap.String(key, data)
}

func Debug(msg string, fields ...zap.Field) {
	GetLog().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetLog().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLog().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLog().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	GetLog().DPanic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLog().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	GetLog().Panic(msg, fields...)
}

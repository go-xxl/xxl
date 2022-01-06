package main

import (
	"fmt"
	"github.com/go-xxl/xxl/utils"
	"go.uber.org/zap"
)

type CLog struct {
}

func (c *CLog) Debug(msg string, fields ...zap.Field) {
	fmt.Println("Debug", msg)
}
func (c *CLog) Info(msg string, fields ...zap.Field) {
	ToString("Info", msg, fields...)
}
func (c *CLog) Warn(msg string, fields ...zap.Field) {
	ToString("Warn", msg, fields...)
}
func (c *CLog) Error(msg string, fields ...zap.Field) {
	ToString("Error", msg, fields...)
}
func (c *CLog) DPanic(msg string, fields ...zap.Field) {
	ToString("DPanic", msg, fields...)
}
func (c *CLog) Fatal(msg string, fields ...zap.Field) {
	ToString("Fatal", msg, fields...)
}
func (c *CLog) Panic(msg string, fields ...zap.Field) {
	ToString("Panic", msg, fields...)
}

func ToString(level string, msg string, fields ...zap.Field) {
	var d string

	for _, field := range fields {
		d = d + "  || " + utils.ObjToStr(field.Interface)
	}
	fmt.Println(level, msg, d)
}

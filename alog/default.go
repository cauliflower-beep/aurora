package alog

import (
	"aurora/aiface"
	"context"
	"fmt"
)

var zLogInstance aiface.ILogger = new(zinxDefaultLog)

type zinxDefaultLog struct{}

func (log *zinxDefaultLog) InfoF(format string, v ...interface{}) {
	StdZinxLog.Infof(format, v...)
}

func (log *zinxDefaultLog) ErrorF(format string, v ...interface{}) {
	StdZinxLog.Errorf(format, v...)
}

func (log *zinxDefaultLog) DebugF(format string, v ...interface{}) {
	StdZinxLog.Debugf(format, v...)
}

func (log *zinxDefaultLog) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdZinxLog.Infof(format, v...)
}

func (log *zinxDefaultLog) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdZinxLog.Errorf(format, v...)
}

func (log *zinxDefaultLog) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdZinxLog.Debugf(format, v...)
}

func SetLogger(newlog aiface.ILogger) {
	zLogInstance = newlog
}

func Ins() aiface.ILogger {
	return zLogInstance
}

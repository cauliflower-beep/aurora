package alog

import (
	"aurora/aiface"
	"context"
	"fmt"
)

var aLogInstance aiface.ILogger = new(auroraDefaultLog)

type auroraDefaultLog struct{}

func (log *auroraDefaultLog) InfoF(format string, v ...interface{}) {
	StdAuroraLog.Infof(format, v...)
}

func (log *auroraDefaultLog) ErrorF(format string, v ...interface{}) {
	StdAuroraLog.Errorf(format, v...)
}

func (log *auroraDefaultLog) DebugF(format string, v ...interface{}) {
	StdAuroraLog.Debugf(format, v...)
}

func (log *auroraDefaultLog) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdAuroraLog.Infof(format, v...)
}

func (log *auroraDefaultLog) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdAuroraLog.Errorf(format, v...)
}

func (log *auroraDefaultLog) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	StdAuroraLog.Debugf(format, v...)
}

func SetLogger(newlog aiface.ILogger) {
	aLogInstance = newlog
}

func Ins() aiface.ILogger {
	return aLogInstance
}

package alog_test

import (
	"aurora/alog"
	"testing"
)

func TestStdZLog(t *testing.T) {

	//测试 默认debug输出
	alog.Debug("zinx debug content1")
	alog.Debug("zinx debug content2")

	alog.Debugf(" zinx debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	alog.ResetFlags(alog.BitDate | alog.BitLongFile | alog.BitLevel)
	alog.Info("zinx info content")

	//设置日志前缀，主要标记当前日志模块
	alog.SetPrefix("MODULE")
	alog.Error("zinx error content")

	//添加标记位
	alog.AddFlag(alog.BitShortFile | alog.BitTime)
	alog.Stack(" Zinx Stack! ")

	//设置日志写入文件
	alog.SetLogFile("./log", "testfile.log")
	alog.Debug("===> zinx debug content ~~666")
	alog.Debug("===> zinx debug content ~~888")
	alog.Error("===> zinx Error!!!! ~~~555~~~")

	//关闭debug调试
	alog.CloseDebug()
	alog.Debug("===> 我不应该出现~！")
	alog.Debug("===> 我不应该出现~！")
	alog.Error("===> zinx Error  after debug close !!!!")
}

func TestZLogger(t *testing.T) {
}

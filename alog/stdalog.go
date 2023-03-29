// Package alog 主要提供Aurora相关日志记录接口
// 包括:
//		stdzlog模块， 提供全局日志方法
//		zlogger模块,  日志内部定义协议，均为对象类方法
//
// 当前文件描述:
// @Title  stdalog.go
// @Description    包裹zlogger日志方法，提供全局方法
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package alog

/*
   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdAuroraLog
   注意：本文件方法不支持自定义，无法替换日志记录模式，如果需要自定义Logger:

   请使用如下方法:
   alog.SetLogger(yourLogger)
   alog.Ins().InfoF()等方法
*/

import "os"

//StdAuroraLog 创建全局标准log
var StdAuroraLog = NewAuroraLog(os.Stderr, "", BitDefault)

//Flags 获取StdAuroraLog 标记位
func Flags() int {
	return StdAuroraLog.Flags()
}

//ResetFlags 设置StdAuroraLog标记位
func ResetFlags(flag int) {
	StdAuroraLog.ResetFlags(flag)
}

//AddFlag 添加flag标记
func AddFlag(flag int) {
	StdAuroraLog.AddFlag(flag)
}

//SetPrefix 设置StdAuroraLog 日志头前缀
func SetPrefix(prefix string) {
	StdAuroraLog.SetPrefix(prefix)
}

//SetLogFile 设置StdAuroraLog绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdAuroraLog.SetLogFile(fileDir, fileName)
}

//CloseDebug 设置关闭debug
func CloseDebug() {
	StdAuroraLog.CloseDebug()
}

//OpenDebug 设置打开debug
func OpenDebug() {
	StdAuroraLog.OpenDebug()
}

//Debugf ====> Debug <====
func Debugf(format string, v ...interface{}) {
	StdAuroraLog.Debugf(format, v...)
}

//Debug Debug
func Debug(v ...interface{}) {
	StdAuroraLog.Debug(v...)
}

//Infof ====> Info <====
func Infof(format string, v ...interface{}) {
	StdAuroraLog.Infof(format, v...)
}

//Info -
func Info(v ...interface{}) {
	StdAuroraLog.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	StdAuroraLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdAuroraLog.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	StdAuroraLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdAuroraLog.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	StdAuroraLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdAuroraLog.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	StdAuroraLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdAuroraLog.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	StdAuroraLog.Stack(v...)
}

func init() {
	//因为StdAuroraLog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的AuroraLogger对象 calldDepth=2, StdAuroraLog的calldDepth=3
	StdAuroraLog.calldDepth = 3
}

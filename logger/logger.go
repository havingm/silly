package logger

import (
	"fmt"
	"runtime/debug"
)

//错误
func Error(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))
	fmt.Println(string(debug.Stack()))
}

//警报
func Alert(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))
}

//警告
func Warning(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))

}

//信息
func Info(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))
}

//调试
func Debug(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))
}

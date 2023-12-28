package utils

import (
	"silly/logger"
	"syscall"
	"unsafe"
)

var _WinFuncSetConsoleTitleW uintptr

func init() {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		logger.Info("load kernel32.dll failed: ", err)
		return
	}
	defer syscall.FreeLibrary(kernel32)
	_WinFuncSetConsoleTitleW, _ = syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
}
func SetConsoleTitle(title string) int {
	ret, _, err := syscall.Syscall(_WinFuncSetConsoleTitleW, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	if err != 0 {
		logger.Info("SetConsoleTitle failed, errCode: ", err)
	}
	return int(ret)
}

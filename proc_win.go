//go:build windows
// +build windows

package sorbet

import (
	"syscall"
	"unsafe"
)

var (
	user31                  = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow = user31.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user31.NewProc("GetWindowTextW")
)

func getForegroundWindow() syscall.Handle {
	hwnd, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(hwnd)
}

func getWindowText(hwnd syscall.Handle) string {
	var buf [255]uint16
	procGetWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	return syscall.UTF16ToString(buf[:])
}

func GetTitle() string {
	hwnd := getForegroundWindow()
	title := getWindowText(hwnd)
	return title
}

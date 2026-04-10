// Package ink 实现终端 UI 渲染循环和实例管理
package ink

import (
	"os"
)

// IsTerminal 检查文件描述符是否为终端 (TTY)
func IsTerminal(file *os.File) bool {
	if file == nil {
		return false
	}

	// 获取文件信息
	fi, err := file.Stat()
	if err != nil {
		return false
	}

	// 检查是否为字符设备 (终端)
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// IsStdoutTerminal 检查 stdout 是否为终端
func IsStdoutTerminal() bool {
	return IsTerminal(os.Stdout)
}

// IsStdinTerminal 检查 stdin 是否为终端
func IsStdinTerminal() bool {
	return IsTerminal(os.Stdin)
}

// IsStderrTerminal 检查 stderr 是否为终端
func IsStderrTerminal() bool {
	return IsTerminal(os.Stderr)
}

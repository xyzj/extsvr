package model

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	Prefix = "\033["
	Suffix = "m"
	Reset  = "\033[0m"

	// 文本样式
	StyleBold      = "1"
	StyleDim       = "2"
	StyleItalic    = "3"
	StyleUnderline = "4"
	StyleBlink     = "5"
	StyleReverse   = "7"

	// 前景色
	FgRed          = "31"
	FgGreen        = "32"
	FgYellow       = "33"
	FgBlue         = "34"
	FgMagenta      = "35"
	FgCyan         = "36"
	FgWhite        = "37"
	FgDefault      = "39"
	FgLightGrey    = "90"
	FgLightRed     = "91"
	FgLightGreen   = "92"
	FgLightYellow  = "93"
	FgLightBlue    = "94"
	FgLightMagenta = "95"
	FgLightCyan    = "96"
	FgLightWhite   = "97"

	// 背景色
	BgBlack        = "40"
	BgRed          = "41"
	BgGreen        = "42"
	BgYellow       = "43"
	BgBlue         = "44"
	BgMagenta      = "45"
	BgCyan         = "46"
	BgWhite        = "47"
	BgLightGrey    = "100"
	BgLightRed     = "101"
	BgLightGreen   = "102"
	BgLightYellow  = "103"
	BgLightBlue    = "104"
	BgLightMagenta = "105"
	BgLightCyan    = "106"
	BgLightWhite   = "107"
	BgDefault      = "49"
)

// Colorize 提供一种动态拼装颜色的快捷方法
func Colorize(text string, codes ...string) string {
	if len(codes) == 0 {
		return text
	}
	// 将多个代码用分号连接，例如 "1;31;40"
	var attrs strings.Builder
	for i, code := range codes {
		if i > 0 {
			attrs.WriteString(";")
		}
		attrs.WriteString(code)
	}
	return Prefix + attrs.String() + Suffix + text + Reset
}

type ProcessInfo struct {
	Name    string
	CmdLine string
	Pid     int
}

// ProcessExist only for linux
func ProcessExist(pid int) bool { // 发送信号 0，检查进程是否存在
	err := syscall.Kill(pid, 0)

	if err == nil {
		return true // 进程存在且你有权限
	}

	if err == syscall.EPERM {
		return true // 进程存在，但你没权限（说明它肯定活着）
	}

	if err == syscall.ESRCH {
		return false // 进程不存在
	}

	return false
}

// QueryProcess only for linux
func QueryProcess(name string) []*ProcessInfo {
	pi := make([]*ProcessInfo, 0)
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return pi
	}
	for _, proc := range procs {
		if !proc.IsDir() {
			continue
		}
		pid, _ := strconv.ParseInt(proc.Name(), 10, 32)
		if pid == 0 {
			continue
		}
		cmd, _ := os.ReadFile("/proc/" + proc.Name() + "/cmdline")
		if len(cmd) == 0 {
			continue
		}
		cl := strings.Split(string(cmd), "\x00")
		if name != filepath.Base(cl[0]) {
			continue
		}
		pi = append(pi, &ProcessInfo{
			Name:    name,
			Pid:     int(pid),
			CmdLine: strings.Join(cl, " "),
		})
	}
	return pi
}

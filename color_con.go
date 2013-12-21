package ops

import (
	"fmt"
	"runtime"
)

func ConRed(s string, a ...interface{}) string {
	red := "\x1b[31m"
	clear := "\x1b[0m"

	switch runtime.GOOS {
	case "linux":
		fallthrough
	case "darwin":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("%s%s%s", red, s, clear)
	case "windows":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("[!!!!!] %s", s)
	}
	return s
}

func ConGreen(s string, a ...interface{}) string {
	green := "\x1b[32m"
	clear := "\x1b[0m"

	switch runtime.GOOS {
	case "linux":
		fallthrough
	case "darwin":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("%s%s%s", green, s, clear)
	case "windows":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("[.....] %s", s)
	}
	return s
}

func ConYellow(s string, a ...interface{}) string {
	green := "\x1b[33m"
	clear := "\x1b[0m"

	switch runtime.GOOS {
	case "linux":
		fallthrough
	case "darwin":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("%s%s%s", green, s, clear)
	case "windows":
		s = fmt.Sprintf(s, a...)
		s = fmt.Sprintf("[-----] %s", s)
	}
	return s
}

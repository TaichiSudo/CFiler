package fileops

import (
	"os/exec"
	"runtime"
)

// OpenFile opens a file with the OS-associated application.
func OpenFile(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}

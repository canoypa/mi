package utils

import (
	"net/url"
	"os/exec"
	"runtime"
)

func OpenUrl(url url.URL) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url.String())
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url.String())
	default:
		cmd = exec.Command("xdg-open", url.String())
	}

	return cmd.Start()
}

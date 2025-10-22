package cmd

import (
	"os"
	"os/exec"
)

// run: 쉘 커맨드 실행 helper
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

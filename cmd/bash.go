//go:build !windows
// +build !windows

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func runCommand(path, name string, arg string, execStr string) (msg string, err error) {
	cmd := exec.Command(name, arg, execStr)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = path
	err = cmd.Run()
	msg = out.String()
	if err != nil {
		msg += stderr.String()
		if len(msg) == 0 {
			msg = fmt.Sprintf("%v", err)
		}
		err = errors.New(msg)
	}
	return
}

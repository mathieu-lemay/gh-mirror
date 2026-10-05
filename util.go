package main

import (
	"bytes"
	"context"
	"io"
	"os/exec"

	"github.com/cli/go-gh/v2"
)

// runCli runs the GitHub CLI, capturing stdOut and stdErr
func runCli(wo, we io.Writer, path *string, args ...string) error {
	ghExe, err := gh.Path()
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(context.Background(), ghExe, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if path != nil {
		cmd.Dir = *path
	}

	err = cmd.Run()

	if wo != nil {
		wo.Write(stdout.Bytes())
	}
	if we != nil {
		we.Write(stderr.Bytes())
	}

	return err
}

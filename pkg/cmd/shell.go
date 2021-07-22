package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type Shell struct {
	Path     string
	Template string
}

func NewShell() *Shell {
	return &Shell{
		Path:     "/bin/sh",
		Template: "%s",
	}
}

func (c *Shell) Exec(script string) (stdout []byte, stderr []byte, exitCode int, err error) {
	finalCmd := fmt.Sprintf(c.Template, script)
	cmd := exec.Command(c.Path, "-c", finalCmd)

	stdoutBuf := bytes.Buffer{}
	cmd.Stdout = &stdoutBuf

	stderrBuf := bytes.Buffer{}
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()

	stdout = stdoutBuf.Bytes()
	stderr = stderrBuf.Bytes()

	log.Printf(
		"[DEBUG] exec \"%s\"\nexit code: %d\nerr: %v\nstderr: %v\n",
		finalCmd,
		exitCode,
		err,
		string(stderr),
	)
	return
}

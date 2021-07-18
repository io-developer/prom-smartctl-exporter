package exporter

import (
	"fmt"
	"log"
	"os/exec"
)

type Shell struct {
	ShellPath string
	Template  string
}

func NewShell() *Shell {
	return &Shell{
		ShellPath: "/bin/sh",
		Template:  "%s",
	}
}

func (c *Shell) Exec(cmd string) ([]byte, error) {
	finalCmd := fmt.Sprintf(c.Template, cmd)
	out, err := exec.Command(c.ShellPath, "-c", finalCmd).CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] failed exec cmd \"%s\" (\"%s\"):\n%v\n%v\n", cmd, finalCmd, string(out), err)
	} else {
		log.Printf("[DEBUG] cmd \"%s\" (\"%s\"):\n%v\n", cmd, finalCmd, string(out))
	}
	return out, err
}

package cmd

import (
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

func (c *Shell) Exec(cmd string) ([]byte, error) {
	finalCmd := fmt.Sprintf(c.Template, cmd)
	out, err := exec.Command(c.Path, "-c", finalCmd).CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] failed exec \"%s\" (\"%s\"):\n%v\n%v\n", cmd, finalCmd, string(out), err)
	} else {
		log.Printf("[DEBUG] exec \"%s\" (\"%s\"):\n%v\n", cmd, finalCmd, string(out))
	}
	return out, err
}

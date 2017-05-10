package main

import (
	"os"
	"os/exec"
)

func nginxReload() error {
	cmd := exec.Command("/bin/sh", "-c", "sudo nginx -s reload")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

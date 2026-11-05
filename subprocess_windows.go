package main

import (
	"os/exec"
	"syscall"
)

func hideEngineWindow(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} }

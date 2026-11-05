//go:build !windows

package main

import "os/exec"

func hideEngineWindow(cmd *exec.Cmd) {}

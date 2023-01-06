package controllers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func CreateBackup(local_patch string, remote_patch string, key_patch string) {

	cmd := exec.Command("restic", "-r", remote_patch, "backup", local_patch, "--json", "-p", key_patch)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	defer cmd.Wait()
	cmd.Start()
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		current_status = scanner.Text()
		fmt.Println(current_status)
	}
}

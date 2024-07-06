package work

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type LocalCmd struct {
}

// pid,error
func (*LocalCmd) Start(shell string, opts ...WorkOptionsFunc) (int, error) {
	// 创建一个带有超时的上下文
	wopt := &WorkOptions{
		TimeOut: 30 * time.Second,
	}
	for _, optfunc := range opts {
		optfunc(wopt)
	}

	ctx, _ := context.WithTimeout(context.Background(), wopt.TimeOut)
	cmd := exec.CommandContext(ctx, "sh", "-c", shell)
	if wopt.SudoPassword != "" && wopt.Username != "root" {
		cmd = exec.CommandContext(ctx, "sudo", "-S", "sh", "-c", shell)
		cmd.Stdin = strings.NewReader(wopt.SudoPassword + "\n")
	}

	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var errFile, outFile *os.File
	var err error
	if wopt.ErrPath != "" {
		errFile, err = os.OpenFile(wopt.ErrPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stderr = errFile
	}

	if wopt.OutPath != "" {
		outFile, err = os.OpenFile(wopt.OutPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stdout = outFile
	}
	// 启动命令
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	// 等待命令完成或上下文取消
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
		if errFile != nil {
			errFile.Close()
		}
		if outFile != nil {
			outFile.Close()
		}
	}()
	return cmd.ProcessState.Pid(), nil
}

// state_code,output,errout,error, 默认超时时间120秒
func (w *LocalCmd) RunOutput(shell string, opts ...WorkOptionsFunc) (int, string, error) {
	// 创建一个带有超时的上下文
	wopt := &WorkOptions{
		TimeOut: 30 * time.Second,
	}
	for _, optfunc := range opts {
		optfunc(wopt)
	}

	ctx, cancel := context.WithTimeout(context.Background(), wopt.TimeOut)
	defer cancel()

	fmt.Println("RunOutput shell:", shell)
	cmd := exec.CommandContext(ctx, "sh", "-c", shell)
	if wopt.SudoPassword != "" && wopt.Username != "root" {
		cmd = exec.CommandContext(ctx, "sudo", "-S", "sh", "-c", shell)
		cmd.Stdin = strings.NewReader(wopt.SudoPassword + "\n")
	}
	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	// 启动命令
	if err := cmd.Start(); err != nil {
		return -1, strings.TrimSpace(stdout.String()), err
	}

	// 等待命令完成或上下文取消
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 上下文取消时终止所有相关进程
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), err
		}
		return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), WorkTimeoutErr
	case err := <-done:
		// 命令完成
		if err != nil {
			return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), err
		}
		return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), nil
	}
}

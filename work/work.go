package work

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	ShellNilErr    = errors.New("work shell is nil")
	WorkTimeoutErr = errors.New("work timeout")
)

type Work struct {
	Args WorkArgs
}

type WorkArgs struct {
	Shell    string
	TimeOut  time.Duration
	OutPath  string
	ErrPaht  string
	Username string
	Become   bool
}

func Get(wa WorkArgs) *Work {
	return &Work{
		Args: wa,
	}
}

// state_code,error
func (w *Work) Run() (int, error) {
	if w.Args.Shell == "" {
		return 0, ShellNilErr
	}

	// 创建一个带有超时的上下文
	timeout := 30 * time.Second
	if w.Args.TimeOut != 0 {
		timeout = w.Args.TimeOut
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	cmd := exec.CommandContext(ctx, "bash", "-c", w.Args.Shell)
	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var errFile, outFile *os.File
	var err error
	if w.Args.ErrPaht != "" {
		errFile, err = os.OpenFile(w.Args.ErrPaht, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stderr = errFile
		defer errFile.Close()
	}

	if w.Args.OutPath != "" {
		outFile, err = os.OpenFile(w.Args.OutPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stdout = outFile
		defer outFile.Close()
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return -1, err
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
			return -1, err
		}
		return -1, WorkTimeoutErr
	case err := <-done:
		// 命令完成
		if err != nil {
			return cmd.ProcessState.ExitCode(), err
		}
		return cmd.ProcessState.ExitCode(), nil
	}
}

// pid,error
func (w *Work) AsyncRun() (int, error) {
	// 创建一个带有超时的上下文
	timeout := 30 * time.Second
	if w.Args.TimeOut != 0 {
		timeout = w.Args.TimeOut
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, "bash", "-c", w.Args.Shell)
	if w.Args.Become {
		cmd = exec.CommandContext(ctx, "sudo", "-S", "bash", "-c", w.Args.Shell)
	}

	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var errFile, outFile *os.File
	var err error
	if w.Args.ErrPaht != "" {
		errFile, err = os.OpenFile(w.Args.ErrPaht, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stderr = errFile
	}

	if w.Args.OutPath != "" {
		outFile, err = os.OpenFile(w.Args.OutPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
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

	go func() {
		select {
		case <-ctx.Done():
			// 上下文取消时终止所有相关进程
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				return
			}
			return
		case err := <-done:
			// 命令完成
			if err != nil {
				return
			}
			return
		}
	}()
	return cmd.ProcessState.Pid(), nil
}

// state_code,output,errout,error, 默认超时时间120秒
func (w *Work) RunOutput() (int, string, string, error) {
	// 创建一个带有超时的上下文
	timeout := 120 * time.Second
	if w.Args.TimeOut != 0 {
		timeout = w.Args.TimeOut
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if !strings.Contains(w.Args.Shell, "#!/bin") {
		w.Args.Shell = "set -e \n" + w.Args.Shell
	}
	fmt.Println("RunOutput shell:", w.Args.Shell)
	cmd := exec.CommandContext(ctx, "bash", "-c", w.Args.Shell)
	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 启动命令
	if err := cmd.Start(); err != nil {
		return -1, strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
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
			return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
		}
		return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), WorkTimeoutErr
	case err := <-done:
		// 命令完成
		if err != nil {
			return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
		}
		return cmd.ProcessState.ExitCode(), strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), nil
	}
}

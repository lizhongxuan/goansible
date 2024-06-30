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
	Shell        string
	TimeOut      time.Duration
	OutPath      string
	ErrPath      string
	Username     string
	Become       bool
	SudoPassword string
	Stdin        string
}

func GetWork(shell string, opts ...WorkOptionsFunc) *Work {
	w := &Work{
		Shell:   shell,
		TimeOut: 60 * time.Second,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

type WorkOptionsFunc func(*Work)

func WithStdin(str string) WorkOptionsFunc {
	return func(w *Work) {
		w.Stdin = str
		return
	}
}

func WithTimeOut(t time.Duration) WorkOptionsFunc {
	return func(w *Work) {
		w.TimeOut = t
		return
	}
}

func WithOutPath(path string) WorkOptionsFunc {
	return func(w *Work) {
		w.OutPath = path
		return
	}
}

func WithErrPath(path string) WorkOptionsFunc {
	return func(w *Work) {
		w.ErrPath = path
		return
	}
}

func WithUsername(username string) WorkOptionsFunc {
	return func(w *Work) {
		w.Username = username
		return
	}
}

func WithBecome(become bool) WorkOptionsFunc {
	return func(w *Work) {
		w.Become = become
		return
	}
}

func WithSudoPassword(sudoPassword string) WorkOptionsFunc {
	return func(w *Work) {
		w.SudoPassword = sudoPassword
		return
	}
}

// pid,error
func (w *Work) AsyncRun() (int, error) {
	// 创建一个带有超时的上下文
	timeout := 30 * time.Second
	if w.TimeOut != 0 {
		timeout = w.TimeOut
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	cmd := exec.CommandContext(ctx, "bash", "-c", w.Shell)
	if w.Become && w.SudoPassword != "" && w.Username != "root" {
		cmd = exec.CommandContext(ctx, "sudo", "-S", "bash", "-c", w.Shell)
		cmd.Stdin = strings.NewReader(w.SudoPassword + "\n")
	}

	// 设置命令的进程组 ID
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var errFile, outFile *os.File
	var err error
	if w.ErrPath != "" {
		errFile, err = os.OpenFile(w.ErrPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		if err != nil {
			return 0, err
		}
		cmd.Stderr = errFile
	}

	if w.OutPath != "" {
		outFile, err = os.OpenFile(w.OutPath, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
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
	if w.TimeOut != 0 {
		timeout = w.TimeOut
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Println("RunOutput shell:", w.Shell)
	cmd := exec.CommandContext(ctx, "bash", "-c", w.Shell)
	if w.Become && w.SudoPassword != "" && w.Username != "root" {
		cmd = exec.CommandContext(ctx, "sudo", "-S", "bash", "-c", w.Shell)
		cmd.Stdin = strings.NewReader(w.SudoPassword + "\n")
	}
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

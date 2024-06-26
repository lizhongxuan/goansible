package work

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// CmdSession ...
type CmdSession struct {
	sess            *exec.Cmd
	ctx             context.Context
	onceStdinCloser sync.Once
	stdin           io.WriteCloser
	output          *bytes.Buffer
}

// Start starts a remote process in the current session
func (s *CmdSession) Start(cmd string, logFunc ...func(scanner *bufio.Scanner)) error {
	var stdout bytes.Buffer
	sess := exec.CommandContext(s.ctx, "sh", "-c", cmd)
	if runtime.GOOS == "windows" {
		sess = exec.CommandContext(s.ctx, "cmd", "/c", cmd)
	}
	sess.Stdout = &stdout
	sess.Stderr = &stdout
	s.output = &stdout
	s.sess = sess
	stdin, err := sess.StdinPipe()
	if err != nil {
		return err
	}
	s.stdin = stdin
	if len(logFunc) > 0 {
		stdout, err := sess.StdoutPipe()
		if err != nil {
			return err
		}
		go logFunc[0](bufio.NewScanner(stdout))
	}
	err = sess.Start()
	if len(logFunc) > 0 {
		time.Sleep(2 * time.Second)
	}
	return err
}

// Wait wait blocks until the remote process completes or is cancelled
func (s *CmdSession) Wait() error {
	return s.sess.Wait()
}

// Output ...
func (s *CmdSession) Output() string {
	return s.output.String()
}

// Stdin returns a pipe to the stdin of the remote process
func (s *CmdSession) Stdin() io.Writer {
	return s.stdin
}

// CloseStdin closes the stdin pipe of the remote process
func (s *CmdSession) CloseStdin() error {
	var err error
	s.onceStdinCloser.Do(func() {
		err = s.stdin.Close()
	})
	return err
}

// Close closes the current session
func (s *CmdSession) Close() error {
	err := s.CloseStdin()
	if err != nil {
		return fmt.Errorf("failed to close stdin: %s", err)
	}
	return nil
}

func newCmdSession(ctx context.Context) *CmdSession {
	return &CmdSession{ctx: ctx}
}

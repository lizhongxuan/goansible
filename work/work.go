package work

import (
	"errors"
	"time"
)

var (
	WorkTimeoutErr = errors.New("work timeout")
)

// Worker is an interface to run a command
type Worker interface {
	RunOutput(shell string, opts ...WorkOptionsFunc) (int, string, error)
	Start(shell string, opts ...WorkOptionsFunc) (int, error)
}

type WorkOptions struct {
	Worker       Worker
	TimeOut      time.Duration
	OutPath      string
	ErrPath      string
	Username     string
	SudoPassword string
	Stdin        string
}
type WorkOptionsFunc func(*WorkOptions)

func WithStdin(str string) WorkOptionsFunc {
	return func(w *WorkOptions) {
		w.Stdin = str
		return
	}
}

func WithTimeOut(t time.Duration) WorkOptionsFunc {
	return func(w *WorkOptions) {
		if t == 0 {
			return
		}
		w.TimeOut = t
		return
	}
}

func WithOutPath(path string) WorkOptionsFunc {
	return func(w *WorkOptions) {
		w.OutPath = path
		return
	}
}

func WithErrPath(path string) WorkOptionsFunc {
	return func(w *WorkOptions) {
		w.ErrPath = path
		return
	}
}

func WithUsername(username string) WorkOptionsFunc {
	return func(w *WorkOptions) {
		w.Username = username
		return
	}
}

func WithSudoPassword(sudoPassword string) WorkOptionsFunc {
	return func(w *WorkOptions) {
		w.SudoPassword = sudoPassword
		return
	}
}

func WithFunc(f func(*WorkOptions)) WorkOptionsFunc {
	return func(w *WorkOptions) {
		f(w)
		return
	}
}

package work

import (
	"io"
)

// Worker is an interface to run a command
type Worker interface {
	CombinedOutput() ([]byte, error)
	Environ() []string
	Output() ([]byte, error)
	Run() error
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	String() string
	Wait() error
}

type CmdWorker struct {
	Exec Worker
}

func NewCmdWorker(w Worker) *CmdWorker {
	return &CmdWorker{
		Exec: w,
	}
}

func (c *CmdWorker) CombinedOutput() ([]byte, error) {
	return c.Exec.CombinedOutput()
}

func (c *CmdWorker) Combined() ([]byte, error) {
	return c.Exec.CombinedOutput()
}
func (c *CmdWorker) Wait() error {
	return c.Exec.Wait()
}

func (c *CmdWorker) SetEnviron() []string {
	return c.Exec.Environ()
}
func (c *CmdWorker) Environ() []string {
	return c.Exec.Environ()
}

func (c *CmdWorker) Output() ([]byte, error) {
	return c.Exec.Output()
}

func (c *CmdWorker) Run() error {
	return c.Exec.Run()
}

func (c *CmdWorker) Start() error {
	return c.Exec.Start()
}

func (c *CmdWorker) StderrPipe() (io.ReadCloser, error) {
	return c.Exec.StderrPipe()
}

func (c *CmdWorker) StdinPipe() (io.WriteCloser, error) {
	return c.Exec.StdinPipe()
}

func (c *CmdWorker) StdoutPipe() (io.ReadCloser, error) {
	return c.Exec.StdoutPipe()
}

func (c *CmdWorker) String() string {
	return c.Exec.String()
}

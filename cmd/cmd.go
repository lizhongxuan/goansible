package cmd

import (
	"bufio"
	"io"
)

type Command interface {
}

type Session interface {
	Start(cmd string, logFunc ...func(scanner *bufio.Scanner)) error
	Stdin() io.Writer
	CloseStdin() error
	Wait() error
	Close() error
	Output() string
}

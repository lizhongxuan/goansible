package work

type SSHCmd struct {
	Url string
}

func (sc *SSHCmd) RunOutput(shell string, opts ...WorkOptionsFunc) (int, string, error) {
	return 0, "", nil
}

func (sc *SSHCmd) Start(shell string, opts ...WorkOptionsFunc) (int, error) {
	return 0, nil
}

package work

type ClientCmd struct {
	Url string
}

func (*ClientCmd) RunOutput(shell string, opts ...WorkOptionsFunc) (int, string, error) {
	return 0, "", nil
}

func (*ClientCmd) Start(shell string, opts ...WorkOptionsFunc) (int, error) {
	return 0, nil
}

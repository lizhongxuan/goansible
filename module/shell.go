package module

type ShellModule struct {
	Shell string `yaml:"shell"`
}

func (*ShellModule) StringShell(m Module, args map[string]interface{}) (string, error) {
	return Template(m.Shell.Shell, args)
}

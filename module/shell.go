package module

type ShellModule struct {
	Shell string `yaml:"shell"`
}

func (m *ShellModule) StringShell(args map[string]interface{}) (string, error) {
	return Template(m.Shell, args)
}

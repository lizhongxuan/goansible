package module

type ShellModule struct {
	Shell string `yaml:"shell"`
}

var _ ModuleInterface = &ShellModule{}

func (m *ShellModule) StringShell(args map[string]interface{}) (string, error) {
	return Template(m.Shell, args)
}
func (m *ShellModule) Show() string {
	return "Shell Module"
}

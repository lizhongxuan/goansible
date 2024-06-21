package module

type ShellModule string

func (cm *ShellModule) shellString(args map[string]interface{}) (string, error) {
	return Template(string(*cm), args)
}

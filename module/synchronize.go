package module

import (
	"fmt"
	"strings"
)

type SynchronizeModule struct {
	Src  string `yaml:"src,omitempty"`
	Dest string `yaml:"dest,omitempty"`
}

func (*SynchronizeModule) StringShell(m Module, args map[string]interface{}) (string, error) {
	src, err := Template(strings.TrimSpace(m.Copy.Src), args)
	if err != nil {
		return "", err
	}
	dest, err := Template(strings.TrimSpace(m.Copy.Dest), args)
	if err != nil {
		return "", err
	}
	mode := "a"
	if checkIsSSH(src) || checkIsSSH(dest) {
		mode = "az"
	}
	return fmt.Sprintf("rsync -%s %s %s", mode, src, dest), nil
}

func checkIsSSH(str string) bool {
	if strings.Count(str, "@") != 0 && strings.Count(str, ":") != 1 {
		return true
	}
	return false
}

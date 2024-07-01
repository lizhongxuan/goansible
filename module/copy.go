package module

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type CopyModule struct {
	Src  string      `yaml:"src,omitempty"`
	Dest string      `yaml:"dest,omitempty"`
	Mode fs.FileMode `yaml:"mode,omitempty"`
}

var _ ModuleInterface = &CopyModule{}

func (m *CopyModule) StringShell(args map[string]interface{}) (string, error) {
	src, err := Template(strings.TrimSpace(m.Src), args)
	if err != nil {
		return "", err
	}
	dest, err := Template(strings.TrimSpace(m.Dest), args)
	if err != nil {
		return "", err
	}

	cpArg := ""
	chmodArg := ""
	srcFile := ""
	if strings.HasSuffix(src, "/") {
		cpArg = "-r"
		chmodArg = "-R"
	} else {
		srcFile = filepath.Base(src)
	}
	chmodDest := dest
	if strings.HasSuffix(dest, "/") {
		chmodDest = dest + srcFile
	}
	return fmt.Sprintf("cp %s %s %s && chmod %s %s", cpArg, src, dest, chmodArg, chmodDest), nil
}

func (m *CopyModule) Show() string {
	return "Copy Module"
}

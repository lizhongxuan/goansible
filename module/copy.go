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

func (cm *CopyModule) shellString(args map[string]interface{}) (string, error) {
	src := strings.TrimSpace(cm.Src)
	src, err := Template(src, args)
	if err != nil {
		return "", err
	}
	dest := strings.TrimSpace(cm.Dest)
	dest, err = Template(dest, args)
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

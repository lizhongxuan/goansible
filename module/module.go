package module

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

type Module struct {
	Shell ShellModule `yaml:"shell,omitempty"`
	Copy  *CopyModule `yaml:"copy,omitempty"`
	File  *FileModule `yaml:"file,omitempty"`
}

func Template(str string, args map[string]interface{}) (string, error) {
	tmpl, err := template.New("task").Parse(str)
	if err != nil {
		return "", err
	}
	std := bytes.Buffer{}
	if err = tmpl.Execute(&std, args); err != nil {
		panic(err)
	}
	return std.String(), nil
}

func ModuleVerify(m Module) (string, error) {
	num := 0
	moduleName := ""
	if m.Shell != "" {
		num++
		moduleName = "shell"
	}
	if m.Copy != nil {
		num++
		moduleName = "copy"
	}
	if num != 1 {
		return "", errors.New("task module count invalid")
	}
	return moduleName, nil
}

func (m Module) ShellString(args map[string]interface{}) (string, error) {
	moduleName, err := ModuleVerify(m)
	if err != nil {
		return "", err
	}
	switch moduleName {
	case "copy":
		return m.Copy.shellString(args)
	case "shell":
		return m.Shell.shellString(args)
	default:
		return "", errors.New(fmt.Sprintf("task module_name:%s invalid.", moduleName))
	}
}

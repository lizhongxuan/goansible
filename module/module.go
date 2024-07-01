package module

import (
	"bytes"
	"html/template"
	"reflect"
)

type Module struct {
	Shell       *ShellModule       `yaml:",inline"`
	Copy        *CopyModule        `yaml:"copy,omitempty"`
	File        *FileModule        `yaml:"file,omitempty"`
	Synchronize *SynchronizeModule `yaml:"synchronize,omitempty"`
}

type ModuleInterface interface {
	StringShell(map[string]interface{}) (string, error)
	Show() string
}

func Find(m Module) ModuleInterface {
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsNil() {
			continue
		}
		if obj, ok := field.Interface().(ModuleInterface); ok {
			return obj
		}
	}
	return nil
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

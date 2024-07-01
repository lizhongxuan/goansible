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

var modules map[string]ModuleInterface

func init() {
	// key需要跟Module下的变量名一样
	modules = map[string]ModuleInterface{
		"Copy":        &CopyModule{},
		"Shell":       &ShellModule{},
		"Synchronize": &SynchronizeModule{},
	}
}

type ModuleInterface interface {
	StringShell(map[string]interface{}) (string, error)
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

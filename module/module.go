package module

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"reflect"
)

type Module struct {
	Shell *ShellModule `yaml:",inline"`
	Copy  *CopyModule  `yaml:"copy,omitempty"`
	File  *FileModule  `yaml:"file,omitempty"`
}

var modules map[string]ModuleInterface

func init() {
	// key需要跟Module下的变量名一样
	modules = map[string]ModuleInterface{
		"Copy":  &CopyModule{},
		"Shell": &ShellModule{},
	}
}

type ModuleInterface interface {
	StringShell(Module, map[string]interface{}) (string, error)
}

func FindAndVerify(m Module) string {
	v := reflect.ValueOf(m)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if field.IsNil() {
			continue
		}
		if _, ok := modules[fieldType.Name]; ok {
			return fieldType.Name
		}
	}
	return ""
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

func (m Module) ShellString(args map[string]interface{}) (string, error) {
	mi, ok := modules[FindAndVerify(m)]
	if !ok {
		return "", errors.New("task module invalid.")
	}
	log.Printf("Module: %+v", m)
	return mi.StringShell(m, args)
}

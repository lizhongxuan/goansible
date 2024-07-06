package checker

import (
	"fmt"
	"go-ansible/work"
	"regexp"
)

type Checker struct {
	Name   string
	CMD    string
	Output string
	Regex  string
	Hooks  map[string]string
}

var checkerPoor map[string]*Checker

func init() {
	if checkerPoor == nil {
		checkerPoor = make(map[string]*Checker)
	}
}

func (c *Checker) check() error {
	defer func() {
		for _, v := range c.Hooks {
			xxxx
		}
	}()

	output, err := w.CombinedOutput()
	if err != nil {
		return err
	}
	if c.Regex != "" {
		// 编译正则表达式
		re, err := regexp.Compile(c.Regex)
		if err != nil {
			return err
		}

		// 使用正则表达式进行匹配
		if !re.MatchString(string(output)) {
			return fmt.Errorf("fail match,regex:%s output:%s", c.Regex, string(output))
		}
	}
	return nil
}

func Check(w work.Worker, list []string) error {
	for _, v := range list {
		cmdkey := v
		c, ok := checkerPoor[cmdkey]
		if !ok {
			continue
		}
		if err := c.check(); err != nil {
			return err
		}
	}
	return nil
}

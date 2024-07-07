package checker

import (
	"fmt"
	"goansible/work"
	"regexp"
)

type Checker struct {
	Name         string   `yaml:"name"`
	Shell        string   `yaml:"cmd"`
	Output       string   `yaml:"output"`
	Regex        string   `yaml:"regex"`
	SuccessHooks []string `yaml:"success_hooks"`
	FailHooks    []string `yaml:"fail_hooks"`
	IgnoreErrors bool     `yaml:"ignore_errors"` // 是否忽略任务错误
	Work         work.Worker
}

var CheckerPoor map[string]*Checker

func init() {
	if CheckerPoor == nil {
		CheckerPoor = make(map[string]*Checker)
	}
}

func (c *Checker) check() error {
	isSuccess := false
	defer func() {
		hooks := make([]string, 0)
		if isSuccess {
			hooks = append(hooks, c.SuccessHooks...)
		} else {
			hooks = append(hooks, c.FailHooks...)
		}
		for _, v := range hooks {
			_, output, err := c.Work.RunOutput(v)
			if err != nil {
				fmt.Println("output:", output)
				return
			}
		}
	}()

	_, output, err := c.Work.RunOutput(c.Shell)
	if err != nil {
		if c.IgnoreErrors {
			return nil
		}
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
	isSuccess = true
	return nil
}

func Check(w work.Worker, list []string) error {
	for _, v := range list {
		cmdkey := v
		c, ok := CheckerPoor[cmdkey]
		if !ok {
			continue
		}
		c.Work = w
		if err := c.check(); err != nil {
			return err
		}
	}
	return nil
}
